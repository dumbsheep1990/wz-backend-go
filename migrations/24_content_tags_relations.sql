-- 内容标签关联表结构

-- 标签表
CREATE TABLE IF NOT EXISTS `tags` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '标签名称',
  `slug` VARCHAR(64) NOT NULL COMMENT '标签别名',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '标签描述',
  `icon` VARCHAR(255) DEFAULT NULL COMMENT '标签图标',
  `color` VARCHAR(32) DEFAULT NULL COMMENT '标签颜色',
  `category` VARCHAR(32) DEFAULT NULL COMMENT '标签分类',
  `count` INT DEFAULT 0 COMMENT '使用次数',
  `is_hot` TINYINT(1) DEFAULT 0 COMMENT '是否热门',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `sort_order` INT DEFAULT 0 COMMENT '排序顺序',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name_tenant` (`name`, `tenant_id`),
  UNIQUE KEY `idx_slug_tenant` (`slug`, `tenant_id`),
  KEY `idx_category` (`category`),
  KEY `idx_is_hot` (`is_hot`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签表';

-- 内容标签关联表
CREATE TABLE IF NOT EXISTS `content_tags` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `content_id` BIGINT NOT NULL COMMENT '内容ID',
  `tag_id` BIGINT NOT NULL COMMENT '标签ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_content_tag` (`content_id`, `tag_id`),
  KEY `idx_content_id` (`content_id`),
  KEY `idx_tag_id` (`tag_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_content_tags_content_id` FOREIGN KEY (`content_id`) REFERENCES `contents` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_content_tags_tag_id` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容标签关联表';

-- 产品标签关联表
CREATE TABLE IF NOT EXISTS `product_tags` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL COMMENT '产品ID',
  `tag_id` BIGINT NOT NULL COMMENT '标签ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_product_tag` (`product_id`, `tag_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_tag_id` (`tag_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_product_tags_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_product_tags_tag_id` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品标签关联表';

-- 用户标签关联表（用户兴趣标签）
CREATE TABLE IF NOT EXISTS `user_tags` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `tag_id` BIGINT NOT NULL COMMENT '标签ID',
  `weight` DECIMAL(5,2) DEFAULT 1.00 COMMENT '权重',
  `source` VARCHAR(32) DEFAULT 'user' COMMENT '来源：user-用户设置 system-系统推荐 analysis-行为分析',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_tag` (`user_id`, `tag_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_tag_id` (`tag_id`),
  KEY `idx_source` (`source`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_user_tags_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_tags_tag_id` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户标签关联表';
