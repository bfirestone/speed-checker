package cmd

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"

	"github.com/bfirestone/speed-checker/internal/api"
	"github.com/bfirestone/speed-checker/internal/database"
	"github.com/bfirestone/speed-checker/internal/handlers"
	"github.com/bfirestone/speed-checker/internal/services"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server only",
	Long: `Start the HTTP API server without background testing.
This mode is perfect for production deployments where you want
dedicated daemon processes for testing.`,
	RunE: runAPI,
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func runAPI(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	log.Printf("Starting API server on %s:%s", cfg.Server.Host, cfg.Server.Port)

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
	legacyHandler := handlers.NewAPIHandler(speedTestService, iperfService)
	openAPIHandler := handlers.NewOpenAPIHandler(speedTestService, iperfService)

	// Initialize Echo
	e := echo.New()

	// Set custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Legacy API routes (preserved for compatibility)
	legacyAPI := e.Group("/api/v1/legacy")
	legacyAPI.GET("/speedtest", legacyHandler.GetSpeedTests)
	legacyAPI.GET("/speedtest/range", legacyHandler.GetSpeedTestsInRange)
	legacyAPI.POST("/speedtest/run", legacyHandler.RunSpeedTest)
	legacyAPI.GET("/iperf", legacyHandler.GetIperfTests)
	legacyAPI.POST("/iperf/run", legacyHandler.RunIperfTests)
	legacyAPI.GET("/hosts", legacyHandler.GetHosts)
	legacyAPI.POST("/hosts", legacyHandler.AddHost)
	legacyAPI.PUT("/hosts/:id", legacyHandler.UpdateHost)
	legacyAPI.DELETE("/hosts/:id", legacyHandler.DeleteHost)
	legacyAPI.GET("/dashboard", legacyHandler.GetDashboard)

	// New OpenAPI v1 routes (contract-first design)
	apiV1 := e.Group("/api/v1")
	api.RegisterHandlers(apiV1, openAPIHandler)

	// Static files (for SvelteKit frontend)
	e.Static("/", "frontend/build")

	// Start server
	log.Printf("API server ready:")
	log.Printf("  Legacy API: http://%s:%s/api/v1/legacy/", cfg.Server.Host, cfg.Server.Port)
	log.Printf("  OpenAPI v1: http://%s:%s/api/v1/", cfg.Server.Host, cfg.Server.Port)
	log.Printf("  Frontend:   http://%s:%s/", cfg.Server.Host, cfg.Server.Port)
	return e.Start(":" + cfg.Server.Port)
}
