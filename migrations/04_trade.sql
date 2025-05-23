-- 交易模块数据表结构

-- 产品表
CREATE TABLE IF NOT EXISTS `products` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL COMMENT '产品ID，业务唯一标识',
  `name` VARCHAR(128) NOT NULL COMMENT '产品名称',
  `company_id` BIGINT NOT NULL COMMENT '公司ID',
  `company_name` VARCHAR(128) NOT NULL COMMENT '公司名称',
  `category` VARCHAR(64) NOT NULL COMMENT '产品分类',
  `price` DECIMAL(10,2) NOT NULL COMMENT '产品价格',
  `specification` VARCHAR(255) DEFAULT NULL COMMENT '规格',
  `material` VARCHAR(255) DEFAULT NULL COMMENT '材质',
  `stock` INT NOT NULL DEFAULT 0 COMMENT '库存',
  `min_order` INT DEFAULT 1 COMMENT '最小订购量',
  `description` TEXT DEFAULT NULL COMMENT '产品描述',
  `images` TEXT DEFAULT NULL COMMENT '产品图片，JSON格式',
  `contact_person` VARCHAR(64) DEFAULT NULL COMMENT '联系人',
  `contact_phone` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
  `contact_email` VARCHAR(128) DEFAULT NULL COMMENT '联系邮箱',
  `address` VARCHAR(255) DEFAULT NULL COMMENT '地址',
  `views` INT NOT NULL DEFAULT 0 COMMENT '浏览量',
  `sales` INT NOT NULL DEFAULT 0 COMMENT '销量',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_product_id` (`product_id`),
  KEY `idx_company_id` (`company_id`),
  KEY `idx_category` (`category`),
  KEY `idx_price` (`price`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品表';

-- 订单表
CREATE TABLE IF NOT EXISTS `orders` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `order_id` VARCHAR(64) NOT NULL COMMENT '订单ID，业务唯一标识',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `product_id` BIGINT NOT NULL COMMENT '产品ID',
  `product_type` VARCHAR(32) NOT NULL COMMENT '产品类型',
  `quantity` INT NOT NULL COMMENT '数量',
  `amount` DECIMAL(10,2) NOT NULL COMMENT '金额',
  `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
  `status` VARCHAR(32) NOT NULL COMMENT '订单状态',
  `payment_id` VARCHAR(64) DEFAULT NULL COMMENT '支付ID',
  `payment_type` VARCHAR(32) DEFAULT NULL COMMENT '支付类型',
  `payment_time` TIMESTAMP NULL DEFAULT NULL COMMENT '支付时间',
  `description` TEXT DEFAULT NULL COMMENT '描述',
  `metadata` TEXT DEFAULT NULL COMMENT '元数据，JSON格式',
  `client_ip` VARCHAR(64) DEFAULT NULL COMMENT '客户端IP',
  `device_id` VARCHAR(64) DEFAULT NULL COMMENT '设备ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `expire_time` TIMESTAMP NULL DEFAULT NULL COMMENT '过期时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_id` (`order_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`),
  KEY `idx_payment_id` (`payment_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_orders_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 订单项表
