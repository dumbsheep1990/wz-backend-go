-- 权限管理模块数据表结构

-- 角色表
CREATE TABLE IF NOT EXISTS `sys_roles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '角色名称',
  `code` VARCHAR(64) NOT NULL COMMENT '角色编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '角色描述',
  `is_default` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否默认角色：0-否 1-是',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 系统权限表
CREATE TABLE IF NOT EXISTS `sys_authorities` (
  `id` VARCHAR(64) NOT NULL COMMENT '权限ID',
  `name` VARCHAR(64) NOT NULL COMMENT '权限名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '权限描述',
  `parent_id` VARCHAR(64) DEFAULT NULL COMMENT '父权限ID',
  `default_router` VARCHAR(255) DEFAULT NULL COMMENT '默认路由',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统权限表';

-- 系统菜单表
CREATE TABLE IF NOT EXISTS `sys_menus` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '菜单名称',
  `path` VARCHAR(255) DEFAULT NULL COMMENT '菜单路径',
  `component` VARCHAR(255) DEFAULT NULL COMMENT '组件路径',
  `redirect` VARCHAR(255) DEFAULT NULL COMMENT '重定向路径',
  `hidden` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否隐藏：0-显示 1-隐藏',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `parent_id` BIGINT DEFAULT 0 COMMENT '父菜单ID',
  `active_menu` VARCHAR(255) DEFAULT NULL COMMENT '激活菜单',
  `is_iframe` TINYINT(1) DEFAULT 0 COMMENT '是否外链：0-否 1-是',
  `is_cache` TINYINT(1) DEFAULT 0 COMMENT '是否缓存：0-否 1-是',
  `close_tab` TINYINT(1) DEFAULT 0 COMMENT '关闭标签：0-否 1-是',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统菜单表';

-- 菜单元数据表
CREATE TABLE IF NOT EXISTS `sys_menu_metas` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `menu_id` BIGINT NOT NULL COMMENT '菜单ID',
  `title` VARCHAR(64) NOT NULL COMMENT '菜单标题',
  `icon` VARCHAR(64) DEFAULT NULL COMMENT '菜单图标',
  `active_icon` VARCHAR(64) DEFAULT NULL COMMENT '激活图标',
  `svg_icon` VARCHAR(64) DEFAULT NULL COMMENT 'SVG图标',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_menu_id` (`menu_id`),
  CONSTRAINT `fk_menu_metas_menu_id` FOREIGN KEY (`menu_id`) REFERENCES `sys_menus` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单元数据表';

-- 系统API表
CREATE TABLE IF NOT EXISTS `sys_apis` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `path` VARCHAR(255) NOT NULL COMMENT 'API路径',
  `method` VARCHAR(20) NOT NULL COMMENT '请求方法：GET, POST, PUT, DELETE',
  `description` VARCHAR(255) DEFAULT NULL COMMENT 'API描述',
  `api_group` VARCHAR(64) DEFAULT NULL COMMENT 'API分组',
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_path_method` (`path`, `method`),
  KEY `idx_api_group` (`api_group`),
  KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统API表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS `sys_role_authorities` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `role_id` BIGINT NOT NULL COMMENT '角色ID',
  `authority_id` VARCHAR(64) NOT NULL COMMENT '权限ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_authority` (`role_id`, `authority_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_authority_id` (`authority_id`),
  CONSTRAINT `fk_role_authorities_role_id` FOREIGN KEY (`role_id`) REFERENCES `sys_roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_role_authorities_authority_id` FOREIGN KEY (`authority_id`) REFERENCES `sys_authorities` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS `sys_user_roles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `role_id` BIGINT NOT NULL COMMENT '角色ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_role` (`user_id`, `role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`),
  CONSTRAINT `fk_user_roles_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `sys_roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 菜单权限关联表
CREATE TABLE IF NOT EXISTS `sys_menu_authorities` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `menu_id` BIGINT NOT NULL COMMENT '菜单ID',
  `authority_id` VARCHAR(64) NOT NULL COMMENT '权限ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_menu_authority` (`menu_id`, `authority_id`),
  KEY `idx_menu_id` (`menu_id`),
  KEY `idx_authority_id` (`authority_id`),
  CONSTRAINT `fk_menu_authorities_menu_id` FOREIGN KEY (`menu_id`) REFERENCES `sys_menus` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_menu_authorities_authority_id` FOREIGN KEY (`authority_id`) REFERENCES `sys_authorities` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单权限关联表';

-- Casbin策略表
CREATE TABLE IF NOT EXISTS `casbin_rules` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `ptype` VARCHAR(10) NOT NULL COMMENT '策略类型',
  `v0` VARCHAR(256) DEFAULT NULL COMMENT '主体',
  `v1` VARCHAR(256) DEFAULT NULL COMMENT '资源',
  `v2` VARCHAR(256) DEFAULT NULL COMMENT '动作',
  `v3` VARCHAR(256) DEFAULT NULL COMMENT '领域',
  `v4` VARCHAR(256) DEFAULT NULL COMMENT '条件1',
  `v5` VARCHAR(256) DEFAULT NULL COMMENT '条件2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rules` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Casbin策略表';

-- JWT黑名单表
CREATE TABLE IF NOT EXISTS `jwt_blacklists` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `jwt` TEXT NOT NULL COMMENT 'JWT令牌',
  `expires_at` TIMESTAMP NOT NULL COMMENT '过期时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='JWT黑名单表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS `sys_operation_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `ip` VARCHAR(64) DEFAULT NULL COMMENT '请求IP',
  `method` VARCHAR(20) DEFAULT NULL COMMENT '请求方法',
  `path` VARCHAR(255) DEFAULT NULL COMMENT '请求路径',
  `status` INT DEFAULT NULL COMMENT '状态码',
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `username` VARCHAR(64) DEFAULT NULL COMMENT '用户名',
  `request_data` TEXT DEFAULT NULL COMMENT '请求参数',
  `response_data` TEXT DEFAULT NULL COMMENT '响应数据',
  `user_agent` VARCHAR(255) DEFAULT NULL COMMENT '用户代理',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `cost_time` BIGINT DEFAULT NULL COMMENT '执行时间(ms)',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_method` (`method`),
  KEY `idx_path` (`path`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- 初始化基础角色数据
INSERT INTO `sys_roles` (`name`, `code`, `description`, `is_default`, `status`) VALUES
('超级管理员', 'superadmin', '系统超级管理员，拥有所有权限', 0, 1),
('管理员', 'admin', '系统管理员，拥有大部分权限', 0, 1),
('普通用户', 'user', '普通用户，拥有基本权限', 1, 1);

-- 初始化基础权限数据
INSERT INTO `sys_authorities` (`id`, `name`, `description`, `parent_id`) VALUES
('888', '超级管理员', '超级管理员权限组', NULL),
('8881', '系统管理', '系统管理权限组', '888'),
('8882', '内容管理', '内容管理权限组', '888'),
('8883', '用户管理', '用户管理权限组', '888'),
('8884', '租户管理', '租户管理权限组', '888'),
('8885', '交易管理', '交易管理权限组', '888');
