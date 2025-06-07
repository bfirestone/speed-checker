package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bfirestone/speed-checker/internal/config"
)

var (
	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "speed-checker",
	Short: "A network speed and performance testing tool",
	Long: `Speed Checker is a comprehensive network testing tool that performs:

• Automated internet speed tests using Ookla Speedtest CLI
• Network performance tests using iperf3 against configurable hosts
• Real-time monitoring dashboard with SvelteKit frontend
• Host management for LAN, VPN, and remote testing targets
• Background scheduled testing with configurable intervals

The tool provides both a web dashboard and CLI interface for managing
network performance monitoring and testing workflows.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration before running any command
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().String("host", "", "server host (overrides config)")
	rootCmd.PersistentFlags().String("port", "", "server port (overrides config)")
	rootCmd.PersistentFlags().String("database-dsn", "", "database connection string (overrides config)")

	// Bind flags to viper
	viper.BindPFlag("server.host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("server.port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("database.dsn", rootCmd.PersistentFlags().Lookup("database-dsn"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}
}

// GetConfig returns the loaded configuration
func GetConfig() *config.Config {
	return cfg
}
