-- 任务管理模块数据表结构

-- 定时任务表
CREATE TABLE IF NOT EXISTS `scheduled_tasks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '任务名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '任务描述',
  `cron_expression` VARCHAR(64) NOT NULL COMMENT 'Cron表达式',
  `job_handler` VARCHAR(128) NOT NULL COMMENT '任务处理器',
  `job_params` TEXT DEFAULT NULL COMMENT '任务参数(JSON格式)',
  `retry_count` INT DEFAULT 0 COMMENT '重试次数',
  `retry_interval` INT DEFAULT 0 COMMENT '重试间隔(秒)',
  `timeout` INT DEFAULT 60 COMMENT '超时时间(秒)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `concurrent` TINYINT(1) DEFAULT 0 COMMENT '是否允许并发执行：0-禁止 1-允许',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示系统任务',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name_tenant` (`name`, `tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='定时任务表';

-- 任务执行日志表
CREATE TABLE IF NOT EXISTS `task_execution_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `task_id` BIGINT NOT NULL COMMENT '任务ID',
  `task_name` VARCHAR(128) NOT NULL COMMENT '任务名称',
  `execution_time` TIMESTAMP NOT NULL COMMENT '执行时间',
  `complete_time` TIMESTAMP NULL COMMENT '完成时间',
  `execution_status` INT DEFAULT 0 COMMENT '执行状态：0-执行中 1-成功 2-失败 3-超时',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `retry_count` INT DEFAULT 0 COMMENT '重试次数',
  `trigger_type` VARCHAR(32) DEFAULT 'cron' COMMENT '触发类型：cron-定时触发 manual-手动触发',
  `result` TEXT DEFAULT NULL COMMENT '执行结果',
  `duration` BIGINT DEFAULT NULL COMMENT '执行时长(毫秒)',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_execution_status` (`execution_status`),
  KEY `idx_execution_time` (`execution_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_task_logs_task_id` FOREIGN KEY (`task_id`) REFERENCES `scheduled_tasks` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务执行日志表';

-- 系统通知表
CREATE TABLE IF NOT EXISTS `system_notifications` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(128) NOT NULL COMMENT '通知标题',
  `content` TEXT NOT NULL COMMENT '通知内容',
  `type` VARCHAR(32) NOT NULL COMMENT '通知类型：system-系统通知 personal-个人通知 broadcast-广播通知',
  `sender_id` BIGINT DEFAULT NULL COMMENT '发送者ID',
  `sender_name` VARCHAR(64) DEFAULT NULL COMMENT '发送者名称',
  `sender_type` VARCHAR(32) DEFAULT 'system' COMMENT '发送者类型：system-系统 admin-管理员 user-用户',
  `importance` INT DEFAULT 0 COMMENT '重要性：0-普通 1-重要 2-紧急',
  `status` INT DEFAULT 0 COMMENT '状态：0-待发送 1-已发送 2-发送失败',
  `scheduled_time` TIMESTAMP NULL COMMENT '计划发送时间，为NULL表示立即发送',
  `expire_time` TIMESTAMP NULL COMMENT '过期时间，为NULL表示永不过期',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示全局通知',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_importance` (`importance`),
  KEY `idx_scheduled_time` (`scheduled_time`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统通知表';

-- 用户通知表
CREATE TABLE IF NOT EXISTS `user_notifications` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `notification_id` BIGINT NOT NULL COMMENT '通知ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `is_read` TINYINT(1) DEFAULT 0 COMMENT '是否已读：0-未读 1-已读',
  `read_time` TIMESTAMP NULL COMMENT '阅读时间',
  `is_deleted` TINYINT(1) DEFAULT 0 COMMENT '是否删除：0-未删除 1-已删除',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_notification_user` (`notification_id`, `user_id`),
  KEY `idx_notification_id` (`notification_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_is_deleted` (`is_deleted`),
  CONSTRAINT `fk_user_notifications_notification_id` FOREIGN KEY (`notification_id`) REFERENCES `system_notifications` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户通知表';

-- 通知接收者组表
CREATE TABLE IF NOT EXISTS `notification_receivers` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `notification_id` BIGINT NOT NULL COMMENT '通知ID',
  `receiver_type` VARCHAR(32) NOT NULL COMMENT '接收者类型：all-所有用户 role-角色 tenant-租户 department-部门 user-指定用户',
  `receiver_id` VARCHAR(128) DEFAULT NULL COMMENT '接收者ID，可以是角色ID、租户ID、部门ID或用户ID',
  `receiver_params` JSON DEFAULT NULL COMMENT '接收者参数，用于存储复杂的接收者规则',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_notification_id` (`notification_id`),
  KEY `idx_receiver_type` (`receiver_type`),
  KEY `idx_receiver_id` (`receiver_id`),
  CONSTRAINT `fk_notification_receivers_notification_id` FOREIGN KEY (`notification_id`) REFERENCES `system_notifications` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知接收者组表';

-- 工作流定义表
CREATE TABLE IF NOT EXISTS `workflow_definitions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(64) NOT NULL COMMENT '工作流编码',
  `name` VARCHAR(128) NOT NULL COMMENT '工作流名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '工作流描述',
  `version` INT NOT NULL DEFAULT 1 COMMENT '版本号',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用 2-草稿',
  `definition` TEXT NOT NULL COMMENT '工作流定义(JSON格式)',
  `form_definition` TEXT DEFAULT NULL COMMENT '表单定义(JSON格式)',
  `category` VARCHAR(64) DEFAULT 'default' COMMENT '分类',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示系统工作流',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建者ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_version_tenant` (`code`, `version`, `tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_category` (`category`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流定义表';

-- 工作流实例表
CREATE TABLE IF NOT EXISTS `workflow_instances` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `definition_id` BIGINT NOT NULL COMMENT '工作流定义ID',
  `instance_code` VARCHAR(64) NOT NULL COMMENT '实例编码',
  `title` VARCHAR(255) NOT NULL COMMENT '实例标题',
  `business_key` VARCHAR(128) DEFAULT NULL COMMENT '业务键',
  `initiator_id` BIGINT NOT NULL COMMENT '发起人ID',
  `initiator_name` VARCHAR(64) DEFAULT NULL COMMENT '发起人名称',
  `status` VARCHAR(32) NOT NULL DEFAULT 'running' COMMENT '状态：draft-草稿 running-运行中 completed-已完成 terminated-已终止 canceled-已取消',
  `start_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `duration` BIGINT DEFAULT NULL COMMENT '持续时间(毫秒)',
  `form_data` TEXT DEFAULT NULL COMMENT '表单数据(JSON格式)',
  `variables` JSON DEFAULT NULL COMMENT '流程变量',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `parent_instance_id` BIGINT DEFAULT NULL COMMENT '父实例ID，用于子流程',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_instance_code` (`instance_code`),
  KEY `idx_definition_id` (`definition_id`),
  KEY `idx_business_key` (`business_key`),
  KEY `idx_initiator_id` (`initiator_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_parent_instance_id` (`parent_instance_id`),
  CONSTRAINT `fk_workflow_instances_definition_id` FOREIGN KEY (`definition_id`) REFERENCES `workflow_definitions` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流实例表';

-- 工作流任务表
CREATE TABLE IF NOT EXISTS `workflow_tasks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `instance_id` BIGINT NOT NULL COMMENT '实例ID',
  `definition_id` BIGINT NOT NULL COMMENT '工作流定义ID',
  `task_code` VARCHAR(64) NOT NULL COMMENT '任务编码',
  `task_name` VARCHAR(128) NOT NULL COMMENT '任务名称',
  `task_type` VARCHAR(32) NOT NULL COMMENT '任务类型：userTask-用户任务 serviceTask-服务任务 scriptTask-脚本任务 等',
  `status` VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '状态：pending-待处理 claimed-已认领 completed-已完成 canceled-已取消',
  `assignee_id` BIGINT DEFAULT NULL COMMENT '受理人ID',
  `assignee_name` VARCHAR(64) DEFAULT NULL COMMENT '受理人名称',
  `assignee_type` VARCHAR(32) DEFAULT NULL COMMENT '受理人类型：user-用户 role-角色 department-部门',
  `priority` INT DEFAULT 0 COMMENT '优先级：0-普通 1-重要 2-紧急',
  `due_date` TIMESTAMP NULL COMMENT '截止日期',
  `claim_time` TIMESTAMP NULL COMMENT '认领时间',
  `start_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `duration` BIGINT DEFAULT NULL COMMENT '持续时间(毫秒)',
  `form_key` VARCHAR(128) DEFAULT NULL COMMENT '表单标识',
  `form_data` TEXT DEFAULT NULL COMMENT '表单数据(JSON格式)',
  `action_result` VARCHAR(32) DEFAULT NULL COMMENT '操作结果：approve-同意 reject-拒绝 transfer-转办 delegate-委派 等',
  `comment` TEXT DEFAULT NULL COMMENT '处理意见',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_task_code` (`task_code`),
  KEY `idx_instance_id` (`instance_id`),
  KEY `idx_definition_id` (`definition_id`),
  KEY `idx_status` (`status`),
  KEY `idx_assignee_id` (`assignee_id`),
  KEY `idx_priority` (`priority`),
  KEY `idx_due_date` (`due_date`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_workflow_tasks_instance_id` FOREIGN KEY (`instance_id`) REFERENCES `workflow_instances` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_workflow_tasks_definition_id` FOREIGN KEY (`definition_id`) REFERENCES `workflow_definitions` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流任务表';

-- 工作流历史记录表
CREATE TABLE IF NOT EXISTS `workflow_histories` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `instance_id` BIGINT NOT NULL COMMENT '实例ID',
  `task_id` BIGINT DEFAULT NULL COMMENT '任务ID',
  `node_id` VARCHAR(64) NOT NULL COMMENT '节点ID',
  `node_name` VARCHAR(128) NOT NULL COMMENT '节点名称',
  `node_type` VARCHAR(32) NOT NULL COMMENT '节点类型',
  `action` VARCHAR(32) NOT NULL COMMENT '操作：start-开始 complete-完成 claim-认领 transfer-转办 delegate-委派 等',
  `action_result` VARCHAR(32) DEFAULT NULL COMMENT '操作结果',
  `operator_id` BIGINT DEFAULT NULL COMMENT '操作人ID',
  `operator_name` VARCHAR(64) DEFAULT NULL COMMENT '操作人名称',
  `comment` TEXT DEFAULT NULL COMMENT '操作意见',
  `data` JSON DEFAULT NULL COMMENT '数据(JSON格式)',
  `timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_instance_id` (`instance_id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_node_id` (`node_id`),
  KEY `idx_operator_id` (`operator_id`),
  KEY `idx_timestamp` (`timestamp`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_workflow_histories_instance_id` FOREIGN KEY (`instance_id`) REFERENCES `workflow_instances` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流历史记录表';

-- 初始化基本定时任务
INSERT INTO `scheduled_tasks` (`name`, `description`, `cron_expression`, `job_handler`, `retry_count`, `status`, `concurrent`) VALUES
('系统状态检查', '定期检查系统状态和性能', '0 0 * * * ?', 'systemStatusCheckJob', 3, 1, 0),
('数据库备份', '定期备份数据库', '0 0 2 * * ?', 'databaseBackupJob', 2, 1, 0),
('临时文件清理', '清理临时文件和缓存', '0 0 3 * * ?', 'tempFileCleanupJob', 1, 1, 0),
('失效TOKEN清理', '清理过期的JWT黑名单记录', '0 0 1 * * ?', 'expiredTokenCleanupJob', 1, 1, 0);
