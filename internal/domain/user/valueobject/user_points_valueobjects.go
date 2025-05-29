package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

// ID 表示唯一标识符值对象
type ID string

// NewID 创建一个新的ID
func NewID(id string) ID {
	return ID(id)
}

// String 返回ID的字符串表示
func (id ID) String() string {
	return string(id)
}

// IsEmpty 检查ID是否为空
func (id ID) IsEmpty() bool {
	return id == ""
}

// UserID 表示用户标识符值对象
type UserID int64

// NewUserID 创建一个新的UserID
func NewUserID(id int64) UserID {
	return UserID(id)
}

// Value 返回UserID的整数值
func (id UserID) Value() int64 {
	return int64(id)
}

// IsEmpty 检查UserID是否为空或无效
func (id UserID) IsEmpty() bool {
	return id <= 0
}

// Points 表示积分值对象
type Points int

// NewPoints 创建一个新的Points值对象
func NewPoints(value int) (Points, error) {
	if value < 0 {
		return 0, errors.New("积分值不能为负数")
	}
	return Points(value), nil
}

// Value 返回Points的整数值
func (p Points) Value() int {
	return int(p)
}

// Add 将积分与另一个积分相加并返回新的积分值对象
func (p Points) Add(other Points) Points {
	return Points(p.Value() + other.Value())
}

// Subtract 将积分减去另一个积分并返回新的积分值对象
func (p Points) Subtract(other Points) (Points, error) {
	result := p.Value() - other.Value()
	if result < 0 {
		return 0, errors.New("积分结果不能为负数")
	}
	return Points(result), nil
}

// PointsType 表示积分类型值对象
type PointsType int

const (
	PointsTypeIncrease PointsType = 1 // 增加积分
	PointsTypeDecrease PointsType = 2 // 减少积分
)

// NewPointsType 创建一个新的PointsType值对象
func NewPointsType(value int) (PointsType, error) {
	if value != int(PointsTypeIncrease) && value != int(PointsTypeDecrease) {
		return 0, errors.New("无效的积分类型，必须是1(增加)或2(减少)")
	}
	return PointsType(value), nil
}

// Value 返回PointsType的整数值
func (pt PointsType) Value() int {
	return int(pt)
}

// String 返回PointsType的字符串表示
func (pt PointsType) String() string {
	switch pt {
	case PointsTypeIncrease:
		return "增加"
	case PointsTypeDecrease:
		return "减少"
	default:
		return "未知"
	}
}

// Source 表示积分来源值对象
type Source string

const (
	SourceSignIn   Source = "sign"     // 签到
	SourcePurchase Source = "purchase" // 购买
	SourceActivity Source = "activity" // 活动
	SourceComment  Source = "comment"  // 评论
	SourceInvite   Source = "invite"   // 邀请
	SourceArticle  Source = "article"  // 发布文章
	SourceShare    Source = "share"    // 分享
	SourceAdmin    Source = "admin"    // 管理员调整
	SourceSystem   Source = "system"   // 系统
	SourceExchange Source = "exchange" // 积分兑换
)

// NewSource 创建一个新的Source值对象
func NewSource(source string) (Source, error) {
	source = strings.TrimSpace(source)
	if source == "" {
		return "", errors.New("积分来源不能为空")
	}

	// 验证来源是否为预定义的来源之一
	validSources := []Source{
		SourceSignIn, SourcePurchase, SourceActivity, SourceComment,
		SourceInvite, SourceArticle, SourceShare, SourceAdmin,
		SourceSystem, SourceExchange,
	}

	for _, validSource := range validSources {
		if Source(source) == validSource {
			return Source(source), nil
		}
	}

	// 也可以接受自定义来源，但需要符合一定格式
	if matched, _ := regexp.MatchString(`^[a-z0-9_-]{1,20}$`, source); !matched {
		return "", errors.New("无效的积分来源格式")
	}

	return Source(source), nil
}

// String 返回Source的字符串表示
func (s Source) String() string {
	return string(s)
}

// DisplayName 返回Source的显示名称
func (s Source) DisplayName() string {
	sourceMap := map[Source]string{
		SourceSignIn:   "签到",
		SourcePurchase: "购买",
		SourceActivity: "活动",
		SourceComment:  "评论",
		SourceInvite:   "邀请",
		SourceArticle:  "发布文章",
		SourceShare:    "分享",
		SourceAdmin:    "管理员调整",
		SourceSystem:   "系统",
		SourceExchange: "积分兑换",
	}

	if name, ok := sourceMap[s]; ok {
		return name
	}

	return string(s)
}

// Description 表示描述值对象
type Description string

// NewDescription 创建一个新的Description
func NewDescription(desc string) (Description, error) {
	desc = strings.TrimSpace(desc)
	if desc == "" {
		return "", errors.New("描述不能为空")
	}

	if len(desc) > 200 {
		return "", errors.New("描述过长（最多200个字符）")
	}

	return Description(desc), nil
}

// String 返回Description的字符串表示
func (d Description) String() string {
	return string(d)
}

// IsEmpty 检查Description是否为空
func (d Description) IsEmpty() bool {
	return d == ""
}

// RelatedType 表示关联类型值对象
type RelatedType string

// NewRelatedType 创建一个新的RelatedType
func NewRelatedType(relatedType string) RelatedType {
	return RelatedType(relatedType)
}

// String 返回RelatedType的字符串表示
func (rt RelatedType) String() string {
	return string(rt)
}

// TenantID 表示租户ID值对象
type TenantID int64

// NewTenantID 创建一个新的TenantID
func NewTenantID(id int64) TenantID {
	return TenantID(id)
}

// Value 返回TenantID的整数值
func (id TenantID) Value() int64 {
	return int64(id)
}

// IsEmpty 检查TenantID是否为空
func (id TenantID) IsEmpty() bool {
	return id <= 0
}
