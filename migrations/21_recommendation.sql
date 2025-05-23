-- 推荐系统数据表结构

-- 推荐策略表
CREATE TABLE IF NOT EXISTS `recommendation_strategies` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '策略名称',
  `code` VARCHAR(64) NOT NULL COMMENT '策略编码',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '策略描述',
  `strategy_type` VARCHAR(32) NOT NULL COMMENT '策略类型：content-内容推荐 product-产品推荐 user-用户推荐 mixed-混合推荐',
  `algorithm` VARCHAR(64) NOT NULL COMMENT '算法：collaborative-协同过滤 content-内容过滤 rule-规则过滤 hybrid-混合',
  `parameters` JSON DEFAULT NULL COMMENT '策略参数(JSON格式)',
  `filter_rules` JSON DEFAULT NULL COMMENT '过滤规则(JSON格式)',
  `sort_rules` JSON DEFAULT NULL COMMENT '排序规则(JSON格式)',
  `weight` INT DEFAULT 100 COMMENT '权重',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_strategy_type` (`strategy_type`),
  KEY `idx_algorithm` (`algorithm`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推荐策略表';

-- 推荐场景表
CREATE TABLE IF NOT EXISTS `recommendation_scenarios` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '场景名称',
  `code` VARCHAR(64) NOT NULL COMMENT '场景编码',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '场景描述',
  `page_position` VARCHAR(128) DEFAULT NULL COMMENT '页面位置',
  `item_count` INT DEFAULT 10 COMMENT '推荐数量',
  `strategy_id` BIGINT NOT NULL COMMENT '策略ID',
  `fallback_strategy_id` BIGINT DEFAULT NULL COMMENT '备选策略ID',
  `refresh_interval` INT DEFAULT 3600 COMMENT '刷新间隔(秒)',
  `cache_ttl` INT DEFAULT 600 COMMENT '缓存有效期(秒)',
  `ab_test_enabled` TINYINT(1) DEFAULT 0 COMMENT '是否启用A/B测试：0-否 1-是',
  `ab_test_config` JSON DEFAULT NULL COMMENT 'A/B测试配置(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_strategy_id` (`strategy_id`),
  KEY `idx_fallback_strategy_id` (`fallback_strategy_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_recommendation_scenarios_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `recommendation_strategies` (`id`),
  CONSTRAINT `fk_recommendation_scenarios_fallback_strategy_id` FOREIGN KEY (`fallback_strategy_id`) REFERENCES `recommendation_strategies` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推荐场景表';

-- 推荐项目表
CREATE TABLE IF NOT EXISTS `recommendation_items` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `item_id` BIGINT NOT NULL COMMENT '项目ID',
  `item_type` VARCHAR(32) NOT NULL COMMENT '项目类型：content-内容 product-产品 user-用户',
  `features` JSON DEFAULT NULL COMMENT '特征向量(JSON格式)',
  `categories` VARCHAR(255) DEFAULT NULL COMMENT '分类，多个用逗号分隔',
  `tags` VARCHAR(255) DEFAULT NULL COMMENT '标签，多个用逗号分隔',
  `popularity` DECIMAL(10,4) DEFAULT 0 COMMENT '流行度',
  `quality_score` DECIMAL(5,4) DEFAULT 0 COMMENT '质量分',
  `relevance_factors` JSON DEFAULT NULL COMMENT '相关性因子(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `last_updated` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_item_id_type` (`item_id`, `item_type`),
  KEY `idx_item_type` (`item_type`),
  KEY `idx_categories` (`categories`),
  KEY `idx_tags` (`tags`),
  KEY `idx_popularity` (`popularity`),
  KEY `idx_quality_score` (`quality_score`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推荐项目表';

-- 用户兴趣表
CREATE TABLE IF NOT EXISTS `user_interests` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `interest_type` VARCHAR(32) NOT NULL COMMENT '兴趣类型：category-分类 tag-标签 topic-话题 entity-实体',
  `interest_key` VARCHAR(128) NOT NULL COMMENT '兴趣键',
  `interest_value` VARCHAR(128) DEFAULT NULL COMMENT '兴趣值',
  `weight` DECIMAL(5,4) DEFAULT 1 COMMENT '权重(0-1)',
  `source` VARCHAR(32) DEFAULT 'system' COMMENT '来源：system-系统 user-用户 algorithm-算法',
  `last_updated` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_interest` (`user_id`, `interest_type`, `interest_key`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_interest_type` (`interest_type`),
  KEY `idx_interest_key` (`interest_key`),
  KEY `idx_weight` (`weight`),
  KEY `idx_source` (`source`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户兴趣表';

-- 用户特征表
CREATE TABLE IF NOT EXISTS `user_profiles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `feature_type` VARCHAR(32) NOT NULL COMMENT '特征类型：demographic-人口统计 behavior-行为 preference-偏好',
  `feature_key` VARCHAR(128) NOT NULL COMMENT '特征键',
  `feature_value` TEXT NOT NULL COMMENT '特征值',
  `confidence` DECIMAL(5,4) DEFAULT 1 COMMENT '置信度(0-1)',
  `last_updated` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_feature` (`user_id`, `feature_type`, `feature_key`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_feature_type` (`feature_type`),
  KEY `idx_feature_key` (`feature_key`),
  KEY `idx_confidence` (`confidence`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户特征表';

-- 相似项目表
CREATE TABLE IF NOT EXISTS `similar_items` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `item_id` BIGINT NOT NULL COMMENT '项目ID',
  `item_type` VARCHAR(32) NOT NULL COMMENT '项目类型',
  `similar_item_id` BIGINT NOT NULL COMMENT '相似项目ID',
  `similar_item_type` VARCHAR(32) NOT NULL COMMENT '相似项目类型',
  `similarity_score` DECIMAL(5,4) NOT NULL COMMENT '相似度(0-1)',
  `similarity_type` VARCHAR(32) DEFAULT 'content' COMMENT '相似度类型：content-内容相似 collaborative-协同相似 hybrid-混合',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_item_similar_item` (`item_id`, `item_type`, `similar_item_id`, `similar_item_type`),
  KEY `idx_item_id_type` (`item_id`, `item_type`),
  KEY `idx_similar_item_id_type` (`similar_item_id`, `similar_item_type`),
  KEY `idx_similarity_score` (`similarity_score`),
  KEY `idx_similarity_type` (`similarity_type`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='相似项目表';

-- 推荐结果表
CREATE TABLE IF NOT EXISTS `recommendation_results` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `scenario_id` BIGINT NOT NULL COMMENT '场景ID',
  `strategy_id` BIGINT NOT NULL COMMENT '策略ID',
  `session_id` VARCHAR(64) DEFAULT NULL COMMENT '会话ID',
  `items` JSON NOT NULL COMMENT '推荐项目(JSON格式)',
  `context` JSON DEFAULT NULL COMMENT '上下文信息(JSON格式)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `expire_at` TIMESTAMP NULL COMMENT '过期时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_scenario_id` (`scenario_id`),
  KEY `idx_strategy_id` (`strategy_id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_expire_at` (`expire_at`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_recommendation_results_scenario_id` FOREIGN KEY (`scenario_id`) REFERENCES `recommendation_scenarios` (`id`),
  CONSTRAINT `fk_recommendation_results_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `recommendation_strategies` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推荐结果表';

-- 推荐反馈表
CREATE TABLE IF NOT EXISTS `recommendation_feedbacks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `item_id` BIGINT NOT NULL COMMENT '项目ID',
  `item_type` VARCHAR(32) NOT NULL COMMENT '项目类型',
  `scenario_id` BIGINT NOT NULL COMMENT '场景ID',
  `session_id` VARCHAR(64) DEFAULT NULL COMMENT '会话ID',
  `action` VARCHAR(32) NOT NULL COMMENT '操作：view-浏览 click-点击 like-点赞 favorite-收藏 share-分享 ignore-忽略 dislike-不喜欢',
  `position` INT DEFAULT NULL COMMENT '位置',
  `duration` INT DEFAULT NULL COMMENT '时长(秒)',
  `feedback_value` DECIMAL(5,4) DEFAULT NULL COMMENT '反馈值(0-1)',
  `properties` JSON DEFAULT NULL COMMENT '属性(JSON格式)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_item_id_type` (`item_id`, `item_type`),
  KEY `idx_scenario_id` (`scenario_id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_action` (`action`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_recommendation_feedbacks_scenario_id` FOREIGN KEY (`scenario_id`) REFERENCES `recommendation_scenarios` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推荐反馈表';

-- 推荐统计表
CREATE TABLE IF NOT EXISTS `recommendation_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `scenario_id` BIGINT NOT NULL COMMENT '场景ID',
  `strategy_id` BIGINT NOT NULL COMMENT '策略ID',
  `date` DATE NOT NULL COMMENT '统计日期',
  `hour` INT DEFAULT NULL COMMENT '小时(0-23)',
  `impression_count` INT DEFAULT 0 COMMENT '展示次数',
  `click_count` INT DEFAULT 0 COMMENT '点击次数',
  `ctr` DECIMAL(5,4) DEFAULT 0 COMMENT '点击率',
  `like_count` INT DEFAULT 0 COMMENT '点赞次数',
  `favorite_count` INT DEFAULT 0 COMMENT '收藏次数',
  `share_count` INT DEFAULT 0 COMMENT '分享次数',
  `user_count` INT DEFAULT 0 COMMENT '用户数',
  `engagement_rate` DECIMAL(5,4) DEFAULT 0 COMMENT '参与率',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_scenario_strategy_date_hour` (`scenario_id`, `strategy_id`, `date`, `hour`),
  KEY `idx_scenario_id` (`scenario_id`),
  KEY `idx_strategy_id` (`strategy_id`),
  KEY `idx_date` (`date`),
  KEY `idx_hour` (`hour`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_recommendation_statistics_scenario_id` FOREIGN KEY (`scenario_id`) REFERENCES `recommendation_scenarios` (`id`),
  CONSTRAINT `fk_recommendation_statistics_strategy_id` FOREIGN KEY (`strategy_id`) REFERENCES `recommendation_strategies` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推荐统计表';

-- A/B测试表
CREATE TABLE IF NOT EXISTS `ab_tests` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '测试名称',
  `code` VARCHAR(64) NOT NULL COMMENT '测试编码',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '测试描述',
  `start_time` TIMESTAMP NULL COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `status` VARCHAR(32) DEFAULT 'draft' COMMENT '状态：draft-草稿 running-运行中 paused-已暂停 completed-已完成 terminated-已终止',
  `traffic_ratio` JSON NOT NULL COMMENT '流量分配(JSON格式)',
  `goal_metrics` JSON NOT NULL COMMENT '目标指标(JSON格式)',
  `variants` JSON NOT NULL COMMENT '变体配置(JSON格式)',
  `audience` JSON DEFAULT NULL COMMENT '受众配置(JSON格式)',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建者ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_end_time` (`end_time`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='A/B测试表';

-- A/B测试用户分组表
CREATE TABLE IF NOT EXISTS `ab_test_users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `test_id` BIGINT NOT NULL COMMENT '测试ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `variant` VARCHAR(64) NOT NULL COMMENT '变体',
  `assignment_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_test_user` (`test_id`, `user_id`),
  KEY `idx_test_id` (`test_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_variant` (`variant`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_ab_test_users_test_id` FOREIGN KEY (`test_id`) REFERENCES `ab_tests` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='A/B测试用户分组表';

-- 初始化基础推荐策略
INSERT INTO `recommendation_strategies` (`name`, `code`, `description`, `strategy_type`, `algorithm`, `weight`, `status`) VALUES
('热门内容推荐', 'popular-content', '基于热度和时间的内容推荐', 'content', 'rule', 100, 1),
('个性化内容推荐', 'personalized-content', '基于用户兴趣的个性化内容推荐', 'content', 'hybrid', 100, 1),
('相似内容推荐', 'similar-content', '基于当前浏览内容的相似推荐', 'content', 'content', 100, 1),
('热销产品推荐', 'hot-products', '热销产品推荐', 'product', 'rule', 100, 1),
('个性化产品推荐', 'personalized-products', '基于用户兴趣的个性化产品推荐', 'product', 'collaborative', 100, 1),
('用户推荐', 'user-recommendation', '推荐可能感兴趣的用户', 'user', 'hybrid', 100, 1);

-- 初始化基础推荐场景
INSERT INTO `recommendation_scenarios` (`name`, `code`, `description`, `page_position`, `item_count`, `strategy_id`, `status`) VALUES
('首页推荐', 'home-recommendation', '首页内容推荐', 'home-main', 10, 1, 1),
('相关内容', 'related-content', '内容详情页底部相关推荐', 'content-detail-bottom', 6, 3, 1),
('猜你喜欢', 'for-you', '个性化推荐内容', 'user-center', 20, 2, 1),
('热门商品', 'hot-products', '商城页面热门商品', 'shop-main', 8, 4, 1),
('相关商品', 'related-products', '商品详情页相关商品', 'product-detail-right', 4, 5, 1);
