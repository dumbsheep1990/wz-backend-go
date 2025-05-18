chon'g'gpackage repository

import (
	"context"
	"wz-backend-go/services/admin-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	GetUserList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.User, int64, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	UpdateUser(ctx context.Context, id int64, user *model.User) error
	DeleteUser(ctx context.Context, id int64) error
	ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, id int64, password string) error
	SetUserAuthority(ctx context.Context, userId int64, authorityId string) error
	SetUserAuthorities(ctx context.Context, userId int64, authorityIds []string) error
	GetUserAuthorities(ctx context.Context, userId int64) ([]string, error)
}

// AuthorityRepository 角色权限仓库接口
type AuthorityRepository interface {
	CreateAuthority(ctx context.Context, authority *model.SysAuthority) error
	UpdateAuthority(ctx context.Context, authority *model.SysAuthority) error
	DeleteAuthority(ctx context.Context, authorityId string) error
	GetAuthorityById(ctx context.Context, authorityId string) (*model.SysAuthority, error)
	GetAuthorityList(ctx context.Context) ([]*model.SysAuthority, error)
	SetDataAuthority(ctx context.Context, authorityId string, dataAuthorityIds []string) error
	GetAuthorityMenus(ctx context.Context, authorityId string) ([]*model.SysMenu, error)
}

// MenuRepository 菜单仓库接口
type MenuRepository interface {
	CreateMenu(ctx context.Context, menu *model.SysMenu, meta *model.SysMenuMeta) (int64, error)
	UpdateMenu(ctx context.Context, menu *model.SysMenu, meta *model.SysMenuMeta) error
	DeleteMenu(ctx context.Context, id int64) error
	GetMenuById(ctx context.Context, id int64) (*model.SysMenu, *model.SysMenuMeta, error)
	GetMenuList(ctx context.Context) ([]*model.SysMenu, error)
	GetMenusByAuthorityId(ctx context.Context, authorityId string) ([]*model.SysMenu, error)
	AddMenuAuthority(ctx context.Context, menuIds []int64, authorityId string) error
}

// ApiRepository API仓库接口
type ApiRepository interface {
	CreateApi(ctx context.Context, api *model.SysApi) error
	UpdateApi(ctx context.Context, api *model.SysApi) error
	DeleteApi(ctx context.Context, id int64) error
	DeleteApisByIds(ctx context.Context, ids []int64) error
	GetApiById(ctx context.Context, id int64) (*model.SysApi, error)
	GetApiList(ctx context.Context, page, pageSize int, path, description, apiGroup, method string) ([]*model.SysApi, int64, error)
	GetAllApis(ctx context.Context) ([]*model.SysApi, error)
}

// CasbinRepository Casbin策略仓库接口
type CasbinRepository interface {
	GetPolicyPathByAuthorityId(ctx context.Context, authorityId string) ([]model.CasbinPolicyPath, error)
	UpdateCasbin(ctx context.Context, authorityId string, policies []model.CasbinPolicyPath) error
}

// DictionaryRepository 字典仓库接口
type DictionaryRepository interface {
	CreateDictionary(ctx context.Context, dictionary *model.SysDictionary) error
	UpdateDictionary(ctx context.Context, dictionary *model.SysDictionary) error
	DeleteDictionary(ctx context.Context, id int64) error
	GetDictionaryById(ctx context.Context, id int64) (*model.SysDictionary, error)
	GetDictionaryList(ctx context.Context, page, pageSize int, name, typeName string, status int) ([]*model.SysDictionary, int64, error)
}

// DictionaryDetailRepository 字典详情仓库接口
type DictionaryDetailRepository interface {
	CreateDictionaryDetail(ctx context.Context, detail *model.SysDictionaryDetail) error
	UpdateDictionaryDetail(ctx context.Context, detail *model.SysDictionaryDetail) error
	DeleteDictionaryDetail(ctx context.Context, id int64) error
	GetDictionaryDetailById(ctx context.Context, id int64) (*model.SysDictionaryDetail, error)
	GetDictionaryDetailList(ctx context.Context, page, pageSize int, dictionaryId int64, label, value string, status int) ([]*model.SysDictionaryDetail, int64, error)
}

