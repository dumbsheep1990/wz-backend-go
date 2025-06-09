-- AI服务数据库表结构

-- 推荐记录表
CREATE TABLE IF NOT EXISTS `recommendations` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_id` bigint NOT NULL COMMENT '用户ID',
    `content_id` bigint NOT NULL COMMENT '内容ID',
    `content_type` varchar(20) NOT NULL COMMENT '内容类型：post/article',
    `score` float NOT NULL COMMENT '推荐分数',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user` (`user_id`),
    KEY `idx_content` (`content_type`, `content_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推荐记录表';

-- 内容审核记录表
CREATE TABLE IF NOT EXISTS `content_reviews` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `content_id` bigint NOT NULL COMMENT '内容ID',
    `content_type` varchar(20) NOT NULL COMMENT '内容类型：text/image',
    `passed` tinyint NOT NULL COMMENT '是否通过',
    `reason` varchar(255) DEFAULT NULL COMMENT '原因',
    `labels` json DEFAULT NULL COMMENT '标签',
    `confidence` float NOT NULL COMMENT '置信度',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_content` (`content_type`, `content_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容审核记录表';

-- 客服对话记录表
CREATE TABLE IF NOT EXISTS `chat_records` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_id` bigint NOT NULL COMMENT '用户ID',
    `message` text NOT NULL COMMENT '消息内容',
    `response` text NOT NULL COMMENT '回复内容',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客服对话记录表'; 