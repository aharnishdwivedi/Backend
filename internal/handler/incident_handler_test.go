package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"incident-triage-assistant/internal/domain"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockIncidentUseCase is a mock implementation of domain.IncidentUseCase
type MockIncidentUseCase struct {
	mock.Mock
}

func (m *MockIncidentUseCase) CreateIncident(req *domain.CreateIncidentRequest) (*domain.Incident, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Incident), args.Error(1)
}

func (m *MockIncidentUseCase) GetIncident(id int) (*domain.Incident, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Incident), args.Error(1)
}

func (m *MockIncidentUseCase) GetAllIncidents() ([]*domain.Incident, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Incident), args.Error(1)
}

func (m *MockIncidentUseCase) UpdateIncident(id int, req *domain.CreateIncidentRequest) (*domain.Incident, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Incident), args.Error(1)
}

func (m *MockIncidentUseCase) DeleteIncident(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateIncident(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		setupMock      func(*MockIncidentUseCase)
	}{
		{
			name: "successful incident creation",
			requestBody: map[string]interface{}{
				"title":            "Test Incident",
				"description":      "Test Description",
				"affected_service": "Test Service",
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mockUC *MockIncidentUseCase) {
				expectedIncident := &domain.Incident{
					ID:              1,
					Title:           "Test Incident",
					Description:     "Test Description",
					AffectedService: "Test Service",
					AISeverity:      "Medium",
					AICategory:      "Software",
				}
				mockUC.On("CreateIncident", mock.AnythingOfType("*domain.CreateIncidentRequest")).
					Return(expectedIncident, nil)
			},
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"title": "Test Incident",
				// missing description and affected_service
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockUC *MockIncidentUseCase) {},
		},
		{
			name:           "invalid JSON",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockUC *MockIncidentUseCase) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			mockUC := new(MockIncidentUseCase)
			handler := NewIncidentHandler(mockUC)

			tt.setupMock(mockUC)

			var req *http.Request
			if tt.requestBody != nil {
				jsonBody, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(http.MethodPost, "/incidents", bytes.NewReader(jsonBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			} else {
				req = httptest.NewRequest(http.MethodPost, "/incidents", bytes.NewReader([]byte("invalid json")))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Test
			err := handler.CreateIncident(c)

			// Assertions
			if err != nil {
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockUC.AssertExpectations(t)
		})
	}
}

func TestGetIncident(t *testing.T) {
	tests := []struct {
		name           string
		incidentID     string
		expectedStatus int
		setupMock      func(*MockIncidentUseCase)
	}{
		{
			name:           "successful incident retrieval",
			incidentID:     "1",
			expectedStatus: http.StatusOK,
			setupMock: func(mockUC *MockIncidentUseCase) {
				expectedIncident := &domain.Incident{
					ID:              1,
					Title:           "Test Incident",
					Description:     "Test Description",
					AffectedService: "Test Service",
					AISeverity:      "Medium",
					AICategory:      "Software",
				}
				mockUC.On("GetIncident", 1).Return(expectedIncident, nil)
			},
		},
		{
			name:           "invalid incident ID",
			incidentID:     "invalid",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockUC *MockIncidentUseCase) {},
		},
		{
			name:           "incident not found",
			incidentID:     "999",
			expectedStatus: http.StatusNotFound,
			setupMock: func(mockUC *MockIncidentUseCase) {
				mockUC.On("GetIncident", 999).Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			mockUC := new(MockIncidentUseCase)
			handler := NewIncidentHandler(mockUC)

			tt.setupMock(mockUC)

			req := httptest.NewRequest(http.MethodGet, "/incidents/"+tt.incidentID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.incidentID)

			// Test
			err := handler.GetIncident(c)

			// Assertions
			if err != nil {
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockUC.AssertExpectations(t)
		})
	}
}

func TestGetAllIncidents(t *testing.T) {
	// Setup
	e := echo.New()
	mockUC := new(MockIncidentUseCase)
	handler := NewIncidentHandler(mockUC)

	expectedIncidents := []*domain.Incident{
		{
			ID:              1,
			Title:           "Test Incident 1",
			Description:     "Test Description 1",
			AffectedService: "Test Service 1",
			AISeverity:      "Medium",
			AICategory:      "Software",
		},
		{
			ID:              2,
			Title:           "Test Incident 2",
			Description:     "Test Description 2",
			AffectedService: "Test Service 2",
			AISeverity:      "High",
			AICategory:      "Network",
		},
	}

	mockUC.On("GetAllIncidents").Return(expectedIncidents, nil)

	req := httptest.NewRequest(http.MethodGet, "/incidents", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	err := handler.GetAllIncidents(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["incidents"])
	assert.Equal(t, float64(2), response["count"])

	mockUC.AssertExpectations(t)
}

func TestHealthCheck(t *testing.T) {
	// Setup
	e := echo.New()
	mockUC := new(MockIncidentUseCase)
	handler := NewIncidentHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	err := handler.HealthCheck(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}
