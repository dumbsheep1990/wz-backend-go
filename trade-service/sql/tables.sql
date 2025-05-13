-- 创建订单表
CREATE TABLE `orders` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` varchar(64) NOT NULL COMMENT '订单编号',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `total_amount` decimal(10,2) NOT NULL COMMENT '订单总金额',
  `pay_amount` decimal(10,2) NOT NULL COMMENT '实付金额',
  `discount_amount` decimal(10,2) DEFAULT '0.00' COMMENT '优惠金额',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '订单状态：0待支付，1已支付，2已发货，3已完成，4已取消，5已退款',
  `pay_type` tinyint(4) DEFAULT NULL COMMENT '支付方式：1支付宝，2微信',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `consignee` varchar(64) DEFAULT NULL COMMENT '收货人',
  `phone` varchar(20) DEFAULT NULL COMMENT '联系电话',
  `address` varchar(255) DEFAULT NULL COMMENT '收货地址',
  `remark` varchar(255) DEFAULT NULL COMMENT '订单备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 创建订单项表
CREATE TABLE `order_items` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '订单项ID',
  `order_id` bigint(20) NOT NULL COMMENT '订单ID',
  `product_id` bigint(20) NOT NULL COMMENT '商品ID',
  `product_name` varchar(128) NOT NULL COMMENT '商品名称',
  `product_image` varchar(255) DEFAULT NULL COMMENT '商品图片',
  `price` decimal(10,2) NOT NULL COMMENT '商品单价',
  `quantity` int(11) NOT NULL COMMENT '商品数量',
  `subtotal` decimal(10,2) NOT NULL COMMENT '小计金额',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单项表';

-- 创建购物车表
CREATE TABLE `carts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '购物车ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `product_id` bigint(20) NOT NULL COMMENT '商品ID',
  `product_name` varchar(128) NOT NULL COMMENT '商品名称',
  `product_image` varchar(255) DEFAULT NULL COMMENT '商品图片',
  `price` decimal(10,2) NOT NULL COMMENT '商品单价',
  `quantity` int(11) NOT NULL COMMENT '商品数量',
  `selected` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否选中：0未选中，1已选中',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_product` (`user_id`,`product_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';

-- 创建支付记录表
CREATE TABLE `payments` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '支付ID',
  `order_id` bigint(20) NOT NULL COMMENT '订单ID',
  `order_no` varchar(64) NOT NULL COMMENT '订单编号',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `pay_type` tinyint(4) NOT NULL COMMENT '支付方式：1支付宝，2微信',
  `trade_no` varchar(64) DEFAULT NULL COMMENT '交易流水号',
  `total_amount` decimal(10,2) NOT NULL COMMENT '支付金额',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '支付状态：0未支付，1支付成功，2支付失败，3已退款',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付记录表';

-- 创建支付配置表
CREATE TABLE `payment_configs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `pay_type` tinyint(4) NOT NULL COMMENT '支付方式：1支付宝，2微信',
  `pay_name` varchar(32) NOT NULL COMMENT '支付名称',
  `app_id` varchar(64) NOT NULL COMMENT '应用ID',
  `merchant_id` varchar(64) DEFAULT NULL COMMENT '商户ID',
  `private_key` text COMMENT '私钥',
  `public_key` text COMMENT '公钥',
  `notify_url` varchar(255) DEFAULT NULL COMMENT '异步通知地址',
  `return_url` varchar(255) DEFAULT NULL COMMENT '同步跳转地址',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：0禁用，1启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_pay_type` (`pay_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付配置表'; 