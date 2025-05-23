-- 用户会员等级与商城功能关系表结构

-- 会员等级表
CREATE TABLE IF NOT EXISTS `user_membership_levels` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '等级名称',
  `code` VARCHAR(32) NOT NULL COMMENT '等级编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '等级描述',
  `icon` VARCHAR(255) DEFAULT NULL COMMENT '等级图标',
  `level` INT NOT NULL DEFAULT 1 COMMENT '等级值',
  `min_points` INT NOT NULL DEFAULT 0 COMMENT '最小积分',
  `discount_rate` DECIMAL(5,2) DEFAULT 1.00 COMMENT '折扣率',
  `benefits` JSON DEFAULT NULL COMMENT '权益列表(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_level` (`level`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员等级表';

-- 用户会员等级关系表
CREATE TABLE IF NOT EXISTS `user_memberships` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `level_id` BIGINT NOT NULL COMMENT '等级ID',
  `start_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `is_permanent` TINYINT(1) DEFAULT 0 COMMENT '是否永久：0-否 1-是',
  `source` VARCHAR(32) DEFAULT 'system' COMMENT '来源：system-系统 purchase-购买 promotion-促销',
  `status` INT DEFAULT 1 COMMENT '状态：0-无效 1-有效',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_level_id` (`level_id`),
  KEY `idx_status` (`status`),
  KEY `idx_end_time` (`end_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_user_memberships_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_memberships_level_id` FOREIGN KEY (`level_id`) REFERENCES `user_membership_levels` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户会员等级关系表';

-- 会员权益表
CREATE TABLE IF NOT EXISTS `membership_benefits` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '权益名称',
  `code` VARCHAR(32) NOT NULL COMMENT '权益编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '权益描述',
  `icon` VARCHAR(255) DEFAULT NULL COMMENT '权益图标',
  `benefit_type` VARCHAR(32) NOT NULL COMMENT '权益类型：discount-折扣 coupon-优惠券 service-服务 access-访问权限',
  `value` VARCHAR(128) DEFAULT NULL COMMENT '权益值',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_benefit_type` (`benefit_type`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员权益表';

-- 会员等级权益关系表
CREATE TABLE IF NOT EXISTS `membership_level_benefits` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `level_id` BIGINT NOT NULL COMMENT '等级ID',
  `benefit_id` BIGINT NOT NULL COMMENT '权益ID',
  `value` VARCHAR(128) DEFAULT NULL COMMENT '权益值（可覆盖权益表中的默认值）',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_level_benefit` (`level_id`, `benefit_id`),
  KEY `idx_level_id` (`level_id`),
  KEY `idx_benefit_id` (`benefit_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_membership_level_benefits_level_id` FOREIGN KEY (`level_id`) REFERENCES `user_membership_levels` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_membership_level_benefits_benefit_id` FOREIGN KEY (`benefit_id`) REFERENCES `membership_benefits` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员等级权益关系表';

-- 初始化会员等级数据
INSERT INTO `user_membership_levels` (`name`, `code`, `description`, `level`, `min_points`, `discount_rate`, `status`) VALUES
('普通会员', 'regular', '普通会员等级', 1, 0, 1.00, 1),
('银卡会员', 'silver', '银卡会员等级', 2, 1000, 0.95, 1),
('金卡会员', 'gold', '金卡会员等级', 3, 5000, 0.90, 1),
('铂金会员', 'platinum', '铂金会员等级', 4, 10000, 0.85, 1),
('钻石会员', 'diamond', '钻石会员等级', 5, 30000, 0.80, 1);
