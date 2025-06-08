package entity

import (
	"time"
)

// Store represents a seller's store in the commerce system
type Store struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	OwnerID     string    `db:"owner_id"`
	LogoURL     string    `db:"logo_url"`
	Province    string    `db:"province"`
	City        string    `db:"city"`
	District    string    `db:"district"`
	Address     string    `db:"address"`
	ContactName string    `db:"contact_name"`
	Phone       string    `db:"phone"`
	Rating      float64   `db:"rating"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// NewStore creates a new store
func NewStore(id, name, description, ownerID, logoURL, province, city, district, address, contactName, phone string) *Store {
	now := time.Now()
	return &Store{
		ID:          id,
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		LogoURL:     logoURL,
		Province:    province,
		City:        city,
		District:    district,
		Address:     address,
		ContactName: contactName,
		Phone:       phone,
		Rating:      0.0,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Activate sets the store as active
func (s *Store) Activate() {
	s.IsActive = true
	s.UpdatedAt = time.Now()
}

// Deactivate sets the store as inactive
func (s *Store) Deactivate() {
	s.IsActive = false
	s.UpdatedAt = time.Now()
}

// UpdateRating updates the store rating
func (s *Store) UpdateRating(rating float64) {
	s.Rating = rating
	s.UpdatedAt = time.Now()
}

// UpdateDetails updates store details
func (s *Store) UpdateDetails(name, description, logoURL, province, city, district, address, contactName, phone string) {
	s.Name = name
	s.Description = description
	s.LogoURL = logoURL
	s.Province = province
	s.City = city
	s.District = district
	s.Address = address
	s.ContactName = contactName
	s.Phone = phone
	s.UpdatedAt = time.Now()
}
