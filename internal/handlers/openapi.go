package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"log"

	"github.com/bfirestone/speed-checker/ent"
	"github.com/bfirestone/speed-checker/internal/api"
	"github.com/bfirestone/speed-checker/internal/services"
)

// OpenAPIHandler implements the generated ServerInterface
type OpenAPIHandler struct {
	speedTestService *services.SpeedTestService
	iperfService     *services.IperfService
}

// NewOpenAPIHandler creates a new OpenAPI handler
func NewOpenAPIHandler(speedTestService *services.SpeedTestService, iperfService *services.IperfService) *OpenAPIHandler {
	return &OpenAPIHandler{
		speedTestService: speedTestService,
		iperfService:     iperfService,
	}
}

// Ensure OpenAPIHandler implements the ServerInterface
var _ api.ServerInterface = (*OpenAPIHandler)(nil)

// Speed Test Endpoints

// GetSpeedTests implements GET /speedtest/results
func (h *OpenAPIHandler) GetSpeedTests(ctx echo.Context, params api.GetSpeedTestsParams) error {
	// Set defaults
	limit := 100
	offset := 0

	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}

	// For now, use the existing service methods - we'll enhance these later
	tests, err := h.speedTestService.GetRecentTests(ctx.Request().Context(), limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Error:   "internal_error",
			Message: "Failed to retrieve speed tests",
		})
	}

	// Convert Ent models to OpenAPI models
	results := make([]api.SpeedTestResult, len(tests))
	for i, test := range tests {
		results[i] = entSpeedTestToAPI(test)
	}

	// Get total count for pagination
	totalCount, err := h.speedTestService.GetTotalCount(ctx.Request().Context())
	if err != nil {
		// Log error but don't fail the request
		totalCount = len(results)
	}

	response := struct {
		Results []api.SpeedTestResult `json:"results"`
		Total   int                   `json:"total"`
		Limit   int                   `json:"limit"`
		Offset  int                   `json:"offset"`
	}{
		Results: results,
		Total:   totalCount,
		Limit:   limit,
		Offset:  offset,
	}

	return ctx.JSON(http.StatusOK, response)
}

// SubmitSpeedTest implements POST /speedtest/results
func (h *OpenAPIHandler) SubmitSpeedTest(ctx echo.Context) error {
	var submission api.SpeedTestSubmission
	if err := ctx.Bind(&submission); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.Error{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Create speed test via service
	speedTest, err := h.speedTestService.CreateFromSubmission(ctx.Request().Context(), submission)
	if err != nil {
		log.Printf("Failed to create speed test from submission: %v", err)
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Error:   "creation_failed",
			Message: "Failed to save speed test result",
		})
	}

	// Return the created speed test result
	result := entSpeedTestToAPI(speedTest)
	return ctx.JSON(http.StatusCreated, result)
}

// DeleteSpeedTest implements DELETE /speedtest/results/{testId}
func (h *OpenAPIHandler) DeleteSpeedTest(ctx echo.Context, testId int) error {
	err := h.speedTestService.DeleteTest(ctx.Request().Context(), testId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.Error{
			Error:   "not_found",
			Message: "Speed test result not found",
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Iperf Test Endpoints

// GetIperfTests implements GET /iperf/results
func (h *OpenAPIHandler) GetIperfTests(ctx echo.Context, params api.GetIperfTestsParams) error {
	// Set defaults
	limit := 100
	offset := 0

	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}

	// For now, use the existing service methods
	tests, err := h.iperfService.GetRecentTests(ctx.Request().Context(), limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Error:   "internal_error",
			Message: "Failed to retrieve iperf tests",
		})
	}

	// Convert Ent models to OpenAPI models
	results := make([]api.IperfTestResult, len(tests))
	for i, test := range tests {
		results[i] = entIperfTestToAPI(test)
	}

	// Get total count for pagination
	totalCount, err := h.iperfService.GetTotalCount(ctx.Request().Context())
	if err != nil {
		// Log error but don't fail the request
		totalCount = len(results)
	}

	response := struct {
		Results []api.IperfTestResult `json:"results"`
		Total   int                   `json:"total"`
		Limit   int                   `json:"limit"`
		Offset  int                   `json:"offset"`
	}{
		Results: results,
		Total:   totalCount,
		Limit:   limit,
		Offset:  offset,
	}

	return ctx.JSON(http.StatusOK, response)
}

