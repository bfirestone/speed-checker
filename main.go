package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SpeedTestResult struct {
	Timestamp    time.Time `json:"timestamp"`
	Download     float64   `json:"download"`
	Upload       float64   `json:"upload"`
	Ping         float64   `json:"ping"`
	Jitter       float64   `json:"jitter"`
	ServerName   string    `json:"server_name"`
	ServerID     string    `json:"server_id"`
	ISP          string    `json:"isp"`
	ExternalIP   string    `json:"external_ip"`
	ResultURL    string    `json:"result_url"`
}

type OoklaResult struct {
	Timestamp time.Time `json:"timestamp"`
	Ping      struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
	} `json:"ping"`
	Download struct {
		Bandwidth int64 `json:"bandwidth"`
	} `json:"download"`
	Upload struct {
		Bandwidth int64 `json:"bandwidth"`
	} `json:"upload"`
	Server struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"server"`
	ISP    string `json:"isp"`
	Result struct {
		URL string `json:"url"`
	} `json:"result"`
	Interface struct {
		ExternalIP string `json:"externalIp"`
	} `json:"interface"`
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./speedtest_results.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS speed_tests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME NOT NULL,
		download_mbps REAL NOT NULL,
		upload_mbps REAL NOT NULL,
		ping_ms REAL NOT NULL,
		jitter_ms REAL,
		server_name TEXT,
		server_id TEXT,
		isp TEXT,
		external_ip TEXT,
		result_url TEXT
	);`

	_, err = db.Exec(createTableSQL)
	return db, err
}

func runSpeedTest() (*SpeedTestResult, error) {
	// Run speedtest-cli with JSON output
	cmd := exec.Command("speedtest", "--format=json", "--accept-license", "--accept-gdpr")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run speedtest: %v", err)
	}

	var ooklaResult OoklaResult
	if err := json.Unmarshal(output, &ooklaResult); err != nil {
		return nil, fmt.Errorf("failed to parse speedtest output: %v", err)
	}

	// Convert bandwidth from bits/second to Mbps
	downloadMbps := float64(ooklaResult.Download.Bandwidth) * 8 / 1000000
	uploadMbps := float64(ooklaResult.Upload.Bandwidth) * 8 / 1000000

	result := &SpeedTestResult{
		Timestamp:    ooklaResult.Timestamp,
		Download:     downloadMbps,
		Upload:       uploadMbps,
		Ping:         ooklaResult.Ping.Latency,
		Jitter:       ooklaResult.Ping.Jitter,
		ServerName:   ooklaResult.Server.Name,
		ServerID:     fmt.Sprintf("%d", ooklaResult.Server.ID),
		ISP:          ooklaResult.ISP,
		ExternalIP:   ooklaResult.Interface.ExternalIP,
		ResultURL:    ooklaResult.Result.URL,
	}

	return result, nil
}

func saveResult(db *sql.DB, result *SpeedTestResult) error {
	insertSQL := `
	INSERT INTO speed_tests (
		timestamp, download_mbps, upload_mbps, ping_ms, jitter_ms,
		server_name, server_id, isp, external_ip, result_url
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(insertSQL,
		result.Timestamp,
		result.Download,
		result.Upload,
		result.Ping,
		result.Jitter,
		result.ServerName,
		result.ServerID,
		result.ISP,
		result.ExternalIP,
		result.ResultURL,
	)
	return err
}

func main() {
	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Set test interval (15-30 minutes)
	interval := 10 * time.Minute // Change this as needed
	log.Printf("Starting speed test monitor with %v intervals", interval)

	// Run initial test
	runTest(db)

	// Set up ticker for regular tests
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run tests on schedule
	for range ticker.C {
		runTest(db)
	}
}

func runTest(db *sql.DB) {
	log.Println("Running speed test...")
	
	result, err := runSpeedTest()
	if err != nil {
		log.Printf("Speed test failed: %v", err)
		return
	}

	err = saveResult(db, result)
	if err != nil {
		log.Printf("Failed to save result: %v", err)
		return
	}

	log.Printf("Test completed - Download: %.2f Mbps, Upload: %.2f Mbps, Ping: %.2f ms",
		result.Download, result.Upload, result.Ping)
}
