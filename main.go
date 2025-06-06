package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/bfirestone/speed-checker/ent"
	"github.com/bfirestone/speed-checker/internal/config"
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

func main() {
	// Load configuration
	cfg := config.Default()

	// Initialize database client
	client, err := ent.Open("sqlite3", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer client.Close()

	// Run database migration
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

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
	if err := e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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
