package service

import (
	"context"
	"errors"

	"wz-backend-go/internal/domain/learn/entity"
	"wz-backend-go/internal/domain/learn/repository"
)

var (
	ErrTeacherExists     = errors.New("讲师已存在")
	ErrUserNotFound      = errors.New("用户不存在")
	ErrInvalidTeacherStatus = errors.New("无效的讲师状态")
)

// TeacherService 讲师领域服务，处理讲师相关的业务逻辑
type TeacherService struct {
	teacherRepo   repository.TeacherRepository
	courseRepo    repository.CourseRepository
	// userRepo可能需要从user微服务获取信息，这里是一个可能的接口
	// userRepo      repository.UserRepository
}

// NewTeacherService 创建讲师服务
func NewTeacherService(
	teacherRepo repository.TeacherRepository,
	courseRepo repository.CourseRepository,
) *TeacherService {
	return &TeacherService{
		teacherRepo: teacherRepo,
		courseRepo:  courseRepo,
	}
}

// CreateTeacher 创建新讲师
func (s *TeacherService) CreateTeacher(ctx context.Context, userID, name, title, introduction, avatar string) (*entity.Teacher, error) {
	// 检查讲师是否已存在
	existingTeacher, err := s.teacherRepo.GetByUserID(ctx, userID)
	if err == nil && existingTeacher != nil {
		return nil, ErrTeacherExists
	}

	// 创建讲师
	teacher := entity.NewTeacher(userID, name)
	teacher.Title = title
	teacher.Introduction = introduction
	teacher.Avatar = avatar

	// 保存讲师
	if err := s.teacherRepo.Create(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

// UpdateTeacherProfile 更新讲师基本资料
func (s *TeacherService) UpdateTeacherProfile(ctx context.Context, id, name, avatar, title, introduction string) (*entity.Teacher, error) {
	// 获取讲师
	teacher, err := s.teacherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTeacherNotFound
	}

	// 更新讲师信息
	teacher.Update(name, avatar, title, introduction)

	// 保存讲师
	if err := s.teacherRepo.Update(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

// UpdateTeacherContact 更新讲师联系信息
func (s *TeacherService) UpdateTeacherContact(ctx context.Context, id, email, phone string) (*entity.Teacher, error) {
	// 获取讲师
	teacher, err := s.teacherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTeacherNotFound
	}

	// 更新联系信息
	teacher.UpdateContact(email, phone)

	// 保存讲师
	if err := s.teacherRepo.Update(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

// SetTeacherSpecialties 设置讲师专长领域
func (s *TeacherService) SetTeacherSpecialties(ctx context.Context, id string, specialties []string) (*entity.Teacher, error) {
	// 获取讲师
	teacher, err := s.teacherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTeacherNotFound
	}

	// 更新讲师专长
	teacher.SetSpecialties(specialties)

	// 保存讲师
	if err := s.teacherRepo.Update(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

// SetTeacherSocialProfiles 设置讲师社交档案
func (s *TeacherService) SetTeacherSocialProfiles(ctx context.Context, id string, profiles []string) (*entity.Teacher, error) {
	// 获取讲师
	teacher, err := s.teacherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTeacherNotFound
	}

	// 更新社交档案
	teacher.SetSocialProfiles(profiles)

	// 保存讲师
	if err := s.teacherRepo.Update(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

// ActivateTeacher 激活讲师
func (s *TeacherService) ActivateTeacher(ctx context.Context, id string) (*entity.Teacher, error) {
	// 获取讲师
	teacher, err := s.teacherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTeacherNotFound
	}

	// 激活讲师
	teacher.Activate()

	// 保存讲师
	if err := s.teacherRepo.Update(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

// DeactivateTeacher 停用讲师
func (s *TeacherService) DeactivateTeacher(ctx context.Context, id string) (*entity.Teacher, error) {
	// 获取讲师
	teacher, err := s.teacherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTeacherNotFound
	}

	// 停用讲师
	teacher.Deactivate()

	// 保存讲师
	if err := s.teacherRepo.Update(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

// SuspendTeacher 暂停讲师
func (s *TeacherService) SuspendTeacher(ctx context.Context, id string) (*entity.Teacher, error) {
	// 获取讲师
	teacher, err := s.teacherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTeacherNotFound
	}

	// 暂停讲师
	teacher.Suspend()

	// 保存讲师
	if err := s.teacherRepo.Update(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

// DeleteTeacher 删除讲师
func (s *TeacherService) DeleteTeacher(ctx context.Context, id string) error {
	// 获取讲师
	teacher, err := s.teacherRepo.GetByID(ctx, id)
	if err != nil {
		return ErrTeacherNotFound
	}

	// 检查讲师是否有课程
	count, err := s.courseRepo.CountByTeacherID(ctx, id)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("讲师有关联课程，无法删除")
	}

	// 删除讲师
	if err := s.teacherRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetTeacherStats 获取讲师统计数据
func (s *TeacherService) GetTeacherStats(ctx context.Context) (totalCount, activeCount, inactiveCount int64, err error) {
	totalCount, err = s.teacherRepo.CountAll(ctx)
	if err != nil {
		return 0, 0, 0, err
	}

	activeCount, err = s.teacherRepo.CountByStatus(ctx, entity.TeacherStatusActive)
	if err != nil {
		return 0, 0, 0, err
	}

	inactiveCount, err = s.teacherRepo.CountByStatus(ctx, entity.TeacherStatusInactive)
	if err != nil {
		return 0, 0, 0, err
	}

	return totalCount, activeCount, inactiveCount, nil
}
