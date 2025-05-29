package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// FavoriteID 收藏ID值对象
type FavoriteID int64

// NewFavoriteID 创建收藏ID值对象
func NewFavoriteID(id int64) (FavoriteID, error) {
	if id <= 0 {
		return 0, errors.New("收藏ID必须大于0")
	}
	return FavoriteID(id), nil
}

// Value 获取ID值
func (id FavoriteID) Value() int64 {
	return int64(id)
}

// String 字符串表示
func (id FavoriteID) String() string {
	return fmt.Sprintf("%d", id)
}

// FavoriteUserID 用户ID值对象
type FavoriteUserID int64

// NewFavoriteUserID 创建用户ID值对象
func NewFavoriteUserID(id int64) (FavoriteUserID, error) {
	if id <= 0 {
		return 0, errors.New("用户ID必须大于0")
	}
	return FavoriteUserID(id), nil
}

// Value 获取ID值
func (id FavoriteUserID) Value() int64 {
	return int64(id)
}

// String 字符串表示
func (id FavoriteUserID) String() string {
	return fmt.Sprintf("%d", id)
}

// FavoriteItemID 内容项ID值对象
type FavoriteItemID int64

// NewFavoriteItemID 创建内容项ID值对象
func NewFavoriteItemID(id int64) (FavoriteItemID, error) {
	if id <= 0 {
		return 0, errors.New("内容项ID必须大于0")
	}
	return FavoriteItemID(id), nil
}

// Value 获取ID值
func (id FavoriteItemID) Value() int64 {
	return int64(id)
}

// String 字符串表示
func (id FavoriteItemID) String() string {
	return fmt.Sprintf("%d", id)
}

// FavoriteItemType 内容类型值对象
type FavoriteItemType string

// 预定义的内容类型常量
const (
	ItemTypeArticle FavoriteItemType = "article" // 文章
	ItemTypeProduct FavoriteItemType = "product" // 商品
	ItemTypeVideo   FavoriteItemType = "video"   // 视频
	ItemTypeAudio   FavoriteItemType = "audio"   // 音频
	ItemTypeImage   FavoriteItemType = "image"   // 图片
	ItemTypeOther   FavoriteItemType = "other"   // 其他
)

// NewFavoriteItemType 创建内容类型值对象
func NewFavoriteItemType(itemType string) (FavoriteItemType, error) {
	if itemType == "" {
		return "", errors.New("内容类型不能为空")
	}

	itemType = strings.ToLower(itemType)
	switch FavoriteItemType(itemType) {
	case ItemTypeArticle, ItemTypeProduct, ItemTypeVideo, ItemTypeAudio, ItemTypeImage, ItemTypeOther:
		return FavoriteItemType(itemType), nil
	default:
		return "", fmt.Errorf("无效的内容类型: %s", itemType)
	}
}

// Value 获取类型值
func (t FavoriteItemType) Value() string {
	return string(t)
}

// String 字符串表示
func (t FavoriteItemType) String() string {
	return string(t)
}

// FavoriteTitle 标题值对象
type FavoriteTitle string

// NewFavoriteTitle 创建标题值对象
func NewFavoriteTitle(title string) (FavoriteTitle, error) {
	if title == "" {
		return "", errors.New("标题不能为空")
	}
	if len(title) > 200 {
		return "", errors.New("标题长度不能超过200个字符")
	}
	return FavoriteTitle(title), nil
}

// Value 获取标题值
func (t FavoriteTitle) Value() string {
	return string(t)
}

// String 字符串表示
func (t FavoriteTitle) String() string {
	return string(t)
}

// FavoriteURL URL值对象
type FavoriteURL string

// urlRegex URL验证正则表达式
var urlRegex = regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:/[a-zA-Z0-9\-\._~:/?#[\]@!$&'\(\)\*\+,;=]*)?$`)

// NewFavoriteURL 创建URL值对象
func NewFavoriteURL(url string) (FavoriteURL, error) {
	if url == "" {
		return "", nil // URL可以为空
	}
	if !urlRegex.MatchString(url) {
		return "", errors.New("无效的URL格式")
	}
	if len(url) > 1000 {
		return "", errors.New("URL长度不能超过1000个字符")
	}
	return FavoriteURL(url), nil
}

// Value 获取URL值
func (u FavoriteURL) Value() string {
	return string(u)
}

// String 字符串表示
func (u FavoriteURL) String() string {
	return string(u)
}

// FavoriteTenantID 租户ID值对象
type FavoriteTenantID int64

// NewFavoriteTenantID 创建租户ID值对象
func NewFavoriteTenantID(id int64) (FavoriteTenantID, error) {
	if id < 0 {
		return 0, errors.New("租户ID不能为负数")
	}
	return FavoriteTenantID(id), nil
}

// Value 获取ID值
func (id FavoriteTenantID) Value() int64 {
	return int64(id)
}

// String 字符串表示
func (id FavoriteTenantID) String() string {
	return fmt.Sprintf("%d", id)
}
