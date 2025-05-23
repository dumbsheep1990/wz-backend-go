-- 链接和积分模块数据表结构

-- 友情链接表
CREATE TABLE IF NOT EXISTS `links` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL COMMENT '链接名称',
  `url` VARCHAR(255) NOT NULL COMMENT '链接URL',
  `logo` VARCHAR(255) DEFAULT NULL COMMENT '链接Logo',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `description` TEXT DEFAULT NULL COMMENT '链接描述',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='友情链接表';

-- 用户消息表（系统通知和私信等）
CREATE TABLE IF NOT EXISTS `user_messages` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `title` VARCHAR(255) NOT NULL COMMENT '消息标题',
  `content` TEXT NOT NULL COMMENT '消息内容',
  `type` INT NOT NULL COMMENT '消息类型：1-系统通知 2-业务通知 3-私信',
  `status` INT NOT NULL DEFAULT 0 COMMENT '状态：0-未读 1-已读',
  `is_important` INT NOT NULL DEFAULT 0 COMMENT '是否重要：0-普通 1-重要',
  `related_id` BIGINT DEFAULT NULL COMMENT '关联ID',
  `related_type` VARCHAR(50) DEFAULT NULL COMMENT '关联类型',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_type` (`type`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_user_messages_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户消息表';

-- 用户积分表（详细记录）
CREATE TABLE IF NOT EXISTS `user_points_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `points` INT NOT NULL COMMENT '积分变动值',
  `total_points` INT NOT NULL COMMENT '总积分（冗余字段）',
  `type` INT NOT NULL COMMENT '类型：1-增加 2-减少',
  `source` VARCHAR(50) NOT NULL COMMENT '来源：sign_in-签到 comment-评论 share-分享 article-发文章 invite-邀请 purchase-购买',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `related_id` BIGINT DEFAULT NULL COMMENT '关联ID',
  `related_type` VARCHAR(50) DEFAULT NULL COMMENT '关联类型',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `operator_id` BIGINT DEFAULT NULL COMMENT '操作员ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_source` (`source`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_user_points_logs_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分记录表';

-- 积分规则表
CREATE TABLE IF NOT EXISTS `points_rules` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `sign_in_points` INT NOT NULL DEFAULT 5 COMMENT '签到积分',
  `comment_points` INT NOT NULL DEFAULT 2 COMMENT '评论积分',
  `share_points` INT NOT NULL DEFAULT 3 COMMENT '分享积分',
  `article_points` INT NOT NULL DEFAULT 10 COMMENT '发布文章积分',
  `invite_points` INT NOT NULL DEFAULT 20 COMMENT '邀请积分',
  `purchase_rate` INT NOT NULL DEFAULT 10 COMMENT '购买积分比例（消费1元获得多少积分）',
  `max_daily_points` INT NOT NULL DEFAULT 100 COMMENT '每日最大获取积分',
  `enable_exchange` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否可兑换商品',
  `exchange_rate` INT NOT NULL DEFAULT 100 COMMENT '兑换比例（多少积分兑换1元）',
  `min_exchange_points` INT NOT NULL DEFAULT 1000 COMMENT '最小兑换积分',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分规则表';

-- 积分兑换记录表
CREATE TABLE IF NOT EXISTS `points_exchanges` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `points` INT NOT NULL COMMENT '兑换积分',
  `amount` DECIMAL(10,2) NOT NULL COMMENT '兑换金额',
  `exchange_type` VARCHAR(50) NOT NULL COMMENT '兑换类型：cash-现金 goods-商品 coupon-优惠券',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending-处理中 approved-已批准 completed-已完成 rejected-已拒绝',
  `remark` VARCHAR(255) DEFAULT NULL COMMENT '备注',
  `operator_id` BIGINT DEFAULT NULL COMMENT '操作员ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `completed_at` TIMESTAMP NULL DEFAULT NULL COMMENT '完成时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_points_exchanges_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分兑换记录表';

-- 积分商城商品表
CREATE TABLE IF NOT EXISTS `points_goods` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL COMMENT '商品名称',
  `points` INT NOT NULL COMMENT '所需积分',
  `original_price` DECIMAL(10,2) DEFAULT NULL COMMENT '原价',
  `stock` INT NOT NULL DEFAULT 0 COMMENT '库存',
  `sold` INT NOT NULL DEFAULT 0 COMMENT '已售数量',
  `image` VARCHAR(255) DEFAULT NULL COMMENT '商品图片',
  `description` TEXT DEFAULT NULL COMMENT '商品描述',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：0-下架 1-上架',
  `start_time` TIMESTAMP NULL DEFAULT NULL COMMENT '开始时间',
  `end_time` TIMESTAMP NULL DEFAULT NULL COMMENT '结束时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_points` (`points`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分商城商品表';

-- 积分商品兑换记录表
CREATE TABLE IF NOT EXISTS `points_goods_orders` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `order_no` VARCHAR(50) NOT NULL COMMENT '订单编号',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `goods_id` BIGINT NOT NULL COMMENT '商品ID',
  `goods_name` VARCHAR(100) NOT NULL COMMENT '商品名称',
  `points` INT NOT NULL COMMENT '兑换积分',
  `quantity` INT NOT NULL DEFAULT 1 COMMENT '数量',
  `total_points` INT NOT NULL COMMENT '总积分',
  `address_info` TEXT DEFAULT NULL COMMENT '收货信息（JSON格式）',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending-待处理 shipped-已发货 completed-已完成 canceled-已取消',
  `tracking_no` VARCHAR(50) DEFAULT NULL COMMENT '物流单号',
  `shipping_company` VARCHAR(50) DEFAULT NULL COMMENT '物流公司',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `completed_at` TIMESTAMP NULL DEFAULT NULL COMMENT '完成时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_goods_id` (`goods_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_points_goods_orders_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_points_goods_orders_goods_id` FOREIGN KEY (`goods_id`) REFERENCES `points_goods` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分商品兑换记录表';

-- 初始化积分规则数据
INSERT INTO `points_rules` (`sign_in_points`, `comment_points`, `share_points`, `article_points`, `invite_points`, 
                         `purchase_rate`, `max_daily_points`, `enable_exchange`, `exchange_rate`, `min_exchange_points`) 
VALUES (5, 2, 3, 10, 20, 10, 100, 1, 100, 1000);
