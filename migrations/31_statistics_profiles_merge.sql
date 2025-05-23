-- 统计服务和用户画像表结构合并

-- 用户行为统计详细表（对20_interaction.sql中user_behaviors表的补充）
CREATE TABLE IF NOT EXISTS `user_behavior_details` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `behavior_id` BIGINT NOT NULL COMMENT '关联的user_behaviors表ID',
  `session_id` VARCHAR(64) DEFAULT NULL COMMENT '会话ID',
  `action_detail` TEXT DEFAULT NULL COMMENT '行为详细描述',
  `context` JSON DEFAULT NULL COMMENT '行为上下文信息',
  `client_info` JSON DEFAULT NULL COMMENT '客户端信息',
  `geo_location` JSON DEFAULT NULL COMMENT '地理位置信息',
  `referrer` VARCHAR(512) DEFAULT NULL COMMENT '来源URL',
  `time_spent` INT DEFAULT NULL COMMENT '停留时间(秒)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_behavior_id` (`behavior_id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_user_behavior_details_behavior_id` FOREIGN KEY (`behavior_id`) REFERENCES `user_behaviors` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户行为统计详细表';

-- 用户画像标签表（对user_profiles的补充）
CREATE TABLE IF NOT EXISTS `user_profile_tags` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `tag_type` VARCHAR(32) NOT NULL COMMENT '标签类型：interest-兴趣 demography-人口特征 behavior-行为 custom-自定义',
  `tag_key` VARCHAR(64) NOT NULL COMMENT '标签键',
  `tag_value` VARCHAR(255) NOT NULL COMMENT '标签值',
  `confidence` DECIMAL(5,4) DEFAULT 1.0000 COMMENT '置信度',
  `source` VARCHAR(32) DEFAULT 'system' COMMENT '来源：system-系统 manual-手动 analysis-分析',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_tag` (`user_id`, `tag_type`, `tag_key`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_tag_type_key` (`tag_type`, `tag_key`),
  KEY `idx_tag_value` (`tag_value`),
  KEY `idx_confidence` (`confidence`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户画像标签表';

-- 用户活跃度统计表
CREATE TABLE IF NOT EXISTS `user_activity_stats` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `date` DATE NOT NULL COMMENT '统计日期',
  `login_count` INT DEFAULT 0 COMMENT '登录次数',
  `view_count` INT DEFAULT 0 COMMENT '浏览次数',
  `interaction_count` INT DEFAULT 0 COMMENT '交互次数(点赞、评论等)',
  `content_count` INT DEFAULT 0 COMMENT '内容产出数',
  `time_spent` INT DEFAULT 0 COMMENT '总时长(秒)',
  `last_active_time` TIMESTAMP NULL COMMENT '最后活跃时间',
  `active_score` DECIMAL(10,2) DEFAULT 0 COMMENT '活跃度得分',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_date` (`user_id`, `date`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_date` (`date`),
  KEY `idx_active_score` (`active_score`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户活跃度统计表';

-- 内容流行度排行表（对content_popularity表的补充）
CREATE TABLE IF NOT EXISTS `content_popularity_rankings` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `content_id` BIGINT NOT NULL COMMENT '内容ID',
  `content_type` VARCHAR(32) NOT NULL COMMENT '内容类型',
  `ranking_type` VARCHAR(32) NOT NULL COMMENT '排行类型：daily-日榜 weekly-周榜 monthly-月榜 all_time-总榜',
  `ranking_date` DATE NOT NULL COMMENT '排行日期',
  `category` VARCHAR(64) DEFAULT NULL COMMENT '分类',
  `rank` INT NOT NULL COMMENT '排名',
  `score` DECIMAL(10,2) NOT NULL COMMENT '得分',
  `previous_rank` INT DEFAULT NULL COMMENT '前一次排名',
  `rank_change` INT DEFAULT NULL COMMENT '排名变化',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_content_ranking_date` (`content_id`, `content_type`, `ranking_type`, `ranking_date`),
  KEY `idx_content` (`content_id`, `content_type`),
  KEY `idx_ranking_date` (`ranking_type`, `ranking_date`),
  KEY `idx_category` (`category`),
  KEY `idx_rank` (`rank`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容流行度排行表';
