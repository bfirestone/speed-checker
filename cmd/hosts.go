package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/bfirestone/speed-checker/internal/database"
	"github.com/bfirestone/speed-checker/internal/services"
)

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Manage iperf test hosts",
	Long: `Manage iperf test hosts for network performance testing:

• List all configured hosts
• Add new hosts for testing
• Update existing host configuration
• Delete hosts from the system

Hosts can be categorized as LAN, VPN, or remote for different
types of network performance testing.`,
}

// hostsListCmd represents the hosts list command
var hostsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured hosts",
	Long:  `List all configured iperf test hosts with their details.`,
	RunE:  listHosts,
}

// hostsAddCmd represents the hosts add command
var hostsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new iperf test host",
	Long: `Add a new iperf test host to the system.

Type must be one of: lan, vpn, remote

Examples:
  speed-checker hosts add --name "Local Server" --hostname 192.168.1.100 --type lan --description "Main server"
  speed-checker hosts add --name "VPN Gateway" --hostname vpn.example.com --type vpn --port 5202
  speed-checker hosts add --name "Remote Test" --hostname remote.example.com --type remote`,
	RunE: addHost,
}

// hostsDeleteCmd represents the hosts delete command
var hostsDeleteCmd = &cobra.Command{
	Use:   "delete <host_id>",
	Short: "Delete a host by ID",
	Long:  `Delete an iperf test host from the system by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE:  deleteHost,
}

var (
	hostPort        int
	hostName        string
	hostHostname    string
	hostType        string
	hostDescription string
)

func init() {
	rootCmd.AddCommand(hostsCmd)
	hostsCmd.AddCommand(hostsListCmd)
	hostsCmd.AddCommand(hostsAddCmd)
	hostsCmd.AddCommand(hostsDeleteCmd)

	// Flags for add command
	hostsAddCmd.Flags().StringVarP(&hostName, "name", "n", "", "Host name (required)")
	hostsAddCmd.Flags().StringVarP(&hostHostname, "hostname", "H", "", "Host hostname/IP address (required)")
	hostsAddCmd.Flags().StringVarP(&hostType, "type", "t", "", "Host type: lan, vpn, or remote (required)")
	hostsAddCmd.Flags().StringVarP(&hostDescription, "description", "d", "", "Host description (optional)")
	hostsAddCmd.Flags().IntVarP(&hostPort, "port", "p", 5201, "Host port")

	// Mark required flags
	hostsAddCmd.MarkFlagRequired("name")
	hostsAddCmd.MarkFlagRequired("hostname")
	hostsAddCmd.MarkFlagRequired("type")
}

func listHosts(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// Initialize database client
	client, err := database.InitializeDatabase(cfg.Database)
	if err != nil {
		return err
	}
	defer client.Close()

	// Initialize service
	iperfService := services.NewIperfService(client)

	hosts, err := iperfService.GetHosts(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get hosts: %w", err)
	}

	if len(hosts) == 0 {
		fmt.Println("No hosts configured")
		return nil
	}

	fmt.Printf("\n🏠 Configured Hosts (%d total):\n", len(hosts))
	fmt.Printf("%-4s %-20s %-25s %-8s %-6s %-8s %s\n",
		"ID", "Name", "Hostname", "Type", "Port", "Active", "Description")
	fmt.Println("────────────────────────────────────────────────────────────────────────────────")

	for _, host := range hosts {
		activeStatus := "✓"
		if !host.Active {
			activeStatus = "✗"
		}

		description := host.Description
		if description == "" {
			description = "-"
		}

		fmt.Printf("%-4d %-20s %-25s %-8s %-6d %-8s %s\n",
			host.ID, host.Name, host.Hostname, host.Type, host.Port, activeStatus, description)
	}

	return nil
}

func addHost(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// Validate host type
	if hostType != "lan" && hostType != "vpn" && hostType != "remote" {
		return fmt.Errorf("invalid host type '%s'. Must be one of: lan, vpn, remote", hostType)
	}

	// Initialize database client
	client, err := database.InitializeDatabase(cfg.Database)
	if err != nil {
		return err
	}
	defer client.Close()

	// Initialize service
	iperfService := services.NewIperfService(client)

	host, err := iperfService.AddHost(context.Background(), hostName, hostHostname, hostType, hostDescription, hostPort)
	if err != nil {
		return fmt.Errorf("failed to add host: %w", err)
	}

	fmt.Printf("✅ Host added successfully:\n")
	fmt.Printf("   ID:          %d\n", host.ID)
	fmt.Printf("   Name:        %s\n", host.Name)
	fmt.Printf("   Hostname:    %s\n", host.Hostname)
	fmt.Printf("   Type:        %s\n", host.Type)
	fmt.Printf("   Port:        %d\n", host.Port)
	fmt.Printf("   Description: %s\n", host.Description)

	return nil
}

func deleteHost(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	hostID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid host ID: %s", args[0])
	}

	// Initialize database client
	client, err := database.InitializeDatabase(cfg.Database)
	if err != nil {
		return err
	}
	defer client.Close()

	// Initialize service
	iperfService := services.NewIperfService(client)

	err = iperfService.DeleteHost(context.Background(), hostID)
	if err != nil {
		return fmt.Errorf("failed to delete host: %w", err)
	}

	fmt.Printf("✅ Host ID %d deleted successfully\n", hostID)

	return nil
}
