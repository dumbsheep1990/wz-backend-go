package service

import (
	"context"
	"errors"
	"time"

	"github.com/wz-project/wz-backend-go/internal/domain"
	"github.com/wz-project/wz-backend-go/internal/domain/model"
)

// EnterpriseRegistrationService implements domain.EnterpriseRegistrationService
type EnterpriseRegistrationService struct {
	repo domain.EnterpriseRegistrationRepository
}

// NewEnterpriseRegistrationService creates a new enterprise registration service
func NewEnterpriseRegistrationService(repo domain.EnterpriseRegistrationRepository) *EnterpriseRegistrationService {
	return &EnterpriseRegistrationService{
		repo: repo,
	}
}

// CreateEnterpriseRegistration creates a new enterprise registration
func (s *EnterpriseRegistrationService) CreateEnterpriseRegistration(ctx context.Context, req *model.EnterpriseRegistrationRequest) (*model.EnterpriseRegistration, error) {
	// Check if registration already exists
	existing, err := s.repo.FindByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("企业入驻申请已存在")
	}

	// Create new registration
	registration := &model.EnterpriseRegistration{
		UserID:             req.UserID,
		CompanyName:        req.CompanyName,
		CompanyType:        req.CompanyType,
		ContactPerson:      req.ContactPerson,
		JobPosition:        req.JobPosition,
		Region:             req.Region,
		VerificationMethod: req.VerificationMethod,
		DetailedAddress:    req.DetailedAddress,
		LocationLatitude:   req.LocationLatitude,
		LocationLongitude:  req.LocationLongitude,
		Status:             0, // Pending
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	err = s.repo.Create(ctx, registration)
	if err != nil {
		return nil, err
	}

	return registration, nil
}

// GetEnterpriseRegistration retrieves enterprise registration by userID
func (s *EnterpriseRegistrationService) GetEnterpriseRegistration(ctx context.Context, userID int64) (*model.EnterpriseRegistration, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// UpdateEnterpriseRegistration updates enterprise registration
func (s *EnterpriseRegistrationService) UpdateEnterpriseRegistration(ctx context.Context, registration *model.EnterpriseRegistration) error {
	// Check if registration exists
	existing, err := s.repo.FindByUserID(ctx, registration.UserID)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("企业入驻申请不存在")
	}

	// Update registration
	return s.repo.Update(ctx, registration)
}

// VerifyEnterprise verifies an enterprise using verification code
func (s *EnterpriseRegistrationService) VerifyEnterprise(ctx context.Context, userID int64, verificationCode string) error {
	// Get registration
	registration, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if registration == nil {
		return errors.New("企业入驻申请不存在")
	}

	// In a real implementation, we would verify the code here
	// For now, we'll just update the status
	
	registration.Status = 1 // Approved
	registration.UpdatedAt = time.Now()

	return s.repo.Update(ctx, registration)
}
