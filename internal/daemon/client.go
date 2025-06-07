package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/bfirestone/speed-checker/internal/client"
	"github.com/bfirestone/speed-checker/internal/config"
)

// APIClient handles communication with the Speed Checker API
type APIClient struct {
	client   *client.ClientWithResponses
	daemonID string
	config   *config.Config
}

// NewAPIClient creates a new API-based daemon client
func NewAPIClient(apiBaseURL string, cfg *config.Config) *APIClient {
	// Create the API client
	apiClient, err := client.NewClientWithResponses(apiBaseURL)
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
	}

	// Generate a unique daemon ID
	hostname, _ := os.Hostname()
	daemonID := fmt.Sprintf("daemon-%s-%d", hostname, os.Getpid())

	return &APIClient{
		client:   apiClient,
		daemonID: daemonID,
		config:   cfg,
	}
}

// StartBackgroundTesting starts the daemon testing loops
func (d *APIClient) StartBackgroundTesting(ctx context.Context) error {
	log.Printf("Starting API-based daemon with ID: %s", d.daemonID)
	log.Printf("API endpoint: %s", d.client.ClientInterface.(*client.Client).Server)

	// Speed test ticker
	speedTestTicker := time.NewTicker(d.config.Testing.SpeedTestInterval)
	defer speedTestTicker.Stop()

	// Iperf test ticker
	iperfTestTicker := time.NewTicker(d.config.Testing.IperfTestInterval)
	defer iperfTestTicker.Stop()

	// Run initial tests
	go func() {
		log.Println("Running initial speed test...")
		if err := d.runSpeedTest(ctx); err != nil {
			log.Printf("Initial speed test failed: %v", err)
		}
	}()

	go func() {
		log.Println("Running initial iperf tests...")
		if err := d.runIperfTests(ctx); err != nil {
			log.Printf("Initial iperf tests failed: %v", err)
		}
	}()

	// Handle scheduled tests
	for {
		select {
		case <-ctx.Done():
			log.Println("API daemon stopped")
			return nil

		case <-speedTestTicker.C:
			go func() {
				log.Println("Running scheduled speed test...")
				if err := d.runSpeedTest(ctx); err != nil {
					log.Printf("Scheduled speed test failed: %v", err)
				}
			}()

		case <-iperfTestTicker.C:
			go func() {
				log.Println("Running scheduled iperf tests...")
				if err := d.runIperfTests(ctx); err != nil {
					log.Printf("Scheduled iperf tests failed: %v", err)
				}
			}()
		}
	}
}

// runSpeedTest executes a speed test and submits results via API
func (d *APIClient) runSpeedTest(ctx context.Context) error {
	log.Println("🚀 Starting Ookla speed test...")

	// Run the speedtest CLI
	cmd := exec.CommandContext(ctx, "speedtest", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("speedtest command failed: %w", err)
	}

	// Parse the speedtest output (we'll implement proper JSON parsing)
	result, err := d.parseSpeedTestJSON(output)
	if err != nil {
		return fmt.Errorf("failed to parse speedtest output: %w", err)
	}

	// Submit results via API
	submission := client.SpeedTestSubmission{
		Timestamp:    result.Timestamp,
		DownloadMbps: result.DownloadMbps,
		UploadMbps:   result.UploadMbps,
		PingMs:       result.PingMs,
		DaemonId:     d.daemonID,
		JitterMs:     result.JitterMs,
		ServerName:   result.ServerName,
		ServerId:     result.ServerID,
		Isp:          result.ISP,
		ExternalIp:   result.ExternalIP,
		ResultUrl:    result.ResultURL,
	}

	// Submit to API (this will return 501 Not Implemented for now)
	resp, err := d.client.SubmitSpeedTestWithResponse(ctx, submission)
	if err != nil {
		return fmt.Errorf("failed to submit speed test: %w", err)
	}

	log.Printf("📊 Speed test submitted - Status: %d, Download: %.2f Mbps, Upload: %.2f Mbps",
		resp.StatusCode(), result.DownloadMbps, result.UploadMbps)

	return nil
}

