package handler

import (
	"net/http"
	"strconv"

	"incident-triage-assistant/internal/domain"
	"incident-triage-assistant/internal/usecase"

	"github.com/labstack/echo/v4"
)

// IncidentHandler handles HTTP requests for incident management
type IncidentHandler struct {
	incidentUseCase *usecase.IncidentUseCase
}

// NewIncidentHandler creates a new incident handler
func NewIncidentHandler(incidentUseCase *usecase.IncidentUseCase) *IncidentHandler {
	return &IncidentHandler{
		incidentUseCase: incidentUseCase,
	}
}

// CreateIncident handles POST /incidents
func (h *IncidentHandler) CreateIncident(c echo.Context) error {
	var req domain.CreateIncidentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Basic validation
	if req.Title == "" || req.Description == "" || req.AffectedService == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Title, description, and affected service are required")
	}

	incident, err := h.incidentUseCase.CreateIncident(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create incident: "+err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":  "Incident created successfully",
		"incident": incident,
	})
}

// GetIncident handles GET /incidents/:id
func (h *IncidentHandler) GetIncident(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid incident ID")
	}

	incident, err := h.incidentUseCase.GetIncident(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Incident not found")
	}

	return c.JSON(http.StatusOK, incident)
}

// GetAllIncidents handles GET /incidents
func (h *IncidentHandler) GetAllIncidents(c echo.Context) error {
	incidents, err := h.incidentUseCase.GetAllIncidents()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve incidents: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"incidents": incidents,
		"count":     len(incidents),
	})
}

// UpdateIncident handles PUT /incidents/:id
func (h *IncidentHandler) UpdateIncident(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid incident ID")
	}

	var req domain.CreateIncidentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Basic validation
	if req.Title == "" || req.Description == "" || req.AffectedService == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Title, description, and affected service are required")
	}

	incident, err := h.incidentUseCase.UpdateIncident(id, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update incident: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Incident updated successfully",
		"incident": incident,
	})
}

// DeleteIncident handles DELETE /incidents/:id
func (h *IncidentHandler) DeleteIncident(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid incident ID")
	}

	err = h.incidentUseCase.DeleteIncident(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete incident: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Incident deleted successfully",
	})
}

// HealthCheck handles GET /health
func (h *IncidentHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}
