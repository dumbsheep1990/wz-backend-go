package dto

import (
	"time"
)

// LearnStats 学习服务统计数据DTO
type LearnStats struct {
	TotalCourses     int64     `json:"totalCourses"`     // 总课程数
	PublishedCourses int64     `json:"publishedCourses"` // 已发布课程
	DraftCourses     int64     `json:"draftCourses"`     // 草稿课程
	TotalTeachers    int64     `json:"totalTeachers"`    // 总讲师数
	TotalEnrollments int64     `json:"totalEnrollments"` // 总报名数
	TotalCategories  int64     `json:"totalCategories"`  // 总分类数
	TotalLessons     int64     `json:"totalLessons"`     // 总课时数
	WeeklyRevenue    float64   `json:"weeklyRevenue"`    // 周收益
	MonthlyRevenue   float64   `json:"monthlyRevenue"`   // 月收益
	TotalRevenue     float64   `json:"totalRevenue"`     // 总收益
	RecentActivity   int64     `json:"recentActivity"`   // 最近活跃度(30天内)
	LastUpdateTime   time.Time `json:"lastUpdateTime"`   // 最后更新时间
}

// CourseBasicDTO 课程基础信息DTO
type CourseBasicDTO struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Subtitle       string    `json:"subtitle"`
	TeacherID      string    `json:"teacherId"`
	TeacherName    string    `json:"teacherName"`
	TeacherAvatar  string    `json:"teacherAvatar"`
	Cover          string    `json:"cover"`
	Price          float64   `json:"price"`
	DiscountPrice  float64   `json:"discountPrice"`
	Level          string    `json:"level"`
	Duration       int       `json:"duration"`
	LessonsCount   int       `json:"lessonsCount"`
	EnrollmentCount int      `json:"enrollmentCount"`
	Rating         float64   `json:"rating"`
	CategoryNames  []string  `json:"categoryNames"`
	Tags           []string  `json:"tags"`
	PublishedAt    time.Time `json:"publishedAt"`
}

// CourseDetailDTO 课程详细信息DTO
type CourseDetailDTO struct {
	CourseBasicDTO
	Description      string           `json:"description"`
	ChapterList      []ChapterDTO     `json:"chapterList"`
	FreeLessons      []LessonBasicDTO `json:"freeLessons"`
	RelatedCourses   []CourseBasicDTO `json:"relatedCourses"`
	TeacherInfo      TeacherBasicDTO  `json:"teacherInfo"`
	RecentEnrollments int             `json:"recentEnrollments"`
	UpdatedAt        time.Time        `json:"updatedAt"`
}

// ChapterDTO 章节DTO
type ChapterDTO struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Order       int            `json:"order"`
	Duration    int            `json:"duration"`
	LessonCount int            `json:"lessonCount"`
	Lessons     []LessonBasicDTO `json:"lessons"`
}

// LessonBasicDTO 课时基础信息DTO
type LessonBasicDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Order       int       `json:"order"`
	Duration    int       `json:"duration"`
	IsFree      bool      `json:"isFree"`
	ChapterID   string    `json:"chapterId"`
	ChapterTitle string   `json:"chapterTitle"`
}

// LessonDetailDTO 课时详细信息DTO
type LessonDetailDTO struct {
	LessonBasicDTO
	VideoURL      string    `json:"videoUrl"`
	ArticleContent string   `json:"articleContent"`
	AudioURL      string    `json:"audioUrl"`
	NextLessonID  string    `json:"nextLessonId"`
	PrevLessonID  string    `json:"prevLessonId"`
}

// TeacherBasicDTO 讲师基础信息DTO
type TeacherBasicDTO struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Avatar        string   `json:"avatar"`
	Title         string   `json:"title"`
	Introduction  string   `json:"introduction"`
	CoursesCount  int      `json:"coursesCount"`
	StudentsCount int      `json:"studentsCount"`
	Rating        float64  `json:"rating"`
	Specialties   []string `json:"specialties"`
}

