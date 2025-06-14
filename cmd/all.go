package cmd

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"

	"github.com/bfirestone/speed-checker/internal/config"
	"github.com/bfirestone/speed-checker/internal/database"
	"github.com/bfirestone/speed-checker/internal/handlers"
	"github.com/bfirestone/speed-checker/internal/services"
)

// CustomValidator wraps the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// allCmd represents the all command (current monolithic behavior)
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run both API server and testing daemon (current behavior)",
	Long: `Runs the complete speed checker application including:

• HTTP API server with web dashboard
• Background speed testing daemon
• Background iperf testing daemon
• Automatic scheduling of tests

This preserves the original monolithic behavior where everything
runs in a single process.`,
	RunE: runAll,
}

func init() {
	rootCmd.AddCommand(allCmd)
}

func runAll(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// Initialize database with proper error handling and directory creation
	client, err := database.InitializeDatabase(cfg.Database)
	if err != nil {
		return err
	}
	defer client.Close()

	// Initialize services
	speedTestService := services.NewSpeedTestService(client)
	iperfService := services.NewIperfService(client)

	// Initialize handlers
	apiHandler := handlers.NewAPIHandler(speedTestService, iperfService)

	// Initialize Echo
	e := echo.New()

	// Set custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// API routes
	api := e.Group("/api/v1")

	// Speed test routes
	api.GET("/speedtest", apiHandler.GetSpeedTests)
	api.GET("/speedtest/range", apiHandler.GetSpeedTestsInRange)
	api.POST("/speedtest/run", apiHandler.RunSpeedTest)

	// Iperf test routes
	api.GET("/iperf", apiHandler.GetIperfTests)
	api.POST("/iperf/run", apiHandler.RunIperfTests)

	// Host management routes
	api.GET("/hosts", apiHandler.GetHosts)
	api.POST("/hosts", apiHandler.AddHost)
	api.PUT("/hosts/:id", apiHandler.UpdateHost)
	api.DELETE("/hosts/:id", apiHandler.DeleteHost)

	// Dashboard route
	api.GET("/dashboard", apiHandler.GetDashboard)

	// Static files (for SvelteKit frontend)
	e.Static("/", "frontend/build")

	// Start background testing goroutines
	go startBackgroundTesting(speedTestService, iperfService, cfg)

	// Start server
	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	return e.Start(":" + cfg.Server.Port)
}

func startBackgroundTesting(speedTestService *services.SpeedTestService, iperfService *services.IperfService, cfg *config.Config) {
	// Speed test ticker
	speedTestTicker := time.NewTicker(cfg.Testing.SpeedTestInterval)
	defer speedTestTicker.Stop()

	// Iperf test ticker
	iperfTestTicker := time.NewTicker(cfg.Testing.IperfTestInterval)
	defer iperfTestTicker.Stop()

	// Run initial tests
	go func() {
		ctx := context.Background()
		log.Println("Running initial speed test...")
		if _, err := speedTestService.RunTest(ctx); err != nil {
			log.Printf("Initial speed test failed: %v", err)
		}
	}()

	go func() {
		ctx := context.Background()
		log.Println("Running initial iperf tests...")
		if err := iperfService.RunRandomTests(ctx, cfg.Testing.IperfTestDuration); err != nil {
			log.Printf("Initial iperf tests failed: %v", err)
		}
	}()

	// Handle scheduled tests
	for {
		select {
		case <-speedTestTicker.C:
			go func() {
				ctx := context.Background()
				log.Println("Running scheduled speed test...")
				if _, err := speedTestService.RunTest(ctx); err != nil {
					log.Printf("Scheduled speed test failed: %v", err)
				}
			}()

		case <-iperfTestTicker.C:
			go func() {
				ctx := context.Background()
				log.Println("Running scheduled iperf tests...")
				if err := iperfService.RunRandomTests(ctx, cfg.Testing.IperfTestDuration); err != nil {
					log.Printf("Scheduled iperf tests failed: %v", err)
				}
			}()
		}
	}
}
