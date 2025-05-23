-- 用户交互服务数据表结构

-- 点赞表
CREATE TABLE IF NOT EXISTS `likes` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `target_id` BIGINT NOT NULL COMMENT '目标ID',
  `target_type` VARCHAR(32) NOT NULL COMMENT '目标类型：content-内容 comment-评论 user-用户 product-产品',
  `status` TINYINT(1) DEFAULT 1 COMMENT '状态：0-取消 1-有效',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_target` (`user_id`, `target_id`, `target_type`),
  KEY `idx_target_id_type` (`target_id`, `target_type`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞表';

-- 收藏表
CREATE TABLE IF NOT EXISTS `favorites` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `target_id` BIGINT NOT NULL COMMENT '目标ID',
  `target_type` VARCHAR(32) NOT NULL COMMENT '目标类型：content-内容 product-产品 activity-活动',
  `collection_id` BIGINT DEFAULT NULL COMMENT '收藏夹ID',
  `remark` VARCHAR(255) DEFAULT NULL COMMENT '备注',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_target` (`user_id`, `target_id`, `target_type`),
  KEY `idx_target_id_type` (`target_id`, `target_type`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_collection_id` (`collection_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏表';

-- 收藏夹表
CREATE TABLE IF NOT EXISTS `favorite_collections` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `name` VARCHAR(64) NOT NULL COMMENT '收藏夹名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '收藏夹描述',
  `cover_image` VARCHAR(255) DEFAULT NULL COMMENT '封面图片',
  `is_default` TINYINT(1) DEFAULT 0 COMMENT '是否默认：0-否 1-是',
  `privacy` TINYINT(1) DEFAULT 0 COMMENT '隐私设置：0-公开 1-私密',
  `item_count` INT DEFAULT 0 COMMENT '收藏数量',
  `sort_order` INT DEFAULT 0 COMMENT '排序顺序',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_is_default` (`is_default`),
  KEY `idx_privacy` (`privacy`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏夹表';

-- 关注表
CREATE TABLE IF NOT EXISTS `follows` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `follow_id` BIGINT NOT NULL COMMENT '关注对象ID',
  `follow_type` VARCHAR(32) NOT NULL COMMENT '关注类型：user-用户 topic-话题 tag-标签',
  `remark` VARCHAR(64) DEFAULT NULL COMMENT '备注',
  `status` TINYINT(1) DEFAULT 1 COMMENT '状态：0-取消 1-有效',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_follow` (`user_id`, `follow_id`, `follow_type`),
  KEY `idx_follow_id_type` (`follow_id`, `follow_type`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='关注表';

-- 分享表
CREATE TABLE IF NOT EXISTS `shares` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `target_id` BIGINT NOT NULL COMMENT '目标ID',
  `target_type` VARCHAR(32) NOT NULL COMMENT '目标类型：content-内容 product-产品 activity-活动',
  `share_type` VARCHAR(32) NOT NULL COMMENT '分享类型：wechat-微信 weibo-微博 qq-QQ link-链接',
  `share_content` TEXT DEFAULT NULL COMMENT '分享内容',
  `share_url` VARCHAR(512) DEFAULT NULL COMMENT '分享链接',
  `share_code` VARCHAR(64) DEFAULT NULL COMMENT '分享码',
  `ip` VARCHAR(64) DEFAULT NULL COMMENT 'IP地址',
  `device` VARCHAR(128) DEFAULT NULL COMMENT '设备信息',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_target_id_type` (`target_id`, `target_type`),
  KEY `idx_share_type` (`share_type`),
  KEY `idx_share_code` (`share_code`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分享表';

-- 评分表
CREATE TABLE IF NOT EXISTS `ratings` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `target_id` BIGINT NOT NULL COMMENT '目标ID',
  `target_type` VARCHAR(32) NOT NULL COMMENT '目标类型：content-内容 product-产品 service-服务',
  `score` DECIMAL(3,1) NOT NULL COMMENT '评分(0-5)',
  `comment` VARCHAR(512) DEFAULT NULL COMMENT '评价内容',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_target` (`user_id`, `target_id`, `target_type`),
  KEY `idx_target_id_type` (`target_id`, `target_type`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_score` (`score`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评分表';

-- 用户行为表
CREATE TABLE IF NOT EXISTS `user_behaviors` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `session_id` VARCHAR(64) DEFAULT NULL COMMENT '会话ID',
  `behavior_type` VARCHAR(32) NOT NULL COMMENT '行为类型：view-浏览 click-点击 search-搜索 stay-停留',
  `target_id` BIGINT DEFAULT NULL COMMENT '目标ID',
  `target_type` VARCHAR(32) DEFAULT NULL COMMENT '目标类型',
  `properties` JSON DEFAULT NULL COMMENT '行为属性(JSON格式)',
  `device_type` VARCHAR(32) DEFAULT NULL COMMENT '设备类型：pc, mobile, tablet, app',
  `device_id` VARCHAR(128) DEFAULT NULL COMMENT '设备ID',
  `ip` VARCHAR(64) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` VARCHAR(512) DEFAULT NULL COMMENT '用户代理',
  `referer` VARCHAR(512) DEFAULT NULL COMMENT '来源URL',
  `url` VARCHAR(512) DEFAULT NULL COMMENT '当前URL',
  `duration` INT DEFAULT NULL COMMENT '持续时间(秒)',
  `behavior_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '行为时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_behavior_type` (`behavior_type`),
  KEY `idx_target_id_type` (`target_id`, `target_type`),
  KEY `idx_behavior_time` (`behavior_time`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户行为表';

-- 用户反馈表
CREATE TABLE IF NOT EXISTS `user_feedbacks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `feedback_type` VARCHAR(32) NOT NULL COMMENT '反馈类型：suggestion-建议 bug-缺陷 complaint-投诉',
  `content` TEXT NOT NULL COMMENT '反馈内容',
  `contact` VARCHAR(128) DEFAULT NULL COMMENT '联系方式',
  `attachments` JSON DEFAULT NULL COMMENT '附件(JSON格式)',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-待处理 processing-处理中 completed-已完成 rejected-已拒绝',
  `priority` INT DEFAULT 0 COMMENT '优先级：0-低 1-中 2-高',
  `handler_id` BIGINT DEFAULT NULL COMMENT '处理人ID',
  `handle_time` TIMESTAMP NULL COMMENT '处理时间',
  `handle_result` VARCHAR(512) DEFAULT NULL COMMENT '处理结果',
  `reply_content` TEXT DEFAULT NULL COMMENT '回复内容',
  `is_public` TINYINT(1) DEFAULT 0 COMMENT '是否公开：0-私密 1-公开',
  `ip` VARCHAR(64) DEFAULT NULL COMMENT 'IP地址',
  `device` VARCHAR(128) DEFAULT NULL COMMENT '设备信息',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_feedback_type` (`feedback_type`),
  KEY `idx_status` (`status`),
  KEY `idx_priority` (`priority`),
  KEY `idx_is_public` (`is_public`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户反馈表';

-- 互动统计表
CREATE TABLE IF NOT EXISTS `interaction_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `target_id` BIGINT NOT NULL COMMENT '目标ID',
  `target_type` VARCHAR(32) NOT NULL COMMENT '目标类型',
  `like_count` INT DEFAULT 0 COMMENT '点赞数',
  `favorite_count` INT DEFAULT 0 COMMENT '收藏数',
  `share_count` INT DEFAULT 0 COMMENT '分享数',
  `view_count` INT DEFAULT 0 COMMENT '浏览数',
  `comment_count` INT DEFAULT 0 COMMENT '评论数',
  `rating_count` INT DEFAULT 0 COMMENT '评分数',
  `rating_avg` DECIMAL(3,2) DEFAULT 0 COMMENT '平均评分',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_target_id_type` (`target_id`, `target_type`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='互动统计表';

-- 用户标签表
CREATE TABLE IF NOT EXISTS `user_tags` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `tag_key` VARCHAR(64) NOT NULL COMMENT '标签键',
  `tag_value` VARCHAR(255) NOT NULL COMMENT '标签值',
  `tag_type` VARCHAR(32) DEFAULT 'interest' COMMENT '标签类型：interest-兴趣 demographic-人口统计 behavior-行为',
  `source` VARCHAR(32) DEFAULT 'system' COMMENT '来源：system-系统 user-用户 algorithm-算法',
  `confidence` DECIMAL(5,4) DEFAULT 1 COMMENT '置信度(0-1)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_tag` (`user_id`, `tag_key`, `tag_value`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_tag_key_value` (`tag_key`, `tag_value`),
  KEY `idx_tag_type` (`tag_type`),
  KEY `idx_source` (`source`),
  KEY `idx_confidence` (`confidence`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户标签表';

-- 话题表
CREATE TABLE IF NOT EXISTS `topics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '话题名称',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '话题描述',
  `icon` VARCHAR(255) DEFAULT NULL COMMENT '话题图标',
  `cover_image` VARCHAR(255) DEFAULT NULL COMMENT '封面图片',
  `parent_id` BIGINT DEFAULT 0 COMMENT '父话题ID',
  `level` INT DEFAULT 1 COMMENT '层级',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `sort_order` INT DEFAULT 0 COMMENT '排序顺序',
  `follower_count` INT DEFAULT 0 COMMENT '关注人数',
  `content_count` INT DEFAULT 0 COMMENT '内容数量',
  `is_recommended` TINYINT(1) DEFAULT 0 COMMENT '是否推荐：0-否 1-是',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建者ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name_tenant` (`name`, `tenant_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_recommended` (`is_recommended`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='话题表';

-- 内容话题关联表
CREATE TABLE IF NOT EXISTS `content_topics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `content_id` BIGINT NOT NULL COMMENT '内容ID',
  `topic_id` BIGINT NOT NULL COMMENT '话题ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_content_topic` (`content_id`, `topic_id`),
  KEY `idx_content_id` (`content_id`),
  KEY `idx_topic_id` (`topic_id`),
  CONSTRAINT `fk_content_topics_topic_id` FOREIGN KEY (`topic_id`) REFERENCES `topics` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容话题关联表';

-- 初始化默认收藏夹
INSERT INTO `favorite_collections` (`user_id`, `name`, `description`, `is_default`, `privacy`) VALUES
(1, '默认收藏夹', '系统默认收藏夹', 1, 0);

-- 初始化热门话题
INSERT INTO `topics` (`name`, `description`, `status`, `is_recommended`, `follower_count`) VALUES
('科技', '关于科技领域的讨论', 1, 1, 0),
('旅行', '分享旅行经验和美景', 1, 1, 0),
('美食', '探索美食世界', 1, 1, 0),
('健康', '健康生活方式讨论', 1, 1, 0),
('时尚', '时尚潮流话题', 1, 1, 0);
