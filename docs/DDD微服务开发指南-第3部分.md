# WanZhi平台 - DDD微服务开发指南

## 四、最佳实践与规范

### 1. 代码规范

#### 1.1 命名约定

| 类型 | 规范 | 示例 |
|------|------|------|
| 实体 | 名词，首字母大写 | `Course`, `Teacher` |
| 值对象 | 描述性名词，首字母大写 | `CourseStatus`, `Address` |
| 仓储接口 | 实体名 + Repository | `CourseRepository` |
| 领域服务 | 功能描述 + Service | `CourseService` |
| 应用服务 | 实体/功能 + AppService | `CourseAppService` |
| HTTP处理器 | 实体名 + Handler | `CourseHandler` |
| 方法名 | 动词开头，驼峰式 | `CreateCourse`, `UpdateTeacherProfile` |
| DTO字段 | 小驼峰，JSON标签 | `courseID`, `teacherName` |
| 包名 | 全小写，简洁 | `entity`, `repository` |

#### 1.2 文件组织

- 每个实体或值对象一个文件
- 一个包内相关功能放在同一目录
- 包结构与领域模型一致
- 测试文件与源文件放在同一目录，以 `_test.go` 结尾

#### 1.3 注释规范

- 所有导出的类型、函数、方法必须有注释
- 注释应说明"做什么"而非"怎么做"
- 领域模型的业务规则必须在注释中说明
- 使用标准格式的注释，便于godoc生成文档

```go
// CourseService 处理课程相关的核心业务逻辑
// 包括课程创建、发布、归档等操作
type CourseService struct {
    // ...
}

// CreateCourse 创建新课程并进行业务规则验证
// 返回创建的课程ID和可能的错误
func (s *CourseService) CreateCourse(ctx context.Context, course *entity.Course) (string, error) {
    // ...
}
```

### 2. DDD实践原则

#### 2.1 领域模型设计原则

1. **实体设计**:
   - 每个实体必须有唯一标识符
   - 实体应包含行为和状态
   - 避免贫血模型，确保实体包含业务逻辑

2. **值对象设计**:
   - 值对象应是不可变的
   - 没有唯一标识的概念
   - 可以被替换而不影响业务逻辑

3. **聚合设计**:
   - 明确定义聚合边界
   - 每个聚合需有一个聚合根
   - 跨聚合引用使用ID，不直接引用对象

4. **仓储设计**:
   - 每个聚合根一个仓储接口
   - 仓储应负责持久化和检索聚合
   - 查询方法应基于业务需求设计

#### 2.2 层间交互原则

1. **依赖方向**:
   - 总是从外部层向内部层依赖
   - 内部层不应知道外部层的存在
   - 使用依赖倒置原则处理依赖

2. **数据传输**:
   - 不同层之间通过DTO传递数据
   - 避免直接传递实体对象到应用层之外
   - 应用层负责实体与DTO的转换

3. **关注点分离**:
   - 领域层关注业务规则
   - 应用层关注用例和流程
   - 基础设施层关注技术实现
   - 交付层关注请求处理和响应

### 3. 常见问题与解决方案

#### 3.1 复杂业务逻辑处理

**问题**: 业务逻辑涉及多个实体和聚合。

**解决方案**: 
- 使用领域服务处理跨实体的业务逻辑
- 应用服务协调多个领域服务的操作
- 使用事件驱动方式处理复杂业务流程

```go
// 在领域服务中处理跨实体业务逻辑
func (s *EnrollmentService) EnrollCourse(ctx context.Context, userID, courseID string) (*entity.Enrollment, error) {
    // 检查课程是否可报名
    course, err := s.courseRepo.FindByID(ctx, courseID)
    if err != nil {
        return nil, err
    }
    
    if !course.IsEnrollable() {
        return nil, errors.New("课程当前不可报名")
    }
    
    // 检查用户是否已报名
    exists, err := s.enrollmentRepo.ExistsByUserAndCourse(ctx, userID, courseID)
    if err != nil {
        return nil, err
    }
    
    if exists {
        return nil, errors.New("用户已报名该课程")
    }
    
    // 创建报名记录
    enrollment := entity.NewEnrollment(userID, courseID)
    
    // 保存报名记录
    if err := s.enrollmentRepo.Save(ctx, enrollment); err != nil {
        return nil, err
    }
    
    return enrollment, nil
}
```

#### 3.2 事务管理

