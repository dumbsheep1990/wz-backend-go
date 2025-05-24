package community

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ID 是领域实体的通用标识符值对象
type ID struct {
	value string
}

// NewID 创建新的ID
func NewID() ID {
	return ID{value: uuid.New().String()}
}

// NewIDFromString 从字符串创建ID
func NewIDFromString(id string) (ID, error) {
	if id == "" {
		return ID{}, errors.New("ID不能为空")
	}
	
	// 检查是否是有效的UUID格式（简化校验）
	if len(id) != 36 {
		return ID{}, errors.New("无效的ID格式")
	}
	
	return ID{value: id}, nil
}

// Value 获取ID值
func (id ID) Value() string {
	return id.value
}

// String 实现Stringer接口
func (id ID) String() string {
	return id.value
}

// 社区名称值对象
type CommunityName struct {
	value string
}

// NewCommunityName 创建社区名称
func NewCommunityName(name string) (CommunityName, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return CommunityName{}, errors.New("社区名称不能为空")
	}
	
	if len(name) > 50 {
		return CommunityName{}, errors.New("社区名称不能超过50个字符")
	}
	
	return CommunityName{value: name}, nil
}

// Value 获取名称值
func (n CommunityName) Value() string {
	return n.value
}

// 描述值对象
type Description struct {
	value string
}

// NewDescription 创建描述
func NewDescription(desc string) (Description, error) {
	desc = strings.TrimSpace(desc)
	if len(desc) > 500 {
		return Description{}, errors.New("描述不能超过500个字符")
	}
	
	return Description{value: desc}, nil
}

// Value 获取描述值
func (d Description) Value() string {
	return d.value
}

// 内容值对象
type Content struct {
	value string
}

// NewContent 创建内容
func NewContent(content string) (Content, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return Content{}, errors.New("内容不能为空")
	}
	
	if len(content) > 10000 {
		return Content{}, errors.New("内容不能超过10000个字符")
	}
	
	return Content{value: content}, nil
}

// Value 获取内容值
func (c Content) Value() string {
	return c.value
}

// 标题值对象
type Title struct {
	value string
}

// NewTitle 创建标题
func NewTitle(title string) (Title, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Title{}, errors.New("标题不能为空")
	}
	
	if len(title) > 100 {
		return Title{}, errors.New("标题不能超过100个字符")
	}
	
	return Title{value: title}, nil
}

// Value 获取标题值
func (t Title) Value() string {
	return t.value
}

// Tag 标签值对象
type Tag struct {
	value string
}

// NewTag 创建标签
func NewTag(tag string) (Tag, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return Tag{}, errors.New("标签不能为空")
	}
	
	if len(tag) > 20 {
		return Tag{}, errors.New("标签不能超过20个字符")
	}
	
	// 标签只能包含字母、数字、中文和连字符
	match, _ := regexp.MatchString(`^[a-zA-Z0-9\p{Han}\-]+$`, tag)
	if !match {
		return Tag{}, errors.New("标签只能包含字母、数字、中文和连字符")
	}
	
	return Tag{value: tag}, nil
}

// Value 获取标签值
func (t Tag) Value() string {
	return t.value
}

// 标签列表值对象
type Tags []Tag

// NewTags 创建标签列表
func NewTags(tags []string) (Tags, error) {
	if len(tags) > 10 {
		return nil, errors.New("标签数量不能超过10个")
	}
	
	result := make(Tags, 0, len(tags))
	for _, t := range tags {
		tag, err := NewTag(t)
		if err != nil {
			return nil, err
		}
		result = append(result, tag)
	}
	
	return result, nil
}

// Values 获取标签值列表
func (ts Tags) Values() []string {
	result := make([]string, len(ts))
	for i, t := range ts {
		result[i] = t.Value()
	}
	return result
}

// 位置值对象
type Location struct {
	value string
}

// NewLocation 创建位置
func NewLocation(location string) (Location, error) {
	location = strings.TrimSpace(location)
	if len(location) > 100 {
		return Location{}, errors.New("位置不能超过100个字符")
	}
	
	return Location{value: location}, nil
}

// Value 获取位置值
func (l Location) Value() string {
	return l.value
}

// 图片URL值对象
type ImageURL struct {
	value string
}

// NewImageURL 创建图片URL
func NewImageURL(url string) (ImageURL, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return ImageURL{}, errors.New("图片URL不能为空")
	}
	
	// 简单的URL格式验证
	match, _ := regexp.MatchString(`^(http|https)://`, url)
	if !match {
		return ImageURL{}, errors.New("无效的图片URL格式")
	}
	
	return ImageURL{value: url}, nil
}

// Value 获取图片URL值
func (i ImageURL) Value() string {
	return i.value
}

// 图片URL列表值对象
type ImageURLs []ImageURL

// NewImageURLs 创建图片URL列表
func NewImageURLs(urls []string) (ImageURLs, error) {
	if len(urls) > 20 {
		return nil, errors.New("图片数量不能超过20张")
	}
	
	result := make(ImageURLs, 0, len(urls))
	for _, u := range urls {
		url, err := NewImageURL(u)
		if err != nil {
			return nil, err
		}
		result = append(result, url)
	}
	
	return result, nil
}

// Values 获取图片URL值列表
func (is ImageURLs) Values() []string {
	result := make([]string, len(is))
	for i, img := range is {
		result[i] = img.Value()
	}
	return result
}

