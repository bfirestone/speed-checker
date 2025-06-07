package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/bfirestone/speed-checker/internal/database"
	"github.com/bfirestone/speed-checker/internal/services"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run individual tests or manage test configuration",
	Long: `Test management commands for running one-off tests and managing configuration:

• Run individual speed tests or iperf tests
• View recent test results  
• Manage iperf test hosts

Use subcommands to perform specific test operations.`,
}

// testSpeedCmd represents the test speed command
var testSpeedCmd = &cobra.Command{
	Use:   "speed",
	Short: "Run a single speed test",
	Long:  `Run a single internet speed test using Ookla Speedtest CLI and display results.`,
	RunE:  runSpeedTest,
}

// testIperfCmd represents the test iperf command
var testIperfCmd = &cobra.Command{
	Use:   "iperf [host_id]",
	Short: "Run iperf test against specific host or random hosts",
	Long: `Run iperf test against a specific host (by ID) or random hosts if no ID provided.
	
Examples:
  speed-checker test iperf           # Test against random hosts
  speed-checker test iperf 1         # Test against host ID 1`,
	RunE: runIperfTest,
}

// testListCmd represents the test list command
var testListCmd = &cobra.Command{
	Use:   "list [speed|iperf]",
	Short: "List recent test results",
	Long: `List recent test results from the database.

Examples:
  speed-checker test list           # List both speed and iperf tests
  speed-checker test list speed     # List only speed tests
  speed-checker test list iperf     # List only iperf tests`,
	RunE: listTests,
}

var (
	iperfDuration time.Duration
	resultCount   int
)

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(testSpeedCmd)
	testCmd.AddCommand(testIperfCmd)
	testCmd.AddCommand(testListCmd)

	// Flags for iperf command
	testIperfCmd.Flags().DurationVarP(&iperfDuration, "duration", "d", 10*time.Second, "Test duration")

	// Flags for list command
	testListCmd.Flags().IntVarP(&resultCount, "count", "c", 10, "Number of results to show")
}

func runSpeedTest(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// Initialize database client
	client, err := database.InitializeDatabase(cfg.Database)
	if err != nil {
		return err
	}
	defer client.Close()

	// Initialize service
	speedTestService := services.NewSpeedTestService(client)

	result, err := speedTestService.RunTest(context.Background())
	if err != nil {
		return fmt.Errorf("speed test failed: %w", err)
	}

	// Display results
	fmt.Printf("\n🚀 Speed Test Results:\n")
	fmt.Printf("   Download: %.2f Mbps\n", result.DownloadMbps)
	fmt.Printf("   Upload:   %.2f Mbps\n", result.UploadMbps)
	fmt.Printf("   Ping:     %.2f ms\n", result.PingMs)
	fmt.Printf("   Server:   %s\n", result.ServerName)
	fmt.Printf("   ISP:      %s\n", result.Isp)
	fmt.Printf("   Tested:   %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))

	return nil
}

func runIperfTest(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// Initialize database client
	client, err := database.InitializeDatabase(cfg.Database)
	if err != nil {
		return err
	}
	defer client.Close()

	// Initialize service
	iperfService := services.NewIperfService(client)

	if len(args) > 0 {
		// TODO: Implement RunTestByHostID method or simplify approach
		return fmt.Errorf("testing specific host by ID not yet implemented - use random tests instead")
	} else {
		// Test random hosts
		log.Println("Running iperf tests against random hosts...")
		err := iperfService.RunRandomTests(context.Background(), int(iperfDuration.Seconds()))
		if err != nil {
			return fmt.Errorf("iperf tests failed: %w", err)
		}
		fmt.Println("✅ Iperf tests completed successfully")
	}

	return nil
}

func listTests(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// Initialize database client
	client, err := database.InitializeDatabase(cfg.Database)
	if err != nil {
		return err
	}
	defer client.Close()

	// Initialize services
	speedTestService := services.NewSpeedTestService(client)
	iperfService := services.NewIperfService(client)

	testType := "all"
	if len(args) > 0 {
		testType = args[0]
	}

	ctx := context.Background()

	switch testType {
	case "speed":
		tests, err := speedTestService.GetRecentTests(ctx, resultCount)
		if err != nil {
			return err
		}

		fmt.Printf("\n🚀 Recent Speed Tests (%d results):\n", len(tests))
		for _, test := range tests {
			fmt.Printf("  %s | ↓%.1f ↑%.1f Mbps | %s\n",
				test.Timestamp.Format("01-02 15:04"),
				test.DownloadMbps, test.UploadMbps, test.ServerName)
		}

	case "iperf":
		tests, err := iperfService.GetRecentTests(ctx, resultCount)
		if err != nil {
			return err
		}

		fmt.Printf("\n⚡ Recent Iperf Tests (%d results):\n", len(tests))
		for _, test := range tests {
			hostName := "Unknown"
			if test.Edges.Host != nil {
				hostName = test.Edges.Host.Name
			}
			fmt.Printf("  %s | %s | ↓%.1f ↑%.1f Mbps | %s\n",
				test.Timestamp.Format("01-02 15:04"),
				test.Protocol, test.ReceivedMbps, test.SentMbps, hostName)
		}

	default:
		// Show both
		speedTests, err := speedTestService.GetRecentTests(ctx, resultCount/2)
		if err != nil {
			return err
		}

		iperfTests, err := iperfService.GetRecentTests(ctx, resultCount/2)
		if err != nil {
			return err
		}

		fmt.Printf("\n📊 Recent Test Results:\n")

		fmt.Printf("\n🚀 Speed Tests (%d results):\n", len(speedTests))
		for _, test := range speedTests {
			fmt.Printf("  %s | ↓%.1f ↑%.1f Mbps | %s\n",
				test.Timestamp.Format("01-02 15:04"),
				test.DownloadMbps, test.UploadMbps, test.ServerName)
		}

		fmt.Printf("\n⚡ Iperf Tests (%d results):\n", len(iperfTests))
		for _, test := range iperfTests {
			hostName := "Unknown"
			if test.Edges.Host != nil {
				hostName = test.Edges.Host.Name
			}
			fmt.Printf("  %s | %s | ↓%.1f ↑%.1f Mbps | %s\n",
				test.Timestamp.Format("01-02 15:04"),
				test.Protocol, test.ReceivedMbps, test.SentMbps, hostName)
		}
	}

	return nil
}
