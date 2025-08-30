package repository

import (
	"database/sql"
	"fmt"
	"incident-triage-assistant/internal/domain"
)

// MySQLIncidentRepository implements the IncidentRepository interface using MySQL
type MySQLIncidentRepository struct {
	db *sql.DB
}

// NewMySQLIncidentRepository creates a new MySQL incident repository
func NewMySQLIncidentRepository(db *sql.DB) *MySQLIncidentRepository {
	return &MySQLIncidentRepository{db: db}
}

// Create inserts a new incident into the database
func (r *MySQLIncidentRepository) Create(incident *domain.Incident) error {
	query := `
		INSERT INTO incidents (title, description, affected_service, ai_severity, ai_category, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.Exec(query,
		incident.Title,
		incident.Description,
		incident.AffectedService,
		incident.AISeverity,
		incident.AICategory,
		incident.CreatedAt,
		incident.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create incident: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	incident.ID = int(id)
	return nil
}

// GetByID retrieves an incident by its ID
func (r *MySQLIncidentRepository) GetByID(id int) (*domain.Incident, error) {
	query := `
		SELECT id, title, description, affected_service, ai_severity, ai_category, created_at, updated_at
		FROM incidents WHERE id = ?
	`
	
	incident := &domain.Incident{}
	err := r.db.QueryRow(query, id).Scan(
		&incident.ID,
		&incident.Title,
		&incident.Description,
		&incident.AffectedService,
		&incident.AISeverity,
		&incident.AICategory,
		&incident.CreatedAt,
		&incident.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("incident not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to get incident: %w", err)
	}

	return incident, nil
}

// GetAll retrieves all incidents from the database
func (r *MySQLIncidentRepository) GetAll() ([]*domain.Incident, error) {
	query := `
		SELECT id, title, description, affected_service, ai_severity, ai_category, created_at, updated_at
		FROM incidents ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query incidents: %w", err)
	}
	defer rows.Close()

	var incidents []*domain.Incident
	for rows.Next() {
		incident := &domain.Incident{}
		err := rows.Scan(
			&incident.ID,
			&incident.Title,
			&incident.Description,
			&incident.AffectedService,
			&incident.AISeverity,
			&incident.AICategory,
			&incident.CreatedAt,
			&incident.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan incident: %w", err)
		}
		incidents = append(incidents, incident)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating incidents: %w", err)
	}

	return incidents, nil
}

// Update updates an existing incident in the database
func (r *MySQLIncidentRepository) Update(incident *domain.Incident) error {
	query := `
		UPDATE incidents 
		SET title = ?, description = ?, affected_service = ?, ai_severity = ?, ai_category = ?, updated_at = ?
		WHERE id = ?
	`
	
	result, err := r.db.Exec(query,
		incident.Title,
		incident.Description,
		incident.AffectedService,
		incident.AISeverity,
		incident.AICategory,
		incident.UpdatedAt,
		incident.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update incident: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("incident not found with id %d", incident.ID)
	}

	return nil
}

// Delete removes an incident from the database
func (r *MySQLIncidentRepository) Delete(id int) error {
	query := `DELETE FROM incidents WHERE id = ?`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete incident: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("incident not found with id %d", id)
	}

	return nil
}
