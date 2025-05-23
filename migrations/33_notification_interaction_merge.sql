-- 通知服务和交互服务表结构合并

-- 通知模板版本表（对22_notification.sql的补充）
CREATE TABLE IF NOT EXISTS `notification_template_versions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `template_id` BIGINT NOT NULL COMMENT '模板ID',
  `version` VARCHAR(32) NOT NULL COMMENT '版本号',
  `title_template` VARCHAR(255) NOT NULL COMMENT '标题模板',
  `content_template` TEXT NOT NULL COMMENT '内容模板',
  `variables` JSON DEFAULT NULL COMMENT '变量定义(JSON格式)',
  `status` INT DEFAULT 0 COMMENT '状态：0-草稿 1-已发布 2-已归档',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `published_at` TIMESTAMP NULL COMMENT '发布时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_template_version` (`template_id`, `version`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_notification_template_versions_template_id` FOREIGN KEY (`template_id`) REFERENCES `notification_templates` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知模板版本表';

-- 通知测试记录表
CREATE TABLE IF NOT EXISTS `notification_test_records` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `template_id` BIGINT NOT NULL COMMENT '模板ID',
  `template_version_id` BIGINT DEFAULT NULL COMMENT '模板版本ID',
  `channel_id` BIGINT DEFAULT NULL COMMENT '渠道ID',
  `test_data` JSON DEFAULT NULL COMMENT '测试数据(JSON格式)',
  `recipient` VARCHAR(255) DEFAULT NULL COMMENT '接收者',
  `result` VARCHAR(32) DEFAULT NULL COMMENT '测试结果：success-成功 failed-失败',
  `error_message` VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
  `rendered_content` TEXT DEFAULT NULL COMMENT '渲染后的内容',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_template_version_id` (`template_version_id`),
  KEY `idx_channel_id` (`channel_id`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_notification_test_records_template_id` FOREIGN KEY (`template_id`) REFERENCES `notification_templates` (`id`),
  CONSTRAINT `fk_notification_test_records_version_id` FOREIGN KEY (`template_version_id`) REFERENCES `notification_template_versions` (`id`),
  CONSTRAINT `fk_notification_test_records_channel_id` FOREIGN KEY (`channel_id`) REFERENCES `notification_channels` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知测试记录表';

-- 通知发送统计表
CREATE TABLE IF NOT EXISTS `notification_send_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `date` DATE NOT NULL COMMENT '统计日期',
  `template_id` BIGINT DEFAULT NULL COMMENT '模板ID',
  `channel_id` BIGINT DEFAULT NULL COMMENT '渠道ID',
  `channel_type` VARCHAR(32) DEFAULT NULL COMMENT '渠道类型',
  `total_count` INT DEFAULT 0 COMMENT '总发送量',
  `success_count` INT DEFAULT 0 COMMENT '成功数量',
  `failed_count` INT DEFAULT 0 COMMENT '失败数量',
  `open_count` INT DEFAULT 0 COMMENT '打开数量',
  `click_count` INT DEFAULT 0 COMMENT '点击数量',
  `response_time_avg` INT DEFAULT 0 COMMENT '平均响应时间(毫秒)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date_template_channel` (`date`, `template_id`, `channel_id`, `tenant_id`),
  KEY `idx_date` (`date`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_channel_id` (`channel_id`),
  KEY `idx_channel_type` (`channel_type`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知发送统计表';

-- 交互评分详情表（对20_interaction.sql的补充）
CREATE TABLE IF NOT EXISTS `rating_details` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `rating_id` BIGINT NOT NULL COMMENT '评分ID',
  `dimension` VARCHAR(64) NOT NULL COMMENT '评分维度',
  `score` DECIMAL(3,1) NOT NULL COMMENT '维度评分',
  `weight` DECIMAL(3,2) DEFAULT 1.00 COMMENT '权重',
  `comment` VARCHAR(512) DEFAULT NULL COMMENT '维度评价',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_rating_id` (`rating_id`),
  KEY `idx_dimension` (`dimension`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_rating_details_rating_id` FOREIGN KEY (`rating_id`) REFERENCES `ratings` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交互评分详情表';

-- 交互内容举报处理表
CREATE TABLE IF NOT EXISTS `report_processing` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `report_id` BIGINT NOT NULL COMMENT '举报ID',
  `processor_id` BIGINT NOT NULL COMMENT '处理人ID',
  `status` VARCHAR(32) NOT NULL COMMENT '处理状态：pending-待处理 processing-处理中 resolved-已解决 rejected-已拒绝',
  `processing_result` VARCHAR(32) DEFAULT NULL COMMENT '处理结果：ignore-忽略 warning-警告 delete-删除 ban-封禁',
  `comment` VARCHAR(512) DEFAULT NULL COMMENT '处理意见',
  `target_action` VARCHAR(64) DEFAULT NULL COMMENT '目标操作',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_report_id` (`report_id`),
  KEY `idx_processor_id` (`processor_id`),
  KEY `idx_status` (`status`),
  KEY `idx_processing_result` (`processing_result`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_report_processing_report_id` FOREIGN KEY (`report_id`) REFERENCES `reports` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交互内容举报处理表';

-- 内容敏感词表
CREATE TABLE IF NOT EXISTS `sensitive_words` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `word` VARCHAR(128) NOT NULL COMMENT '敏感词',
  `category` VARCHAR(32) NOT NULL COMMENT '分类：politics-政治 porn-色情 abuse-辱骂 spam-垃圾信息 custom-自定义',
  `level` INT DEFAULT 1 COMMENT '级别：1-低 2-中 3-高',
  `replacement` VARCHAR(128) DEFAULT NULL COMMENT '替换词',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_word_tenant` (`word`, `tenant_id`),
  KEY `idx_category` (`category`),
  KEY `idx_level` (`level`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容敏感词表';
