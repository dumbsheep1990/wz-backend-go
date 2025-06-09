package learn

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
)

// MySQLTeacherRepository implements the TeacherRepository interface for MySQL
type MySQLTeacherRepository struct {
	db *sqlx.DB
}

// NewMySQLTeacherRepository creates a new instance of MySQLTeacherRepository
func NewMySQLTeacherRepository(db *sqlx.DB) *MySQLTeacherRepository {
	return &MySQLTeacherRepository{
		db: db,
	}
}

// Create adds a new teacher to the database
func (r *MySQLTeacherRepository) Create(ctx context.Context, teacher *entity.Teacher) error {
	query := `
		INSERT INTO learn_teachers (
			id, user_id, name, avatar, title, introduction, specialties,
			status, courses_count, students_count, rating, rating_count, 
			contact_email, contact_phone, social_profiles, created_at, updated_at
		) VALUES (
			:id, :user_id, :name, :avatar, :title, :introduction, :specialties,
			:status, :courses_count, :students_count, :rating, :rating_count, 
			:contact_email, :contact_phone, :social_profiles, :created_at, :updated_at
		)
	`
	
	// Convert slice to string for MySQL
	specialtiesStr := strings.Join(teacher.Specialties, ",")
	socialProfilesStr := strings.Join(teacher.SocialProfiles, ",")
	
	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":              teacher.ID,
		"user_id":         teacher.UserID,
		"name":            teacher.Name,
		"avatar":          teacher.Avatar,
		"title":           teacher.Title,
		"introduction":    teacher.Introduction,
		"specialties":     specialtiesStr,
		"status":          string(teacher.Status),
		"courses_count":   teacher.CoursesCount,
		"students_count":  teacher.StudentsCount,
		"rating":          teacher.Rating,
		"rating_count":    teacher.RatingCount,
		"contact_email":   teacher.ContactEmail,
		"contact_phone":   teacher.ContactPhone,
		"social_profiles": socialProfilesStr,
		"created_at":      teacher.CreatedAt,
		"updated_at":      teacher.UpdatedAt,
	})
	
	return err
}

