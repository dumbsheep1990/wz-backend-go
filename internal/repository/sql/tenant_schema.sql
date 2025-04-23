-- 租户表
CREATE TABLE IF NOT EXISTS tenants (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '租户名称',
    subdomain VARCHAR(50) NOT NULL COMMENT '子域名',
    tenant_type INT NOT NULL COMMENT '租户类型：1-企业，2-个人，3-教育机构',
    description TEXT COMMENT '租户描述',
    logo VARCHAR(255) COMMENT '租户Logo URL',
    creator_user_id BIGINT NOT NULL COMMENT '创建者用户ID',
    status INT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，2-禁用',
    expiration_date TIMESTAMP NULL COMMENT '过期时间',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_tenant_subdomain (subdomain),
    INDEX idx_tenant_creator (creator_user_id),
    INDEX idx_tenant_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表';

-- 租户-用户关联表
CREATE TABLE IF NOT EXISTS tenant_users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role VARCHAR(20) NOT NULL COMMENT '角色：tenant_admin-租户管理员，tenant_user-租户普通用户',
    status INT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，2-禁用',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_tenant_user (tenant_id, user_id),
    INDEX idx_tenant_users_tenant (tenant_id),
    INDEX idx_tenant_users_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户-用户关联表';

-- 租户配置表
CREATE TABLE IF NOT EXISTS tenant_configs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL COMMENT '租户ID',
    config_key VARCHAR(50) NOT NULL COMMENT '配置键',
    config_value TEXT NOT NULL COMMENT '配置值',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_tenant_config (tenant_id, config_key),
    INDEX idx_tenant_config_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户配置表';

-- 租户导航分类表
CREATE TABLE IF NOT EXISTS tenant_categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL COMMENT '租户ID',
    name VARCHAR(50) NOT NULL COMMENT '分类名称',
    url VARCHAR(100) NOT NULL COMMENT '分类URL',
    icon VARCHAR(50) COMMENT '分类图标',
    description VARCHAR(255) COMMENT '分类描述',
    parent_id BIGINT DEFAULT 0 COMMENT '父分类ID，0表示顶级分类',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
    status INT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，2-禁用',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant_category_tenant (tenant_id),
    INDEX idx_tenant_category_parent (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户导航分类表';

-- 租户域名绑定表
CREATE TABLE IF NOT EXISTS tenant_domains (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL COMMENT '租户ID',
    domain VARCHAR(100) NOT NULL COMMENT '自定义域名',
    ssl_enabled TINYINT NOT NULL DEFAULT 0 COMMENT '是否启用SSL',
    certificate_id BIGINT COMMENT 'SSL证书ID',
    status INT NOT NULL DEFAULT 1 COMMENT '状态：1-待验证，2-已验证，3-验证失败',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_tenant_domain (domain),
    INDEX idx_tenant_domain_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户域名绑定表';

-- 租户静态页面表
CREATE TABLE IF NOT EXISTS tenant_static_pages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL COMMENT '租户ID',
    page_key VARCHAR(50) NOT NULL COMMENT '页面标识符，如privacy, terms',
    title VARCHAR(100) NOT NULL COMMENT '页面标题',
    content TEXT NOT NULL COMMENT '页面内容，支持HTML',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_tenant_page (tenant_id, page_key),
    INDEX idx_tenant_page_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户静态页面表';

-- 更新用户表，添加默认租户ID和角色字段
ALTER TABLE users 
    ADD COLUMN default_tenant_id BIGINT NULL COMMENT '默认租户ID' AFTER is_company_verified,
    ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'personal_user' COMMENT '用户角色：platform_admin-平台管理员，tenant_admin-租户管理员，tenant_user-租户普通用户，personal_user-个人用户' AFTER default_tenant_id,
    ADD INDEX idx_user_default_tenant (default_tenant_id);
