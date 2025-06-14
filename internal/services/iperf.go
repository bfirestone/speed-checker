package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/bfirestone/speed-checker/ent"
	"github.com/bfirestone/speed-checker/ent/host"
	"github.com/bfirestone/speed-checker/internal/api"
)

type IperfService struct {
	client *ent.Client
}

type IperfResult struct {
	Start struct {
		Connected []struct {
			Socket     int    `json:"socket"`
			LocalHost  string `json:"local_host"`
			LocalPort  int    `json:"local_port"`
			RemoteHost string `json:"remote_host"`
			RemotePort int    `json:"remote_port"`
		} `json:"connected"`
	} `json:"start"`
	Intervals []struct {
		Streams []struct {
			Socket        int     `json:"socket"`
			Start         float64 `json:"start"`
			End           float64 `json:"end"`
			Seconds       float64 `json:"seconds"`
			Bytes         int64   `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
			Retransmits   int     `json:"retransmits"`
		} `json:"streams"`
		Sum struct {
			Start         float64 `json:"start"`
			End           float64 `json:"end"`
			Seconds       float64 `json:"seconds"`
			Bytes         int64   `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
			Retransmits   int     `json:"retransmits"`
		} `json:"sum"`
	} `json:"intervals"`
	End struct {
		Streams []struct {
			Sender struct {
				Socket        int     `json:"socket"`
				Start         float64 `json:"start"`
				End           float64 `json:"end"`
				Seconds       float64 `json:"seconds"`
				Bytes         int64   `json:"bytes"`
				BitsPerSecond float64 `json:"bits_per_second"`
				Retransmits   int     `json:"retransmits"`
				MeanRtt       int     `json:"mean_rtt"`
			} `json:"sender"`
			Receiver struct {
				Socket        int     `json:"socket"`
				Start         float64 `json:"start"`
				End           float64 `json:"end"`
				Seconds       float64 `json:"seconds"`
				Bytes         int64   `json:"bytes"`
				BitsPerSecond float64 `json:"bits_per_second"`
			} `json:"receiver"`
		} `json:"streams"`
		SumSent struct {
			Start         float64 `json:"start"`
			End           float64 `json:"end"`
			Seconds       float64 `json:"seconds"`
			Bytes         int64   `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
			Retransmits   int     `json:"retransmits"`
		} `json:"sum_sent"`
		SumReceived struct {
			Start         float64 `json:"start"`
			End           float64 `json:"end"`
			Seconds       float64 `json:"seconds"`
			Bytes         int64   `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
		} `json:"sum_received"`
	} `json:"end"`
}

func NewIperfService(client *ent.Client) *IperfService {
	return &IperfService{
		client: client,
	}
}

func (s *IperfService) RunRandomTests(ctx context.Context, duration int) error {
	// Test against random LAN hosts
	if err := s.runTestsForType(ctx, "lan", duration); err != nil {
		log.Printf("LAN tests failed: %v", err)
	}

	// Test against random VPN hosts
	if err := s.runTestsForType(ctx, "vpn", duration); err != nil {
		log.Printf("VPN tests failed: %v", err)
	}

	// Test against random remote hosts
	if err := s.runTestsForType(ctx, "remote", duration); err != nil {
		log.Printf("Remote tests failed: %v", err)
	}

	return nil
}

func (s *IperfService) runTestsForType(ctx context.Context, hostType string, duration int) error {
	// Get active hosts of the specified type
	hosts, err := s.client.Host.
		Query().
		Where(
			host.And(
				host.TypeEQ(host.Type(hostType)),
				host.ActiveEQ(true),
			),
		).
		All(ctx)

	if err != nil {
		return fmt.Errorf("failed to query %s hosts: %v", hostType, err)
	}

	if len(hosts) == 0 {
		log.Printf("No active %s hosts found", hostType)
		return nil
	}

	// Select a random host
	selectedHost := hosts[rand.Intn(len(hosts))]
	log.Printf("Running iperf3 test against %s host: %s (%s:%d)",
		hostType, selectedHost.Name, selectedHost.Hostname, selectedHost.Port)

	// Run the test
	return s.runTest(ctx, selectedHost, duration)
}