// GetByID retrieves a teacher by their ID
func (r *MySQLTeacherRepository) GetByID(ctx context.Context, id string) (*entity.Teacher, error) {
	query := `
		SELECT 
			id, user_id, name, avatar, title, introduction, specialties,
			status, courses_count, students_count, rating, rating_count, 
			contact_email, contact_phone, social_profiles, created_at, updated_at
		FROM learn_teachers
		WHERE id = ?
	`
	
	var teacherDB struct {
		ID             string    `db:"id"`
		UserID         string    `db:"user_id"`
		Name           string    `db:"name"`
		Avatar         string    `db:"avatar"`
		Title          string    `db:"title"`
		Introduction   string    `db:"introduction"`
		Specialties    string    `db:"specialties"`
		Status         string    `db:"status"`
		CoursesCount   int       `db:"courses_count"`
		StudentsCount  int       `db:"students_count"`
		Rating         float64   `db:"rating"`
		RatingCount    int       `db:"rating_count"`
		ContactEmail   string    `db:"contact_email"`
		ContactPhone   string    `db:"contact_phone"`
		SocialProfiles string    `db:"social_profiles"`
		CreatedAt      time.Time `db:"created_at"`
		UpdatedAt      time.Time `db:"updated_at"`
	}
	
	err := r.db.GetContext(ctx, &teacherDB, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("teacher not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}
	
	// Convert strings back to slices
	var specialties []string
	if teacherDB.Specialties != "" {
		specialties = strings.Split(teacherDB.Specialties, ",")
	} else {
		specialties = []string{}
	}
	
	var socialProfiles []string
	if teacherDB.SocialProfiles != "" {
		socialProfiles = strings.Split(teacherDB.SocialProfiles, ",")
	} else {
		socialProfiles = []string{}
	}
	
	teacher := &entity.Teacher{
		ID:             teacherDB.ID,
		UserID:         teacherDB.UserID,
		Name:           teacherDB.Name,
		Avatar:         teacherDB.Avatar,
		Title:          teacherDB.Title,
		Introduction:   teacherDB.Introduction,
		Specialties:    specialties,
		Status:         entity.TeacherStatus(teacherDB.Status),
		CoursesCount:   teacherDB.CoursesCount,
		StudentsCount:  teacherDB.StudentsCount,
		Rating:         teacherDB.Rating,
		RatingCount:    teacherDB.RatingCount,
		ContactEmail:   teacherDB.ContactEmail,
		ContactPhone:   teacherDB.ContactPhone,
		SocialProfiles: socialProfiles,
		CreatedAt:      teacherDB.CreatedAt,
		UpdatedAt:      teacherDB.UpdatedAt,
	}
	
	return teacher, nil
}

// GetByUserID retrieves a teacher by their user ID
func (r *MySQLTeacherRepository) GetByUserID(ctx context.Context, userID string) (*entity.Teacher, error) {
	query := `
		SELECT 
			id, user_id, name, avatar, title, introduction, specialties,
			status, courses_count, students_count, rating, rating_count, 
			contact_email, contact_phone, social_profiles, created_at, updated_at
		FROM learn_teachers
		WHERE user_id = ?
	`
	
	var teacherDB struct {
		ID             string    `db:"id"`
		UserID         string    `db:"user_id"`
		Name           string    `db:"name"`
		Avatar         string    `db:"avatar"`
		Title          string    `db:"title"`
		Introduction   string    `db:"introduction"`
		Specialties    string    `db:"specialties"`
		Status         string    `db:"status"`
		CoursesCount   int       `db:"courses_count"`
		StudentsCount  int       `db:"students_count"`
		Rating         float64   `db:"rating"`
		RatingCount    int       `db:"rating_count"`
		ContactEmail   string    `db:"contact_email"`
		ContactPhone   string    `db:"contact_phone"`
		SocialProfiles string    `db:"social_profiles"`
		CreatedAt      time.Time `db:"created_at"`
		UpdatedAt      time.Time `db:"updated_at"`
	}
	
	err := r.db.GetContext(ctx, &teacherDB, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("teacher not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}
	
	// Convert strings back to slices
	var specialties []string
	if teacherDB.Specialties != "" {
		specialties = strings.Split(teacherDB.Specialties, ",")
	} else {
		specialties = []string{}
	}
	
	var socialProfiles []string
	if teacherDB.SocialProfiles != "" {
		socialProfiles = strings.Split(teacherDB.SocialProfiles, ",")
	} else {
		socialProfiles = []string{}
	}
	
	teacher := &entity.Teacher{
		ID:             teacherDB.ID,
		UserID:         teacherDB.UserID,
		Name:           teacherDB.Name,
		Avatar:         teacherDB.Avatar,
		Title:          teacherDB.Title,
		Introduction:   teacherDB.Introduction,
		Specialties:    specialties,
		Status:         entity.TeacherStatus(teacherDB.Status),
		CoursesCount:   teacherDB.CoursesCount,
		StudentsCount:  teacherDB.StudentsCount,
		Rating:         teacherDB.Rating,
		RatingCount:    teacherDB.RatingCount,
		ContactEmail:   teacherDB.ContactEmail,
		ContactPhone:   teacherDB.ContactPhone,
		SocialProfiles: socialProfiles,
		CreatedAt:      teacherDB.CreatedAt,
		UpdatedAt:      teacherDB.UpdatedAt,
	}
	
	return teacher, nil
}

// Update updates an existing teacher in the database
func (r *MySQLTeacherRepository) Update(ctx context.Context, teacher *entity.Teacher) error {
	query := `
		UPDATE learn_teachers
		SET 
			name = :name,
			avatar = :avatar,
			title = :title,
			introduction = :introduction,
			specialties = :specialties,
			status = :status,
			courses_count = :courses_count,
			students_count = :students_count,
			rating = :rating,
			rating_count = :rating_count,
			contact_email = :contact_email,
			contact_phone = :contact_phone,
			social_profiles = :social_profiles,
			updated_at = :updated_at
		WHERE id = :id
	`
	
	// Convert slices to strings for MySQL
	specialtiesStr := strings.Join(teacher.Specialties, ",")
	socialProfilesStr := strings.Join(teacher.SocialProfiles, ",")
	
	result, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":              teacher.ID,
		"name":            teacher.Name,
		"avatar":          teacher.Avatar,
		"title":           teacher.Title,
		"introduction":    teacher.Introduction,
		"specialties":     specialtiesStr,
		"status":          string(teacher.Status),
		"courses_count":   teacher.CoursesCount,
		"students_count":  teacher.StudentsCount,
		"rating":          teacher.Rating,
		"rating_count":    teacher.RatingCount,
		"contact_email":   teacher.ContactEmail,
		"contact_phone":   teacher.ContactPhone,
		"social_profiles": socialProfilesStr,
		"updated_at":      time.Now(),
	})
	
	if err != nil {
		return fmt.Errorf("failed to update teacher: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("no teacher found with id: %s", teacher.ID)
	}
	
	return nil
}

// Delete removes a teacher from the database
func (r *MySQLTeacherRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM learn_teachers WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete teacher: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("no teacher found with id: %s", id)
	}
	
	return nil
}

// List retrieves a paginated list of teachers based on query parameters
func (r *MySQLTeacherRepository) List(ctx context.Context, params repository.TeacherQueryParams) ([]*entity.Teacher, int64, error) {
	baseQuery := `
		SELECT 
			id, user_id, name, avatar, title, introduction, specialties,
			status, courses_count, students_count, rating, rating_count, 
			contact_email, contact_phone, social_profiles, created_at, updated_at
		FROM learn_teachers
		WHERE 1=1
	`
	
	countQuery := `SELECT COUNT(*) FROM learn_teachers WHERE 1=1`
	
	args := []interface{}{}
	
	// Apply filters
	if params.Status != nil {
		baseQuery += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, string(*params.Status))
	}
	
	if params.Specialty != nil && *params.Specialty != "" {
		baseQuery += " AND specialties LIKE ?"
		countQuery += " AND specialties LIKE ?"
		args = append(args, "%"+*params.Specialty+"%")
	}
	
	// Apply sorting
	if params.SortBy != "" {
		var direction string
		if strings.ToLower(params.SortOrder) == "desc" {
			direction = "DESC"
		} else {
			direction = "ASC"
		}
		
		// Validate sort field to prevent SQL injection
		allowedSortFields := map[string]bool{
			"name": true, "rating": true, "courses_count": true,
			"students_count": true, "created_at": true,
		}
		
		if allowedSortFields[params.SortBy] {
			baseQuery += fmt.Sprintf(" ORDER BY %s %s", params.SortBy, direction)
		} else {
			baseQuery += " ORDER BY created_at DESC"
		}
	} else {
		baseQuery += " ORDER BY created_at DESC"
	}
	
	// Apply pagination
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		baseQuery += " LIMIT ? OFFSET ?"
		args = append(args, params.PageSize, offset)
	}
	
	// Get total count
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count teachers: %w", err)
	}
	
	// Get teachers
	var teachersDB []struct {
		ID             string    `db:"id"`
		UserID         string    `db:"user_id"`
		Name           string    `db:"name"`
		Avatar         string    `db:"avatar"`
		Title          string    `db:"title"`
		Introduction   string    `db:"introduction"`
		Specialties    string    `db:"specialties"`
		Status         string    `db:"status"`
		CoursesCount   int       `db:"courses_count"`
		StudentsCount  int       `db:"students_count"`
		Rating         float64   `db:"rating"`
		RatingCount    int       `db:"rating_count"`
		ContactEmail   string    `db:"contact_email"`
		ContactPhone   string    `db:"contact_phone"`
		SocialProfiles string    `db:"social_profiles"`
		CreatedAt      time.Time `db:"created_at"`
		UpdatedAt      time.Time `db:"updated_at"`
	}
	
	err = r.db.SelectContext(ctx, &teachersDB, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list teachers: %w", err)
	}
	
	// Convert to domain entities
	teachers := make([]*entity.Teacher, 0, len(teachersDB))
	for _, t := range teachersDB {
		// Convert strings to slices
		var specialties []string
		if t.Specialties != "" {
			specialties = strings.Split(t.Specialties, ",")
		} else {
			specialties = []string{}
		}
		
		var socialProfiles []string
		if t.SocialProfiles != "" {
			socialProfiles = strings.Split(t.SocialProfiles, ",")
		} else {
			socialProfiles = []string{}
		}
		
		teachers = append(teachers, &entity.Teacher{
			ID:             t.ID,
			UserID:         t.UserID,
			Name:           t.Name,
			Avatar:         t.Avatar,
			Title:          t.Title,
			Introduction:   t.Introduction,
			Specialties:    specialties,
			Status:         entity.TeacherStatus(t.Status),
			CoursesCount:   t.CoursesCount,
			StudentsCount:  t.StudentsCount,
			Rating:         t.Rating,
			RatingCount:    t.RatingCount,
			ContactEmail:   t.ContactEmail,
			ContactPhone:   t.ContactPhone,
			SocialProfiles: socialProfiles,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
		})
	}
	
	return teachers, total, nil
}