// SubmitIperfTest implements POST /iperf/results
func (h *OpenAPIHandler) SubmitIperfTest(ctx echo.Context) error {
	var submission api.IperfTestSubmission
	if err := ctx.Bind(&submission); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.Error{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Create iperf test via service
	iperfTest, err := h.iperfService.CreateFromSubmission(ctx.Request().Context(), submission)
	if err != nil {
		log.Printf("Failed to create iperf test from submission: %v", err)
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Error:   "creation_failed",
			Message: "Failed to save iperf test result",
		})
	}

	// Return the created iperf test result
	result := entIperfTestToAPI(iperfTest)
	return ctx.JSON(http.StatusCreated, result)
}

// DeleteIperfTest implements DELETE /iperf/results/{testId}
func (h *OpenAPIHandler) DeleteIperfTest(ctx echo.Context, testId int) error {
	err := h.iperfService.DeleteTest(ctx.Request().Context(), testId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.Error{
			Error:   "not_found",
			Message: "Iperf test result not found",
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Host Management Endpoints

// GetHosts implements GET /hosts
func (h *OpenAPIHandler) GetHosts(ctx echo.Context, params api.GetHostsParams) error {
	hosts, err := h.iperfService.GetHosts(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Error:   "internal_error",
			Message: "Failed to retrieve hosts",
		})
	}

	// Convert Ent models to OpenAPI models
	results := make([]api.Host, len(hosts))
	for i, host := range hosts {
		results[i] = entHostToAPI(host)
	}

	return ctx.JSON(http.StatusOK, results)
}

// AddHost implements POST /hosts
func (h *OpenAPIHandler) AddHost(ctx echo.Context) error {
	var hostCreation api.HostCreation
	if err := ctx.Bind(&hostCreation); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.Error{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Create host via service
	host, err := h.iperfService.AddHost(
		ctx.Request().Context(),
		hostCreation.Name,
		hostCreation.Hostname,
		string(hostCreation.Type),
		derefString(hostCreation.Description, ""),
		hostCreation.Port,
	)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.Error{
			Error:   "creation_failed",
			Message: "Failed to create host",
		})
	}

	result := entHostToAPI(host)
	return ctx.JSON(http.StatusCreated, result)
}

// GetHost implements GET /hosts/{hostId}
func (h *OpenAPIHandler) GetHost(ctx echo.Context, hostId int) error {
	// TODO: Implement GetHost method in service
	// For now, return not implemented
	return ctx.JSON(http.StatusNotImplemented, api.Error{
		Error:   "not_implemented",
		Message: "Individual host retrieval not yet implemented",
	})
}

// UpdateHost implements PUT /hosts/{hostId}
func (h *OpenAPIHandler) UpdateHost(ctx echo.Context, hostId int) error {
	var hostUpdate api.HostUpdate
	if err := ctx.Bind(&hostUpdate); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.Error{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Update host via service
	host, err := h.iperfService.UpdateHost(
		ctx.Request().Context(),
		hostId,
		hostUpdate.Name,
		hostUpdate.Hostname,
		string(hostUpdate.Type),
		derefString(hostUpdate.Description, ""),
		hostUpdate.Port,
		derefBool(hostUpdate.Active, true),
	)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.Error{
			Error:   "update_failed",
			Message: "Failed to update host",
		})
	}

	result := entHostToAPI(host)
	return ctx.JSON(http.StatusOK, result)
}