func (s *IperfService) runTest(ctx context.Context, testHost *ent.Host, duration int) error {
	// Create context with timeout
	testCtx, cancel := context.WithTimeout(ctx, time.Duration(duration+30)*time.Second)
	defer cancel()

	// Run iperf3 client
	cmd := exec.CommandContext(testCtx, "iperf3",
		"-c", testHost.Hostname,
		"-p", fmt.Sprintf("%d", testHost.Port),
		"-t", fmt.Sprintf("%d", duration),
		"-J", // JSON output
	)

	output, err := cmd.Output()
	if err != nil {
		// Save failed test result
		_, saveErr := s.client.IperfTest.
			Create().
			SetHost(testHost).
			SetSuccess(false).
			SetErrorMessage(err.Error()).
			SetDurationSeconds(duration).
			Save(ctx)
		if saveErr != nil {
			log.Printf("Failed to save error result: %v", saveErr)
		}
		return fmt.Errorf("iperf3 test failed: %v", err)
	}

	// Parse JSON output
	var result IperfResult
	if err := json.Unmarshal(output, &result); err != nil {
		return fmt.Errorf("failed to parse iperf3 output: %v", err)
	}

	// Extract metrics
	sentMbps := result.End.SumSent.BitsPerSecond / 1000000
	receivedMbps := result.End.SumReceived.BitsPerSecond / 1000000
	retransmits := float64(result.End.SumSent.Retransmits)

	var meanRtt float64
	if len(result.End.Streams) > 0 {
		meanRtt = float64(result.End.Streams[0].Sender.MeanRtt) / 1000 // Convert to ms
	}

	// Save successful test result
	iperfTest, err := s.client.IperfTest.
		Create().
		SetHost(testHost).
		SetSentMbps(sentMbps).
		SetReceivedMbps(receivedMbps).
		SetRetransmits(retransmits).
		SetMeanRttMs(meanRtt).
		SetDurationSeconds(duration).
		SetProtocol("TCP").
		SetSuccess(true).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to save iperf test result: %v", err)
	}

	log.Printf("Iperf3 test completed - Sent: %.2f Mbps, Received: %.2f Mbps, RTT: %.2f ms",
		sentMbps, receivedMbps, meanRtt)

	_ = iperfTest
	return nil
}

func (s *IperfService) GetRecentTests(ctx context.Context, limit int) ([]*ent.IperfTest, error) {
	return s.client.IperfTest.
		Query().
		WithHost().
		Order(ent.Desc("timestamp")).
		Limit(limit).
		All(ctx)
}

func (s *IperfService) GetTestsInRange(ctx context.Context, start, end time.Time) ([]*ent.IperfTest, error) {
	// TODO: Implement timestamp filtering once we add the iperftest predicate import
	return s.client.IperfTest.
		Query().
		WithHost().
		Order(ent.Desc("timestamp")).
		All(ctx)
}

// Host management methods
func (s *IperfService) AddHost(ctx context.Context, name, hostname, hostType, description string, port int) (*ent.Host, error) {
	return s.client.Host.
		Create().
		SetName(name).
		SetHostname(hostname).
		SetPort(port).
		SetType(host.Type(hostType)).
		SetDescription(description).
		SetActive(true).
		Save(ctx)
}

func (s *IperfService) GetHosts(ctx context.Context) ([]*ent.Host, error) {
	return s.client.Host.
		Query().
		Order(ent.Asc("name")).
		All(ctx)
}

func (s *IperfService) GetActiveHosts(ctx context.Context) ([]*ent.Host, error) {
	return s.client.Host.
		Query().
		Where(host.ActiveEQ(true)).
		Order(ent.Asc("name")).
		All(ctx)
}

