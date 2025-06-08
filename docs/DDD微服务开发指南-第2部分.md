# WanZhi平台 - DDD微服务开发指南

## 三、微服务开发流程

### 1. 领域建模

领域建模是整个DDD开发流程的核心和起点，它直接影响后续所有开发工作。

#### 1.1 关键步骤

1. **识别限界上下文**：
   - 确定微服务边界，如学习微服务、交易微服务等
   - 定义上下文地图，明确各微服务之间的关系

2. **提取核心实体**：
   - 分析业务需求，识别核心业务概念
   - 对于学习微服务，核心实体包括：Course（课程）、Teacher（教师）、Enrollment（报名）、Category（分类）、Chapter（章节）、Lesson（课时）

3. **建立领域模型**：
   - 确定实体间的关系（如一对多、多对多）
   - 识别聚合根（如Course可作为聚合根，包含多个Chapter）
   - 定义值对象（如Address、ContactInfo等）

4. **设计仓储接口**：
   - 为每个聚合根定义仓储接口
   - 设计符合业务需求的查询方法
   - 示例：`CourseRepository`、`TeacherRepository`等

#### 1.2 命名规范

- **实体**：使用名词，如`Course`、`Teacher`
- **值对象**：使用描述性名词，如`CourseStatus`、`TeacherContactInfo`
- **仓储接口**：实体名 + Repository，如`CourseRepository`
- **领域服务**：当操作涉及多个实体或复杂业务规则时，使用动词 + 名词，如`ChapterLessonService`

### 2. 领域层实现

#### 2.1 实体定义

创建各实体结构体，包含其属性和行为：

```go
// internal/domain/learn/entity/course.go
package entity

type Course struct {
    ID          string
    Title       string
    Description string
    // 其他属性...
    Chapters    []*Chapter
    Status      CourseStatus
}

// 实体方法 - 行为
func (c *Course) Publish() error {
    // 业务规则：必须有至少一章内容才能发布
    if len(c.Chapters) == 0 {
        return errors.New("课程必须有至少一章内容才能发布")
    }
    c.Status = StatusPublished
    return nil
}

// 其他领域行为方法...
```

#### 2.2 仓储接口

定义实体持久化和检索的接口：

```go
// internal/domain/learn/repository/course_repository.go
package repository

import (
    "context"
    "wz-backend-go/internal/domain/learn/entity"
)

type CourseRepository interface {
    Save(ctx context.Context, course *entity.Course) error
    FindByID(ctx context.Context, id string) (*entity.Course, error)
    // 其他查询方法...
}
```

#### 2.3 领域服务

处理跨实体的复杂领域逻辑：

```go
// internal/domain/learn/service/chapter_lesson_service.go
package service

import (
    "context"
    "wz-backend-go/internal/domain/learn/entity"
    "wz-backend-go/internal/domain/learn/repository"
)

type ChapterLessonService struct {
    chapterRepo repository.ChapterRepository
    lessonRepo  repository.LessonRepository
    courseRepo  repository.CourseRepository
}

func NewChapterLessonService(
    chapterRepo repository.ChapterRepository,
    lessonRepo repository.LessonRepository,
    courseRepo repository.CourseRepository,
) *ChapterLessonService {
    return &ChapterLessonService{
        chapterRepo: chapterRepo,
        lessonRepo:  lessonRepo,
        courseRepo:  courseRepo,
    }
}

// 领域服务方法
func (s *ChapterLessonService) AddLessonToChapter(
    ctx context.Context, 
    chapterID string, 
    lesson *entity.Lesson,
) error {
    // 实现复杂业务逻辑...
    return nil
}

// 其他服务方法...
```

#### 2.4 领域层DTO

定义领域层数据传输对象，用于跨边界的数据交换：

```go
// internal/domain/learn/dto/course_dto.go
package dto

type CourseDTO struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    // 其他必要的字段...
}
```

### 3. 应用层实现

应用层是连接领域层和交付层的桥梁，负责用例编排、事务管理和DTO转换。

#### 3.1 应用服务

开发处理特定用例的应用服务：

```go
// internal/application/learn/course_app_service.go
package learn

import (
    "context"
    "wz-backend-go/internal/domain/learn/dto"
    "wz-backend-go/internal/domain/learn/service"
)

type CourseAppService struct {
    courseService       *service.CourseService
    teacherService      *service.TeacherService
    categoryService     *service.CategoryService
    chapterLessonService *service.ChapterLessonService
    enrollmentService   *service.EnrollmentService
}

func NewCourseAppService(
    courseService *service.CourseService,
    teacherService *service.TeacherService,
    categoryService *service.CategoryService,
    chapterLessonService *service.ChapterLessonService,
    enrollmentService *service.EnrollmentService,
) *CourseAppService {
    return &CourseAppService{
        courseService:       courseService,
        teacherService:      teacherService,
        categoryService:     categoryService,
        chapterLessonService: chapterLessonService,
        enrollmentService:   enrollmentService,
    }
}

// 应用服务方法
func (s *CourseAppService) CreateCourse(
    ctx context.Context, 
    req *dto.CreateCourseRequest,
) (*dto.CourseDTO, error) {
    // 1. 参数验证（应在交付层处理）
    // 2. 调用领域服务创建课程
    // 3. 转换为DTO返回
    // ...实现业务流程
    return &dto.CourseDTO{}, nil
}

// 其他应用服务方法...
```