// ListByIDs retrieves teachers by their IDs
func (r *MySQLTeacherRepository) ListByIDs(ctx context.Context, ids []string) ([]*entity.Teacher, error) {
	if len(ids) == 0 {
		return []*entity.Teacher{}, nil
	}
	
	query := `
		SELECT 
			id, user_id, name, avatar, title, introduction, specialties,
			status, courses_count, students_count, rating, rating_count, 
			contact_email, contact_phone, social_profiles, created_at, updated_at
		FROM learn_teachers
		WHERE id IN (?)
	`
	
	// Use sqlx.In to properly handle the IN clause
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	
	// Rebind the query for the specific database
	query = r.db.Rebind(query)
	
	var teachersDB []struct {
		ID             string    `db:"id"`
		UserID         string    `db:"user_id"`
		Name           string    `db:"name"`
		Avatar         string    `db:"avatar"`
		Title          string    `db:"title"`
		Introduction   string    `db:"introduction"`
		Specialties    string    `db:"specialties"`
		Status         string    `db:"status"`
		CoursesCount   int       `db:"courses_count"`
		StudentsCount  int       `db:"students_count"`
		Rating         float64   `db:"rating"`
		RatingCount    int       `db:"rating_count"`
		ContactEmail   string    `db:"contact_email"`
		ContactPhone   string    `db:"contact_phone"`
		SocialProfiles string    `db:"social_profiles"`
		CreatedAt      time.Time `db:"created_at"`
		UpdatedAt      time.Time `db:"updated_at"`
	}
	
	err = r.db.SelectContext(ctx, &teachersDB, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list teachers by IDs: %w", err)
	}
	
	// Convert to domain entities
	teachers := make([]*entity.Teacher, 0, len(teachersDB))
	for _, t := range teachersDB {
		// Convert strings to slices
		var specialties []string
		if t.Specialties != "" {
			specialties = strings.Split(t.Specialties, ",")
		} else {
			specialties = []string{}
		}
		
		var socialProfiles []string
		if t.SocialProfiles != "" {
			socialProfiles = strings.Split(t.SocialProfiles, ",")
		} else {
			socialProfiles = []string{}
		}
		
		teachers = append(teachers, &entity.Teacher{
			ID:             t.ID,
			UserID:         t.UserID,
			Name:           t.Name,
			Avatar:         t.Avatar,
			Title:          t.Title,
			Introduction:   t.Introduction,
			Specialties:    specialties,
			Status:         entity.TeacherStatus(t.Status),
			CoursesCount:   t.CoursesCount,
			StudentsCount:  t.StudentsCount,
			Rating:         t.Rating,
			RatingCount:    t.RatingCount,
			ContactEmail:   t.ContactEmail,
			ContactPhone:   t.ContactPhone,
			SocialProfiles: socialProfiles,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
		})
	}
	
	return teachers, nil
}