func (s *IperfService) UpdateHost(ctx context.Context, id int, name, hostname, hostType, description string, port int, active bool) (*ent.Host, error) {
	return s.client.Host.
		UpdateOneID(id).
		SetName(name).
		SetHostname(hostname).
		SetPort(port).
		SetType(host.Type(hostType)).
		SetDescription(description).
		SetActive(active).
		Save(ctx)
}

func (s *IperfService) DeleteHost(ctx context.Context, id int) error {
	return s.client.Host.DeleteOneID(id).Exec(ctx)
}

func (s *IperfService) GetTestsByHostName(ctx context.Context, hostName string, limit int) ([]*ent.IperfTest, error) {
	return s.client.IperfTest.
		Query().
		WithHost(func(q *ent.HostQuery) {
			q.Where(host.NameContains(hostName))
		}).
		Order(ent.Desc("timestamp")).
		Limit(limit).
		All(ctx)
}

func (s *IperfService) GetTestsByHostType(ctx context.Context, hostType string, limit int) ([]*ent.IperfTest, error) {
	return s.client.IperfTest.
		Query().
		WithHost(func(q *ent.HostQuery) {
			q.Where(host.TypeEQ(host.Type(hostType)))
		}).
		Order(ent.Desc("timestamp")).
		Limit(limit).
		All(ctx)
}

func (s *IperfService) GetSlowestTests(ctx context.Context, limit int) ([]*ent.IperfTest, error) {
	return s.client.IperfTest.
		Query().
		WithHost().
		Order(ent.Asc("received_mbps")). // Ascending order to get slowest first
		Limit(limit).
		All(ctx)
}

func (s *IperfService) GetTotalCount(ctx context.Context) (int, error) {
	return s.client.IperfTest.Query().Count(ctx)
}

// CreateFromSubmission creates an iperf test record from an API submission
func (s *IperfService) CreateFromSubmission(ctx context.Context, submission api.IperfTestSubmission) (*ent.IperfTest, error) {
	// Get the host by ID
	targetHost, err := s.client.Host.Get(ctx, submission.HostId)
	if err != nil {
		return nil, fmt.Errorf("failed to find host with ID %d: %w", submission.HostId, err)
	}

	// Create the iperf test record using Ent
	builder := s.client.IperfTest.
		Create().
		SetHost(targetHost).
		SetTimestamp(submission.Timestamp).
		SetSentMbps(submission.SentMbps).
		SetReceivedMbps(submission.ReceivedMbps).
		SetProtocol(string(submission.Protocol)).
		SetDurationSeconds(submission.DurationSeconds).
		SetDaemonID(submission.DaemonId)

	// Determine success based on whether we have meaningful throughput
	success := submission.SentMbps > 0 || submission.ReceivedMbps > 0
	builder.SetSuccess(success)

	// Set optional fields if provided
	if submission.MeanRttMs != nil {
		builder.SetMeanRttMs(*submission.MeanRttMs)
	}
	if submission.Retransmits != nil {
		builder.SetRetransmits(*submission.Retransmits)
	}

	iperfTest, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save iperf test submission: %w", err)
	}

	// Set the host edge manually since we already have the host
	iperfTest.Edges.Host = targetHost

	log.Printf("Iperf test submission saved - ID: %d, Daemon: %s, Host: %s, Sent: %.2f Mbps, Received: %.2f Mbps",
		iperfTest.ID, submission.DaemonId, targetHost.Name, submission.SentMbps, submission.ReceivedMbps)

	return iperfTest, nil
}

// DeleteTest deletes an iperf test by ID
func (s *IperfService) DeleteTest(ctx context.Context, testID int) error {
	err := s.client.IperfTest.DeleteOneID(testID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete iperf test %d: %w", testID, err)
	}

	log.Printf("Iperf test deleted - ID: %d", testID)
	return nil
}