// TeacherDetailDTO 讲师详细信息DTO
type TeacherDetailDTO struct {
	TeacherBasicDTO
	Courses       []CourseBasicDTO `json:"courses"`
	ContactEmail  string           `json:"contactEmail"`
	SocialProfiles []string        `json:"socialProfiles"`
}

// EnrollmentDTO 报名信息DTO
type EnrollmentDTO struct {
	ID             string    `json:"id"`
	CourseID       string    `json:"courseId"`
	CourseTitle    string    `json:"courseTitle"`
	CourseCover    string    `json:"courseCover"`
	TeacherName    string    `json:"teacherName"`
	Progress       float64   `json:"progress"`
	Status         string    `json:"status"`
	EnrolledAt     time.Time `json:"enrolledAt"`
	ExpiresAt      *time.Time `json:"expiresAt"`
	LastLearnTime  *time.Time `json:"lastLearnTime"`
	CompletedCount int       `json:"completedCount"`
	TotalCount     int       `json:"totalCount"`
}

// CategoryDTO 分类信息DTO
type CategoryDTO struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Icon         string       `json:"icon"`
	ParentID     *string      `json:"parentId"`
	Level        int          `json:"level"`
	Order        int          `json:"order"`
	CoursesCount int          `json:"coursesCount"`
	Children     []CategoryDTO `json:"children,omitempty"`
}

// ReviewDTO 评价信息DTO
type ReviewDTO struct {
	ID         string     `json:"id"`
	UserID     string     `json:"userId"`
	CourseID   string     `json:"courseId"`
	Rating     int        `json:"rating"`
	Content    string     `json:"content"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	ApprovedAt *time.Time `json:"approvedAt"`
}

// CourseRatingStatsDTO 课程评分统计DTO
type CourseRatingStatsDTO struct {
	CourseID      string             `json:"courseId"`
	AverageRating float64            `json:"averageRating"`
	TotalCount    int64              `json:"totalCount"`
	Distribution  map[int]int64      `json:"distribution"` // 评分分布：1-5星的数量
}

// ProgressDTO 学习进度DTO
type ProgressDTO struct {
	ID              string     `json:"id"`
	UserID          string     `json:"userId"`
	CourseID        string     `json:"courseId"`
	LessonID        string     `json:"lessonId"`
	Status          string     `json:"status"`
	WatchedDuration int        `json:"watchedDuration"` // 已观看时长（秒）
	TotalDuration   int        `json:"totalDuration"`   // 总时长（秒）
	CompletionRate  float64    `json:"completionRate"`  // 完成率（0-1）
	ProgressPercent int        `json:"progressPercent"` // 进度百分比（0-100）
	LastWatchedAt   *time.Time `json:"lastWatchedAt"`
	CompletedAt     *time.Time `json:"completedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// CourseProgressStatsDTO 课程学习进度统计DTO
type CourseProgressStatsDTO struct {
	CourseID        string  `json:"courseId"`
	UserID          string  `json:"userId"`
	OverallProgress float64 `json:"overallProgress"` // 整体进度（0-1）
	CompletedCount  int64   `json:"completedCount"`  // 已完成课时数
	InProgressCount int64   `json:"inProgressCount"` // 学习中课时数
	TotalCount      int64   `json:"totalCount"`      // 总课时数
	CompletionRate  float64 `json:"completionRate"`  // 完成率（0-1）
}

// UserProgressStatsDTO 用户整体学习进度统计DTO
type UserProgressStatsDTO struct {
	UserID          string  `json:"userId"`
	OverallProgress float64 `json:"overallProgress"` // 整体学习进度（0-1）
	CompletedCount  int64   `json:"completedCount"`  // 已完成课时总数
	InProgressCount int64   `json:"inProgressCount"` // 学习中课时总数
	TotalCount      int64   `json:"totalCount"`      // 总学习课时数
}
