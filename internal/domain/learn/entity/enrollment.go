package entity

import (
	"time"

	"github.com/google/uuid"
)

// EnrollmentStatus 报名状态
type EnrollmentStatus string

const (
	EnrollmentStatusActive    EnrollmentStatus = "active"    // 活跃
	EnrollmentStatusCompleted EnrollmentStatus = "completed" // 已完成
	EnrollmentStatusExpired   EnrollmentStatus = "expired"   // 已过期
	EnrollmentStatusRefunded  EnrollmentStatus = "refunded"  // 已退款
)

// Enrollment 课程报名实体
type Enrollment struct {
	ID             string           `json:"id"`
	CourseID       string           `json:"courseId"`       // 课程ID
	UserID         string           `json:"userId"`         // 用户ID
	OrderID        string           `json:"orderId"`        // 订单ID
	Status         EnrollmentStatus `json:"status"`         // 报名状态
	Progress       float64          `json:"progress"`       // 学习进度(百分比)
	CompletedCount int              `json:"completedCount"` // 已完成课时数
	TotalCount     int              `json:"totalCount"`     // 总课时数
	LastLearnTime  *time.Time       `json:"lastLearnTime"`  // 最后学习时间
	Rating         *float64         `json:"rating"`         // 用户评分
	Comment        string           `json:"comment"`        // 用户评价
	ExpiresAt      *time.Time       `json:"expiresAt"`      // 过期时间
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	CompletedAt    *time.Time       `json:"completedAt"`
}

// NewEnrollment 创建新报名记录
func NewEnrollment(courseID, userID, orderID string, lessonCount int, expiresAt *time.Time) *Enrollment {
	now := time.Now()
	return &Enrollment{
		ID:             uuid.New().String(),
		CourseID:       courseID,
		UserID:         userID,
		OrderID:        orderID,
		Status:         EnrollmentStatusActive,
		Progress:       0,
		CompletedCount: 0,
		TotalCount:     lessonCount,
		ExpiresAt:      expiresAt,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// UpdateProgress 更新学习进度
func (e *Enrollment) UpdateProgress(completedCount int) {
	e.CompletedCount = completedCount
	if e.TotalCount > 0 {
		e.Progress = float64(completedCount) / float64(e.TotalCount) * 100
	}
	now := time.Now()
	e.LastLearnTime = &now
	e.UpdatedAt = now

	// 如果全部完成，更新状态
	if e.Progress >= 100 {
		e.Status = EnrollmentStatusCompleted
		e.CompletedAt = &now
	}
}

// Complete 标记为已完成
func (e *Enrollment) Complete() {
	e.Status = EnrollmentStatusCompleted
	e.Progress = 100
	now := time.Now()
	e.CompletedAt = &now
	e.UpdatedAt = now
}

// Expire 标记为已过期
func (e *Enrollment) Expire() {
	e.Status = EnrollmentStatusExpired
	e.UpdatedAt = time.Now()
}

// Refund 标记为已退款
func (e *Enrollment) Refund() {
	e.Status = EnrollmentStatusRefunded
	e.UpdatedAt = time.Now()
}

// AddRating 添加评分和评价
func (e *Enrollment) AddRating(score float64, comment string) {
	e.Rating = &score
	e.Comment = comment
	e.UpdatedAt = time.Now()
}