// OperationRepository 操作日志仓库接口
type OperationRepository interface {
	CreateOperation(ctx context.Context, record *model.SysOperationRecord) error
	DeleteOperation(ctx context.Context, id int64) error
	DeleteOperationsByIds(ctx context.Context, ids []int64) error
	GetOperationById(ctx context.Context, id int64) (*model.SysOperationRecord, error)
	GetOperationList(ctx context.Context, page, pageSize int, method, path string, status int, userId int64, startTime, endTime string) ([]*model.SysOperationRecord, int64, error)
}

// ParamsRepository 系统参数仓库接口
type ParamsRepository interface {
	CreateParams(ctx context.Context, params *model.SysParams) error
	UpdateParams(ctx context.Context, params *model.SysParams) error
	DeleteParams(ctx context.Context, id int64) error
	GetParamsById(ctx context.Context, id int64) (*model.SysParams, error)
	GetParamsByKey(ctx context.Context, key string) (*model.SysParams, error)
	GetParamsList(ctx context.Context, page, pageSize int, paramName, paramKey, paramValue, paramType, paramGroup string, status int) ([]*model.SysParams, int64, error)
}

// SystemConfigRepository 系统配置仓库接口
type SystemConfigRepository interface {
	GetSystemConfig(ctx context.Context) (*model.SystemConfig, error)
	UpdateSystemConfig(ctx context.Context, config *model.SystemConfig) error
}

// JwtRepository JWT黑名单仓库接口
type JwtRepository interface {
	CreateJwtBlacklist(ctx context.Context, jwt *model.JwtBlacklist) error
	IsJwtInBlacklist(ctx context.Context, jwt string) (bool, error)
	CleanupExpiredJwt(ctx context.Context) error
}

// TenantRepository 租户仓库接口
type TenantRepository interface {
	GetTenantList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.Tenant, int64, error)
	GetTenantByID(ctx context.Context, id int64) (*model.Tenant, error)
	CreateTenant(ctx context.Context, tenant *model.Tenant) (int64, error)
	UpdateTenant(ctx context.Context, id int64, tenant *model.Tenant) error
	DeleteTenant(ctx context.Context, id int64) error
}

// ContentRepository 内容仓库接口
type ContentRepository interface {
	GetContentList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.Content, int64, error)
	GetContentByID(ctx context.Context, id int64) (*model.Content, error)
	UpdateContentStatus(ctx context.Context, id int64, status int, reason string) error
	DeleteContent(ctx context.Context, id int64) error
	RecommendContent(ctx context.Context, contentId int64, priority int) error
	CancelRecommendation(ctx context.Context, id int64) error
}

// TradeRepository 交易仓库接口
type TradeRepository interface {
	GetOrderList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.Order, int64, error)
	GetOrderByID(ctx context.Context, id int64) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, status int, remark string) error
	DeleteOrder(ctx context.Context, id int64) error
}

// SettingsRepository 系统设置仓库接口
type SettingsRepository interface {
	GetSettings(ctx context.Context) (map[string]string, error)
	UpdateSetting(ctx context.Context, key, value string) error
}

// AdminRepository 管理员仓库接口
type AdminRepository interface {
	GetAdminList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.Admin, int64, error)
	GetAdminByID(ctx context.Context, id int64) (*model.Admin, error)
	GetAdminByUsername(ctx context.Context, username string) (*model.Admin, error)
	CreateAdmin(ctx context.Context, admin *model.Admin) (int64, error)
	UpdateAdmin(ctx context.Context, id int64, admin *model.Admin) error
	DeleteAdmin(ctx context.Context, id int64) error
	UpdateLastLogin(ctx context.Context, id int64) error
	VerifyLogin(ctx context.Context, username, password string) (*model.Admin, error)
}

// RoleRepository 角色仓库接口
type RoleRepository interface {
	GetRoleList(ctx context.Context, page, pageSize int) ([]*model.Role, int64, error)
	GetRoleByID(ctx context.Context, id int64) (*model.Role, error)
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role) (int64, error)
	UpdateRole(ctx context.Context, id int64, role *model.Role) error
	DeleteRole(ctx context.Context, id int64) error
}

// OperationLogRepository 操作日志仓库接口
type OperationLogRepository interface {
	GetOperationLogList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.OperationLog, int64, error)
	GetOperationLogByID(ctx context.Context, id int64) (*model.OperationLog, error)
	CreateOperationLog(ctx context.Context, log *model.OperationLog) error
}

// 以下是仓库实现的工厂函数

// NewUserRepository 创建用户仓库
func NewUserRepository(conn sqlx.SqlConn) UserRepository {
	// TODO: 实现用户仓库
	return &userRepository{conn}
}

