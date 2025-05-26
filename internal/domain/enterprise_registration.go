package domain

import (
	"context"
	"time"

	"github.com/wz-project/wz-backend-go/internal/domain/model"
)

// EnterpriseRegistrationService defines the service interface for enterprise registration
type EnterpriseRegistrationService interface {
	// CreateEnterpriseRegistration creates a new enterprise registration
	CreateEnterpriseRegistration(ctx context.Context, req *model.EnterpriseRegistrationRequest) (*model.EnterpriseRegistration, error)
	
	// GetEnterpriseRegistration retrieves enterprise registration by userID
	GetEnterpriseRegistration(ctx context.Context, userID int64) (*model.EnterpriseRegistration, error)
	
	// UpdateEnterpriseRegistration updates enterprise registration
	UpdateEnterpriseRegistration(ctx context.Context, registration *model.EnterpriseRegistration) error
	
	// VerifyEnterprise verifies an enterprise using verification code
	VerifyEnterprise(ctx context.Context, userID int64, verificationCode string) error
}

// EnterpriseRegistrationRepository defines the repository interface for enterprise registration
type EnterpriseRegistrationRepository interface {
	// Create creates a new enterprise registration
	Create(ctx context.Context, registration *model.EnterpriseRegistration) error
	
	// FindByUserID finds enterprise registration by userID
	FindByUserID(ctx context.Context, userID int64) (*model.EnterpriseRegistration, error)
	
	// Update updates enterprise registration
	Update(ctx context.Context, registration *model.EnterpriseRegistration) error
}
