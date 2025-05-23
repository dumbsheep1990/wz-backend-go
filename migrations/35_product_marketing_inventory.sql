-- 产品管理、电商营销和库存管理补充表结构

-- 产品分类表（补充04_trade.sql中的产品表）
CREATE TABLE IF NOT EXISTS `product_categories` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `parent_id` BIGINT DEFAULT NULL COMMENT '父分类ID',
  `name` VARCHAR(64) NOT NULL COMMENT '分类名称',
  `code` VARCHAR(32) NOT NULL COMMENT '分类编码',
  `level` INT NOT NULL DEFAULT 1 COMMENT '层级，1为顶级',
  `path` VARCHAR(255) DEFAULT NULL COMMENT '分类路径，例如：1,2,3',
  `icon` VARCHAR(255) DEFAULT NULL COMMENT '分类图标',
  `banner` VARCHAR(255) DEFAULT NULL COMMENT '分类banner',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `is_visible` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否可见',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '描述',
  `seo_title` VARCHAR(128) DEFAULT NULL COMMENT 'SEO标题',
  `seo_keywords` VARCHAR(255) DEFAULT NULL COMMENT 'SEO关键词',
  `seo_description` VARCHAR(512) DEFAULT NULL COMMENT 'SEO描述',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_code` (`tenant_id`, `code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_is_visible` (`is_visible`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品分类表';

-- 产品属性表
CREATE TABLE IF NOT EXISTS `product_attributes` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '属性名称',
  `code` VARCHAR(32) NOT NULL COMMENT '属性编码',
  `attribute_type` VARCHAR(16) NOT NULL COMMENT '属性类型：text-文本 number-数字 enum-枚举 date-日期',
  `input_type` VARCHAR(16) NOT NULL COMMENT '输入类型：input-输入框 select-下拉 radio-单选 checkbox-多选 date-日期',
  `is_required` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否必填',
  `is_multiple` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否多选',
  `is_filterable` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否可筛选',
  `is_searchable` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否可搜索',
  `is_comparable` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否可比较',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `default_value` VARCHAR(255) DEFAULT NULL COMMENT '默认值',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '描述',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_code` (`tenant_id`, `code`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品属性表';

-- 产品属性值表
CREATE TABLE IF NOT EXISTS `product_attribute_values` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `attribute_id` BIGINT NOT NULL COMMENT '属性ID',
  `value` VARCHAR(255) NOT NULL COMMENT '属性值',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_attribute_id` (`attribute_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_sort_order` (`sort_order`),
  CONSTRAINT `fk_product_attribute_values_attribute_id` FOREIGN KEY (`attribute_id`) REFERENCES `product_attributes` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品属性值表';

-- 产品规格表
CREATE TABLE IF NOT EXISTS `product_specifications` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL COMMENT '产品ID',
  `spec_code` VARCHAR(64) NOT NULL COMMENT '规格编码，如SKU',
  `spec_values` JSON NOT NULL COMMENT '规格值，JSON格式，如：{"颜色":"红色","尺寸":"XL"}',
  `price` DECIMAL(10,2) NOT NULL COMMENT '价格',
  `original_price` DECIMAL(10,2) DEFAULT NULL COMMENT '原价',
  `cost_price` DECIMAL(10,2) DEFAULT NULL COMMENT '成本价',
  `weight` DECIMAL(10,2) DEFAULT NULL COMMENT '重量(g)',
  `volume` DECIMAL(10,2) DEFAULT NULL COMMENT '体积(cm³)',
  `stock` INT NOT NULL DEFAULT 0 COMMENT '库存',
  `sold` INT NOT NULL DEFAULT 0 COMMENT '已售',
  `is_default` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否默认',
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `images` TEXT DEFAULT NULL COMMENT '图片，JSON格式',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_product_spec_code` (`product_id`, `spec_code`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_is_enabled` (`is_enabled`),
  CONSTRAINT `fk_product_specifications_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品规格表';

-- 产品分类关联表
CREATE TABLE IF NOT EXISTS `product_category_relations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL COMMENT '产品ID',
  `category_id` BIGINT NOT NULL COMMENT '分类ID',
  `is_primary` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否主分类',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_product_category` (`product_id`, `category_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_is_primary` (`is_primary`),
  CONSTRAINT `fk_product_category_relations_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_product_category_relations_category_id` FOREIGN KEY (`category_id`) REFERENCES `product_categories` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品分类关联表';

-- 产品属性关联表
CREATE TABLE IF NOT EXISTS `product_attribute_relations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL COMMENT '产品ID',
  `attribute_id` BIGINT NOT NULL COMMENT '属性ID',
  `attribute_value` VARCHAR(512) DEFAULT NULL COMMENT '属性值',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_product_attribute` (`product_id`, `attribute_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_attribute_id` (`attribute_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_product_attribute_relations_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_product_attribute_relations_attribute_id` FOREIGN KEY (`attribute_id`) REFERENCES `product_attributes` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品属性关联表';

-- 购物车表
CREATE TABLE IF NOT EXISTS `shopping_carts` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `product_id` BIGINT NOT NULL COMMENT '产品ID',
  `spec_id` BIGINT DEFAULT NULL COMMENT '规格ID',
  `quantity` INT NOT NULL DEFAULT 1 COMMENT '数量',
  `selected` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否选中',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_product_spec` (`user_id`, `product_id`, `spec_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_spec_id` (`spec_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_shopping_carts_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_shopping_carts_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_shopping_carts_spec_id` FOREIGN KEY (`spec_id`) REFERENCES `product_specifications` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';

-- 优惠券表
CREATE TABLE IF NOT EXISTS `coupons` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(32) NOT NULL COMMENT '优惠码',
  `name` VARCHAR(64) NOT NULL COMMENT '名称',
  `type` VARCHAR(16) NOT NULL COMMENT '类型：amount-金额 percent-百分比',
  `value` DECIMAL(10,2) NOT NULL COMMENT '优惠值',
  `min_amount` DECIMAL(10,2) DEFAULT NULL COMMENT '最低消费金额',
  `max_discount` DECIMAL(10,2) DEFAULT NULL COMMENT '最大优惠金额',
  `quantity` INT NOT NULL COMMENT '发放数量',
  `used_quantity` INT NOT NULL DEFAULT 0 COMMENT '已使用数量',
  `start_time` TIMESTAMP NULL COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `description` TEXT DEFAULT NULL COMMENT '描述',
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_code` (`tenant_id`, `code`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_type` (`type`),
  KEY `idx_is_enabled` (`is_enabled`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_end_time` (`end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券表';

-- 用户优惠券表
CREATE TABLE IF NOT EXISTS `user_coupons` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `coupon_id` BIGINT NOT NULL COMMENT '优惠券ID',
  `coupon_code` VARCHAR(32) NOT NULL COMMENT '优惠码',
  `status` VARCHAR(16) NOT NULL COMMENT '状态：unused-未使用 used-已使用 expired-已过期',
  `get_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '获取时间',
  `use_time` TIMESTAMP NULL COMMENT '使用时间',
  `order_id` VARCHAR(64) DEFAULT NULL COMMENT '订单ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_coupon_id` (`coupon_id`),
  KEY `idx_status` (`status`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_user_coupons_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_coupons_coupon_id` FOREIGN KEY (`coupon_id`) REFERENCES `coupons` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户优惠券表';

-- 促销活动表
CREATE TABLE IF NOT EXISTS `promotions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '活动名称',
  `code` VARCHAR(32) NOT NULL COMMENT '活动编码',
  `type` VARCHAR(16) NOT NULL COMMENT '活动类型：discount-折扣 flash-秒杀 bundle-捆绑',
  `start_time` TIMESTAMP NOT NULL COMMENT '开始时间',
  `end_time` TIMESTAMP NOT NULL COMMENT '结束时间',
  `status` VARCHAR(16) NOT NULL COMMENT '状态：pending-未开始 active-进行中 ended-已结束 canceled-已取消',
  `priority` INT NOT NULL DEFAULT 0 COMMENT '优先级',
  `rules` JSON DEFAULT NULL COMMENT '规则，JSON格式',
  `actions` JSON DEFAULT NULL COMMENT '动作，JSON格式',
  `description` TEXT DEFAULT NULL COMMENT '描述',
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_code` (`tenant_id`, `code`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_end_time` (`end_time`),
  KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='促销活动表';

-- 促销产品关联表
CREATE TABLE IF NOT EXISTS `promotion_products` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `promotion_id` BIGINT NOT NULL COMMENT '促销ID',
  `product_id` BIGINT NOT NULL COMMENT '产品ID',
  `spec_id` BIGINT DEFAULT NULL COMMENT '规格ID',
  `promotion_price` DECIMAL(10,2) DEFAULT NULL COMMENT '促销价格',
  `promotion_stock` INT DEFAULT NULL COMMENT '促销库存',
  `limit_per_user` INT DEFAULT NULL COMMENT '每人限购',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_promotion_product_spec` (`promotion_id`, `product_id`, `spec_id`),
  KEY `idx_promotion_id` (`promotion_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_spec_id` (`spec_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_promotion_products_promotion_id` FOREIGN KEY (`promotion_id`) REFERENCES `promotions` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_promotion_products_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_promotion_products_spec_id` FOREIGN KEY (`spec_id`) REFERENCES `product_specifications` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='促销产品关联表';

-- 库存变动记录表
CREATE TABLE IF NOT EXISTS `inventory_changes` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL COMMENT '产品ID',
  `spec_id` BIGINT DEFAULT NULL COMMENT '规格ID',
  `change_type` VARCHAR(16) NOT NULL COMMENT '变动类型：in-入库 out-出库 adjust-调整',
  `change_quantity` INT NOT NULL COMMENT '变动数量',
  `before_quantity` INT NOT NULL COMMENT '变动前数量',
  `after_quantity` INT NOT NULL COMMENT '变动后数量',
  `related_type` VARCHAR(16) DEFAULT NULL COMMENT '关联类型：order-订单 return-退货 adjust-调整',
  `related_id` VARCHAR(64) DEFAULT NULL COMMENT '关联ID',
  `operator_id` BIGINT DEFAULT NULL COMMENT '操作人ID',
  `operator_name` VARCHAR(64) DEFAULT NULL COMMENT '操作人名称',
  `remark` VARCHAR(255) DEFAULT NULL COMMENT '备注',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_spec_id` (`spec_id`),
  KEY `idx_change_type` (`change_type`),
  KEY `idx_related` (`related_type`, `related_id`),
  KEY `idx_operator_id` (`operator_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_inventory_changes_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_inventory_changes_spec_id` FOREIGN KEY (`spec_id`) REFERENCES `product_specifications` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存变动记录表';

-- 商品评价表
CREATE TABLE IF NOT EXISTS `product_reviews` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `product_id` BIGINT NOT NULL COMMENT '产品ID',
  `order_id` VARCHAR(64) DEFAULT NULL COMMENT '订单ID',
  `spec_id` BIGINT DEFAULT NULL COMMENT '规格ID',
  `rating` INT NOT NULL COMMENT '评分(1-5星)',
  `content` TEXT DEFAULT NULL COMMENT '评价内容',
  `images` TEXT DEFAULT NULL COMMENT '图片，JSON格式',
  `is_anonymous` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否匿名',
  `status` VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT '状态：pending-待审核 approved-已通过 rejected-已拒绝',
  `reply_content` TEXT DEFAULT NULL COMMENT '回复内容',
  `reply_time` TIMESTAMP NULL COMMENT '回复时间',
  `likes` INT NOT NULL DEFAULT 0 COMMENT '点赞数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_spec_id` (`spec_id`),
  KEY `idx_rating` (`rating`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_product_reviews_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_product_reviews_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_product_reviews_spec_id` FOREIGN KEY (`spec_id`) REFERENCES `product_specifications` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品评价表';

-- 收货地址表
CREATE TABLE IF NOT EXISTS `shipping_addresses` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `receiver_name` VARCHAR(64) NOT NULL COMMENT '收货人姓名',
  `receiver_phone` VARCHAR(20) NOT NULL COMMENT '收货人电话',
  `province` VARCHAR(32) NOT NULL COMMENT '省份',
  `city` VARCHAR(32) NOT NULL COMMENT '城市',
  `district` VARCHAR(32) NOT NULL COMMENT '区县',
  `detail_address` VARCHAR(255) NOT NULL COMMENT '详细地址',
  `postal_code` VARCHAR(16) DEFAULT NULL COMMENT '邮政编码',
  `is_default` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否默认',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_is_default` (`is_default`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_shipping_addresses_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收货地址表';
