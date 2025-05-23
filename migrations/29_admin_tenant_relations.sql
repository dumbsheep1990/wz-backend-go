-- 管理员租户关联表结构

-- 管理员表
CREATE TABLE IF NOT EXISTS `admin_users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(64) NOT NULL COMMENT '用户名',
  `password` VARCHAR(128) NOT NULL COMMENT '密码',
  `real_name` VARCHAR(64) DEFAULT NULL COMMENT '真实姓名',
  `avatar` VARCHAR(255) DEFAULT NULL COMMENT '头像',
  `email` VARCHAR(128) DEFAULT NULL COMMENT '邮箱',
  `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `is_super_admin` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否超级管理员：0-否 1-是',
  `last_login_time` TIMESTAMP NULL COMMENT '最后登录时间',
  `last_login_ip` VARCHAR(64) DEFAULT NULL COMMENT '最后登录IP',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '所属租户ID，超级管理员为NULL',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

-- 管理员租户关联表（多对多关系，一个管理员可以管理多个租户）
CREATE TABLE IF NOT EXISTS `admin_tenant_relations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `admin_id` BIGINT NOT NULL COMMENT '管理员ID',
  `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
  `role` VARCHAR(32) NOT NULL DEFAULT 'manager' COMMENT '角色：owner-拥有者 manager-管理员 operator-操作员',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admin_tenant` (`admin_id`, `tenant_id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_role` (`role`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_admin_tenant_relations_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `admin_users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_admin_tenant_relations_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员租户关联表';

-- 管理员角色表
CREATE TABLE IF NOT EXISTS `admin_roles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '角色名称',
  `code` VARCHAR(64) NOT NULL COMMENT '角色编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `sort_order` INT DEFAULT 0 COMMENT '排序',
  `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统角色：0-否 1-是',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，NULL表示通用角色',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员角色表';

-- 管理员角色关联表
CREATE TABLE IF NOT EXISTS `admin_user_roles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `admin_id` BIGINT NOT NULL COMMENT '管理员ID',
  `role_id` BIGINT NOT NULL COMMENT '角色ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admin_role` (`admin_id`, `role_id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_admin_user_roles_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `admin_users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_admin_user_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `admin_roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员角色关联表';

-- 管理员权限表
CREATE TABLE IF NOT EXISTS `admin_permissions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '权限名称',
  `code` VARCHAR(128) NOT NULL COMMENT '权限编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `type` VARCHAR(32) NOT NULL COMMENT '类型：menu-菜单 button-按钮 api-接口',
  `parent_id` BIGINT DEFAULT NULL COMMENT '父级权限ID',
  `path` VARCHAR(255) DEFAULT NULL COMMENT '路径',
  `component` VARCHAR(255) DEFAULT NULL COMMENT '组件',
  `redirect` VARCHAR(255) DEFAULT NULL COMMENT '重定向',
  `icon` VARCHAR(64) DEFAULT NULL COMMENT '图标',
  `sort_order` INT DEFAULT 0 COMMENT '排序',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统权限：0-否 1-是',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员权限表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS `admin_role_permissions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `role_id` BIGINT NOT NULL COMMENT '角色ID',
  `permission_id` BIGINT NOT NULL COMMENT '权限ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_permission` (`role_id`, `permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_admin_role_permissions_role_id` FOREIGN KEY (`role_id`) REFERENCES `admin_roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_admin_role_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `admin_permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- 初始化超级管理员
INSERT INTO `admin_users` (`username`, `password`, `real_name`, `status`, `is_super_admin`, `created_at`, `updated_at`)
VALUES ('admin', '$2a$10$V4rGBoWMRKbpODRtVuUH8eITfQ7QlZ24riKw9s3RpkQm/xrMJU.8C', '系统管理员', 1, 1, NOW(), NOW());

-- 初始化基础角色
INSERT INTO `admin_roles` (`name`, `code`, `description`, `status`, `is_system`, `created_at`, `updated_at`)
VALUES 
('超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 1, 1, NOW(), NOW()),
('平台管理员', 'platform_admin', '平台管理员，管理所有租户', 1, 1, NOW(), NOW()),
('租户管理员', 'tenant_admin', '租户管理员，管理单个租户', 1, 1, NOW(), NOW()),
('运营人员', 'operator', '运营人员，负责日常运营', 1, 1, NOW(), NOW()),
('内容管理员', 'content_manager', '内容管理员，负责内容审核', 1, 1, NOW(), NOW());
