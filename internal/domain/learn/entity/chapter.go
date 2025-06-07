package entity

import (
	"time"

	"github.com/google/uuid"
)

// Chapter 课程章节实体
type Chapter struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"courseId"`    // 所属课程ID
	Title       string    `json:"title"`       // 章节标题
	Description string    `json:"description"` // 章节描述
	Order       int       `json:"order"`       // 章节顺序
	LessonCount int       `json:"lessonCount"` // 课时数量
	Duration    int       `json:"duration"`    // 章节总时长(分钟)
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewChapter 创建新章节
func NewChapter(courseID, title string, order int) *Chapter {
	now := time.Now()
	return &Chapter{
		ID:          uuid.New().String(),
		CourseID:    courseID,
		Title:       title,
		Order:       order,
		LessonCount: 0,
		Duration:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update 更新章节基本信息
func (c *Chapter) Update(title, description string, order int) {
	c.Title = title
	c.Description = description
	c.Order = order
	c.UpdatedAt = time.Now()
}

// UpdateLessonStats 更新课时统计信息
func (c *Chapter) UpdateLessonStats(lessonCount, duration int) {
	c.LessonCount = lessonCount
	c.Duration = duration
	c.UpdatedAt = time.Now()
}
