package user

import "context"

// Repository 用户仓储接口
// 遵循DDD原则，仓储负责领域对象的持久化，但不包含业务逻辑
type Repository interface {
	// 基本CRUD操作
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id UserID) (*User, error)
	FindByUsername(ctx context.Context, username Username) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
	FindByPhone(ctx context.Context, phone Phone) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id UserID) error
	
	// 业务查询方法
	ExistsByUsername(ctx context.Context, username Username) (bool, error)
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
	ExistsByPhone(ctx context.Context, phone Phone) (bool, error)
	FindUsersByTenant(ctx context.Context, tenantID TenantID, offset, limit int) ([]*User, error)
	CountUsersByTenant(ctx context.Context, tenantID TenantID) (int, error)
	FindUsersByRole(ctx context.Context, role UserRole, offset, limit int) ([]*User, error)
	CountUsersByRole(ctx context.Context, role UserRole) (int, error)
	FindActiveUsers(ctx context.Context, offset, limit int) ([]*User, error)
	CountActiveUsers(ctx context.Context) (int, error)
}

// ProfileRepository 用户详情仓储接口
type ProfileRepository interface {
	// 基本CRUD操作
	SaveProfile(ctx context.Context, profile *UserProfile) error
	FindProfileByUserID(ctx context.Context, userID UserID) (*UserProfile, error)
	UpdateProfile(ctx context.Context, profile *UserProfile) error
	DeleteProfile(ctx context.Context, userID UserID) error
}

// UserProfile 用户详细资料实体
type UserProfile struct {
	id          UserProfileID
	userID      UserID
	realName    string
	idCardType  string
	idCardNo    string
	address     Address
	company     string
	position    string
	interests   []string
	education   []EducationExperience
	workHistory []WorkExperience
	isVerified  bool
	verifiedAt  time.Time
	tenantID    TenantID
	createdAt   time.Time
	updatedAt   time.Time
}

// UserProfileID 用户详情ID值对象
type UserProfileID int64

// NewUserProfileID 创建用户详情ID值对象
func NewUserProfileID(id int64) UserProfileID {
	return UserProfileID(id)
}

// Value 获取用户详情ID值
func (id UserProfileID) Value() int64 {
	return int64(id)
}

// EducationExperience 教育经历值对象
type EducationExperience struct {
	School     string
	Degree     string
	Major      string
	StartDate  time.Time
	EndDate    time.Time
	Achievement string
}

// WorkExperience 工作经历值对象
type WorkExperience struct {
	Company     string
	Position    string
	Department  string
	StartDate   time.Time
	EndDate     time.Time
	Description string
}

// LoginRepository 登录日志仓储接口
type LoginRepository interface {
	// 保存登录日志
	SaveLoginLog(ctx context.Context, log *UserLoginLog) error
	// 查询用户登录日志
	FindLoginLogsByUserID(ctx context.Context, userID UserID, offset, limit int) ([]*UserLoginLog, error)
	// 统计用户登录次数
	CountLoginLogsByUserID(ctx context.Context, userID UserID) (int, error)
	// 查询最近一次登录记录
	FindLastLoginByUserID(ctx context.Context, userID UserID) (*UserLoginLog, error)
}

// UserLoginLog 用户登录日志实体
type UserLoginLog struct {
	id          UserLoginLogID
	userID      UserID
	loginType   string
	ip          string
	userAgent   string
	device      string
	os          string
	browser     string
	location    string
	status      int
	errorMsg    string
	loginAt     time.Time
	tenantID    TenantID
}

// UserLoginLogID 登录日志ID值对象
type UserLoginLogID int64

// NewUserLoginLogID 创建登录日志ID值对象
func NewUserLoginLogID(id int64) UserLoginLogID {
	return UserLoginLogID(id)
}

// Value 获取登录日志ID值
func (id UserLoginLogID) Value() int64 {
	return int64(id)
}

// BehaviorRepository 用户行为仓储接口
type BehaviorRepository interface {
	// 保存用户行为日志
	SaveBehaviorLog(ctx context.Context, log *UserBehaviorLog) error
	// 查询用户行为日志
	FindBehaviorLogsByUserID(ctx context.Context, userID UserID, offset, limit int) ([]*UserBehaviorLog, error)
	// 统计用户行为次数
	CountBehaviorLogsByUserID(ctx context.Context, userID UserID) (int, error)
	// 查询特定类型的用户行为
	FindBehaviorLogsByAction(ctx context.Context, userID UserID, action string, offset, limit int) ([]*UserBehaviorLog, error)
}

// UserBehaviorLog 用户行为日志实体
type UserBehaviorLog struct {
	id           UserBehaviorLogID
	userID       UserID
	action       string
	resourceType string
	resourceID   int64
	ip           string
	userAgent    string
	createdAt    time.Time
}

// UserBehaviorLogID 行为日志ID值对象
type UserBehaviorLogID int64

// NewUserBehaviorLogID 创建行为日志ID值对象
func NewUserBehaviorLogID(id int64) UserBehaviorLogID {
	return UserBehaviorLogID(id)
}

// Value 获取行为日志ID值
func (id UserBehaviorLogID) Value() int64 {
	return int64(id)
}
