package entity

import (
	"time"

	"wz-backend-go/internal/domain/user/valueobject"
)

// UserFavorite 用户收藏实体
type UserFavorite struct {
	id        valueobject.FavoriteID
	userID    valueobject.FavoriteUserID
	itemID    valueobject.FavoriteItemID
	itemType  valueobject.FavoriteItemType
	title     valueobject.FavoriteTitle
	cover     string
	summary   string
	url       valueobject.FavoriteURL
	remark    string
	tenantID  valueobject.FavoriteTenantID
	username  string
	createdAt time.Time
	updatedAt time.Time
}

// NewUserFavorite 创建新的用户收藏实体
func NewUserFavorite(
	id valueobject.FavoriteID,
	userID valueobject.FavoriteUserID,
	itemID valueobject.FavoriteItemID,
	itemType valueobject.FavoriteItemType,
	title valueobject.FavoriteTitle,
	cover string,
	summary string,
	url valueobject.FavoriteURL,
	remark string,
	tenantID valueobject.FavoriteTenantID,
) (*UserFavorite, error) {
	now := time.Now()

	favorite := &UserFavorite{
		id:        id,
		userID:    userID,
		itemID:    itemID,
		itemType:  itemType,
		title:     title,
		cover:     cover,
		summary:   summary,
		url:       url,
		remark:    remark,
		tenantID:  tenantID,
		createdAt: now,
		updatedAt: now,
	}

	return favorite, nil
}

// ID 获取收藏ID
func (f *UserFavorite) ID() valueobject.FavoriteID {
	return f.id
}

// UserID 获取用户ID
func (f *UserFavorite) UserID() valueobject.FavoriteUserID {
	return f.userID
}

// ItemID 获取内容项ID
func (f *UserFavorite) ItemID() valueobject.FavoriteItemID {
	return f.itemID
}

// ItemType 获取内容类型
func (f *UserFavorite) ItemType() valueobject.FavoriteItemType {
	return f.itemType
}

// Title 获取标题
func (f *UserFavorite) Title() valueobject.FavoriteTitle {
	return f.title
}

// Cover 获取封面
func (f *UserFavorite) Cover() string {
	return f.cover
}

// SetCover 设置封面
func (f *UserFavorite) SetCover(cover string) {
	f.cover = cover
	f.updatedAt = time.Now()
}

// Summary 获取摘要
func (f *UserFavorite) Summary() string {
	return f.summary
}

// SetSummary 设置摘要
func (f *UserFavorite) SetSummary(summary string) {
	f.summary = summary
	f.updatedAt = time.Now()
}

// URL 获取URL
func (f *UserFavorite) URL() valueobject.FavoriteURL {
	return f.url
}

// SetURL 设置URL
func (f *UserFavorite) SetURL(url valueobject.FavoriteURL) {
	f.url = url
	f.updatedAt = time.Now()
}

// Remark 获取备注
func (f *UserFavorite) Remark() string {
	return f.remark
}

// SetRemark 设置备注
func (f *UserFavorite) SetRemark(remark string) {
	f.remark = remark
	f.updatedAt = time.Now()
}

// TenantID 获取租户ID
func (f *UserFavorite) TenantID() valueobject.FavoriteTenantID {
	return f.tenantID
}

// Username 获取用户名
func (f *UserFavorite) Username() string {
	return f.username
}

// SetUsername 设置用户名
func (f *UserFavorite) SetUsername(username string) {
	f.username = username
}

// CreatedAt 获取创建时间
func (f *UserFavorite) CreatedAt() time.Time {
	return f.createdAt
}

// UpdatedAt 获取更新时间
func (f *UserFavorite) UpdatedAt() time.Time {
	return f.updatedAt
}

// IsBelongsToUser 检查收藏是否属于指定用户
func (f *UserFavorite) IsBelongsToUser(userID valueobject.FavoriteUserID) bool {
	return f.userID == userID
}

// Update 更新收藏信息（仅允许更新部分字段）
func (f *UserFavorite) Update(cover string, summary string, url valueobject.FavoriteURL, remark string) error {
	f.cover = cover
	f.summary = summary
	f.url = url
	f.remark = remark
	f.updatedAt = time.Now()
	return nil
}

// ToDTO 将实体转换为数据传输对象
func (f *UserFavorite) ToDTO() map[string]interface{} {
	return map[string]interface{}{
		"id":         f.id.Value(),
		"user_id":    f.userID.Value(),
		"username":   f.username,
		"item_id":    f.itemID.Value(),
		"item_type":  f.itemType.Value(),
		"title":      f.title.Value(),
		"cover":      f.cover,
		"summary":    f.summary,
		"url":        f.url.Value(),
		"remark":     f.remark,
		"tenant_id":  f.tenantID.Value(),
		"created_at": f.createdAt,
		"updated_at": f.updatedAt,
	}
}
