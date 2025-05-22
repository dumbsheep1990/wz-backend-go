package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SimilarCategoryType 同X分类类型
type SimilarCategoryType string

// 同X分类常量
const (
	SimilarCategoryTongYong  SimilarCategoryType = "同用" // 同用
	SimilarCategoryTongHao   SimilarCategoryType = "同好" // 同好
	SimilarCategoryTongGou   SimilarCategoryType = "同购" // 同购
	SimilarCategoryTongNian  SimilarCategoryType = "同年" // 同年
	SimilarCategoryTongYou   SimilarCategoryType = "同游" // 同游
	SimilarCategoryTongZai   SimilarCategoryType = "同在" // 同在
	SimilarCategoryTongShi   SimilarCategoryType = "同市" // 同市
	SimilarCategoryTongQi    SimilarCategoryType = "同企" // 同企
	SimilarCategoryTongQin   SimilarCategoryType = "同亲" // 同亲
	SimilarCategoryTongBan   SimilarCategoryType = "同班" // 同班
	SimilarCategoryTongShi2  SimilarCategoryType = "同师" // 同师
	SimilarCategoryTongYe    SimilarCategoryType = "同业" // 同业
	SimilarCategoryTongWang  SimilarCategoryType = "同网" // 同网
	SimilarCategoryTongGong  SimilarCategoryType = "同工" // 同工
	SimilarCategoryTongWu    SimilarCategoryType = "同务" // 同务
	SimilarCategoryTongYi    SimilarCategoryType = "同艺" // 同艺
	SimilarCategoryTongWan   SimilarCategoryType = "同玩" // 同玩
	SimilarCategoryTongXian  SimilarCategoryType = "同闲" // 同闲
	SimilarCategoryTongPai   SimilarCategoryType = "同拍" // 同拍
	SimilarCategoryTongXiang SimilarCategoryType = "同乡" // 同乡
	SimilarCategoryTongXue   SimilarCategoryType = "同学" // 同学
)

// GetAllSimilarCategories 获取所有同X分类
func GetAllSimilarCategories() []SimilarCategoryType {
	return []SimilarCategoryType{
		SimilarCategoryTongYong,
		SimilarCategoryTongHao,
		SimilarCategoryTongGou,
		SimilarCategoryTongNian,
		SimilarCategoryTongYou,
		SimilarCategoryTongZai,
		SimilarCategoryTongShi,
		SimilarCategoryTongQi,
		SimilarCategoryTongQin,
		SimilarCategoryTongBan,
		SimilarCategoryTongShi2,
		SimilarCategoryTongYe,
		SimilarCategoryTongWang,
		SimilarCategoryTongGong,
		SimilarCategoryTongWu,
		SimilarCategoryTongYi,
		SimilarCategoryTongWan,
		SimilarCategoryTongXian,
		SimilarCategoryTongPai,
		SimilarCategoryTongXiang,
		SimilarCategoryTongXue,
	}
}

// SimilarApplication 入同申请模型
type SimilarApplication struct {
	ID             string    `json:"id" gorm:"type:varchar(36);primary_key"`
	UserID         string    `json:"user_id" gorm:"type:varchar(36);index"`
	ApplicationType string    `json:"application_type" gorm:"type:varchar(20);comment:申请类型(同用/同好/同购等)"`
	Name           string    `json:"name" gorm:"type:varchar(50);comment:姓名"`
	Gender         string    `json:"gender" gorm:"type:varchar(10);comment:性别"`
	Birthplace     string    `json:"birthplace" gorm:"type:varchar(100);comment:出生地点"`
	Occupation     string    `json:"occupation" gorm:"type:varchar(50);comment:职业"`
	Education      string    `json:"education" gorm:"type:varchar(50);comment:学历"`
	WorkPosition   string    `json:"work_position" gorm:"type:varchar(50);comment:工作岗位"`
	WorkPlace      string    `json:"work_place" gorm:"type:varchar(100);comment:工作地点"`
	Hobby          string    `json:"hobby" gorm:"type:text;comment:爱好"`
	Address        string    `json:"address" gorm:"type:varchar(200);comment:地址"`
	ContactType    string    `json:"contact_type" gorm:"type:varchar(20);comment:联系方式类型"`
	ContactValue   string    `json:"contact_value" gorm:"type:varchar(100);comment:联系方式值"`
	Description    string    `json:"description" gorm:"type:text;comment:个人简介"`
	Status         string    `json:"status" gorm:"type:varchar(20);default:pending;comment:状态(pending/approved/rejected)"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BeforeCreate 创建记录前自动生成UUID
func (s *SimilarApplication) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// TableName 指定表名
func (SimilarApplication) TableName() string {
	return "similar_applications"
}
