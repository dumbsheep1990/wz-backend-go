package entity

import (
	"errors"
	"time"

	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/domain/user/valueobject"
)

// PointsRules 表示积分规则领域实体
type PointsRules struct {
	id                valueobject.ID
	signInPoints      int
	commentPoints     int
	sharePoints       int
	articlePoints     int
	invitePoints      int
	purchaseRate      int
	maxDailyPoints    int
	enableExchange    bool
	exchangeRate      int
	minExchangePoints int
	tenantID          valueobject.TenantID
	updatedAt         time.Time

	domainEvents []event.DomainEvent
}

// NewPointsRules 创建一个新的积分规则实体
func NewPointsRules(
	signInPoints int,
	commentPoints int,
	sharePoints int,
	articlePoints int,
	invitePoints int,
	purchaseRate int,
	maxDailyPoints int,
	enableExchange bool,
	exchangeRate int,
	minExchangePoints int,
	tenantID valueobject.TenantID,
) (*PointsRules, error) {
	// 验证积分规则参数的合理性
	if signInPoints < 0 || commentPoints < 0 || sharePoints < 0 ||
		articlePoints < 0 || invitePoints < 0 || purchaseRate < 0 ||
		maxDailyPoints < 0 || exchangeRate < 0 || minExchangePoints < 0 {
		return nil, errors.New("积分规则参数不能为负数")
	}

	rules := &PointsRules{
		id:                valueobject.NewID(""), // 由仓储层生成
		signInPoints:      signInPoints,
		commentPoints:     commentPoints,
		sharePoints:       sharePoints,
		articlePoints:     articlePoints,
		invitePoints:      invitePoints,
		purchaseRate:      purchaseRate,
		maxDailyPoints:    maxDailyPoints,
		enableExchange:    enableExchange,
		exchangeRate:      exchangeRate,
		minExchangePoints: minExchangePoints,
		tenantID:          tenantID,
		updatedAt:         time.Now(),
	}

	// 添加领域事件
	event := NewPointsRulesCreatedEvent(rules)
	rules.addDomainEvent(event)

	return rules, nil
}

// ID 获取积分规则ID
func (pr *PointsRules) ID() valueobject.ID {
	return pr.id
}

// SignInPoints 获取签到积分
func (pr *PointsRules) SignInPoints() int {
	return pr.signInPoints
}

// CommentPoints 获取评论积分
func (pr *PointsRules) CommentPoints() int {
	return pr.commentPoints
}

// SharePoints 获取分享积分
func (pr *PointsRules) SharePoints() int {
	return pr.sharePoints
}

// ArticlePoints 获取发布文章积分
func (pr *PointsRules) ArticlePoints() int {
	return pr.articlePoints
}

// InvitePoints 获取邀请积分
func (pr *PointsRules) InvitePoints() int {
	return pr.invitePoints
}

// PurchaseRate 获取购买积分比例
func (pr *PointsRules) PurchaseRate() int {
	return pr.purchaseRate
}

// MaxDailyPoints 获取每日最大获取积分
func (pr *PointsRules) MaxDailyPoints() int {
	return pr.maxDailyPoints
}

// EnableExchange 获取是否可兑换商品
func (pr *PointsRules) EnableExchange() bool {
	return pr.enableExchange
}

// ExchangeRate 获取兑换比例
func (pr *PointsRules) ExchangeRate() int {
	return pr.exchangeRate
}

// MinExchangePoints 获取最小兑换积分
func (pr *PointsRules) MinExchangePoints() int {
	return pr.minExchangePoints
}

// TenantID 获取租户ID
func (pr *PointsRules) TenantID() valueobject.TenantID {
	return pr.tenantID
}

// UpdatedAt 获取更新时间
func (pr *PointsRules) UpdatedAt() time.Time {
	return pr.updatedAt
}

// Update 更新积分规则
func (pr *PointsRules) Update(
	signInPoints int,
	commentPoints int,
	sharePoints int,
	articlePoints int,
	invitePoints int,
	purchaseRate int,
	maxDailyPoints int,
	enableExchange bool,
	exchangeRate int,
	minExchangePoints int,
) error {
	// 验证积分规则参数的合理性
	if signInPoints < 0 || commentPoints < 0 || sharePoints < 0 ||
		articlePoints < 0 || invitePoints < 0 || purchaseRate < 0 ||
		maxDailyPoints < 0 || exchangeRate < 0 || minExchangePoints < 0 {
		return errors.New("积分规则参数不能为负数")
	}

	// 检查是否有变化
	if pr.signInPoints == signInPoints &&
		pr.commentPoints == commentPoints &&
		pr.sharePoints == sharePoints &&
		pr.articlePoints == articlePoints &&
		pr.invitePoints == invitePoints &&
		pr.purchaseRate == purchaseRate &&
		pr.maxDailyPoints == maxDailyPoints &&
		pr.enableExchange == enableExchange &&
		pr.exchangeRate == exchangeRate &&
		pr.minExchangePoints == minExchangePoints {
		return nil
	}

	// 更新规则
	pr.signInPoints = signInPoints
	pr.commentPoints = commentPoints
	pr.sharePoints = sharePoints
	pr.articlePoints = articlePoints
	pr.invitePoints = invitePoints
	pr.purchaseRate = purchaseRate
	pr.maxDailyPoints = maxDailyPoints
	pr.enableExchange = enableExchange
	pr.exchangeRate = exchangeRate
	pr.minExchangePoints = minExchangePoints
	pr.updatedAt = time.Now()

	// 添加领域事件
	event := NewPointsRulesUpdatedEvent(pr)
	pr.addDomainEvent(event)

	return nil
}

// SetID 设置ID（通常由仓储层调用）
func (pr *PointsRules) SetID(id valueobject.ID) {
	pr.id = id
}

// addDomainEvent 添加领域事件
func (pr *PointsRules) addDomainEvent(event event.DomainEvent) {
	pr.domainEvents = append(pr.domainEvents, event)
}

// GetDomainEvents 获取所有领域事件
func (pr *PointsRules) GetDomainEvents() []event.DomainEvent {
	return pr.domainEvents
}

// ClearDomainEvents 清除所有领域事件
func (pr *PointsRules) ClearDomainEvents() {
	pr.domainEvents = []event.DomainEvent{}
}

// ReconstructPointsRules 从存储中重建积分规则实体（仅供仓储层使用）
func ReconstructPointsRules(
	id valueobject.ID,
	signInPoints int,
	commentPoints int,
	sharePoints int,
	articlePoints int,
	invitePoints int,
	purchaseRate int,
	maxDailyPoints int,
	enableExchange bool,
	exchangeRate int,
	minExchangePoints int,
	tenantID valueobject.TenantID,
	updatedAt time.Time,
) (*PointsRules, error) {
	return &PointsRules{
		id:                id,
		signInPoints:      signInPoints,
		commentPoints:     commentPoints,
		sharePoints:       sharePoints,
		articlePoints:     articlePoints,
		invitePoints:      invitePoints,
		purchaseRate:      purchaseRate,
		maxDailyPoints:    maxDailyPoints,
		enableExchange:    enableExchange,
		exchangeRate:      exchangeRate,
		minExchangePoints: minExchangePoints,
		tenantID:          tenantID,
		updatedAt:         updatedAt,
	}, nil
}