// runIperfTests executes iperf tests against random hosts and submits results via API
func (d *APIClient) runIperfTests(ctx context.Context) error {
	// Get available hosts from API
	hostsResp, err := d.client.GetHostsWithResponse(ctx, &client.GetHostsParams{
		Active: &[]bool{true}[0], // Only get active hosts
	})
	if err != nil {
		return fmt.Errorf("failed to get hosts: %w", err)
	}

	if hostsResp.StatusCode() != 200 || hostsResp.JSON200 == nil {
		return fmt.Errorf("unexpected response getting hosts: %d", hostsResp.StatusCode())
	}

	hosts := *hostsResp.JSON200
	if len(hosts) == 0 {
		log.Println("⚠️  No active hosts available for iperf testing")
		return nil
	}

	// Select a random host
	host := hosts[time.Now().Unix()%int64(len(hosts))]

	log.Printf("🔗 Running iperf test against %s (%s:%d)", host.Name, host.Hostname, host.Port)

	// Run iperf test
	result, err := d.runSingleIperfTest(ctx, host)
	if err != nil {
		log.Printf("❌ Iperf test failed against %s: %v", host.Name, err)

		// Submit failed test result
		submission := client.IperfTestSubmission{
			Timestamp:       time.Now(),
			HostId:          host.Id,
			SentMbps:        0,
			ReceivedMbps:    0,
			Protocol:        client.IperfTestSubmissionProtocolTCP,
			DurationSeconds: d.config.Testing.IperfTestDuration,
			DaemonId:        d.daemonID,
		}

		_, submitErr := d.client.SubmitIperfTestWithResponse(ctx, submission)
		if submitErr != nil {
			log.Printf("Failed to submit failed iperf test: %v", submitErr)
		}

		return err
	}

	// Submit successful test result
	submission := client.IperfTestSubmission{
		Timestamp:       result.Timestamp,
		HostId:          host.Id,
		SentMbps:        result.SentMbps,
		ReceivedMbps:    result.ReceivedMbps,
		Protocol:        client.IperfTestSubmissionProtocol(result.Protocol),
		DurationSeconds: result.Duration,
		DaemonId:        d.daemonID,
		MeanRttMs:       result.MeanRTT,
		Retransmits:     result.Retransmits,
	}

	resp, err := d.client.SubmitIperfTestWithResponse(ctx, submission)
	if err != nil {
		return fmt.Errorf("failed to submit iperf test: %w", err)
	}

	log.Printf("📈 Iperf test submitted - Status: %d, Sent: %.2f Mbps, Received: %.2f Mbps",
		resp.StatusCode(), result.SentMbps, result.ReceivedMbps)

	return nil
}

// Placeholder structures for parsing (we'll implement proper parsing)
type SpeedTestResult struct {
	Timestamp    time.Time
	DownloadMbps float64
	UploadMbps   float64
	PingMs       float64
	JitterMs     *float64
	ServerName   *string
	ServerID     *string
	ISP          *string
	ExternalIP   *string
	ResultURL    *string
}

