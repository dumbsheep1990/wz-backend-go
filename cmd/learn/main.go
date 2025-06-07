package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"wz-backend-go/internal/application/learn"
	"wz-backend-go/internal/delivery/http/internal/router"
	"wz-backend-go/internal/domain/learn/service"
	"wz-backend-go/internal/infrastructure/repository/mysql"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到 .env 文件")
	}

	// 初始化数据库连接
	db, err := mysql.NewMySQLConnection()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 初始化仓储
	courseRepo := mysql.NewCourseRepository(db)
	teacherRepo := mysql.NewTeacherRepository(db)
	categoryRepo := mysql.NewCategoryRepository(db)
	chapterRepo := mysql.NewChapterRepository(db)
	lessonRepo := mysql.NewLessonRepository(db)
	enrollmentRepo := mysql.NewEnrollmentRepository(db)

	// 初始化领域服务
	courseService := service.NewCourseService(courseRepo)
	teacherService := service.NewTeacherService(teacherRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	chapterLessonService := service.NewChapterLessonService(chapterRepo, lessonRepo, courseRepo)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseRepo)

	// 初始化应用服务
	courseAppService := learn.NewCourseAppService(
		courseService, teacherService, categoryService, chapterLessonService, enrollmentService)
	teacherAppService := learn.NewTeacherAppService(teacherService, courseService)
	enrollmentAppService := learn.NewEnrollmentAppService(enrollmentService, courseService)
	categoryAppService := learn.NewCategoryAppService(categoryService)
	chapterLessonAppService := learn.NewChapterLessonAppService(chapterLessonService)

	// 设置 Gin 路由
	r := gin.Default()

	// 配置中间件
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// 注册路由
	router.SetupLearnRoutes(r, courseAppService, teacherAppService, enrollmentAppService, 
		categoryAppService, chapterLessonAppService)

	// 获取服务端口
	port := os.Getenv("LEARN_SERVICE_PORT")
	if port == "" {
		port = "8082" // 默认端口
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// 在单独的goroutine中启动服务器
	go func() {
		log.Printf("学习微服务启动于端口 %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听失败: %v\n", err)
		}
	}()

	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭学习微服务...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v\n", err)
	}

	log.Println("学习微服务已关闭")
}
