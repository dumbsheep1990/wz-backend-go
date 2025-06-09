-- 统计服务数据库表结构

-- 用户行为统计表
CREATE TABLE IF NOT EXISTS `user_behaviors` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_id` bigint NOT NULL COMMENT '用户ID',
    `type` varchar(20) NOT NULL COMMENT '行为类型：view/like/comment',
    `target_type` varchar(20) NOT NULL COMMENT '目标类型：post/article',
    `target_id` bigint NOT NULL COMMENT '目标ID',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user` (`user_id`),
    KEY `idx_target` (`target_type`, `target_id`),
    KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户行为统计表';

-- 内容流行度统计表
CREATE TABLE IF NOT EXISTS `content_popularity` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `content_id` bigint NOT NULL COMMENT '内容ID',
    `content_type` varchar(20) NOT NULL COMMENT '内容类型：post/article',
    `views` bigint NOT NULL DEFAULT '0' COMMENT '浏览量',
    `likes` bigint NOT NULL DEFAULT '0' COMMENT '点赞数',
    `comments` bigint NOT NULL DEFAULT '0' COMMENT '评论数',
    `shares` bigint NOT NULL DEFAULT '0' COMMENT '分享数',
    `score` float NOT NULL DEFAULT '0' COMMENT '热度分数',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_content` (`content_type`, `content_id`),
    KEY `idx_score` (`score`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容流行度统计表';

-- 用户画像表
CREATE TABLE IF NOT EXISTS `user_profiles` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_id` bigint NOT NULL COMMENT '用户ID',
    `interests` json DEFAULT NULL COMMENT '兴趣标签',
    `active_time` varchar(20) DEFAULT NULL COMMENT '活跃时间段',
    `content_types` json DEFAULT NULL COMMENT '常浏览的内容类型',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户画像表'; 