CREATE TABLE IF NOT EXISTS `order_items` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
  `product_id` BIGINT NOT NULL COMMENT '产品ID',
  `product_type` VARCHAR(32) NOT NULL COMMENT '产品类型',
  `product_name` VARCHAR(128) NOT NULL COMMENT '产品名称',
  `quantity` INT NOT NULL COMMENT '数量',
  `unit_price` DECIMAL(10,2) NOT NULL COMMENT '单价',
  `total_price` DECIMAL(10,2) NOT NULL COMMENT '总价',
  `discount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '折扣',
  `metadata` TEXT DEFAULT NULL COMMENT '元数据，JSON格式',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_product_id` (`product_id`),
  CONSTRAINT `fk_order_items_order_id` FOREIGN KEY (`order_id`) REFERENCES `orders` (`order_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单项表';

-- 支付表
CREATE TABLE IF NOT EXISTS `payments` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `payment_id` VARCHAR(64) NOT NULL COMMENT '支付ID，业务唯一标识',
  `order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `amount` DECIMAL(10,2) NOT NULL COMMENT '支付金额',
  `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
  `payment_type` VARCHAR(32) NOT NULL COMMENT '支付类型',
  `status` VARCHAR(32) NOT NULL COMMENT '支付状态',
  `transaction_id` VARCHAR(128) DEFAULT NULL COMMENT '第三方交易ID',
  `payment_time` TIMESTAMP NULL DEFAULT NULL COMMENT '支付时间',
  `callback_time` TIMESTAMP NULL DEFAULT NULL COMMENT '回调时间',
  `callback_data` TEXT DEFAULT NULL COMMENT '回调原始数据',
  `client_ip` VARCHAR(64) DEFAULT NULL COMMENT '客户端IP',
  `metadata` TEXT DEFAULT NULL COMMENT '元数据，JSON格式',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_payment_id` (`payment_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_payments_order_id` FOREIGN KEY (`order_id`) REFERENCES `orders` (`order_id`),
  CONSTRAINT `fk_payments_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付表';

-- 退款表
CREATE TABLE IF NOT EXISTS `refunds` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `refund_id` VARCHAR(64) NOT NULL COMMENT '退款ID，业务唯一标识',
  `order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
  `payment_id` VARCHAR(64) NOT NULL COMMENT '支付ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `amount` DECIMAL(10,2) NOT NULL COMMENT '退款金额',
  `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
  `status` VARCHAR(32) NOT NULL COMMENT '退款状态',
  `reason` VARCHAR(255) DEFAULT NULL COMMENT '退款原因',
  `description` TEXT DEFAULT NULL COMMENT '描述',
  `processed_by` VARCHAR(64) DEFAULT NULL COMMENT '处理人',
  `process_time` TIMESTAMP NULL DEFAULT NULL COMMENT '处理时间',
  `refund_transaction_id` VARCHAR(128) DEFAULT NULL COMMENT '退款交易ID',
  `metadata` TEXT DEFAULT NULL COMMENT '元数据，JSON格式',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_refund_id` (`refund_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_payment_id` (`payment_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_refunds_order_id` FOREIGN KEY (`order_id`) REFERENCES `orders` (`order_id`),
  CONSTRAINT `fk_refunds_payment_id` FOREIGN KEY (`payment_id`) REFERENCES `payments` (`payment_id`),
  CONSTRAINT `fk_refunds_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款表';

-- 账户余额表
CREATE TABLE IF NOT EXISTS `account_balances` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
  `available` DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '可用余额',
  `pending` DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '待结算余额',
  `frozen` DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '冻结余额',
  `total` DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '总余额',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_currency` (`user_id`, `currency`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_account_balances_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='账户余额表';

-- 交易记录表
CREATE TABLE IF NOT EXISTS `transactions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `transaction_id` VARCHAR(64) NOT NULL COMMENT '交易ID，业务唯一标识',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `related_id` VARCHAR(64) NOT NULL COMMENT '关联ID（订单ID或退款ID）',
  `related_type` VARCHAR(32) NOT NULL COMMENT '关联类型',
  `type` VARCHAR(32) NOT NULL COMMENT '交易类型',
  `amount` DECIMAL(12,2) NOT NULL COMMENT '金额',
  `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
  `balance_before` DECIMAL(12,2) NOT NULL COMMENT '交易前余额',
  `balance_after` DECIMAL(12,2) NOT NULL COMMENT '交易后余额',
  `status` VARCHAR(32) NOT NULL COMMENT '交易状态',
  `description` TEXT DEFAULT NULL COMMENT '描述',
  `metadata` TEXT DEFAULT NULL COMMENT '元数据，JSON格式',
  `operator_id` VARCHAR(64) DEFAULT NULL COMMENT '操作员ID',
  `client_ip` VARCHAR(64) DEFAULT NULL COMMENT '客户端IP',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_related` (`related_type`, `related_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_transactions_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易记录表';

-- 财务日报表
CREATE TABLE IF NOT EXISTS `financial_daily_reports` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `report_date` DATE NOT NULL COMMENT '报表日期',
  `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
  `income` DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '收入',
  `refund` DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '退款',
  `net` DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '净收入',
  `order_count` INT NOT NULL DEFAULT 0 COMMENT '订单数',
  `payment_count` INT NOT NULL DEFAULT 0 COMMENT '支付数',
  `refund_count` INT NOT NULL DEFAULT 0 COMMENT '退款数',
  `user_count` INT NOT NULL DEFAULT 0 COMMENT '用户数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date_currency` (`report_date`, `currency`),
  KEY `idx_report_date` (`report_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='财务日报表';

-- 财务月报表
CREATE TABLE IF NOT EXISTS `financial_monthly_reports` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `report_year` INT NOT NULL COMMENT '报表年份',
  `report_month` INT NOT NULL COMMENT '报表月份',
  `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
  `income` DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '收入',
  `refund` DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '退款',
  `net` DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '净收入',
  `order_count` INT NOT NULL DEFAULT 0 COMMENT '订单数',
  `payment_count` INT NOT NULL DEFAULT 0 COMMENT '支付数',
  `refund_count` INT NOT NULL DEFAULT 0 COMMENT '退款数',
  `user_count` INT NOT NULL DEFAULT 0 COMMENT '用户数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_year_month_currency` (`report_year`, `report_month`, `currency`),
  KEY `idx_report_year_month` (`report_year`, `report_month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='财务月报表';

-- 支付方式表
CREATE TABLE IF NOT EXISTS `payment_methods` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `method_code` VARCHAR(32) NOT NULL COMMENT '方式代码',
  `method_name` VARCHAR(64) NOT NULL COMMENT '方式名称',
  `method_type` VARCHAR(32) NOT NULL COMMENT '方式类型',
  `config` TEXT DEFAULT NULL COMMENT '配置，JSON格式',
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_method_code` (`method_code`),
  KEY `idx_is_enabled` (`is_enabled`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付方式表';
