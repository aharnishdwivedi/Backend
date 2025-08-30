package repository

import (
	"database/sql"
	"testing"
	"time"

	"incident-triage-assistant/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMySQLIncidentRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLIncidentRepository(db)

	incident := &domain.Incident{
		Title:           "Test Incident",
		Description:     "Test Description",
		AffectedService: "Test Service",
		AISeverity:      "Medium",
		AICategory:      "Software",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mock.ExpectExec("INSERT INTO incidents").
		WithArgs(incident.Title, incident.Description, incident.AffectedService, incident.AISeverity, incident.AICategory, incident.CreatedAt, incident.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(incident)
	assert.NoError(t, err)
	assert.Equal(t, 1, incident.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLIncidentRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLIncidentRepository(db)

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

	rows := sqlmock.NewRows([]string{"id", "title", "description", "affected_service", "ai_severity", "ai_category", "created_at", "updated_at"}).
		AddRow(expectedIncident.ID, expectedIncident.Title, expectedIncident.Description, expectedIncident.AffectedService, expectedIncident.AISeverity, expectedIncident.AICategory, expectedIncident.CreatedAt, expectedIncident.UpdatedAt)

	mock.ExpectQuery("SELECT id, title, description, affected_service, ai_severity, ai_category, created_at, updated_at FROM incidents WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	incident, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedIncident, incident)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLIncidentRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLIncidentRepository(db)

	mock.ExpectQuery("SELECT id, title, description, affected_service, ai_severity, ai_category, created_at, updated_at FROM incidents WHERE id = ?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	incident, err := repo.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, incident)
	assert.Contains(t, err.Error(), "incident not found with id 999")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLIncidentRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLIncidentRepository(db)

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

	rows := sqlmock.NewRows([]string{"id", "title", "description", "affected_service", "ai_severity", "ai_category", "created_at", "updated_at"})
	for _, incident := range expectedIncidents {
		rows.AddRow(incident.ID, incident.Title, incident.Description, incident.AffectedService, incident.AISeverity, incident.AICategory, incident.CreatedAt, incident.UpdatedAt)
	}

	mock.ExpectQuery("SELECT id, title, description, affected_service, ai_severity, ai_category, created_at, updated_at FROM incidents ORDER BY created_at DESC").
		WillReturnRows(rows)

	incidents, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, expectedIncidents, incidents)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLIncidentRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLIncidentRepository(db)

	incident := &domain.Incident{
		ID:              1,
		Title:           "Updated Incident",
		Description:     "Updated Description",
		AffectedService: "Updated Service",
		AISeverity:      "High",
		AICategory:      "Network",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mock.ExpectExec("UPDATE incidents SET title = \\?, description = \\?, affected_service = \\?, ai_severity = \\?, ai_category = \\?, updated_at = \\? WHERE id = \\?").
		WithArgs(incident.Title, incident.Description, incident.AffectedService, incident.AISeverity, incident.AICategory, incident.UpdatedAt, incident.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(incident)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLIncidentRepository_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLIncidentRepository(db)

	incident := &domain.Incident{
		ID:              999,
		Title:           "Updated Incident",
		Description:     "Updated Description",
		AffectedService: "Updated Service",
		AISeverity:      "High",
		AICategory:      "Network",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mock.ExpectExec("UPDATE incidents SET title = \\?, description = \\?, affected_service = \\?, ai_severity = \\?, ai_category = \\?, updated_at = \\? WHERE id = \\?").
		WithArgs(incident.Title, incident.Description, incident.AffectedService, incident.AISeverity, incident.AICategory, incident.UpdatedAt, incident.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Update(incident)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incident not found with id 999")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLIncidentRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLIncidentRepository(db)

	mock.ExpectExec("DELETE FROM incidents WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLIncidentRepository_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLIncidentRepository(db)

	mock.ExpectExec("DELETE FROM incidents WHERE id = ?").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incident not found with id 999")
	assert.NoError(t, mock.ExpectationsWereMet())
}
