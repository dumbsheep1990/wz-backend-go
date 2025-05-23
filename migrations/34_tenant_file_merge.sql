-- 租户管理和文件存储表结构合并

-- 租户信息变更记录表（对02_tenants.sql的补充）
CREATE TABLE IF NOT EXISTS `tenant_change_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
  `change_type` VARCHAR(32) NOT NULL COMMENT '变更类型：create-创建 update-更新 status-状态变更 expire-过期变更',
  `field_name` VARCHAR(64) DEFAULT NULL COMMENT '字段名称',
  `old_value` TEXT DEFAULT NULL COMMENT '旧值',
  `new_value` TEXT DEFAULT NULL COMMENT '新值',
  `change_reason` VARCHAR(255) DEFAULT NULL COMMENT '变更原因',
  `operator_id` BIGINT DEFAULT NULL COMMENT '操作人ID',
  `operator_name` VARCHAR(64) DEFAULT NULL COMMENT '操作人名称',
  `ip` VARCHAR(64) DEFAULT NULL COMMENT '操作IP',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id_copy` BIGINT DEFAULT NULL COMMENT '租户ID副本（用于跨租户查询）',
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_change_type` (`change_type`),
  KEY `idx_operator_id` (`operator_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id_copy` (`tenant_id_copy`),
  CONSTRAINT `fk_tenant_change_logs_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户信息变更记录表';

-- 租户配额表
CREATE TABLE IF NOT EXISTS `tenant_quotas` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
  `resource_type` VARCHAR(32) NOT NULL COMMENT '资源类型：user-用户数 storage-存储空间 api-API调用 bandwidth-带宽',
  `quota_limit` BIGINT NOT NULL COMMENT '配额限制',
  `used_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '已使用数量',
  `unit` VARCHAR(16) DEFAULT NULL COMMENT '单位',
  `reset_cycle` VARCHAR(16) DEFAULT NULL COMMENT '重置周期：none-不重置 daily-每日 monthly-每月 yearly-每年',
  `last_reset_time` TIMESTAMP NULL COMMENT '上次重置时间',
  `next_reset_time` TIMESTAMP NULL COMMENT '下次重置时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id_copy` BIGINT DEFAULT NULL COMMENT '租户ID副本（用于跨租户查询）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_resource` (`tenant_id`, `resource_type`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_resource_type` (`resource_type`),
  KEY `idx_tenant_id_copy` (`tenant_id_copy`),
  CONSTRAINT `fk_tenant_quotas_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户配额表';

-- 租户计费表
CREATE TABLE IF NOT EXISTS `tenant_billings` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
  `plan_id` BIGINT DEFAULT NULL COMMENT '套餐ID',
  `billing_cycle` VARCHAR(16) NOT NULL COMMENT '计费周期：monthly-月付 quarterly-季付 yearly-年付',
  `amount` DECIMAL(10,2) NOT NULL COMMENT '金额',
  `currency` VARCHAR(16) NOT NULL DEFAULT 'CNY' COMMENT '货币',
  `status` VARCHAR(16) NOT NULL COMMENT '状态：active-活跃 pending-待支付 expired-已过期 canceled-已取消',
  `payment_method` VARCHAR(32) DEFAULT NULL COMMENT '支付方式',
  `start_date` DATE NOT NULL COMMENT '开始日期',
  `end_date` DATE NOT NULL COMMENT '结束日期',
  `last_payment_time` TIMESTAMP NULL COMMENT '最后支付时间',
  `next_payment_time` TIMESTAMP NULL COMMENT '下次支付时间',
  `auto_renew` TINYINT(1) DEFAULT 0 COMMENT '是否自动续费',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id_copy` BIGINT DEFAULT NULL COMMENT '租户ID副本（用于跨租户查询）',
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_plan_id` (`plan_id`),
  KEY `idx_status` (`status`),
  KEY `idx_start_date` (`start_date`),
  KEY `idx_end_date` (`end_date`),
  KEY `idx_tenant_id_copy` (`tenant_id_copy`),
  CONSTRAINT `fk_tenant_billings_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户计费表';

