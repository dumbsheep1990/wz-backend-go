package router

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/learn"
	learnHandler "wz-backend-go/internal/delivery/http/internal/handler/learn"
)

// SetupLearnRoutes 设置学习微服务的所有路由
func SetupLearnRoutes(
	router *gin.Engine,
	courseAppService *learn.CourseAppService,
	teacherAppService *learn.TeacherAppService,
	enrollmentAppService *learn.EnrollmentAppService,
	categoryAppService *learn.CategoryAppService,
	chapterLessonAppService *learn.ChapterLessonAppService,
) {
	apiGroup := router.Group("/api/v1/learn")

	// 注册课程相关路由
	courseHandler := learnHandler.NewCourseHandler(courseAppService)
	courseHandler.RegisterRoutes(apiGroup)

	// 注册讲师相关路由
	teacherHandler := learnHandler.NewTeacherHandler(teacherAppService)
	teacherHandler.RegisterRoutes(apiGroup)

	// 注册报名相关路由
	enrollmentHandler := learnHandler.NewEnrollmentHandler(enrollmentAppService)
	enrollmentHandler.RegisterRoutes(apiGroup)

	// 注册分类相关路由
	categoryHandler := learnHandler.NewCategoryHandler(categoryAppService)
	categoryHandler.RegisterRoutes(apiGroup)

	// 注册章节和课时相关路由
	chapterLessonHandler := learnHandler.NewChapterLessonHandler(chapterLessonAppService)
	chapterLessonHandler.RegisterRoutes(apiGroup)
}
