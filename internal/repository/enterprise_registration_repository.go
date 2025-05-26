package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/wz-project/wz-backend-go/internal/domain/model"
)

// EnterpriseRegistrationRepository implements domain.EnterpriseRegistrationRepository
type EnterpriseRegistrationRepository struct {
	db *sql.DB
}

// NewEnterpriseRegistrationRepository creates a new enterprise registration repository
func NewEnterpriseRegistrationRepository(db *sql.DB) *EnterpriseRegistrationRepository {
	return &EnterpriseRegistrationRepository{
		db: db,
	}
}

// Create inserts a new enterprise registration
func (r *EnterpriseRegistrationRepository) Create(ctx context.Context, registration *model.EnterpriseRegistration) error {
	query := `
		INSERT INTO enterprise_registrations (
			user_id, company_name, company_type, contact_person, job_position, 
			region, verification_method, detailed_address, location_latitude, 
			location_longitude, status, remark, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	registration.CreatedAt = now
	registration.UpdatedAt = now

	_, err := r.db.ExecContext(
		ctx, query,
		registration.UserID, registration.CompanyName, registration.CompanyType,
		registration.ContactPerson, registration.JobPosition, registration.Region,
		registration.VerificationMethod, registration.DetailedAddress,
		registration.LocationLatitude, registration.LocationLongitude,
		registration.Status, registration.Remark,
		registration.CreatedAt, registration.UpdatedAt,
	)

	return err
}

// FindByUserID finds enterprise registration by userID
func (r *EnterpriseRegistrationRepository) FindByUserID(ctx context.Context, userID int64) (*model.EnterpriseRegistration, error) {
	query := `
		SELECT id, user_id, company_name, company_type, contact_person, 
		       job_position, region, verification_method, detailed_address, 
		       location_latitude, location_longitude, status, remark, 
		       created_at, updated_at
		FROM enterprise_registrations
		WHERE user_id = ?
	`

	var registration model.EnterpriseRegistration
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&registration.ID, &registration.UserID, &registration.CompanyName,
		&registration.CompanyType, &registration.ContactPerson, &registration.JobPosition,
		&registration.Region, &registration.VerificationMethod, &registration.DetailedAddress,
		&registration.LocationLatitude, &registration.LocationLongitude,
		&registration.Status, &registration.Remark,
		&registration.CreatedAt, &registration.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &registration, nil
}

// Update updates enterprise registration
func (r *EnterpriseRegistrationRepository) Update(ctx context.Context, registration *model.EnterpriseRegistration) error {
	query := `
		UPDATE enterprise_registrations
		SET company_name = ?, company_type = ?, contact_person = ?, 
		    job_position = ?, region = ?, verification_method = ?, 
		    detailed_address = ?, location_latitude = ?, location_longitude = ?, 
		    status = ?, remark = ?, updated_at = ?
		WHERE id = ?
	`

	registration.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(
		ctx, query,
		registration.CompanyName, registration.CompanyType, registration.ContactPerson,
		registration.JobPosition, registration.Region, registration.VerificationMethod,
		registration.DetailedAddress, registration.LocationLatitude, registration.LocationLongitude,
		registration.Status, registration.Remark, registration.UpdatedAt,
		registration.ID,
	)

	return err
}
