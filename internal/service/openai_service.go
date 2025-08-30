package service

import (
	"encoding/json"
	"fmt"
	"incident-triage-assistant/internal/domain"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// OpenAIService implements the AIService interface using OpenAI API
type OpenAIService struct {
	client *openai.Client
}

// NewOpenAIService creates a new OpenAI service instance
func NewOpenAIService() *OpenAIService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient(apiKey)
	return &OpenAIService{client: client}
}

// AnalyzeIncident analyzes an incident using OpenAI to determine severity and category
func (s *OpenAIService) AnalyzeIncident(title, description, affectedService string) (*domain.IncidentAnalysis, error) {
	prompt := fmt.Sprintf(`
Analyze the following IT incident and provide:
1. Severity level (Low, Medium, High, Critical)
2. Category (Network, Software, Hardware, Security, Database, Application, Infrastructure)

Incident Details:
- Title: %s
- Description: %s
- Affected Service: %s

Please respond with only a JSON object in this exact format:
{
  "severity": "Low|Medium|High|Critical",
  "category": "Network|Software|Hardware|Security|Database|Application|Infrastructure"
}
`, title, description, affectedService)

	resp, err := s.client.CreateChatCompletion(
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an IT incident triage assistant. Analyze incidents and provide severity and category classifications. Respond only with valid JSON.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.1, // Low temperature for consistent classification
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get AI analysis: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI service")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	
	// Parse JSON response
	var analysis domain.IncidentAnalysis
	err = json.Unmarshal([]byte(content), &analysis)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Validate severity
	validSeverities := []string{"Low", "Medium", "High", "Critical"}
	if !contains(validSeverities, analysis.Severity) {
		analysis.Severity = "Medium" // Default fallback
	}

	// Validate category
	validCategories := []string{"Network", "Software", "Hardware", "Security", "Database", "Application", "Infrastructure"}
	if !contains(validCategories, analysis.Category) {
		analysis.Category = "Software" // Default fallback
	}

	return &analysis, nil
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
