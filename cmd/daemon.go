package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/bfirestone/speed-checker/internal/config"
	"github.com/bfirestone/speed-checker/internal/daemon"
	"github.com/bfirestone/speed-checker/internal/database"
	"github.com/bfirestone/speed-checker/internal/services"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run only the testing daemon (background tests)",
	Long: `Runs only the background testing daemon including:

• Scheduled internet speed tests using Ookla Speedtest CLI
• Scheduled iperf network performance tests  
• Configurable test intervals and duration
• Automatic random host selection for iperf tests

This command supports two modes:
• API mode (default): Submit results via HTTP API (recommended)
• Legacy mode: Direct database access (--legacy flag)

The daemon will run until interrupted (Ctrl+C).`,
	RunE: runDaemon,
}

var (
	useLegacyMode bool
	apiEndpoint   string
)

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Add flags
	daemonCmd.Flags().BoolVar(&useLegacyMode, "legacy", false, "Use legacy direct database access instead of API")
	daemonCmd.Flags().StringVar(&apiEndpoint, "api-endpoint", "http://localhost:8080", "API endpoint URL for API mode")
}

func runDaemon(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	if useLegacyMode {
		log.Println("Starting testing daemon in LEGACY mode (direct database access)...")
		return runLegacyDaemon(cfg)
	}

	log.Println("Starting testing daemon in API mode...")
	return runAPIDaemon(cfg)
}

// runAPIDaemon runs the daemon using API communication
func runAPIDaemon(cfg *config.Config) error {
	// Create API client
	apiBaseURL := fmt.Sprintf("%s/api/v1", apiEndpoint)
	daemonClient := daemon.NewAPIClient(apiBaseURL, cfg)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received interrupt signal, shutting down gracefully...")
		cancel()
	}()

	// Start background testing
	log.Printf("API daemon started - Endpoint: %s", apiBaseURL)
	log.Printf("Test intervals - Speed: %v, Iperf: %v",
		cfg.Testing.SpeedTestInterval, cfg.Testing.IperfTestInterval)

	return daemonClient.StartBackgroundTesting(ctx)
}

// runLegacyDaemon runs the daemon using direct database access (original behavior)
func runLegacyDaemon(cfg *config.Config) error {
	// Initialize database with proper error handling and directory creation
	client, err := database.InitializeDatabase(cfg.Database)
	if err != nil {
		return err
	}
	defer client.Close()

	// Initialize services
	speedTestService := services.NewSpeedTestService(client)
	iperfService := services.NewIperfService(client)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received interrupt signal, shutting down gracefully...")
		cancel()
	}()

	// Start background testing
	log.Printf("Legacy daemon started with intervals - Speed tests: %v, Iperf tests: %v",
		cfg.Testing.SpeedTestInterval, cfg.Testing.IperfTestInterval)

	return runBackgroundTesting(ctx, speedTestService, iperfService, cfg)
}

func runBackgroundTesting(ctx context.Context, speedTestService *services.SpeedTestService, iperfService *services.IperfService, cfg *config.Config) error {
	// Speed test ticker
	speedTestTicker := time.NewTicker(cfg.Testing.SpeedTestInterval)
	defer speedTestTicker.Stop()

	// Iperf test ticker
	iperfTestTicker := time.NewTicker(cfg.Testing.IperfTestInterval)
	defer iperfTestTicker.Stop()

	// Run initial tests
	go func() {
		log.Println("Running initial speed test...")
		if _, err := speedTestService.RunTest(ctx); err != nil {
			log.Printf("Initial speed test failed: %v", err)
		}
	}()

	go func() {
		log.Println("Running initial iperf tests...")
		if err := iperfService.RunRandomTests(ctx, cfg.Testing.IperfTestDuration); err != nil {
			log.Printf("Initial iperf tests failed: %v", err)
		}
	}()

	// Handle scheduled tests
	for {
		select {
		case <-ctx.Done():
			log.Println("Testing daemon stopped")
			return nil

		case <-speedTestTicker.C:
			go func() {
				log.Println("Running scheduled speed test...")
				if _, err := speedTestService.RunTest(ctx); err != nil {
					log.Printf("Scheduled speed test failed: %v", err)
				}
			}()

		case <-iperfTestTicker.C:
			go func() {
				log.Println("Running scheduled iperf tests...")
				if err := iperfService.RunRandomTests(ctx, cfg.Testing.IperfTestDuration); err != nil {
					log.Printf("Scheduled iperf tests failed: %v", err)
				}
			}()
		}
	}
}
