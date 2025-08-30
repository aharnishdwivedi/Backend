package usecase

import (
	"incident-triage-assistant/internal/domain"
	"time"
)

// IncidentUseCase implements the business logic for incident management
type IncidentUseCase struct {
	incidentRepo domain.IncidentRepository
	aiService    domain.AIService
}

// NewIncidentUseCase creates a new instance of IncidentUseCase
func NewIncidentUseCase(incidentRepo domain.IncidentRepository, aiService domain.AIService) *IncidentUseCase {
	return &IncidentUseCase{
		incidentRepo: incidentRepo,
		aiService:    aiService,
	}
}

// CreateIncident creates a new incident with AI analysis
func (uc *IncidentUseCase) CreateIncident(req *domain.CreateIncidentRequest) (*domain.Incident, error) {
	// Analyze incident using AI
	analysis, err := uc.aiService.AnalyzeIncident(req.Title, req.Description, req.AffectedService)
	if err != nil {
		return nil, err
	}

	// Create incident with AI insights
	incident := &domain.Incident{
		Title:           req.Title,
		Description:     req.Description,
		AffectedService: req.AffectedService,
		AISeverity:      analysis.Severity,
		AICategory:      analysis.Category,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save to repository
	err = uc.incidentRepo.Create(incident)
	if err != nil {
		return nil, err
	}

	return incident, nil
}

// GetIncident retrieves an incident by ID
func (uc *IncidentUseCase) GetIncident(id int) (*domain.Incident, error) {
	return uc.incidentRepo.GetByID(id)
}

// GetAllIncidents retrieves all incidents
func (uc *IncidentUseCase) GetAllIncidents() ([]*domain.Incident, error) {
	return uc.incidentRepo.GetAll()
}

// UpdateIncident updates an existing incident
func (uc *IncidentUseCase) UpdateIncident(id int, req *domain.CreateIncidentRequest) (*domain.Incident, error) {
	// Get existing incident
	incident, err := uc.incidentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Re-analyze with AI if content changed
	analysis, err := uc.aiService.AnalyzeIncident(req.Title, req.Description, req.AffectedService)
	if err != nil {
		return nil, err
	}

	// Update fields
	incident.Title = req.Title
	incident.Description = req.Description
	incident.AffectedService = req.AffectedService
	incident.AISeverity = analysis.Severity
	incident.AICategory = analysis.Category
	incident.UpdatedAt = time.Now()

	// Save to repository
	err = uc.incidentRepo.Update(incident)
	if err != nil {
		return nil, err
	}

	return incident, nil
}

// DeleteIncident deletes an incident by ID
func (uc *IncidentUseCase) DeleteIncident(id int) error {
	return uc.incidentRepo.Delete(id)
}
