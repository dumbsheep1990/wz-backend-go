package entity

import (
	"time"

	"github.com/google/uuid"
)

// TeacherStatus 讲师状态
type TeacherStatus string

const (
	TeacherStatusActive    TeacherStatus = "active"    // 活跃
	TeacherStatusInactive  TeacherStatus = "inactive"  // 非活跃
	TeacherStatusSuspended TeacherStatus = "suspended" // 已暂停
)

// Teacher 讲师实体
type Teacher struct {
	ID             string        `json:"id"`
	UserID         string        `json:"userId"`         // 关联的用户ID
	Name           string        `json:"name"`           // 讲师姓名
	Avatar         string        `json:"avatar"`         // 头像
	Title          string        `json:"title"`          // 职称/头衔
	Introduction   string        `json:"introduction"`   // 简介
	Specialties    []string      `json:"specialties"`    // 专长领域
	Status         TeacherStatus `json:"status"`         // 讲师状态
	CoursesCount   int           `json:"coursesCount"`   // 课程数量
	StudentsCount  int           `json:"studentsCount"`  // 学生数量
	Rating         float64       `json:"rating"`         // 评分
	RatingCount    int           `json:"ratingCount"`    // 评价数量
	ContactEmail   string        `json:"contactEmail"`   // 联系邮箱
	ContactPhone   string        `json:"contactPhone"`   // 联系电话
	SocialProfiles []string      `json:"socialProfiles"` // 社交档案
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
}

// NewTeacher 创建新讲师
func NewTeacher(userID, name string) *Teacher {
	now := time.Now()
	return &Teacher{
		ID:             uuid.New().String(),
		UserID:         userID,
		Name:           name,
		Status:         TeacherStatusActive,
		Specialties:    make([]string, 0),
		CoursesCount:   0,
		StudentsCount:  0,
		Rating:         0,
		RatingCount:    0,
		SocialProfiles: make([]string, 0),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Update 更新讲师基本信息
func (t *Teacher) Update(name, avatar, title, introduction string) {
	t.Name = name
	t.Avatar = avatar
	t.Title = title
	t.Introduction = introduction
	t.UpdatedAt = time.Now()
}

// UpdateContact 更新联系信息
func (t *Teacher) UpdateContact(email, phone string) {
	t.ContactEmail = email
	t.ContactPhone = phone
	t.UpdatedAt = time.Now()
}

// SetSpecialties 设置专长领域
func (t *Teacher) SetSpecialties(specialties []string) {
	t.Specialties = specialties
	t.UpdatedAt = time.Now()
}

// SetSocialProfiles 设置社交档案
func (t *Teacher) SetSocialProfiles(profiles []string) {
	t.SocialProfiles = profiles
	t.UpdatedAt = time.Now()
}

// Activate 激活讲师
func (t *Teacher) Activate() {
	t.Status = TeacherStatusActive
	t.UpdatedAt = time.Now()
}

// Deactivate 停用讲师
func (t *Teacher) Deactivate() {
	t.Status = TeacherStatusInactive
	t.UpdatedAt = time.Now()
}

// Suspend 暂停讲师
func (t *Teacher) Suspend() {
	t.Status = TeacherStatusSuspended
	t.UpdatedAt = time.Now()
}

// IncrementCourseCount 增加课程数
func (t *Teacher) IncrementCourseCount() {
	t.CoursesCount++
	t.UpdatedAt = time.Now()
}

// DecrementCourseCount 减少课程数
func (t *Teacher) DecrementCourseCount() {
	if t.CoursesCount > 0 {
		t.CoursesCount--
		t.UpdatedAt = time.Now()
	}
}

// AddStudentsCount 增加学生数
func (t *Teacher) AddStudentsCount(count int) {
	t.StudentsCount += count
	t.UpdatedAt = time.Now()
}

// AddRating 添加评分
func (t *Teacher) AddRating(score float64) {
	totalScore := t.Rating * float64(t.RatingCount)
	t.RatingCount++
	t.Rating = (totalScore + score) / float64(t.RatingCount)
	t.UpdatedAt = time.Now()
}