// ListPopular retrieves the most popular teachers
func (r *MySQLTeacherRepository) ListPopular(ctx context.Context, limit int) ([]*entity.Teacher, error) {
	query := `
		SELECT 
			id, user_id, name, avatar, title, introduction, specialties,
			status, courses_count, students_count, rating, rating_count, 
			contact_email, contact_phone, social_profiles, created_at, updated_at
		FROM learn_teachers
		WHERE status = ?
		ORDER BY rating DESC, students_count DESC
		LIMIT ?
	`
	
	var teachersDB []struct {
		ID             string    `db:"id"`
		UserID         string    `db:"user_id"`
		Name           string    `db:"name"`
		Avatar         string    `db:"avatar"`
		Title          string    `db:"title"`
		Introduction   string    `db:"introduction"`
		Specialties    string    `db:"specialties"`
		Status         string    `db:"status"`
		CoursesCount   int       `db:"courses_count"`
		StudentsCount  int       `db:"students_count"`
		Rating         float64   `db:"rating"`
		RatingCount    int       `db:"rating_count"`
		ContactEmail   string    `db:"contact_email"`
		ContactPhone   string    `db:"contact_phone"`
		SocialProfiles string    `db:"social_profiles"`
		CreatedAt      time.Time `db:"created_at"`
		UpdatedAt      time.Time `db:"updated_at"`
	}
	
	err := r.db.SelectContext(ctx, &teachersDB, query, string(entity.TeacherStatusActive), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list popular teachers: %w", err)
	}
	
	// Convert to domain entities
	teachers := make([]*entity.Teacher, 0, len(teachersDB))
	for _, t := range teachersDB {
		// Convert strings to slices
		var specialties []string
		if t.Specialties != "" {
			specialties = strings.Split(t.Specialties, ",")
		} else {
			specialties = []string{}
		}
		
		var socialProfiles []string
		if t.SocialProfiles != "" {
			socialProfiles = strings.Split(t.SocialProfiles, ",")
		} else {
			socialProfiles = []string{}
		}
		
		teachers = append(teachers, &entity.Teacher{
			ID:             t.ID,
			UserID:         t.UserID,
			Name:           t.Name,
			Avatar:         t.Avatar,
			Title:          t.Title,
			Introduction:   t.Introduction,
			Specialties:    specialties,
			Status:         entity.TeacherStatus(t.Status),
			CoursesCount:   t.CoursesCount,
			StudentsCount:  t.StudentsCount,
			Rating:         t.Rating,
			RatingCount:    t.RatingCount,
			ContactEmail:   t.ContactEmail,
			ContactPhone:   t.ContactPhone,
			SocialProfiles: socialProfiles,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
		})
	}
	
	return teachers, nil
}