// DeleteHost implements DELETE /hosts/{hostId}
func (h *OpenAPIHandler) DeleteHost(ctx echo.Context, hostId int) error {
	err := h.iperfService.DeleteHost(ctx.Request().Context(), hostId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.Error{
			Error:   "not_found",
			Message: "Host not found",
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Dashboard Endpoint

// GetDashboard implements GET /dashboard
func (h *OpenAPIHandler) GetDashboard(ctx echo.Context) error {
	// Get recent speed tests
	speedTests, err := h.speedTestService.GetRecentTests(ctx.Request().Context(), 10)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Error:   "internal_error",
			Message: "Failed to retrieve dashboard data",
		})
	}

	// Get recent iperf tests
	iperfTests, err := h.iperfService.GetRecentTests(ctx.Request().Context(), 10)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Error:   "internal_error",
			Message: "Failed to retrieve dashboard data",
		})
	}

	// Get active hosts
	hosts, err := h.iperfService.GetActiveHosts(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.Error{
			Error:   "internal_error",
			Message: "Failed to retrieve dashboard data",
		})
	}

	// Get total counts
	totalSpeedTests, _ := h.speedTestService.GetTotalCount(ctx.Request().Context())
	totalIperfTests, _ := h.iperfService.GetTotalCount(ctx.Request().Context())

	// Convert to API models
	recentSpeedTests := make([]api.SpeedTestResult, len(speedTests))
	for i, test := range speedTests {
		recentSpeedTests[i] = entSpeedTestToAPI(test)
	}

	recentIperfTests := make([]api.IperfTestResult, len(iperfTests))
	for i, test := range iperfTests {
		recentIperfTests[i] = entIperfTestToAPI(test)
	}

	// Create dashboard response
	dashboard := api.DashboardData{
		RecentSpeedTests: recentSpeedTests,
		RecentIperfTests: recentIperfTests,
		Statistics: struct {
			ActiveHosts     *int     `json:"active_hosts,omitempty"`
			AvgDownloadMbps *float64 `json:"avg_download_mbps,omitempty"`
			AvgUploadMbps   *float64 `json:"avg_upload_mbps,omitempty"`
			TotalIperfTests *int     `json:"total_iperf_tests,omitempty"`
			TotalSpeedTests *int     `json:"total_speed_tests,omitempty"`
		}{
			ActiveHosts:     &[]int{len(hosts)}[0],
			TotalSpeedTests: &totalSpeedTests,
			TotalIperfTests: &totalIperfTests,
		},
	}

	return ctx.JSON(http.StatusOK, dashboard)
}

// Helper conversion functions

func entSpeedTestToAPI(test *ent.SpeedTest) api.SpeedTestResult {
	daemonId := test.DaemonID
	if daemonId == "" {
		daemonId = "daemon-legacy" // Fallback for tests without daemon_id
	}

	return api.SpeedTestResult{
		Id:           test.ID,
		Timestamp:    test.Timestamp,
		DownloadMbps: test.DownloadMbps,
		UploadMbps:   test.UploadMbps,
		PingMs:       test.PingMs,
		DaemonId:     daemonId,
		CreatedAt:    test.Timestamp, // Use timestamp as created_at for now
		ServerName:   &test.ServerName,
		ServerId:     &test.ServerID,
		Isp:          &test.Isp,
		ExternalIp:   &test.ExternalIP,
		ResultUrl:    &test.ResultURL,
	}
}

func entIperfTestToAPI(test *ent.IperfTest) api.IperfTestResult {
	host := entHostToAPI(test.Edges.Host)

	daemonId := test.DaemonID
	if daemonId == "" {
		daemonId = "daemon-legacy" // Fallback for tests without daemon_id
	}

	return api.IperfTestResult{
		Id:              test.ID,
		Timestamp:       test.Timestamp,
		HostId:          test.Edges.Host.ID,
		SentMbps:        test.SentMbps,
		ReceivedMbps:    test.ReceivedMbps,
		Protocol:        api.IperfTestResultProtocol(test.Protocol),
		DurationSeconds: test.DurationSeconds,
		DaemonId:        daemonId,
		CreatedAt:       test.Timestamp, // Use timestamp as created_at for now
		Host:            host,
		Success:         &test.Success,
		MeanRttMs:       &test.MeanRttMs,
		Retransmits:     &test.Retransmits,
	}
}

func entHostToAPI(host *ent.Host) api.Host {
	// TODO: Add timestamp fields to Host schema
	now := time.Now()
	return api.Host{
		Id:          host.ID,
		Name:        host.Name,
		Hostname:    host.Hostname,
		Type:        api.HostType(host.Type),
		Port:        host.Port,
		Description: &host.Description,
		Active:      &host.Active,
		CreatedAt:   now, // Placeholder until we add timestamps to schema
		UpdatedAt:   now, // Placeholder until we add timestamps to schema
	}
}

// Helper functions for pointer dereferencing
func derefString(ptr *string, defaultValue string) string {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

func derefBool(ptr *bool, defaultValue bool) bool {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}
