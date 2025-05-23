-- 通知服务数据表结构

-- 通知模板表
CREATE TABLE IF NOT EXISTS `notification_templates` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '模板名称',
  `code` VARCHAR(64) NOT NULL COMMENT '模板编码',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '模板描述',
  `type` VARCHAR(32) NOT NULL COMMENT '通知类型：system-系统通知 marketing-营销通知 account-账户通知 transaction-交易通知 social-社交通知',
  `title_template` VARCHAR(255) NOT NULL COMMENT '标题模板',
  `content_template` TEXT NOT NULL COMMENT '内容模板',
  `variables` JSON DEFAULT NULL COMMENT '变量定义(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `category` VARCHAR(64) DEFAULT 'default' COMMENT '分类',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_category` (`category`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知模板表';

-- 通知渠道表
CREATE TABLE IF NOT EXISTS `notification_channels` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '渠道名称',
  `code` VARCHAR(64) NOT NULL COMMENT '渠道编码',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '渠道描述',
  `type` VARCHAR(32) NOT NULL COMMENT '渠道类型：app-应用内 email-邮件 sms-短信 push-推送 wechat-微信 webhook-Webhook',
  `config` JSON DEFAULT NULL COMMENT '渠道配置(JSON格式)',
  `provider` VARCHAR(64) DEFAULT NULL COMMENT '服务提供商',
  `credentials` TEXT DEFAULT NULL COMMENT '认证凭证（加密存储）',
  `template_mapping` JSON DEFAULT NULL COMMENT '模板映射(JSON格式)',
  `rate_limit` INT DEFAULT 0 COMMENT '速率限制(每分钟)，0表示不限制',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_type` (`type`),
  KEY `idx_provider` (`provider`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知渠道表';

-- 通知策略表
CREATE TABLE IF NOT EXISTS `notification_policies` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '策略名称',
  `code` VARCHAR(64) NOT NULL COMMENT '策略编码',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '策略描述',
  `template_id` BIGINT NOT NULL COMMENT '模板ID',
  `channels` JSON NOT NULL COMMENT '渠道配置(JSON格式)',
  `priority` INT DEFAULT 0 COMMENT '优先级，数字越大优先级越高',
  `conditions` JSON DEFAULT NULL COMMENT '触发条件(JSON格式)',
  `schedule` JSON DEFAULT NULL COMMENT '发送时间计划(JSON格式)',
  `throttle_config` JSON DEFAULT NULL COMMENT '频率控制配置(JSON格式)',
  `retry_config` JSON DEFAULT NULL COMMENT '重试配置(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_priority` (`priority`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_notification_policies_template_id` FOREIGN KEY (`template_id`) REFERENCES `notification_templates` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知策略表';

-- 通知表
CREATE TABLE IF NOT EXISTS `notifications` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `notification_id` VARCHAR(64) NOT NULL COMMENT '通知ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `template_id` BIGINT NOT NULL COMMENT '模板ID',
  `policy_id` BIGINT DEFAULT NULL COMMENT '策略ID',
  `title` VARCHAR(255) NOT NULL COMMENT '通知标题',
  `content` TEXT NOT NULL COMMENT '通知内容',
  `content_type` VARCHAR(32) DEFAULT 'text' COMMENT '内容类型：text-文本 html-HTML markdown-Markdown',
  `type` VARCHAR(32) NOT NULL COMMENT '通知类型：system-系统通知 marketing-营销通知 account-账户通知 transaction-交易通知 social-社交通知',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-待发送 sending-发送中 sent-已发送 failed-发送失败',
  `priority` INT DEFAULT 0 COMMENT '优先级：0-普通 1-重要 2-紧急',
  `payload` JSON DEFAULT NULL COMMENT '附加数据(JSON格式)',
  `metadata` JSON DEFAULT NULL COMMENT '元数据(JSON格式)',
  `source` VARCHAR(64) DEFAULT NULL COMMENT '来源',
  `source_id` VARCHAR(64) DEFAULT NULL COMMENT '来源ID',
  `is_read` TINYINT(1) DEFAULT 0 COMMENT '是否已读：0-未读 1-已读',
  `read_time` TIMESTAMP NULL COMMENT '阅读时间',
  `is_deleted` TINYINT(1) DEFAULT 0 COMMENT '是否删除：0-未删除 1-已删除',
  `delete_time` TIMESTAMP NULL COMMENT '删除时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `expire_at` TIMESTAMP NULL COMMENT '过期时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_notification_id` (`notification_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_policy_id` (`policy_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_priority` (`priority`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_is_deleted` (`is_deleted`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_expire_at` (`expire_at`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_notifications_template_id` FOREIGN KEY (`template_id`) REFERENCES `notification_templates` (`id`),
  CONSTRAINT `fk_notifications_policy_id` FOREIGN KEY (`policy_id`) REFERENCES `notification_policies` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知表';

-- 通知发送记录表
CREATE TABLE IF NOT EXISTS `notification_delivery_records` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `notification_id` VARCHAR(64) NOT NULL COMMENT '通知ID',
  `channel_id` BIGINT NOT NULL COMMENT '渠道ID',
  `channel_type` VARCHAR(32) NOT NULL COMMENT '渠道类型',
  `recipient` VARCHAR(255) NOT NULL COMMENT '接收者',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-待发送 sending-发送中 delivered-已送达 failed-失败',
  `sent_time` TIMESTAMP NULL COMMENT '发送时间',
  `delivered_time` TIMESTAMP NULL COMMENT '送达时间',
  `error_message` VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
  `retry_count` INT DEFAULT 0 COMMENT '重试次数',
  `next_retry_time` TIMESTAMP NULL COMMENT '下次重试时间',
  `provider_response` TEXT DEFAULT NULL COMMENT '服务商响应',
  `tracking_info` JSON DEFAULT NULL COMMENT '跟踪信息(JSON格式)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_notification_id` (`notification_id`),
  KEY `idx_channel_id` (`channel_id`),
  KEY `idx_channel_type` (`channel_type`),
  KEY `idx_status` (`status`),
  KEY `idx_sent_time` (`sent_time`),
  KEY `idx_next_retry_time` (`next_retry_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_notification_delivery_records_channel_id` FOREIGN KEY (`channel_id`) REFERENCES `notification_channels` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知发送记录表';

-- 用户通知偏好表
CREATE TABLE IF NOT EXISTS `user_notification_preferences` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `notification_type` VARCHAR(32) NOT NULL COMMENT '通知类型',
  `channel_type` VARCHAR(32) NOT NULL COMMENT '渠道类型',
  `is_enabled` TINYINT(1) DEFAULT 1 COMMENT '是否启用：0-禁用 1-启用',
  `quiet_start` TIME DEFAULT NULL COMMENT '免打扰开始时间',
  `quiet_end` TIME DEFAULT NULL COMMENT '免打扰结束时间',
  `priority_threshold` INT DEFAULT 0 COMMENT '优先级阈值，只接收高于该阈值的通知',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_type_channel` (`user_id`, `notification_type`, `channel_type`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_notification_type` (`notification_type`),
  KEY `idx_channel_type` (`channel_type`),
  KEY `idx_is_enabled` (`is_enabled`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户通知偏好表';

-- 用户设备表
CREATE TABLE IF NOT EXISTS `user_devices` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `device_id` VARCHAR(128) NOT NULL COMMENT '设备ID',
  `device_type` VARCHAR(32) NOT NULL COMMENT '设备类型：ios, android, web, desktop',
  `device_name` VARCHAR(128) DEFAULT NULL COMMENT '设备名称',
  `device_model` VARCHAR(64) DEFAULT NULL COMMENT '设备型号',
  `os_version` VARCHAR(32) DEFAULT NULL COMMENT '操作系统版本',
  `app_version` VARCHAR(32) DEFAULT NULL COMMENT '应用版本',
  `push_token` VARCHAR(255) DEFAULT NULL COMMENT '推送令牌',
  `push_enabled` TINYINT(1) DEFAULT 1 COMMENT '是否启用推送：0-禁用 1-启用',
  `last_active_time` TIMESTAMP NULL COMMENT '最后活跃时间',
  `status` VARCHAR(32) DEFAULT 'active' COMMENT '状态：active-活跃 inactive-不活跃 disabled-禁用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_device` (`user_id`, `device_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_device_type` (`device_type`),
  KEY `idx_push_enabled` (`push_enabled`),
  KEY `idx_status` (`status`),
  KEY `idx_last_active_time` (`last_active_time`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户设备表';

-- 消息中心统计表
CREATE TABLE IF NOT EXISTS `notification_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `unread_count` INT DEFAULT 0 COMMENT '未读数量',
  `system_unread` INT DEFAULT 0 COMMENT '系统通知未读数量',
  `marketing_unread` INT DEFAULT 0 COMMENT '营销通知未读数量',
  `account_unread` INT DEFAULT 0 COMMENT '账户通知未读数量',
  `transaction_unread` INT DEFAULT 0 COMMENT '交易通知未读数量',
  `social_unread` INT DEFAULT 0 COMMENT '社交通知未读数量',
  `last_read_time` TIMESTAMP NULL COMMENT '最后阅读时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息中心统计表';

-- 通知分类表
CREATE TABLE IF NOT EXISTS `notification_categories` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '分类名称',
  `code` VARCHAR(64) NOT NULL COMMENT '分类编码',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '分类描述',
  `parent_id` BIGINT DEFAULT 0 COMMENT '父分类ID',
  `icon` VARCHAR(255) DEFAULT NULL COMMENT '图标',
  `sort_order` INT DEFAULT 0 COMMENT '排序顺序',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知分类表';

-- 批量通知任务表
CREATE TABLE IF NOT EXISTS `notification_batch_tasks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `task_id` VARCHAR(64) NOT NULL COMMENT '任务ID',
  `name` VARCHAR(128) NOT NULL COMMENT '任务名称',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '任务描述',
  `template_id` BIGINT NOT NULL COMMENT '模板ID',
  `policy_id` BIGINT DEFAULT NULL COMMENT '策略ID',
  `audience` JSON NOT NULL COMMENT '受众定义(JSON格式)',
  `variables` JSON DEFAULT NULL COMMENT '变量数据(JSON格式)',
  `scheduled_time` TIMESTAMP NULL COMMENT '计划发送时间',
  `expire_time` TIMESTAMP NULL COMMENT '过期时间',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-待处理 processing-处理中 completed-已完成 failed-失败 canceled-已取消',
  `progress` INT DEFAULT 0 COMMENT '进度(0-100)',
  `total_count` INT DEFAULT 0 COMMENT '总数量',
  `processed_count` INT DEFAULT 0 COMMENT '已处理数量',
  `success_count` INT DEFAULT 0 COMMENT '成功数量',
  `failure_count` INT DEFAULT 0 COMMENT '失败数量',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建者ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_task_id` (`task_id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_policy_id` (`policy_id`),
  KEY `idx_status` (`status`),
  KEY `idx_scheduled_time` (`scheduled_time`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_notification_batch_tasks_template_id` FOREIGN KEY (`template_id`) REFERENCES `notification_templates` (`id`),
  CONSTRAINT `fk_notification_batch_tasks_policy_id` FOREIGN KEY (`policy_id`) REFERENCES `notification_policies` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='批量通知任务表';

-- 初始化基础通知模板
INSERT INTO `notification_templates` (`name`, `code`, `description`, `type`, `title_template`, `content_template`, `status`) VALUES
('账户注册成功', 'account-register-success', '用户注册成功后发送的通知', 'account', '欢迎加入{{app_name}}', '尊敬的{{user_name}}，恭喜您成功注册{{app_name}}账号，您现在可以使用平台的所有功能了。', 1),
('密码重置', 'password-reset', '密码重置通知', 'account', '密码重置通知', '尊敬的{{user_name}}，您的密码已成功重置。如非本人操作，请立即联系客服。', 1),
('新消息提醒', 'new-message', '收到新消息的通知', 'social', '您收到一条新消息', '{{sender_name}}给您发送了一条新消息：{{message_content}}', 1),
('订单创建通知', 'order-created', '创建新订单的通知', 'transaction', '订单创建成功', '您的订单({{order_number}})已创建成功，订单金额：{{order_amount}}元，请及时支付。', 1),
('订单支付成功', 'order-paid', '订单支付成功通知', 'transaction', '订单支付成功', '您的订单({{order_number}})已支付成功，我们将尽快为您处理。', 1),
('系统公告', 'system-announcement', '系统公告通知', 'system', '系统公告：{{title}}', '{{content}}', 1);

-- 初始化基础通知渠道
INSERT INTO `notification_channels` (`name`, `code`, `description`, `type`, `status`) VALUES
('应用内通知', 'app-notification', '应用内推送通知', 'app', 1),
('邮件通知', 'email-notification', '通过邮件发送通知', 'email', 1),
('短信通知', 'sms-notification', '通过短信发送通知', 'sms', 1),
('移动推送', 'mobile-push', '移动应用推送通知', 'push', 1),
('微信通知', 'wechat-notification', '通过微信发送通知', 'wechat', 1);

-- 初始化基础通知策略
INSERT INTO `notification_policies` (`name`, `code`, `description`, `template_id`, `channels`, `priority`, `status`) VALUES
('账户注册通知策略', 'account-register-policy', '用户注册后通知策略', 1, '[{"channel_id":1,"enabled":true},{"channel_id":2,"enabled":true}]', 1, 1),
('订单通知策略', 'order-notification-policy', '订单相关通知策略', 4, '[{"channel_id":1,"enabled":true},{"channel_id":3,"enabled":true},{"channel_id":4,"enabled":true}]', 2, 1),
('系统通知策略', 'system-notification-policy', '系统通知策略', 6, '[{"channel_id":1,"enabled":true}]', 0, 1);

-- 初始化通知分类
INSERT INTO `notification_categories` (`name`, `code`, `description`, `parent_id`, `sort_order`, `status`) VALUES
('系统通知', 'system-notification', '系统相关的通知', 0, 1, 1),
('账户通知', 'account-notification', '账户相关的通知', 0, 2, 1),
('交易通知', 'transaction-notification', '交易相关的通知', 0, 3, 1),
('社交通知', 'social-notification', '社交互动相关的通知', 0, 4, 1),
('营销通知', 'marketing-notification', '营销活动相关的通知', 0, 5, 1);
