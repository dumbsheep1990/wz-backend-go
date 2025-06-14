package entity

import (
	"time"

	"github.com/google/uuid"
)

// ReviewStatus 评价状态
type ReviewStatus string

const (
	ReviewStatusPending   ReviewStatus = "pending"   // 待审核
	ReviewStatusApproved  ReviewStatus = "approved"  // 已通过
	ReviewStatusRejected  ReviewStatus = "rejected"  // 已拒绝
)

// Review 课程评价实体
type Review struct {
	ID        string       `json:"id"`
	UserID    string       `json:"userId"`    // 评价用户ID
	CourseID  string       `json:"courseId"`  // 课程ID
	Rating    int          `json:"rating"`    // 评分(1-5)
	Content   string       `json:"content"`   // 评价内容
	Status    ReviewStatus `json:"status"`    // 评价状态
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	ApprovedAt *time.Time  `json:"approvedAt"`
}

// NewReview 创建新评价
func NewReview(userID, courseID string, rating int, content string) *Review {
	now := time.Now()
	return &Review{
		ID:        uuid.New().String(),
		UserID:    userID,
		CourseID:  courseID,
		Rating:    rating,
		Content:   content,
		Status:    ReviewStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update 更新评价内容
func (r *Review) Update(rating int, content string) {
	r.Rating = rating
	r.Content = content
	r.UpdatedAt = time.Now()
}

// Approve 通过评价
func (r *Review) Approve() {
	r.Status = ReviewStatusApproved
	now := time.Now()
	r.ApprovedAt = &now
	r.UpdatedAt = now
}

// Reject 拒绝评价
func (r *Review) Reject() {
	r.Status = ReviewStatusRejected
	r.UpdatedAt = time.Now()
}

// IsApproved 是否已通过
func (r *Review) IsApproved() bool {
	return r.Status == ReviewStatusApproved
}

// IsValidRating 验证评分是否有效
func (r *Review) IsValidRating() bool {
	return r.Rating >= 1 && r.Rating <= 5
}
