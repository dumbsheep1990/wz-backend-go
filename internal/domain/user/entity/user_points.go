package entity

import (
	"errors"
	"time"

	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/domain/user/valueobject"
)

// ErrInsufficientPoints 表示用户积分不足的错误
var ErrInsufficientPoints = errors.New("用户积分不足")

// UserPoints 表示用户积分领域实体
type UserPoints struct {
	id          valueobject.ID
	userID      valueobject.UserID
	points      valueobject.Points
	totalPoints valueobject.Points
	pointsType  valueobject.PointsType
	source      valueobject.Source
	description valueobject.Description
	relatedID   int64
	relatedType valueobject.RelatedType
	tenantID    valueobject.TenantID
	operatorID  valueobject.UserID
	isRevoked   bool
	createdAt   time.Time
	updatedAt   time.Time

	domainEvents []event.DomainEvent
}

// NewUserPoints 创建一个新的用户积分实体
func NewUserPoints(
	userID valueobject.UserID,
	points valueobject.Points,
	pointsType valueobject.PointsType,
	source valueobject.Source,
	description valueobject.Description,
	relatedID int64,
	relatedType valueobject.RelatedType,
	operatorID valueobject.UserID,
	tenantID valueobject.TenantID,
) (*UserPoints, error) {
	// 验证必要参数
	if userID.IsEmpty() {
		return nil, errors.New("用户ID不能为空")
	}

	if points.Value() <= 0 {
		return nil, errors.New("积分值必须大于0")
	}

	now := time.Now()

	userPoints := &UserPoints{
		id:          valueobject.NewID(""), // 由仓储层生成
		userID:      userID,
		points:      points,
		totalPoints: valueobject.Points(0), // 初始化为0，后续由仓储层更新
		pointsType:  pointsType,
		source:      source,
		description: description,
		relatedID:   relatedID,
		relatedType: relatedType,
		tenantID:    tenantID,
		operatorID:  operatorID,
		isRevoked:   false,
		createdAt:   now,
		updatedAt:   now,
	}

	// 添加领域事件
	event := NewUserPointsCreatedEvent(userPoints)
	userPoints.addDomainEvent(event)

	return userPoints, nil
}

// ID 获取积分记录ID
func (up *UserPoints) ID() valueobject.ID {
	return up.id
}

// UserID 获取用户ID
func (up *UserPoints) UserID() valueobject.UserID {
	return up.userID
}

// Points 获取积分值
func (up *UserPoints) Points() valueobject.Points {
	return up.points
}

// TotalPoints 获取积分总值
func (up *UserPoints) TotalPoints() valueobject.Points {
	return up.totalPoints
}

// SetTotalPoints 设置积分总值（通常由仓储层调用）
func (up *UserPoints) SetTotalPoints(total valueobject.Points) {
	up.totalPoints = total
}

// PointsType 获取积分类型
func (up *UserPoints) PointsType() valueobject.PointsType {
	return up.pointsType
}

// Source 获取积分来源
func (up *UserPoints) Source() valueobject.Source {
	return up.source
}

// Description 获取描述
func (up *UserPoints) Description() valueobject.Description {
	return up.description
}

// RelatedID 获取关联ID
func (up *UserPoints) RelatedID() int64 {
	return up.relatedID
}

// RelatedType 获取关联类型
func (up *UserPoints) RelatedType() valueobject.RelatedType {
	return up.relatedType
}

// TenantID 获取租户ID
func (up *UserPoints) TenantID() valueobject.TenantID {
	return up.tenantID
}

// OperatorID 获取操作者ID
func (up *UserPoints) OperatorID() valueobject.UserID {
	return up.operatorID
}

// IsRevoked 检查积分记录是否已撤销
func (up *UserPoints) IsRevoked() bool {
	return up.isRevoked
}

// CreatedAt 获取创建时间
func (up *UserPoints) CreatedAt() time.Time {
	return up.createdAt
}

// UpdatedAt 获取更新时间
func (up *UserPoints) UpdatedAt() time.Time {
	return up.updatedAt
}

// Revoke 撤销积分记录
func (up *UserPoints) Revoke(operatorID valueobject.UserID) error {
	if up.isRevoked {
		return errors.New("积分记录已经被撤销")
	}

	up.isRevoked = true
	up.updatedAt = time.Now()

	// 添加领域事件
	event := NewUserPointsRevokedEvent(up, operatorID)
	up.addDomainEvent(event)

	return nil
}

// SetID 设置ID（通常由仓储层调用）
func (up *UserPoints) SetID(id valueobject.ID) {
	up.id = id
}

// addDomainEvent 添加领域事件
func (up *UserPoints) addDomainEvent(event event.DomainEvent) {
	up.domainEvents = append(up.domainEvents, event)
}

// GetDomainEvents 获取所有领域事件
func (up *UserPoints) GetDomainEvents() []event.DomainEvent {
	return up.domainEvents
}

// ClearDomainEvents 清除所有领域事件
func (up *UserPoints) ClearDomainEvents() {
	up.domainEvents = []event.DomainEvent{}
}

// ReconstructUserPoints 从存储中重建用户积分实体（仅供仓储层使用）
func ReconstructUserPoints(
	id valueobject.ID,
	userID valueobject.UserID,
	points valueobject.Points,
	totalPoints valueobject.Points,
	pointsType valueobject.PointsType,
	source valueobject.Source,
	description valueobject.Description,
	relatedID int64,
	relatedType valueobject.RelatedType,
	tenantID valueobject.TenantID,
	operatorID valueobject.UserID,
	isRevoked bool,
	createdAt time.Time,
	updatedAt time.Time,
) (*UserPoints, error) {
	userPoints := &UserPoints{
		id:          id,
		userID:      userID,
		points:      points,
		totalPoints: totalPoints,
		pointsType:  pointsType,
		source:      source,
		description: description,
		relatedID:   relatedID,
		relatedType: relatedType,
		tenantID:    tenantID,
		operatorID:  operatorID,
		isRevoked:   isRevoked,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}

	return userPoints, nil
}