-- 租户邀请表
CREATE TABLE IF NOT EXISTS `tenant_invitations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
  `inviter_id` BIGINT NOT NULL COMMENT '邀请人ID',
  `email` VARCHAR(128) NOT NULL COMMENT '被邀请人邮箱',
  `role` VARCHAR(32) NOT NULL COMMENT '角色',
  `invitation_code` VARCHAR(64) NOT NULL COMMENT '邀请码',
  `message` VARCHAR(512) DEFAULT NULL COMMENT '邀请消息',
  `status` VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT '状态：pending-待接受 accepted-已接受 rejected-已拒绝 expired-已过期',
  `accepted_user_id` BIGINT DEFAULT NULL COMMENT '接受邀请的用户ID',
  `expire_time` TIMESTAMP NOT NULL COMMENT '过期时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id_copy` BIGINT DEFAULT NULL COMMENT '租户ID副本（用于跨租户查询）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_invitation_code` (`invitation_code`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_inviter_id` (`inviter_id`),
  KEY `idx_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id_copy` (`tenant_id_copy`),
  CONSTRAINT `fk_tenant_invitations_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户邀请表';

-- 文件分片上传表（对10_files.sql的补充）
CREATE TABLE IF NOT EXISTS `file_chunks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `file_id` BIGINT NOT NULL COMMENT '文件ID',
  `chunk_number` INT NOT NULL COMMENT '分片序号',
  `chunk_size` BIGINT NOT NULL COMMENT '分片大小',
  `chunk_path` VARCHAR(255) NOT NULL COMMENT '分片路径',
  `status` VARCHAR(16) NOT NULL DEFAULT 'uploading' COMMENT '状态：uploading-上传中 uploaded-已上传 merged-已合并',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_chunk` (`file_id`, `chunk_number`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_file_chunks_file_id` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件分片上传表';

-- 文件处理任务表
CREATE TABLE IF NOT EXISTS `file_processing_tasks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `file_id` BIGINT NOT NULL COMMENT '文件ID',
  `task_type` VARCHAR(32) NOT NULL COMMENT '任务类型：compress-压缩 extract-解压 convert-转换 analyze-分析 resize-调整大小',
  `status` VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT '状态：pending-待处理 processing-处理中 completed-已完成 failed-失败',
  `parameters` JSON DEFAULT NULL COMMENT '处理参数',
  `result` JSON DEFAULT NULL COMMENT '处理结果',
  `error_message` VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
  `output_file_id` BIGINT DEFAULT NULL COMMENT '输出文件ID',
  `progress` DECIMAL(5,2) DEFAULT 0 COMMENT '进度',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `started_at` TIMESTAMP NULL COMMENT '开始时间',
  `completed_at` TIMESTAMP NULL COMMENT '完成时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_task_type` (`task_type`),
  KEY `idx_status` (`status`),
  KEY `idx_output_file_id` (`output_file_id`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_file_processing_tasks_file_id` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_file_processing_tasks_output_file_id` FOREIGN KEY (`output_file_id`) REFERENCES `files` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件处理任务表';

-- 文件共享表
CREATE TABLE IF NOT EXISTS `file_shares` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `file_id` BIGINT NOT NULL COMMENT '文件ID',
  `creator_id` BIGINT NOT NULL COMMENT '创建人ID',
  `share_code` VARCHAR(32) NOT NULL COMMENT '分享码',
  `share_type` VARCHAR(16) NOT NULL DEFAULT 'link' COMMENT '分享类型：link-链接 email-邮件',
  `access_type` VARCHAR(16) NOT NULL DEFAULT 'view' COMMENT '访问类型：view-查看 download-下载 edit-编辑',
  `password` VARCHAR(64) DEFAULT NULL COMMENT '访问密码',
  `expire_time` TIMESTAMP NULL COMMENT '过期时间',
  `max_views` INT DEFAULT NULL COMMENT '最大查看次数',
  `view_count` INT DEFAULT 0 COMMENT '已查看次数',
  `download_count` INT DEFAULT 0 COMMENT '下载次数',
  `status` VARCHAR(16) NOT NULL DEFAULT 'active' COMMENT '状态：active-活跃 expired-已过期 disabled-已禁用',
  `recipients` JSON DEFAULT NULL COMMENT '接收人列表',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_share_code` (`share_code`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_creator_id` (`creator_id`),
  KEY `idx_expire_time` (`expire_time`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_file_shares_file_id` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件共享表';