// OoklaSpeedTestResult represents the JSON structure returned by Ookla speedtest CLI
type OoklaSpeedTestResult struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Ping      struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
		Low     float64 `json:"low"`
		High    float64 `json:"high"`
	} `json:"ping"`
	Download struct {
		Bandwidth int64 `json:"bandwidth"` // bits per second
		Bytes     int64 `json:"bytes"`
		Elapsed   int   `json:"elapsed"`
		Latency   struct {
			Iqm    float64 `json:"iqm"`
			Low    float64 `json:"low"`
			High   float64 `json:"high"`
			Jitter float64 `json:"jitter"`
		} `json:"latency"`
	} `json:"download"`
	Upload struct {
		Bandwidth int64 `json:"bandwidth"` // bits per second
		Bytes     int64 `json:"bytes"`
		Elapsed   int   `json:"elapsed"`
		Latency   struct {
			Iqm    float64 `json:"iqm"`
			Low    float64 `json:"low"`
			High   float64 `json:"high"`
			Jitter float64 `json:"jitter"`
		} `json:"latency"`
	} `json:"upload"`
	PacketLoss float64 `json:"packetLoss"`
	ISP        string  `json:"isp"`
	Interface  struct {
		InternalIP string `json:"internalIp"`
		Name       string `json:"name"`
		MacAddr    string `json:"macAddr"`
		IsVpn      bool   `json:"isVpn"`
		ExternalIP string `json:"externalIp"`
	} `json:"interface"`
	Server struct {
		ID       int    `json:"id"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
		IP       string `json:"ip"`
	} `json:"server"`
	Result struct {
		ID        string `json:"id"`
		URL       string `json:"url"`
		Persisted bool   `json:"persisted"`
	} `json:"result"`
}

type IperfTestResult struct {
	Timestamp    time.Time
	SentMbps     float64
	ReceivedMbps float64
	Protocol     string
	Duration     int
	MeanRTT      *float64
	Retransmits  *float64
}

// IperfJSONResult represents the JSON structure returned by iperf3 CLI
// This is a simplified version focusing on the essential fields we need
type IperfJSONResult struct {
	Start struct {
		TestStart struct {
			Protocol string `json:"protocol"`
			Duration int    `json:"duration"`
		} `json:"test_start"`
	} `json:"start"`
	End struct {
		Streams []struct {
			Sender struct {
				MeanRtt     int `json:"mean_rtt"` // microseconds
				Retransmits int `json:"retransmits"`
			} `json:"sender"`
		} `json:"streams"`
		SumSent struct {
			BitsPerSecond float64 `json:"bits_per_second"`
			Retransmits   int     `json:"retransmits"`
		} `json:"sum_sent"`
		SumReceived struct {
			BitsPerSecond float64 `json:"bits_per_second"`
		} `json:"sum_received"`
	} `json:"end"`
}

// parseSpeedTestJSON parses the speedtest CLI JSON output
func (d *APIClient) parseSpeedTestJSON(jsonOutput []byte) (*SpeedTestResult, error) {
	var ooklaResult OoklaSpeedTestResult
	if err := json.Unmarshal(jsonOutput, &ooklaResult); err != nil {
		// Log the raw output for debugging
		log.Printf("Failed to parse speedtest JSON: %v", err)
		log.Printf("Raw speedtest output: %s", string(jsonOutput))
		return nil, fmt.Errorf("failed to parse speedtest JSON: %w", err)
	}

	// Convert Ookla format to our internal format
	// Bandwidth is in bits per second, convert to Mbps
	downloadMbps := float64(ooklaResult.Download.Bandwidth) / 1000000 // bits/s to Mbps
	uploadMbps := float64(ooklaResult.Upload.Bandwidth) / 1000000     // bits/s to Mbps

	// Create result with proper values
	result := &SpeedTestResult{
		Timestamp:    ooklaResult.Timestamp,
		DownloadMbps: downloadMbps,
		UploadMbps:   uploadMbps,
		PingMs:       ooklaResult.Ping.Latency,
	}

	// Set optional fields
	if ooklaResult.Ping.Jitter > 0 {
		result.JitterMs = &ooklaResult.Ping.Jitter
	}

	if ooklaResult.Server.Name != "" {
		result.ServerName = &ooklaResult.Server.Name
	}

	if ooklaResult.Server.ID > 0 {
		serverID := fmt.Sprintf("%d", ooklaResult.Server.ID)
		result.ServerID = &serverID
	}

	if ooklaResult.ISP != "" {
		result.ISP = &ooklaResult.ISP
	}

	if ooklaResult.Interface.ExternalIP != "" {
		result.ExternalIP = &ooklaResult.Interface.ExternalIP
	}

	if ooklaResult.Result.URL != "" {
		result.ResultURL = &ooklaResult.Result.URL
	}

	log.Printf("✅ Parsed speedtest results - Download: %.2f Mbps, Upload: %.2f Mbps, Ping: %.2f ms, Server: %s",
		downloadMbps, uploadMbps, ooklaResult.Ping.Latency, ooklaResult.Server.Name)

	return result, nil
}

// runSingleIperfTest runs an iperf test against a specific host
func (d *APIClient) runSingleIperfTest(ctx context.Context, host client.Host) (*IperfTestResult, error) {
	// Build iperf command
	cmd := exec.CommandContext(ctx, "iperf3",
		"-c", host.Hostname,
		"-p", strconv.Itoa(host.Port),
		"-t", strconv.Itoa(d.config.Testing.IperfTestDuration),
		"-J", // JSON output
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("iperf3 command failed: %w", err)
	}

	// Parse iperf output
	result, err := d.parseIperfJSON(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse iperf output: %w", err)
	}

	return result, nil
}

// parseIperfJSON parses the iperf3 JSON output
func (d *APIClient) parseIperfJSON(jsonOutput []byte) (*IperfTestResult, error) {
	var iperfResult IperfJSONResult
	if err := json.Unmarshal(jsonOutput, &iperfResult); err != nil {
		// Log the raw output for debugging
		log.Printf("Failed to parse iperf JSON: %v", err)
		log.Printf("Raw iperf output: %s", string(jsonOutput))
		return nil, fmt.Errorf("failed to parse iperf JSON: %w", err)
	}

	// Convert bits per second to Mbps
	sentMbps := iperfResult.End.SumSent.BitsPerSecond / 1000000
	receivedMbps := iperfResult.End.SumReceived.BitsPerSecond / 1000000

	// Create result with proper values
	result := &IperfTestResult{
		Timestamp:    time.Now(), // iperf doesn't provide timestamp, use current time
		SentMbps:     sentMbps,
		ReceivedMbps: receivedMbps,
		Protocol:     "TCP", // Default to TCP, could be extracted from Start.TestStart.Protocol
		Duration:     d.config.Testing.IperfTestDuration,
	}

	// Set optional fields
	if len(iperfResult.End.Streams) > 0 && iperfResult.End.Streams[0].Sender.MeanRtt > 0 {
		// Convert from microseconds to milliseconds
		meanRttMs := float64(iperfResult.End.Streams[0].Sender.MeanRtt) / 1000
		result.MeanRTT = &meanRttMs
	}

	if iperfResult.End.SumSent.Retransmits > 0 {
		retransmits := float64(iperfResult.End.SumSent.Retransmits)
		result.Retransmits = &retransmits
	}

	log.Printf("✅ Parsed iperf results - Sent: %.2f Mbps, Received: %.2f Mbps, RTT: %.2f ms",
		sentMbps, receivedMbps,
		func() float64 {
			if result.MeanRTT != nil {
				return *result.MeanRTT
			} else {
				return 0
			}
		}())

	return result, nil
}
