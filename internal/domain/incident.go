package domain

import (
	"time"
)

// Incident represents an IT incident with AI-generated insights
type Incident struct {
	ID              int       `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	AffectedService string    `json:"affected_service" db:"affected_service"`
	AISeverity      string    `json:"ai_severity" db:"ai_severity"`
	AICategory      string    `json:"ai_category" db:"ai_category"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// CreateIncidentRequest represents the request to create a new incident
type CreateIncidentRequest struct {
	Title           string `json:"title" validate:"required"`
	Description     string `json:"description" validate:"required"`
	AffectedService string `json:"affected_service" validate:"required"`
}

// IncidentRepository defines the interface for incident data operations
type IncidentRepository interface {
	Create(incident *Incident) error
	GetByID(id int) (*Incident, error)
	GetAll() ([]*Incident, error)
	Update(incident *Incident) error
	Delete(id int) error
}

// AIService defines the interface for AI-powered incident analysis
type AIService interface {
	AnalyzeIncident(title, description, affectedService string) (*IncidentAnalysis, error)
}

// IncidentAnalysis represents the AI-generated analysis of an incident
type IncidentAnalysis struct {
	Severity string `json:"severity"`
	Category string `json:"category"`
}
