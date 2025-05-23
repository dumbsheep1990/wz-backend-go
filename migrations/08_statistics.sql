-- 统计模块数据表结构

-- 统计数据表
CREATE TABLE IF NOT EXISTS `statistics_data` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `date` DATE NOT NULL COMMENT '统计日期',
  `type` VARCHAR(32) NOT NULL COMMENT '统计类型(pv, uv, content_view, user_register等)',
  `value` BIGINT NOT NULL DEFAULT 0 COMMENT '统计值',
  `item_id` BIGINT DEFAULT NULL COMMENT '关联项目ID',
  `item_type` VARCHAR(32) DEFAULT NULL COMMENT '关联项目类型',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_date` (`date`),
  KEY `idx_type` (`type`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_item` (`item_type`, `item_id`),
  KEY `idx_date_type_tenant` (`date`, `type`, `tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='统计数据表';

-- 日访问统计表
CREATE TABLE IF NOT EXISTS `daily_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `date` DATE NOT NULL COMMENT '统计日期',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `pv` BIGINT NOT NULL DEFAULT 0 COMMENT '页面浏览量',
  `uv` BIGINT NOT NULL DEFAULT 0 COMMENT '独立访客数',
  `ip_count` BIGINT NOT NULL DEFAULT 0 COMMENT '独立IP数',
  `new_user_count` BIGINT NOT NULL DEFAULT 0 COMMENT '新用户数',
  `new_content_count` BIGINT NOT NULL DEFAULT 0 COMMENT '新内容数',
  `avg_visit_time` INT DEFAULT 0 COMMENT '平均访问时长(秒)',
  `bounce_rate` DECIMAL(5,2) DEFAULT 0 COMMENT '跳出率(%)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date_tenant` (`date`, `tenant_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日访问统计表';

-- 内容排行表
CREATE TABLE IF NOT EXISTS `content_rankings` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `content_id` BIGINT NOT NULL COMMENT '内容ID',
  `title` VARCHAR(255) NOT NULL COMMENT '标题',
  `cover` VARCHAR(255) DEFAULT NULL COMMENT '封面',
  `type` VARCHAR(32) NOT NULL COMMENT '内容类型',
  `view_count` BIGINT NOT NULL DEFAULT 0 COMMENT '浏览量',
  `like_count` BIGINT NOT NULL DEFAULT 0 COMMENT '点赞量',
  `comment_count` BIGINT NOT NULL DEFAULT 0 COMMENT '评论量',
  `url` VARCHAR(255) DEFAULT NULL COMMENT '链接',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `date` DATE NOT NULL COMMENT '统计日期',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date_content` (`date`, `content_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_date` (`date`),
  KEY `idx_view_count` (`view_count`),
  KEY `idx_like_count` (`like_count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容排行表';

-- 用户行为统计表
CREATE TABLE IF NOT EXISTS `user_behavior_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `date` DATE NOT NULL COMMENT '统计日期',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `pv` BIGINT NOT NULL DEFAULT 0 COMMENT '页面浏览量',
  `session_count` INT NOT NULL DEFAULT 0 COMMENT '会话数',
  `total_time` INT NOT NULL DEFAULT 0 COMMENT '总停留时间(秒)',
  `action_count` INT NOT NULL DEFAULT 0 COMMENT '操作次数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date_user` (`date`, `user_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户行为统计表';

-- 设备统计表
CREATE TABLE IF NOT EXISTS `device_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `date` DATE NOT NULL COMMENT '统计日期',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `device_type` VARCHAR(32) NOT NULL COMMENT '设备类型：desktop, mobile, tablet, other',
  `browser` VARCHAR(32) DEFAULT NULL COMMENT '浏览器：chrome, firefox, safari, edge, other',
  `os` VARCHAR(32) DEFAULT NULL COMMENT '操作系统：windows, macos, ios, android, linux, other',
  `count` BIGINT NOT NULL DEFAULT 0 COMMENT '访问次数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date_tenant_device` (`date`, `tenant_id`, `device_type`, `browser`, `os`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_date` (`date`),
  KEY `idx_device_type` (`device_type`),
  KEY `idx_browser` (`browser`),
  KEY `idx_os` (`os`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设备统计表';

-- 来源统计表
CREATE TABLE IF NOT EXISTS `referrer_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `date` DATE NOT NULL COMMENT '统计日期',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `referrer_type` VARCHAR(32) NOT NULL COMMENT '来源类型：direct, search, social, referral, other',
  `referrer_domain` VARCHAR(255) DEFAULT NULL COMMENT '来源域名',
  `count` BIGINT NOT NULL DEFAULT 0 COMMENT '访问次数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date_tenant_referrer` (`date`, `tenant_id`, `referrer_type`, `referrer_domain`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_date` (`date`),
  KEY `idx_referrer_type` (`referrer_type`),
  KEY `idx_referrer_domain` (`referrer_domain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='来源统计表';

-- 热门搜索词表
CREATE TABLE IF NOT EXISTS `search_keywords` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `date` DATE NOT NULL COMMENT '统计日期',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `keyword` VARCHAR(255) NOT NULL COMMENT '搜索关键词',
  `count` BIGINT NOT NULL DEFAULT 0 COMMENT '搜索次数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date_tenant_keyword` (`date`, `tenant_id`, `keyword`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_date` (`date`),
  KEY `idx_count` (`count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='热门搜索词表';
