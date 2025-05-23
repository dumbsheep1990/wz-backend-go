-- 用户模块数据表结构

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(64) NOT NULL COMMENT '用户名',
  `password` VARCHAR(128) NOT NULL COMMENT '密码',
  `email` VARCHAR(128) NOT NULL COMMENT '邮箱',
  `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1-正常 2-禁用',
  `is_verified` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已验证',
  `is_company_verified` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已企业认证',
  `default_tenant_id` BIGINT DEFAULT NULL COMMENT '默认租户ID',
  `role` VARCHAR(20) NOT NULL COMMENT '用户角色：platform_admin, tenant_admin, tenant_user, personal_user',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`),
  UNIQUE KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表';

-- 用户详情表
CREATE TABLE IF NOT EXISTS `user_details` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `real_name` VARCHAR(64) DEFAULT NULL COMMENT '真实姓名',
  `id_card` VARCHAR(32) DEFAULT NULL COMMENT '身份证号',
  `avatar` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
  `gender` INT DEFAULT 0 COMMENT '性别：0-未知 1-男 2-女',
  `birthday` DATE DEFAULT NULL COMMENT '生日',
  `address` VARCHAR(255) DEFAULT NULL COMMENT '地址',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_user_details_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户详情表';

-- 企业认证表
CREATE TABLE IF NOT EXISTS `company_verifications` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `company_type` INT NOT NULL COMMENT '企业类型：1-企业 2-集团 3-政府机构 4-科研所',
  `company_name` VARCHAR(128) NOT NULL COMMENT '企业名称',
  `business_license` VARCHAR(255) DEFAULT NULL COMMENT '营业执照',
  `committee_letter` VARCHAR(255) DEFAULT NULL COMMENT '委托书',
  `org_code_cert` VARCHAR(255) DEFAULT NULL COMMENT '组织机构代码证',
  `agency_cert` VARCHAR(255) DEFAULT NULL COMMENT '代理证明',
  `org_structure` VARCHAR(255) DEFAULT NULL COMMENT '组织结构',
  `unified_social_credit` VARCHAR(64) DEFAULT NULL COMMENT '统一社会信用代码',
  `listing_cert` VARCHAR(255) DEFAULT NULL COMMENT '上市证明',
  `contact_person` VARCHAR(64) NOT NULL COMMENT '联系人',
  `contact_phone` VARCHAR(20) NOT NULL COMMENT '联系电话',
  `uploaded_document` VARCHAR(255) DEFAULT NULL COMMENT '上传的证明文件',
  `status` INT NOT NULL DEFAULT 0 COMMENT '状态：0-待审核 1-已通过 2-已拒绝',
  `remark` TEXT DEFAULT NULL COMMENT '备注',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_company_verifications_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业认证表';

-- 用户登录日志表
CREATE TABLE IF NOT EXISTS `user_login_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `login_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `login_ip` VARCHAR(64) DEFAULT NULL COMMENT '登录IP',
  `user_agent` VARCHAR(255) DEFAULT NULL COMMENT '用户代理',
  `device_type` INT DEFAULT 0 COMMENT '设备类型：0-未知 1-PC 2-移动端 3-平板',
  `login_status` INT NOT NULL COMMENT '登录状态：0-失败 1-成功',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_login_time` (`login_time`),
  CONSTRAINT `fk_user_login_logs_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户登录日志表';

-- 用户行为日志表
CREATE TABLE IF NOT EXISTS `user_behavior_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `action` VARCHAR(64) NOT NULL COMMENT '行为动作',
  `resource_type` VARCHAR(64) DEFAULT NULL COMMENT '资源类型',
  `resource_id` BIGINT DEFAULT NULL COMMENT '资源ID',
  `ip` VARCHAR(64) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` VARCHAR(255) DEFAULT NULL COMMENT '用户代理',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_user_behavior_logs_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户行为日志表';

-- 企业入驻信息表
CREATE TABLE IF NOT EXISTS `enterprise_registrations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `company_name` VARCHAR(128) NOT NULL COMMENT '公司名称',
  `company_type` INT NOT NULL COMMENT '企业类型：1-企业 2-集团 3-政府机构 4-科研所',
  `contact_person` VARCHAR(64) NOT NULL COMMENT '联系人',
  `job_position` VARCHAR(64) NOT NULL COMMENT '职位',
  `region` VARCHAR(128) NOT NULL COMMENT '地区',
  `verification_method` VARCHAR(32) NOT NULL COMMENT '验证方式',
  `detailed_address` VARCHAR(255) NOT NULL COMMENT '详细地址',
  `location_latitude` DECIMAL(10,6) DEFAULT NULL COMMENT '位置纬度',
  `location_longitude` DECIMAL(10,6) DEFAULT NULL COMMENT '位置经度',
  `status` INT NOT NULL DEFAULT 0 COMMENT '状态：0-待审核 1-已通过 2-已拒绝',
  `remark` TEXT DEFAULT NULL COMMENT '备注',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_enterprise_registrations_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业入驻信息表';

-- 用户积分表
CREATE TABLE IF NOT EXISTS `user_points` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `points` INT NOT NULL DEFAULT 0 COMMENT '积分总额',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_user_points_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分表';

-- 用户积分流水表
CREATE TABLE IF NOT EXISTS `user_point_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `action` VARCHAR(64) NOT NULL COMMENT '积分动作',
  `points` INT NOT NULL COMMENT '积分变动值',
  `balance` INT NOT NULL COMMENT '积分余额',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_user_point_logs_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分流水表';

-- 用户收藏表
CREATE TABLE IF NOT EXISTS `user_favorites` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `resource_type` VARCHAR(64) NOT NULL COMMENT '资源类型',
  `resource_id` BIGINT NOT NULL COMMENT '资源ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_resource` (`user_id`, `resource_type`, `resource_id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_user_favorites_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收藏表';

-- 用户消息表
CREATE TABLE IF NOT EXISTS `user_messages` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `title` VARCHAR(128) NOT NULL COMMENT '消息标题',
  `content` TEXT NOT NULL COMMENT '消息内容',
  `message_type` VARCHAR(32) NOT NULL COMMENT '消息类型',
  `is_read` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已读',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `read_at` TIMESTAMP NULL DEFAULT NULL COMMENT '阅读时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_user_messages_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户消息表';
