package entity

import (
	"time"

	"github.com/google/uuid"
)

// LessonType 课时类型
type LessonType string

const (
	LessonTypeVideo   LessonType = "video"   // 视频
	LessonTypeArticle LessonType = "article" // 文章
	LessonTypeAudio   LessonType = "audio"   // 音频
)

// LessonStatus 课时状态
type LessonStatus string

const (
	LessonStatusDraft     LessonStatus = "draft"     // 草稿
	LessonStatusPublished LessonStatus = "published" // 已发布
)

// Lesson 课时实体
type Lesson struct {
	ID          string       `json:"id"`
	CourseID    string       `json:"courseId"`    // 所属课程ID
	ChapterID   string       `json:"chapterId"`   // 所属章节ID
	Title       string       `json:"title"`       // 课时标题
	Description string       `json:"description"` // 课时描述
	Type        LessonType   `json:"type"`        // 课时类型
	Status      LessonStatus `json:"status"`      // 课时状态
	Order       int          `json:"order"`       // 课时顺序
	Duration    int          `json:"duration"`    // 课时时长(分钟)
	VideoURL    string       `json:"videoUrl"`    // 视频URL
	VideoSize   int64        `json:"videoSize"`   // 视频大小(字节)
	ArticleContent string    `json:"articleContent"` // 文章内容
	AudioURL    string       `json:"audioUrl"`    // 音频URL
	AudioSize   int64        `json:"audioSize"`   // 音频大小(字节)
	IsFree      bool         `json:"isFree"`      // 是否免费
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	PublishedAt *time.Time   `json:"publishedAt"`
}

// NewLesson 创建新课时
func NewLesson(courseID, chapterID, title string, order int, lessonType LessonType) *Lesson {
	now := time.Now()
	return &Lesson{
		ID:        uuid.New().String(),
		CourseID:  courseID,
		ChapterID: chapterID,
		Title:     title,
		Type:      lessonType,
		Status:    LessonStatusDraft,
		Order:     order,
		IsFree:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update 更新课时基本信息
func (l *Lesson) Update(title, description string, order int, isFree bool) {
	l.Title = title
	l.Description = description
	l.Order = order
	l.IsFree = isFree
	l.UpdatedAt = time.Now()
}

// Publish 发布课时
func (l *Lesson) Publish() {
	l.Status = LessonStatusPublished
	now := time.Now()
	l.PublishedAt = &now
	l.UpdatedAt = now
}

// SetVideo 设置视频信息
func (l *Lesson) SetVideo(videoURL string, duration int, size int64) {
	l.VideoURL = videoURL
	l.Duration = duration
	l.VideoSize = size
	l.UpdatedAt = time.Now()
}

// SetAudio 设置音频信息
func (l *Lesson) SetAudio(audioURL string, duration int, size int64) {
	l.AudioURL = audioURL
	l.Duration = duration
	l.AudioSize = size
	l.UpdatedAt = time.Now()
}

// SetArticle 设置文章内容
func (l *Lesson) SetArticle(content string) {
	l.ArticleContent = content
	l.UpdatedAt = time.Now()
}
