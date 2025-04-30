package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pquerna/otp/totp"
)

// MFAService 多因素认证服务接口
type MFAService interface {
	// 生成TOTP密钥
	GenerateSecret(ctx context.Context, userID int64) (string, string, error)
	// 验证TOTP代码
	ValidateCode(ctx context.Context, userID int64, code string) (bool, error)
	// 启用MFA
	EnableMFA(ctx context.Context, userID int64, secret string) error
	// 禁用MFA
	DisableMFA(ctx context.Context, userID int64) error
	// 检查MFA是否已启用
	IsMFAEnabled(ctx context.Context, userID int64) (bool, error)
}

type mfaService struct {
	redis        *redis.Client
	issuerName   string
	keyValidTime time.Duration
}

// NewMFAService 创建MFA服务
func NewMFAService(redis *redis.Client, issuerName string) MFAService {
	return &mfaService{
		redis:        redis,
		issuerName:   issuerName,
		keyValidTime: 30 * time.Second, // TOTP标准验证窗口
	}
}

// GenerateSecret 生成TOTP密钥
func (s *mfaService) GenerateSecret(ctx context.Context, userID int64) (string, string, error) {
	// 生成TOTP密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuerName,
		AccountName: fmt.Sprintf("user:%d", userID),
		SecretSize:  20, // 推荐的密钥长度
	})
	if err != nil {
		return "", "", err
	}

	// 返回密钥和二维码URL
	return key.Secret(), key.URL(), nil
}

// ValidateCode 验证TOTP代码
func (s *mfaService) ValidateCode(ctx context.Context, userID int64, code string) (bool, error) {
	// 从Redis获取用户的MFA密钥
	secretKey := fmt.Sprintf("mfa:secret:%d", userID)
	secret, err := s.redis.Get(ctx, secretKey).Result()
	if err != nil {
		return false, err
	}

	// 验证代码
	valid := totp.Validate(code, secret)
	if !valid {
		return false, nil
	}

	// 检查代码是否已使用（防重放攻击）
	codeKey := fmt.Sprintf("mfa:used:%d:%s", userID, code)
	used, _ := s.redis.Exists(ctx, codeKey).Result()
	if used > 0 {
		return false, nil
	}

	// 标记代码为已使用
	s.redis.SetEX(ctx, codeKey, "1", time.Minute*5)

	return true, nil
}

// EnableMFA 启用MFA
func (s *mfaService) EnableMFA(ctx context.Context, userID int64, secret string) error {
	// 存储MFA密钥
	secretKey := fmt.Sprintf("mfa:secret:%d", userID)
	enabledKey := fmt.Sprintf("mfa:enabled:%d", userID)

	// 存储密钥和启用状态
	pipe := s.redis.Pipeline()
	pipe.Set(ctx, secretKey, secret, 0) // 永不过期
	pipe.Set(ctx, enabledKey, "1", 0)   // 永不过期
	_, err := pipe.Exec(ctx)

	return err
}

// DisableMFA 禁用MFA
func (s *mfaService) DisableMFA(ctx context.Context, userID int64) error {
	// 删除MFA相关键
	secretKey := fmt.Sprintf("mfa:secret:%d", userID)
	enabledKey := fmt.Sprintf("mfa:enabled:%d", userID)

	pipe := s.redis.Pipeline()
	pipe.Del(ctx, secretKey)
	pipe.Del(ctx, enabledKey)
	_, err := pipe.Exec(ctx)

	return err
}

// IsMFAEnabled 检查MFA是否已启用
func (s *mfaService) IsMFAEnabled(ctx context.Context, userID int64) (bool, error) {
	enabledKey := fmt.Sprintf("mfa:enabled:%d", userID)
	result, err := s.redis.Exists(ctx, enabledKey).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}
