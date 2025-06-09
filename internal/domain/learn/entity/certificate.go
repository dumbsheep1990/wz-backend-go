package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CertificateStatus 证书状态
type CertificateStatus string

const (
	CertificateStatusIssued     CertificateStatus = "issued"     // 已颁发
	CertificateStatusRevoked    CertificateStatus = "revoked"    // 已撤销
	CertificateStatusExpired    CertificateStatus = "expired"    // 已过期
	CertificateStatusInProgress CertificateStatus = "inprogress" // 处理中
)

// Certificate 课程证书实体
type Certificate struct {
	ID              string            `json:"id"`
	UserID          string            `json:"userId"`          // 用户ID
	CourseID        string            `json:"courseId"`        // 课程ID
	EnrollmentID    string            `json:"enrollmentId"`    // 报名ID
	Title           string            `json:"title"`           // 证书标题
	Description     string            `json:"description"`     // 证书描述
	IssueDate       time.Time         `json:"issueDate"`       // 颁发日期
	ExpiryDate      *time.Time        `json:"expiryDate"`      // 过期日期，为空表示永久有效
	Status          CertificateStatus `json:"status"`          // 证书状态
	CertificateCode string            `json:"certificateCode"` // 证书编码/序列号
	VerifyURL       string            `json:"verifyUrl"`       // 验证链接
	FileURL         string            `json:"fileUrl"`         // 证书PDF文件链接
	ImageURL        string            `json:"imageUrl"`        // 证书图片链接
	Metadata        map[string]string `json:"metadata"`        // 元数据
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

// NewCertificate 创建新证书
func NewCertificate(userID, courseID, enrollmentID, title, description string) *Certificate {
	now := time.Now()
	return &Certificate{
		ID:              uuid.New().String(),
		UserID:          userID,
		CourseID:        courseID,
		EnrollmentID:    enrollmentID,
		Title:           title,
		Description:     description,
		IssueDate:       now,
		Status:          CertificateStatusInProgress,
		CertificateCode: generateCertificateCode(userID, courseID),
		Metadata:        make(map[string]string),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Issue 颁发证书
func (c *Certificate) Issue(fileURL, imageURL, verifyURL string) {
	now := time.Now()
	c.Status = CertificateStatusIssued
	c.IssueDate = now
	c.FileURL = fileURL
	c.ImageURL = imageURL
	c.VerifyURL = verifyURL
	c.UpdatedAt = now
}

// Revoke 撤销证书
func (c *Certificate) Revoke(reason string) {
	c.Status = CertificateStatusRevoked
	c.Metadata["revoke_reason"] = reason
	c.Metadata["revoke_date"] = time.Now().Format(time.RFC3339)
	c.UpdatedAt = time.Now()
}

// SetExpiry 设置过期日期
func (c *Certificate) SetExpiry(expiryDate time.Time) {
	c.ExpiryDate = &expiryDate
	c.UpdatedAt = time.Now()
	
	// 如果当前日期已经超过过期日期，则更新状态为已过期
	if time.Now().After(expiryDate) {
		c.Status = CertificateStatusExpired
	}
}

// UpdateMetadata 更新元数据
func (c *Certificate) UpdateMetadata(key, value string) {
	c.Metadata[key] = value
	c.UpdatedAt = time.Now()
}

// IsValid 检查证书是否有效
func (c *Certificate) IsValid() bool {
	// 证书必须是已颁发状态
	if c.Status != CertificateStatusIssued {
		return false
	}
	
	// 如果设置了过期日期，检查是否已过期
	if c.ExpiryDate != nil && time.Now().After(*c.ExpiryDate) {
		return false
	}
	
	return true
}

// 生成证书编码
func generateCertificateCode(userID, courseID string) string {
	timestamp := time.Now().Unix()
	prefix := "CERT"
	// 简单生成一个基于用户ID、课程ID和时间戳的唯一编码
	return fmt.Sprintf("%s-%s-%s-%d", prefix, userID[0:8], courseID[0:8], timestamp)
}
