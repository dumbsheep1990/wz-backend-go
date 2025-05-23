-- 内容模块数据表结构

-- 内容分类表
CREATE TABLE IF NOT EXISTS `categories` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '分类名称',
  `description` TEXT DEFAULT NULL COMMENT '分类描述',
  `parent_id` BIGINT DEFAULT 0 COMMENT '父分类ID，0表示顶级分类',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1-正常 2-禁用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容分类表';

-- 帖子表
CREATE TABLE IF NOT EXISTS `posts` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(255) NOT NULL COMMENT '帖子标题',
  `content` TEXT NOT NULL COMMENT '帖子内容',
  `user_id` BIGINT NOT NULL COMMENT '发布用户ID',
  `category_id` BIGINT NOT NULL COMMENT '分类ID',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1-正常 2-待审核 3-禁用',
  `view_count` INT NOT NULL DEFAULT 0 COMMENT '浏览次数',
  `like_count` INT NOT NULL DEFAULT 0 COMMENT '点赞次数',
  `comment_count` INT NOT NULL DEFAULT 0 COMMENT '评论次数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_posts_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_posts_category_id` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子表';

-- 评论表
CREATE TABLE IF NOT EXISTS `reviews` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `post_id` BIGINT NOT NULL COMMENT '帖子ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `content` TEXT NOT NULL COMMENT '评论内容',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1-正常 2-待审核 3-禁用',
  `like_count` INT NOT NULL DEFAULT 0 COMMENT '点赞次数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_reviews_post_id` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_reviews_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论表';

-- 内容状态变更日志表
CREATE TABLE IF NOT EXISTS `content_status_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `resource_type` VARCHAR(32) NOT NULL COMMENT '资源类型',
  `resource_id` BIGINT NOT NULL COMMENT '资源ID',
  `status` INT NOT NULL COMMENT '状态值',
  `reason` TEXT DEFAULT NULL COMMENT '原因',
  `operator_id` BIGINT NOT NULL COMMENT '操作人ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_resource` (`resource_type`, `resource_id`),
  KEY `idx_operator_id` (`operator_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_content_status_logs_operator_id` FOREIGN KEY (`operator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容状态变更日志表';

-- 热门内容表
CREATE TABLE IF NOT EXISTS `hot_contents` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `resource_type` VARCHAR(32) NOT NULL COMMENT '资源类型',
  `resource_id` BIGINT NOT NULL COMMENT '资源ID',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `operator_id` BIGINT NOT NULL COMMENT '操作人ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_resource` (`resource_type`, `resource_id`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_operator_id` (`operator_id`),
  CONSTRAINT `fk_hot_contents_operator_id` FOREIGN KEY (`operator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='热门内容表';

-- 点赞表
CREATE TABLE IF NOT EXISTS `likes` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `resource_type` VARCHAR(32) NOT NULL COMMENT '资源类型',
  `resource_id` BIGINT NOT NULL COMMENT '资源ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_resource` (`user_id`, `resource_type`, `resource_id`),
  KEY `idx_resource` (`resource_type`, `resource_id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_likes_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞表';

-- 标签表
CREATE TABLE IF NOT EXISTS `tags` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '标签名称',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签表';

-- 内容标签关联表
CREATE TABLE IF NOT EXISTS `content_tags` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `resource_type` VARCHAR(32) NOT NULL COMMENT '资源类型',
  `resource_id` BIGINT NOT NULL COMMENT '资源ID',
  `tag_id` BIGINT NOT NULL COMMENT '标签ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_resource_tag` (`resource_type`, `resource_id`, `tag_id`),
  KEY `idx_resource` (`resource_type`, `resource_id`),
  KEY `idx_tag_id` (`tag_id`),
  CONSTRAINT `fk_content_tags_tag_id` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容标签关联表';
