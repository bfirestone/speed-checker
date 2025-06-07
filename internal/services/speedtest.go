package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/bfirestone/speed-checker/ent"
	"github.com/bfirestone/speed-checker/ent/speedtest"
)

type SpeedTestService struct {
	client *ent.Client
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

func NewSpeedTestService(client *ent.Client) *SpeedTestService {
	return &SpeedTestService{
		client: client,
	}
}

func (s *SpeedTestService) RunTest(ctx context.Context) (*ent.SpeedTest, error) {
	log.Println("Running speed test...")

	// Run speedtest-cli with JSON output
	cmd := exec.CommandContext(ctx, "speedtest", "--format=json", "--accept-license", "--accept-gdpr")
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

	// Save to database using Ent
	speedTest, err := s.client.SpeedTest.
		Create().
		SetTimestamp(ooklaResult.Timestamp).
		SetDownloadMbps(downloadMbps).
		SetUploadMbps(uploadMbps).
		SetPingMs(ooklaResult.Ping.Latency).
		SetJitterMs(ooklaResult.Ping.Jitter).
		SetServerName(ooklaResult.Server.Name).
		SetServerID(fmt.Sprintf("%d", ooklaResult.Server.ID)).
		SetIsp(ooklaResult.ISP).
		SetExternalIP(ooklaResult.Interface.ExternalIP).
		SetResultURL(ooklaResult.Result.URL).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to save speed test result: %v", err)
	}

	log.Printf("Speed test completed - Download: %.2f Mbps, Upload: %.2f Mbps, Ping: %.2f ms",
		downloadMbps, uploadMbps, ooklaResult.Ping.Latency)

	return speedTest, nil
}

func (s *SpeedTestService) GetRecentTests(ctx context.Context, limit int) ([]*ent.SpeedTest, error) {
	return s.client.SpeedTest.
		Query().
		Order(ent.Desc("timestamp")).
		Limit(limit).
		All(ctx)
}

func (s *SpeedTestService) GetTestsInRange(ctx context.Context, start, end time.Time) ([]*ent.SpeedTest, error) {
	return s.client.SpeedTest.
		Query().
		Where(
			speedtest.And(
				speedtest.TimestampGTE(start),
				speedtest.TimestampLTE(end),
			),
		).
		Order(ent.Desc("timestamp")).
		All(ctx)
}

func (s *SpeedTestService) GetTestsByServerName(ctx context.Context, serverName string, limit int) ([]*ent.SpeedTest, error) {
	return s.client.SpeedTest.
		Query().
		Where(speedtest.ServerNameContains(serverName)).
		Order(ent.Desc("timestamp")).
		Limit(limit).
		All(ctx)
}

func (s *SpeedTestService) GetSlowestTests(ctx context.Context, limit int) ([]*ent.SpeedTest, error) {
	return s.client.SpeedTest.
		Query().
		Order(ent.Asc("download_mbps")). // Ascending order to get slowest first
		Limit(limit).
		All(ctx)
}

func (s *SpeedTestService) GetTotalCount(ctx context.Context) (int, error) {
	return s.client.SpeedTest.Query().Count(ctx)
}