// NewTenantRepository 创建租户仓库
func NewTenantRepository(conn sqlx.SqlConn) TenantRepository {
	// TODO: 实现租户仓库
	return &tenantRepository{conn}
}

// NewContentRepository 创建内容仓库
func NewContentRepository(conn sqlx.SqlConn) ContentRepository {
	// TODO: 实现内容仓库
	return &contentRepository{conn}
}

// NewTradeRepository 创建交易仓库
func NewTradeRepository(conn sqlx.SqlConn) TradeRepository {
	// TODO: 实现交易仓库
	return &tradeRepository{conn}
}

// NewSettingsRepository 创建系统设置仓库
func NewSettingsRepository(conn sqlx.SqlConn) SettingsRepository {
	// TODO: 实现系统设置仓库
	return &settingsRepository{conn}
}

// NewAdminRepository 创建管理员仓库
func NewAdminRepository(conn sqlx.SqlConn) AdminRepository {
	// TODO: 实现管理员仓库
	return &adminRepository{conn}
}

// NewRoleRepository 创建角色仓库
func NewRoleRepository(conn sqlx.SqlConn) RoleRepository {
	// TODO: 实现角色仓库
	return nil
}

// NewOperationLogRepository 创建操作日志仓库
func NewOperationLogRepository(conn sqlx.SqlConn) OperationLogRepository {
	// TODO: 实现操作日志仓库
	return nil
}

// NewAuthorityRepository 创建角色权限仓库
func NewAuthorityRepository(conn sqlx.SqlConn) AuthorityRepository {
	return &authorityRepository{conn}
}

// NewMenuRepository 创建菜单仓库
func NewMenuRepository(conn sqlx.SqlConn) MenuRepository {
	return &menuRepository{conn}
}

// NewApiRepository 创建API仓库
func NewApiRepository(conn sqlx.SqlConn) ApiRepository {
	return &apiRepository{conn}
}

// NewCasbinRepository 创建Casbin策略仓库
func NewCasbinRepository(conn sqlx.SqlConn) CasbinRepository {
	return &casbinRepository{conn}
}

// NewDictionaryRepository 创建字典仓库
func NewDictionaryRepository(conn sqlx.SqlConn) DictionaryRepository {
	return &dictionaryRepository{conn}
}

// NewDictionaryDetailRepository 创建字典详情仓库
func NewDictionaryDetailRepository(conn sqlx.SqlConn) DictionaryDetailRepository {
	return &dictionaryDetailRepository{conn}
}

// NewOperationRepository 创建操作日志仓库
func NewOperationRepository(conn sqlx.SqlConn) OperationRepository {
	return &operationRepository{conn}
}

// NewParamsRepository 创建系统参数仓库
func NewParamsRepository(conn sqlx.SqlConn) ParamsRepository {
	return &paramsRepository{conn}
}

// NewSystemConfigRepository 创建系统配置仓库
func NewSystemConfigRepository(conn sqlx.SqlConn) SystemConfigRepository {
	return &systemConfigRepository{conn}
}

// NewJwtRepository 创建JWT黑名单仓库
func NewJwtRepository(conn sqlx.SqlConn) JwtRepository {
	return &jwtRepository{conn}
}

// 存储库基础实现

type userRepository struct {
	conn sqlx.SqlConn
}

type tenantRepository struct {
	conn sqlx.SqlConn
}

type contentRepository struct {
	conn sqlx.SqlConn
}

type tradeRepository struct {
	conn sqlx.SqlConn
}

type settingsRepository struct {
	conn sqlx.SqlConn
}

type adminRepository struct {
	conn sqlx.SqlConn
}

type authorityRepository struct {
	conn sqlx.SqlConn
}

type menuRepository struct {
	conn sqlx.SqlConn
}

type apiRepository struct {
	conn sqlx.SqlConn
}

type casbinRepository struct {
	conn sqlx.SqlConn
}

type dictionaryRepository struct {
	conn sqlx.SqlConn
}

type dictionaryDetailRepository struct {
	conn sqlx.SqlConn
}

type operationRepository struct {
	conn sqlx.SqlConn
}

type paramsRepository struct {
	conn sqlx.SqlConn
}

type systemConfigRepository struct {
	conn sqlx.SqlConn
}

type jwtRepository struct {
	conn sqlx.SqlConn
}