#### 3.2 应用层DTO

定义应用层数据传输对象，用于API请求和响应：

```go
// internal/application/learn/dto/course_dto.go
package dto

// 请求DTO
type CreateCourseRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description" binding:"required"`
    TeacherID   string `json:"teacher_id" binding:"required"`
    CategoryID  string `json:"category_id" binding:"required"`
    // 其他字段...
}

// 响应DTO
type CourseResponse struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    TeacherName string `json:"teacher_name"`
    Category    string `json:"category"`
    // 其他字段...
}

// 其他DTO定义...
```

### 4. 交付层实现

交付层处理外部请求并调用应用服务，提供API接口。

#### 4.1 HTTP处理器

实现REST API的处理函数：

```go
// internal/delivery/http/internal/handler/learn/course_handler.go
package learn

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    
    "wz-backend-go/internal/application/learn"
    "wz-backend-go/internal/application/learn/dto"
)

type CourseHandler struct {
    courseAppService *learn.CourseAppService
}

func NewCourseHandler(courseAppService *learn.CourseAppService) *CourseHandler {
    return &CourseHandler{
        courseAppService: courseAppService,
    }
}

// CreateCourse 创建课程
func (h *CourseHandler) CreateCourse(c *gin.Context) {
    var req dto.CreateCourseRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    courseDTO, err := h.courseAppService.CreateCourse(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, courseDTO)
}

// 其他处理方法...
```

#### 4.2 路由配置

注册HTTP路由：

```go
// internal/delivery/http/internal/router/learn_router.go
package router

import (
    "github.com/gin-gonic/gin"
    
    "wz-backend-go/internal/application/learn"
    "wz-backend-go/internal/delivery/http/internal/handler/learn"
    "wz-backend-go/internal/delivery/http/internal/middleware"
)

// SetupLearnRoutes 配置学习微服务的路由
func SetupLearnRoutes(
    r *gin.Engine,
    courseAppService *learn.CourseAppService,
    teacherAppService *learn.TeacherAppService,
    enrollmentAppService *learn.EnrollmentAppService,
    categoryAppService *learn.CategoryAppService,
    chapterLessonAppService *learn.ChapterLessonAppService,
) {
    // 初始化处理器
    courseHandler := learn.NewCourseHandler(courseAppService)
    teacherHandler := learn.NewTeacherHandler(teacherAppService)
    enrollmentHandler := learn.NewEnrollmentHandler(enrollmentAppService)
    categoryHandler := learn.NewCategoryHandler(categoryAppService)
    chapterLessonHandler := learn.NewChapterLessonHandler(chapterLessonAppService)
    
    // API路由组
    apiV1 := r.Group("/api/v1/learn")
    
    // 课程相关路由
    courses := apiV1.Group("/courses")
    {
        courses.POST("", middleware.AdminAuth(), courseHandler.CreateCourse)
        courses.GET("", courseHandler.ListCourses)
        courses.GET("/:id", courseHandler.GetCourseByID)
        // 其他路由...
    }
    
    // 其他路由组配置...
}
```

### 5. 基础设施层实现

基础设施层提供技术实现，包括仓储实现、第三方集成等。

#### 5.1 仓储实现

基于具体存储技术实现仓储接口：

```go
// internal/infrastructure/repository/mysql/course_repository.go
package mysql

import (
    "context"
    "database/sql"
    
    "wz-backend-go/internal/domain/learn/entity"
    "wz-backend-go/internal/domain/learn/repository"
)

type CourseRepository struct {
    db *sql.DB
}

func NewCourseRepository(db *sql.DB) repository.CourseRepository {
    return &CourseRepository{
        db: db,
    }
}

func (r *CourseRepository) Save(ctx context.Context, course *entity.Course) error {
    // 实现保存逻辑
    return nil
}

func (r *CourseRepository) FindByID(ctx context.Context, id string) (*entity.Course, error) {
    // 实现查询逻辑
    return nil, nil
}

// 其他仓储方法实现...
```

### 6. 微服务入口实现

创建微服务启动入口，进行依赖注入和配置：

```go
// cmd/learn/main.go
package main

import (
    "context"
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
    
    // 设置Gin路由
    r := gin.Default()
    
    // 配置中间件
    r.Use(gin.Recovery())
    r.Use(gin.Logger())
    
    // 注册路由
    router.SetupLearnRoutes(r, courseAppService, teacherAppService, enrollmentAppService, 
        categoryAppService, chapterLessonAppService)
    
    // 启动服务器
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }
    
    // 优雅关闭
    // ...实现优雅关闭逻辑
}
```
