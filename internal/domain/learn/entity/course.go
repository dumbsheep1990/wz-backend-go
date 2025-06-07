package entity

import (
	"time"

	"github.com/google/uuid"
)

// CourseStatus 课程状态
type CourseStatus string

const (
	CourseStatusDraft     CourseStatus = "draft"     // 草稿
	CourseStatusPublished CourseStatus = "published" // 已发布
	CourseStatusArchived  CourseStatus = "archived"  // 已归档
)

// CourseLevel 课程难度等级
type CourseLevel string

const (
	CourseLevelBeginner     CourseLevel = "beginner"     // 初级
	CourseLevelIntermediate CourseLevel = "intermediate" // 中级
	CourseLevelAdvanced     CourseLevel = "advanced"     // 高级
)

// Course 课程实体
type Course struct {
	ID              string       `json:"id"`
	Title           string       `json:"title"`            // 课程标题
	Subtitle        string       `json:"subtitle"`         // 副标题
	Description     string       `json:"description"`      // 课程描述
	TeacherID       string       `json:"teacherId"`        // 讲师ID
	Cover           string       `json:"cover"`            // 课程封面
	Level           CourseLevel  `json:"level"`            // 课程难度
	Duration        int          `json:"duration"`         // 总时长(分钟)
	Price           float64      `json:"price"`            // 价格
	DiscountPrice   float64      `json:"discountPrice"`    // 折扣价
	Status          CourseStatus `json:"status"`           // 课程状态
	CategoryIDs     []string     `json:"categoryIds"`      // 分类ID列表
	Tags            []string     `json:"tags"`             // 标签列表
	ChaptersCount   int          `json:"chaptersCount"`    // 章节数
	LessonsCount    int          `json:"lessonsCount"`     // 课时数
	EnrollmentsCount int         `json:"enrollmentsCount"` // 报名人数
	Rating          float64      `json:"rating"`           // 评分
	RatingCount     int          `json:"ratingCount"`      // 评分人数
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	PublishedAt     *time.Time   `json:"publishedAt"`
}

// NewCourse 创建新课程
func NewCourse(title, teacherID string) *Course {
	now := time.Now()
	return &Course{
		ID:            uuid.New().String(),
		Title:         title,
		TeacherID:     teacherID,
		Status:        CourseStatusDraft,
		Level:         CourseLevelBeginner,
		CreatedAt:     now,
		UpdatedAt:     now,
		CategoryIDs:   make([]string, 0),
		Tags:          make([]string, 0),
		ChaptersCount: 0,
		LessonsCount:  0,
		Rating:        0,
		RatingCount:   0,
	}
}

// Publish 发布课程
func (c *Course) Publish() {
	c.Status = CourseStatusPublished
	now := time.Now()
	c.PublishedAt = &now
	c.UpdatedAt = now
}

// Archive 归档课程
func (c *Course) Archive() {
	c.Status = CourseStatusArchived
	c.UpdatedAt = time.Now()
}

// Update 更新课程基本信息
func (c *Course) Update(title, subtitle, description, cover string, level CourseLevel, price, discountPrice float64) {
	c.Title = title
	c.Subtitle = subtitle
	c.Description = description
	c.Cover = cover
	c.Level = level
	c.Price = price
	c.DiscountPrice = discountPrice
	c.UpdatedAt = time.Now()
}

// SetCategories 设置课程分类
func (c *Course) SetCategories(categoryIDs []string) {
	c.CategoryIDs = categoryIDs
	c.UpdatedAt = time.Now()
}

// SetTags 设置课程标签
func (c *Course) SetTags(tags []string) {
	c.Tags = tags
	c.UpdatedAt = time.Now()
}

// AddEnrollment 增加报名人数
func (c *Course) AddEnrollment() {
	c.EnrollmentsCount++
	c.UpdatedAt = time.Now()
}

// AddRating 添加评分
func (c *Course) AddRating(score float64) {
	totalScore := c.Rating * float64(c.RatingCount)
	c.RatingCount++
	c.Rating = (totalScore + score) / float64(c.RatingCount)
	c.UpdatedAt = time.Now()
}

// UpdateLessonCounts 更新课时计数
func (c *Course) UpdateLessonCounts(chaptersCount, lessonsCount int) {
	c.ChaptersCount = chaptersCount
	c.LessonsCount = lessonsCount
	c.UpdatedAt = time.Now()
}
