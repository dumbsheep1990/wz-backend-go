package entity

import (
	"time"

	"github.com/google/uuid"
)

// ProgressStatus 学习进度状态
type ProgressStatus string

const (
	ProgressStatusNotStarted ProgressStatus = "not_started" // 未开始
	ProgressStatusInProgress ProgressStatus = "in_progress" // 学习中
	ProgressStatusCompleted  ProgressStatus = "completed"   // 已完成
)

// Progress 学习进度实体
type Progress struct {
	ID               string         `json:"id"`
	UserID           string         `json:"userId"`           // 用户ID
	CourseID         string         `json:"courseId"`         // 课程ID
	LessonID         string         `json:"lessonId"`         // 课时ID
	Status           ProgressStatus `json:"status"`           // 进度状态
	WatchedDuration  int            `json:"watchedDuration"`  // 已观看时长(秒)
	TotalDuration    int            `json:"totalDuration"`    // 总时长(秒)
	CompletionRate   float64        `json:"completionRate"`   // 完成率(0-1)
	LastWatchedAt    *time.Time     `json:"lastWatchedAt"`    // 最后观看时间
	CompletedAt      *time.Time     `json:"completedAt"`      // 完成时间
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

// NewProgress 创建新学习进度
func NewProgress(userID, courseID, lessonID string, totalDuration int) *Progress {
	now := time.Now()
	return &Progress{
		ID:              uuid.New().String(),
		UserID:          userID,
		CourseID:        courseID,
		LessonID:        lessonID,
		Status:          ProgressStatusNotStarted,
		WatchedDuration: 0,
		TotalDuration:   totalDuration,
		CompletionRate:  0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// UpdateProgress 更新学习进度
func (p *Progress) UpdateProgress(watchedDuration int) {
	p.WatchedDuration = watchedDuration
	if p.TotalDuration > 0 {
		p.CompletionRate = float64(watchedDuration) / float64(p.TotalDuration)
		if p.CompletionRate > 1 {
			p.CompletionRate = 1
		}
	}
	
	now := time.Now()
	p.LastWatchedAt = &now
	p.UpdatedAt = now
	
	// 更新状态
	if p.WatchedDuration > 0 && p.Status == ProgressStatusNotStarted {
		p.Status = ProgressStatusInProgress
	}
	
	// 如果完成率达到90%以上，标记为已完成
	if p.CompletionRate >= 0.9 && p.Status != ProgressStatusCompleted {
		p.Complete()
	}
}

// Complete 标记为已完成
func (p *Progress) Complete() {
	p.Status = ProgressStatusCompleted
	p.CompletionRate = 1
	now := time.Now()
	p.CompletedAt = &now
	p.UpdatedAt = now
}

// Reset 重置进度
func (p *Progress) Reset() {
	p.Status = ProgressStatusNotStarted
	p.WatchedDuration = 0
	p.CompletionRate = 0
	p.LastWatchedAt = nil
	p.CompletedAt = nil
	p.UpdatedAt = time.Now()
}

// IsCompleted 是否已完成
func (p *Progress) IsCompleted() bool {
	return p.Status == ProgressStatusCompleted
}

// IsStarted 是否已开始
func (p *Progress) IsStarted() bool {
	return p.Status != ProgressStatusNotStarted
}

// GetProgressPercentage 获取进度百分比
func (p *Progress) GetProgressPercentage() int {
	return int(p.CompletionRate * 100)
}
