package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/bfirestone/speed-checker/internal/services"
)

type APIHandler struct {
	speedTestService *services.SpeedTestService
	iperfService     *services.IperfService
}

func NewAPIHandler(speedTestService *services.SpeedTestService, iperfService *services.IperfService) *APIHandler {
	return &APIHandler{
		speedTestService: speedTestService,
		iperfService:     iperfService,
	}
}

// Speed Test endpoints
func (h *APIHandler) GetSpeedTests(c echo.Context) error {
	limit := 50 // default
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	tests, err := h.speedTestService.GetRecentTests(c.Request().Context(), limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  tests,
		"count": len(tests),
	})
}

func (h *APIHandler) GetSpeedTestsInRange(c echo.Context) error {
	startStr := c.QueryParam("start")
	endStr := c.QueryParam("end")

	if startStr == "" || endStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "start and end parameters are required")
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid start time format")
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid end time format")
	}

	tests, err := h.speedTestService.GetTestsInRange(c.Request().Context(), start, end)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  tests,
		"count": len(tests),
		"range": map[string]string{
			"start": start.Format(time.RFC3339),
			"end":   end.Format(time.RFC3339),
		},
	})
}

func (h *APIHandler) RunSpeedTest(c echo.Context) error {
	test, err := h.speedTestService.RunTest(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Speed test completed successfully",
		"data":    test,
	})
}

// Iperf Test endpoints
func (h *APIHandler) GetIperfTests(c echo.Context) error {
	limit := 50 // default
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	tests, err := h.iperfService.GetRecentTests(c.Request().Context(), limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  tests,
		"count": len(tests),
	})
}

func (h *APIHandler) RunIperfTests(c echo.Context) error {
	duration := 10 // default
	if d := c.QueryParam("duration"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			duration = parsed
		}
	}

	err := h.iperfService.RunRandomTests(c.Request().Context(), duration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Iperf tests completed successfully",
	})
}

// Host management endpoints
func (h *APIHandler) GetHosts(c echo.Context) error {
	hosts, err := h.iperfService.GetHosts(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  hosts,
		"count": len(hosts),
	})
}

type AddHostRequest struct {
	Name        string `json:"name" validate:"required"`
	Hostname    string `json:"hostname" validate:"required"`
	Port        int    `json:"port" validate:"required,min=1,max=65535"`
	Type        string `json:"type" validate:"required,oneof=lan vpn remote"`
	Description string `json:"description"`
}

type UpdateHostRequest struct {
	Name        string `json:"name" validate:"required"`
	Hostname    string `json:"hostname" validate:"required"`
	Port        int    `json:"port" validate:"required,min=1,max=65535"`
	Type        string `json:"type" validate:"required,oneof=lan vpn remote"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

func (h *APIHandler) AddHost(c echo.Context) error {
	var req AddHostRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	host, err := h.iperfService.AddHost(
		c.Request().Context(),
		req.Name,
		req.Hostname,
		req.Type,
		req.Description,
		req.Port,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Host added successfully",
		"data":    host,
	})
}

func (h *APIHandler) UpdateHost(c echo.Context) error {
	id := c.Param("id")
	hostID, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid host ID")
	}

	var req UpdateHostRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	host, err := h.iperfService.UpdateHost(
		c.Request().Context(),
		hostID,
		req.Name,
		req.Hostname,
		req.Type,
		req.Description,
		req.Port,
		req.Active,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Host updated successfully",
		"data":    host,
	})
}

func (h *APIHandler) DeleteHost(c echo.Context) error {
	id := c.Param("id")
	hostID, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid host ID")
	}

	err = h.iperfService.DeleteHost(c.Request().Context(), hostID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Host deleted successfully",
	})
}

// Dashboard endpoint
func (h *APIHandler) GetDashboard(c echo.Context) error {
	ctx := c.Request().Context()

	// Get recent speed tests
	speedTests, err := h.speedTestService.GetRecentTests(ctx, 10)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Get recent iperf tests
	iperfTests, err := h.iperfService.GetRecentTests(ctx, 10)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Get active hosts
	hosts, err := h.iperfService.GetActiveHosts(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"speed_tests": speedTests,
		"iperf_tests": iperfTests,
		"hosts":       hosts,
		"summary": map[string]interface{}{
			"total_speed_tests": len(speedTests),
			"total_iperf_tests": len(iperfTests),
			"active_hosts":      len(hosts),
		},
	})
}
