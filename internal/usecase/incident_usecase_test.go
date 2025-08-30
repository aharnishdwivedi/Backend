package usecase

import (
	"errors"
	"testing"
	"time"

	"incident-triage-assistant/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockIncidentRepository is a mock implementation of IncidentRepository
type MockIncidentRepository struct {
	mock.Mock
}

func (m *MockIncidentRepository) Create(incident *domain.Incident) error {
	args := m.Called(incident)
	return args.Error(0)
}

func (m *MockIncidentRepository) GetByID(id int) (*domain.Incident, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Incident), args.Error(1)
}

func (m *MockIncidentRepository) GetAll() ([]*domain.Incident, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Incident), args.Error(1)
}

func (m *MockIncidentRepository) Update(incident *domain.Incident) error {
	args := m.Called(incident)
	return args.Error(0)
}

func (m *MockIncidentRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockAIService is a mock implementation of AIService
type MockAIService struct {
	mock.Mock
}

func (m *MockAIService) AnalyzeIncident(title, description, affectedService string) (*domain.IncidentAnalysis, error) {
	args := m.Called(title, description, affectedService)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IncidentAnalysis), args.Error(1)
}

func TestCreateIncident(t *testing.T) {
	tests := []struct {
		name           string
		request        *domain.CreateIncidentRequest
		aiAnalysis     *domain.IncidentAnalysis
		aiError        error
		repoError      error
		expectedError  bool
		expectedResult *domain.Incident
	}{
		{
			name: "successful incident creation",
			request: &domain.CreateIncidentRequest{
				Title:           "Test Incident",
				Description:     "Test Description",
				AffectedService: "Test Service",
			},
			aiAnalysis: &domain.IncidentAnalysis{
				Severity: "Medium",
				Category: "Software",
			},
			aiError:       nil,
			repoError:     nil,
			expectedError: false,
			expectedResult: &domain.Incident{
				Title:           "Test Incident",
				Description:     "Test Description",
				AffectedService: "Test Service",
				AISeverity:      "Medium",
				AICategory:      "Software",
			},
		},
		{
			name: "AI service error",
			request: &domain.CreateIncidentRequest{
				Title:           "Test Incident",
				Description:     "Test Description",
				AffectedService: "Test Service",
			},
			aiAnalysis:    nil,
			aiError:       errors.New("AI service unavailable"),
			repoError:     nil,
			expectedError: true,
		},
		{
			name: "repository error",
			request: &domain.CreateIncidentRequest{
				Title:           "Test Incident",
				Description:     "Test Description",
				AffectedService: "Test Service",
			},
			aiAnalysis: &domain.IncidentAnalysis{
				Severity: "High",
				Category: "Network",
			},
			aiError:       nil,
			repoError:     errors.New("database error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockIncidentRepository)
			mockAI := new(MockAIService)

			useCase := NewIncidentUseCase(mockRepo, mockAI)

			if tt.aiAnalysis != nil {
				mockAI.On("AnalyzeIncident", tt.request.Title, tt.request.Description, tt.request.AffectedService).
					Return(tt.aiAnalysis, tt.aiError)
			} else {
				mockAI.On("AnalyzeIncident", tt.request.Title, tt.request.Description, tt.request.AffectedService).
					Return(nil, tt.aiError)
			}

			if tt.aiError == nil && tt.repoError == nil {
				mockRepo.On("Create", mock.AnythingOfType("*domain.Incident")).Return(tt.repoError)
			}

			result, err := useCase.CreateIncident(tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Title, result.Title)
				assert.Equal(t, tt.expectedResult.Description, result.Description)
				assert.Equal(t, tt.expectedResult.AffectedService, result.AffectedService)
				assert.Equal(t, tt.expectedResult.AISeverity, result.AISeverity)
				assert.Equal(t, tt.expectedResult.AICategory, result.AICategory)
			}

			mockRepo.AssertExpectations(t)
			mockAI.AssertExpectations(t)
		})
	}
}

func TestGetIncident(t *testing.T) {
	mockRepo := new(MockIncidentRepository)
	mockAI := new(MockAIService)
	useCase := NewIncidentUseCase(mockRepo, mockAI)

	expectedIncident := &domain.Incident{
		ID:              1,
		Title:           "Test Incident",
		Description:     "Test Description",
		AffectedService: "Test Service",
		AISeverity:      "Medium",
		AICategory:      "Software",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockRepo.On("GetByID", 1).Return(expectedIncident, nil)

	result, err := useCase.GetIncident(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedIncident, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllIncidents(t *testing.T) {
	mockRepo := new(MockIncidentRepository)
	mockAI := new(MockAIService)
	useCase := NewIncidentUseCase(mockRepo, mockAI)

	expectedIncidents := []*domain.Incident{
		{
			ID:              1,
			Title:           "Test Incident 1",
			Description:     "Test Description 1",
			AffectedService: "Test Service 1",
			AISeverity:      "Medium",
			AICategory:      "Software",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              2,
			Title:           "Test Incident 2",
			Description:     "Test Description 2",
			AffectedService: "Test Service 2",
			AISeverity:      "High",
			AICategory:      "Network",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	mockRepo.On("GetAll").Return(expectedIncidents, nil)

	result, err := useCase.GetAllIncidents()

	assert.NoError(t, err)
	assert.Equal(t, expectedIncidents, result)
	mockRepo.AssertExpectations(t)
}
