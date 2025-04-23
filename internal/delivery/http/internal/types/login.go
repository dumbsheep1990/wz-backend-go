package types

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	TenantID int64  `json:"tenant_id,omitempty"` // 可选的租户ID，如果是租户用户登录需要指定
}

// LoginResp 登录响应
type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
}