// 时间戳值对象
type Timestamp struct {
	value time.Time
}

// NewTimestamp 创建当前时间戳
func NewTimestamp() Timestamp {
	return Timestamp{value: time.Now()}
}

// NewTimestampFromTime 从time.Time创建时间戳
func NewTimestampFromTime(t time.Time) Timestamp {
	return Timestamp{value: t}
}

// Value 获取时间值
func (t Timestamp) Value() time.Time {
	return t.value
}

// String 实现Stringer接口，返回格式化的时间
func (t Timestamp) String() string {
	return t.value.Format(time.RFC3339)
}

// 社区状态枚举
type CommunityStatus int

const (
	CommunityStatusUnspecified CommunityStatus = iota
	CommunityStatusActive
	CommunityStatusInactive
	CommunityStatusDeleted
)

// String 实现Stringer接口
func (s CommunityStatus) String() string {
	switch s {
	case CommunityStatusActive:
		return "ACTIVE"
	case CommunityStatusInactive:
		return "INACTIVE"
	case CommunityStatusDeleted:
		return "DELETED"
	default:
		return "UNSPECIFIED"
	}
}

// 小组状态枚举
type GroupStatus int

const (
	GroupStatusUnspecified GroupStatus = iota
	GroupStatusActive
	GroupStatusInactive
	GroupStatusDeleted
)

// String 实现Stringer接口
func (s GroupStatus) String() string {
	switch s {
	case GroupStatusActive:
		return "ACTIVE"
	case GroupStatusInactive:
		return "INACTIVE"
	case GroupStatusDeleted:
		return "DELETED"
	default:
		return "UNSPECIFIED"
	}
}

// 帖子状态枚举
type PostStatus int

const (
	PostStatusUnspecified PostStatus = iota
	PostStatusActive
	PostStatusInactive
	PostStatusDeleted
)

// String 实现Stringer接口
func (s PostStatus) String() string {
	switch s {
	case PostStatusActive:
		return "ACTIVE"
	case PostStatusInactive:
		return "INACTIVE"
	case PostStatusDeleted:
		return "DELETED"
	default:
		return "UNSPECIFIED"
	}
}

// 评论状态枚举
type CommentStatus int

const (
	CommentStatusUnspecified CommentStatus = iota
	CommentStatusActive
	CommentStatusDeleted
)

// String 实现Stringer接口
func (s CommentStatus) String() string {
	switch s {
	case CommentStatusActive:
		return "ACTIVE"
	case CommentStatusDeleted:
		return "DELETED"
	default:
		return "UNSPECIFIED"
	}
}

// 根据proto定义中的同X分类，创建社区类型枚举
type CommunityType string

// 社区类型常量
const (
	CommunityTypeTongYong  CommunityType = "同用" // 同用
	CommunityTypeTongHao   CommunityType = "同好" // 同好
	CommunityTypeTongGou   CommunityType = "同购" // 同购
	CommunityTypeTongNian  CommunityType = "同年" // 同年
	CommunityTypeTongYou   CommunityType = "同游" // 同游
	CommunityTypeTongZai   CommunityType = "同在" // 同在
	CommunityTypeTongShi   CommunityType = "同市" // 同市
	CommunityTypeTongQi    CommunityType = "同企" // 同企
	CommunityTypeTongQin   CommunityType = "同亲" // 同亲
	CommunityTypeTongBan   CommunityType = "同班" // 同班
	CommunityTypeTongShi2  CommunityType = "同师" // 同师
	CommunityTypeTongYe    CommunityType = "同业" // 同业
	CommunityTypeTongWang  CommunityType = "同网" // 同网
	CommunityTypeTongGong  CommunityType = "同工" // 同工
	CommunityTypeTongWu    CommunityType = "同务" // 同务
	CommunityTypeTongYi    CommunityType = "同艺" // 同艺
	CommunityTypeTongWan   CommunityType = "同玩" // 同玩
	CommunityTypeTongXian  CommunityType = "同闲" // 同闲
	CommunityTypeTongPai   CommunityType = "同拍" // 同拍
	CommunityTypeTongXiang CommunityType = "同乡" // 同乡
	CommunityTypeTongXue   CommunityType = "同学" // 同学
)

// AllCommunityTypes 获取所有社区类型
func AllCommunityTypes() []CommunityType {
	return []CommunityType{
		CommunityTypeTongYong,
		CommunityTypeTongHao,
		CommunityTypeTongGou,
		CommunityTypeTongNian,
		CommunityTypeTongYou,
		CommunityTypeTongZai,
		CommunityTypeTongShi,
		CommunityTypeTongQi,
		CommunityTypeTongQin,
		CommunityTypeTongBan,
		CommunityTypeTongShi2,
		CommunityTypeTongYe,
		CommunityTypeTongWang,
		CommunityTypeTongGong,
		CommunityTypeTongWu,
		CommunityTypeTongYi,
		CommunityTypeTongWan,
		CommunityTypeTongXian,
		CommunityTypeTongPai,
		CommunityTypeTongXiang,
		CommunityTypeTongXue,
	}
}

// IsValid 检查社区类型是否有效
func (c CommunityType) IsValid() bool {
	for _, t := range AllCommunityTypes() {
		if c == t {
			return true
		}
	}
	return false
}

// String 实现Stringer接口
func (c CommunityType) String() string {
	return string(c)
}