**问题**: 操作需要跨多个仓储的事务一致性。

**解决方案**:
- 在应用层管理事务
- 使用工作单元模式（Unit of Work）
- 领域事件处理最终一致性

```go
// 在应用服务中管理事务
func (s *CourseAppService) CreateCourseWithChapters(ctx context.Context, req *dto.CreateCourseWithChaptersRequest) (*dto.CourseDTO, error) {
    // 开启事务
    tx, err := s.transactionManager.Begin(ctx)
    if err != nil {
        return nil, err
    }
    
    // 确保事务结束时提交或回滚
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r) // 重新panic
        }
    }()
    
    // 创建课程
    course, err := s.courseService.CreateCourse(tx.Context(), /* 参数 */)
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // 创建章节
    for _, chapterReq := range req.Chapters {
        _, err := s.chapterLessonService.CreateChapter(tx.Context(), course.ID, /* 参数 */)
        if err != nil {
            tx.Rollback()
            return nil, err
        }
    }
    
    // 提交事务
    if err := tx.Commit(); err != nil {
        return nil, err
    }
    
    // 转换为DTO返回
    return s.toCourseDTO(course), nil
}
```

#### 3.3 性能优化

**问题**: 领域模型操作可能导致性能问题。

**解决方案**:
- 使用CQRS分离读写操作
- 针对复杂查询使用专门的查询服务
- 利用缓存提高读取性能

```go
// 针对复杂查询使用专门的查询服务
type CourseQueryService struct {
    db *sql.DB // 直接使用数据库连接进行查询优化
}

func (s *CourseQueryService) GetCourseWithDetails(ctx context.Context, courseID string) (*dto.CourseDetailDTO, error) {
    // 使用优化的SQL查询直接获取所需数据
    query := `
        SELECT c.id, c.title, c.description, c.status,
               t.id as teacher_id, t.name as teacher_name,
               cat.id as category_id, cat.name as category_name
        FROM courses c
        JOIN teachers t ON c.teacher_id = t.id
        JOIN categories cat ON c.category_id = cat.id
        WHERE c.id = ?
    `
    
    // 直接映射到DTO而非实体
    var dto dto.CourseDetailDTO
    err := s.db.QueryRowContext(ctx, query, courseID).Scan(
        &dto.ID, &dto.Title, &dto.Description, &dto.Status,
        &dto.Teacher.ID, &dto.Teacher.Name,
        &dto.Category.ID, &dto.Category.Name,
    )
    
    // 处理结果...
    
    return &dto, nil
}
```

## 五、开发工作流程

### 1. 新微服务开发流程

按照以下步骤开发新的微服务：

1. **领域分析与建模**:
   - 识别业务概念和实体
   - 确定聚合边界
   - 设计领域模型

2. **领域层实现**:
   - 实现实体与值对象
   - 定义仓储接口
   - 开发领域服务

3. **应用层实现**:
   - 定义应用服务接口
   - 实现用例编排
   - 设计DTO对象

4. **基础设施层实现**:
   - 实现仓储接口
   - 集成外部系统
   - 配置数据库连接

5. **交付层实现**:
   - 开发HTTP处理器
   - 配置路由
   - 实现中间件

6. **入口程序实现**:
   - 配置依赖注入
   - 设置服务器参数
   - 实现优雅启动和关闭

7. **测试与部署**:
   - 编写单元测试
   - 进行集成测试
   - 配置部署脚本

### 2. 微服务集成规范

在与其他微服务集成时，遵循以下原则：

1. **客户端接口**:
   - 使用客户端接口隔离微服务调用
   - 实现错误处理和重试策略
   - 支持服务发现和负载均衡

2. **数据一致性**:
   - 采用最终一致性模型
   - 使用领域事件进行状态同步
   - 实现补偿事务处理失败场景

3. **服务边界**:
   - 明确定义服务API契约
   - 版本化API接口
   - 文档化所有公开接口

## 六、总结

本指南旨在规范化万智平台微服务的开发流程，确保所有微服务遵循相同的架构模式和代码规范。通过领域驱动设计，我们将复杂业务逻辑封装在领域模型中，保持代码与业务的一致性，提高系统的可维护性和可扩展性。

开发团队在实施过程中应密切关注本指南，同时结合实际项目需求进行灵活调整，不断优化架构和开发流程。随着项目的进展，本指南也将持续更新，以反映最新的最佳实践和经验教训。