// CountAll returns the total number of teachers
func (r *MySQLTeacherRepository) CountAll(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM learn_teachers`
	
	var count int64
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count teachers: %w", err)
	}
	
	return count, nil
}

// CountByStatus returns the number of teachers with a specific status
func (r *MySQLTeacherRepository) CountByStatus(ctx context.Context, status entity.TeacherStatus) (int64, error) {
	query := `SELECT COUNT(*) FROM learn_teachers WHERE status = ?`
	
	var count int64
	err := r.db.GetContext(ctx, &count, query, string(status))
	if err != nil {
		return 0, fmt.Errorf("failed to count teachers by status: %w", err)
	}
	
	return count, nil
}

// Search searches for teachers based on a keyword and filters
func (r *MySQLTeacherRepository) Search(ctx context.Context, keyword string, params repository.TeacherQueryParams) ([]*entity.Teacher, int64, error) {
	baseQuery := `
		SELECT 
			id, user_id, name, avatar, title, introduction, specialties,
			status, courses_count, students_count, rating, rating_count, 
			contact_email, contact_phone, social_profiles, created_at, updated_at
		FROM learn_teachers
		WHERE (name LIKE ? OR title LIKE ? OR introduction LIKE ? OR specialties LIKE ?)
	`
	
	countQuery := `
		SELECT COUNT(*)
		FROM learn_teachers
		WHERE (name LIKE ? OR title LIKE ? OR introduction LIKE ? OR specialties LIKE ?)
	`
	
	// Prepare search terms
	searchTerm := "%" + keyword + "%"
	args := []interface{}{searchTerm, searchTerm, searchTerm, searchTerm}
	countArgs := []interface{}{searchTerm, searchTerm, searchTerm, searchTerm}
	
	// Apply filters
	if params.Status != nil {
		baseQuery += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, string(*params.Status))
		countArgs = append(countArgs, string(*params.Status))
	}
	
	if params.Specialty != nil && *params.Specialty != "" {
		baseQuery += " AND specialties LIKE ?"
		countQuery += " AND specialties LIKE ?"
		specialtyTerm := "%" + *params.Specialty + "%"
		args = append(args, specialtyTerm)
		countArgs = append(countArgs, specialtyTerm)
	}
	
	// Apply sorting
	if params.SortBy != "" {
		var direction string
		if strings.ToLower(params.SortOrder) == "desc" {
			direction = "DESC"
		} else {
			direction = "ASC"
		}
		
		// Validate sort field to prevent SQL injection
		allowedSortFields := map[string]bool{
			"name": true, "rating": true, "courses_count": true,
			"students_count": true, "created_at": true,
		}
		
		if allowedSortFields[params.SortBy] {
			baseQuery += fmt.Sprintf(" ORDER BY %s %s", params.SortBy, direction)
		} else {
			baseQuery += " ORDER BY created_at DESC"
		}
	} else {
		baseQuery += " ORDER BY created_at DESC"
	}
	
	// Apply pagination
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		baseQuery += " LIMIT ? OFFSET ?"
		args = append(args, params.PageSize, offset)
	}
	
	// Get total count
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, countArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}
	
	// Get teachers
	var teachersDB []struct {
		ID             string    `db:"id"`
		UserID         string    `db:"user_id"`
		Name           string    `db:"name"`
		Avatar         string    `db:"avatar"`
		Title          string    `db:"title"`
		Introduction   string    `db:"introduction"`
		Specialties    string    `db:"specialties"`
		Status         string    `db:"status"`
		CoursesCount   int       `db:"courses_count"`
		StudentsCount  int       `db:"students_count"`
		Rating         float64   `db:"rating"`
		RatingCount    int       `db:"rating_count"`
		ContactEmail   string    `db:"contact_email"`
		ContactPhone   string    `db:"contact_phone"`
		SocialProfiles string    `db:"social_profiles"`
		CreatedAt      time.Time `db:"created_at"`
		UpdatedAt      time.Time `db:"updated_at"`
	}
	
	err = r.db.SelectContext(ctx, &teachersDB, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search teachers: %w", err)
	}
	
	// Convert to domain entities
	teachers := make([]*entity.Teacher, 0, len(teachersDB))
	for _, t := range teachersDB {
		// Convert strings to slices
		var specialties []string
		if t.Specialties != "" {
			specialties = strings.Split(t.Specialties, ",")
		} else {
			specialties = []string{}
		}
		
		var socialProfiles []string
		if t.SocialProfiles != "" {
			socialProfiles = strings.Split(t.SocialProfiles, ",")
		} else {
			socialProfiles = []string{}
		}
		
		teachers = append(teachers, &entity.Teacher{
			ID:             t.ID,
			UserID:         t.UserID,
			Name:           t.Name,
			Avatar:         t.Avatar,
			Title:          t.Title,
			Introduction:   t.Introduction,
			Specialties:    specialties,
			Status:         entity.TeacherStatus(t.Status),
			CoursesCount:   t.CoursesCount,
			StudentsCount:  t.StudentsCount,
			Rating:         t.Rating,
			RatingCount:    t.RatingCount,
			ContactEmail:   t.ContactEmail,
			ContactPhone:   t.ContactPhone,
			SocialProfiles: socialProfiles,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
		})
	}
	
	return teachers, total, nil
}
