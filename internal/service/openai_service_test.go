package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"incident-triage-assistant/internal/domain"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOpenAIClient is a mock implementation of the OpenAI client
type MockOpenAIClient struct {
	mock.Mock
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(openai.ChatCompletionResponse), args.Error(1)
}

func TestOpenAIService_AnalyzeIncident(t *testing.T) {
	// Set a dummy API key for testing
	os.Setenv("OPENAI_API_KEY", "test-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	tests := []struct {
		name            string
		title           string
		description     string
		affectedService string
		aiResponse      string
		aiError         error
		expectedResult  *domain.IncidentAnalysis
		expectedError   bool
	}{
		{
			name:            "successful analysis",
			title:           "Database timeout",
			description:     "Users unable to login",
			affectedService: "Auth Service",
			aiResponse:      `{"severity": "High", "category": "Database"}`,
			aiError:         nil,
			expectedResult: &domain.IncidentAnalysis{
				Severity: "High",
				Category: "Database",
			},
			expectedError: false,
		},
		{
			name:            "invalid severity fallback",
			title:           "Minor issue",
			description:     "Small bug",
			affectedService: "UI Service",
			aiResponse:      `{"severity": "Invalid", "category": "Software"}`,
			aiError:         nil,
			expectedResult: &domain.IncidentAnalysis{
				Severity: "Medium", // Should fallback to Medium
				Category: "Software",
			},
			expectedError: false,
		},
		{
			name:            "invalid category fallback",
			title:           "Network issue",
			description:     "Connection lost",
			affectedService: "Network Service",
			aiResponse:      `{"severity": "High", "category": "Invalid"}`,
			aiError:         nil,
			expectedResult: &domain.IncidentAnalysis{
				Severity: "High",
				Category: "Software", // Should fallback to Software
			},
			expectedError: false,
		},
		{
			name:            "AI service error",
			title:           "Test incident",
			description:     "Test description",
			affectedService: "Test Service",
			aiResponse:      "",
			aiError:         errors.New("API error"),
			expectedResult:  nil,
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockOpenAIClient)

			service := &OpenAIService{
				client: mockClient,
			}

			if tt.aiError == nil {
				response := openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{
							Message: openai.ChatCompletionMessage{
								Content: tt.aiResponse,
							},
						},
					},
				}
				mockClient.On("CreateChatCompletion", mock.Anything, mock.AnythingOfType("openai.ChatCompletionRequest")).
					Return(response, nil)
			} else {
				mockClient.On("CreateChatCompletion", mock.Anything, mock.AnythingOfType("openai.ChatCompletionRequest")).
					Return(openai.ChatCompletionResponse{}, tt.aiError)
			}

			result, err := service.AnalyzeIncident(tt.title, tt.description, tt.affectedService)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Severity, result.Severity)
				assert.Equal(t, tt.expectedResult.Category, result.Category)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestOpenAIService_NewOpenAIService(t *testing.T) {
	// Test with missing API key
	os.Unsetenv("OPENAI_API_KEY")

	assert.Panics(t, func() {
		NewOpenAIService()
	})

	// Test with valid API key
	os.Setenv("OPENAI_API_KEY", "test-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	service := NewOpenAIService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.client)
}

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	assert.True(t, contains(slice, "a"))
	assert.True(t, contains(slice, "b"))
	assert.True(t, contains(slice, "c"))
	assert.False(t, contains(slice, "d"))
	assert.False(t, contains(slice, ""))
}
