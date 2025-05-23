-- 租户模块数据表结构

-- 租户表
CREATE TABLE IF NOT EXISTS `tenants` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '租户名称',
  `subdomain` VARCHAR(64) NOT NULL COMMENT '子域名',
  `tenant_type` INT NOT NULL COMMENT '租户类型：1-企业 2-个人 3-教育机构',
  `description` TEXT DEFAULT NULL COMMENT '租户描述',
  `logo` VARCHAR(255) DEFAULT NULL COMMENT '租户Logo',
  `creator_user_id` BIGINT NOT NULL COMMENT '创建者用户ID',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1-正常 2-禁用',
  `expiration_date` TIMESTAMP NULL DEFAULT NULL COMMENT '过期时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_subdomain` (`subdomain`),
  KEY `idx_creator_user_id` (`creator_user_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_tenants_creator_user_id` FOREIGN KEY (`creator_user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户信息表';

-- 租户用户关联表
CREATE TABLE IF NOT EXISTS `tenant_users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `role` VARCHAR(32) NOT NULL COMMENT '角色：admin(租户管理员), user(普通用户)',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1-正常 2-禁用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_user` (`tenant_id`, `user_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_tenant_users_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_tenant_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户用户关联表';

-- 租户配置表
CREATE TABLE IF NOT EXISTS `tenant_configs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
  `config_key` VARCHAR(64) NOT NULL COMMENT '配置键',
  `config_value` TEXT DEFAULT NULL COMMENT '配置值',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '配置描述',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_key` (`tenant_id`, `config_key`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_tenant_configs_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户配置表';

-- 租户计划表
CREATE TABLE IF NOT EXISTS `tenant_plans` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '计划名称',
  `description` TEXT DEFAULT NULL COMMENT '计划描述',
  `price` DECIMAL(10,2) NOT NULL COMMENT '价格',
  `duration_days` INT NOT NULL COMMENT '有效期（天）',
  `features` TEXT DEFAULT NULL COMMENT '功能列表（JSON）',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1-启用 2-禁用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户计划表';

-- 租户订阅表
CREATE TABLE IF NOT EXISTS `tenant_subscriptions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
  `plan_id` BIGINT NOT NULL COMMENT '计划ID',
  `start_date` TIMESTAMP NOT NULL COMMENT '开始时间',
  `end_date` TIMESTAMP NOT NULL COMMENT '结束时间',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1-有效 2-过期 3-取消',
  `payment_id` VARCHAR(64) DEFAULT NULL COMMENT '支付ID',
  `auto_renew` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否自动续费',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_plan_id` (`plan_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_tenant_subscriptions_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_tenant_subscriptions_plan_id` FOREIGN KEY (`plan_id`) REFERENCES `tenant_plans` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户订阅表';

-- 租户统计表
CREATE TABLE IF NOT EXISTS `tenant_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
  `user_count` INT NOT NULL DEFAULT 0 COMMENT '用户数',
  `product_count` INT NOT NULL DEFAULT 0 COMMENT '产品数',
  `order_count` INT NOT NULL DEFAULT 0 COMMENT '订单数',
  `post_count` INT NOT NULL DEFAULT 0 COMMENT '文章数',
  `view_count` INT NOT NULL DEFAULT 0 COMMENT '浏览数',
  `date` DATE NOT NULL COMMENT '统计日期',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_date` (`tenant_id`, `date`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_tenant_statistics_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户统计表';
