-- 系统运维、安全审计和日志监控表结构

-- 系统监控指标表
CREATE TABLE IF NOT EXISTS `system_metrics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `metric_name` VARCHAR(64) NOT NULL COMMENT '指标名称',
  `metric_type` VARCHAR(16) NOT NULL COMMENT '指标类型：counter-计数器 gauge-仪表盘 histogram-直方图',
  `value` DECIMAL(18,4) NOT NULL COMMENT '指标值',
  `unit` VARCHAR(16) DEFAULT NULL COMMENT '单位',
  `labels` JSON DEFAULT NULL COMMENT '标签，JSON格式',
  `timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间戳',
  `hostname` VARCHAR(64) DEFAULT NULL COMMENT '主机名',
  `instance_id` VARCHAR(64) DEFAULT NULL COMMENT '实例ID',
  `service_name` VARCHAR(64) DEFAULT NULL COMMENT '服务名称',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_metric_name` (`metric_name`),
  KEY `idx_metric_type` (`metric_type`),
  KEY `idx_timestamp` (`timestamp`),
  KEY `idx_hostname` (`hostname`),
  KEY `idx_service_name` (`service_name`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统监控指标表';

-- 系统告警规则表
CREATE TABLE IF NOT EXISTS `alert_rules` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `rule_name` VARCHAR(64) NOT NULL COMMENT '规则名称',
  `rule_description` VARCHAR(255) DEFAULT NULL COMMENT '规则描述',
  `metric_name` VARCHAR(64) NOT NULL COMMENT '指标名称',
  `condition` VARCHAR(32) NOT NULL COMMENT '条件：gt-大于 lt-小于 eq-等于 neq-不等于',
  `threshold` DECIMAL(18,4) NOT NULL COMMENT '阈值',
  `duration` INT NOT NULL DEFAULT 0 COMMENT '持续时间(秒)',
  `severity` VARCHAR(16) NOT NULL COMMENT '严重程度：info-信息 warning-警告 error-错误 critical-严重',
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `labels` JSON DEFAULT NULL COMMENT '标签，JSON格式',
  `notification_channels` JSON DEFAULT NULL COMMENT '通知渠道，JSON格式',
  `silence_period` INT NOT NULL DEFAULT 300 COMMENT '静默期(秒)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_rule_name` (`tenant_id`, `rule_name`),
  KEY `idx_metric_name` (`metric_name`),
  KEY `idx_severity` (`severity`),
  KEY `idx_is_enabled` (`is_enabled`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统告警规则表';

-- 告警事件表
CREATE TABLE IF NOT EXISTS `alert_events` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `rule_id` BIGINT NOT NULL COMMENT '规则ID',
  `event_id` VARCHAR(64) NOT NULL COMMENT '事件ID',
  `metric_name` VARCHAR(64) NOT NULL COMMENT '指标名称',
  `metric_value` DECIMAL(18,4) NOT NULL COMMENT '指标值',
  `threshold` DECIMAL(18,4) NOT NULL COMMENT '阈值',
  `severity` VARCHAR(16) NOT NULL COMMENT '严重程度：info-信息 warning-警告 error-错误 critical-严重',
  `status` VARCHAR(16) NOT NULL COMMENT '状态：firing-触发中 resolved-已解决 silenced-已静默',
  `start_time` TIMESTAMP NOT NULL COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `duration` INT DEFAULT NULL COMMENT '持续时间(秒)',
  `labels` JSON DEFAULT NULL COMMENT '标签，JSON格式',
  `annotations` JSON DEFAULT NULL COMMENT '注释，JSON格式',
  `hostname` VARCHAR(64) DEFAULT NULL COMMENT '主机名',
  `service_name` VARCHAR(64) DEFAULT NULL COMMENT '服务名称',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_event_id` (`event_id`),
  KEY `idx_rule_id` (`rule_id`),
  KEY `idx_metric_name` (`metric_name`),
  KEY `idx_severity` (`severity`),
  KEY `idx_status` (`status`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_end_time` (`end_time`),
  KEY `idx_hostname` (`hostname`),
  KEY `idx_service_name` (`service_name`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_alert_events_rule_id` FOREIGN KEY (`rule_id`) REFERENCES `alert_rules` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警事件表';

-- 服务健康检查表
CREATE TABLE IF NOT EXISTS `service_health_checks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `service_name` VARCHAR(64) NOT NULL COMMENT '服务名称',
  `instance_id` VARCHAR(64) NOT NULL COMMENT '实例ID',
  `check_type` VARCHAR(16) NOT NULL COMMENT '检查类型：http-HTTP ping-PING tcp-TCP',
  `endpoint` VARCHAR(255) NOT NULL COMMENT '检查端点',
  `interval` INT NOT NULL DEFAULT 60 COMMENT '检查间隔(秒)',
  `timeout` INT NOT NULL DEFAULT 5 COMMENT '超时时间(秒)',
  `retry_count` INT NOT NULL DEFAULT 3 COMMENT '重试次数',
  `expected_response` TEXT DEFAULT NULL COMMENT '期望响应',
  `status` VARCHAR(16) NOT NULL COMMENT '状态：up-正常 down-异常 warning-警告',
  `last_check_time` TIMESTAMP NULL COMMENT '最后检查时间',
  `next_check_time` TIMESTAMP NULL COMMENT '下次检查时间',
  `consecutive_failures` INT NOT NULL DEFAULT 0 COMMENT '连续失败次数',
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_service_instance` (`service_name`, `instance_id`),
  KEY `idx_service_name` (`service_name`),
  KEY `idx_status` (`status`),
  KEY `idx_is_enabled` (`is_enabled`),
  KEY `idx_last_check_time` (`last_check_time`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务健康检查表';

-- 健康检查日志表
CREATE TABLE IF NOT EXISTS `health_check_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `check_id` BIGINT NOT NULL COMMENT '健康检查ID',
  `service_name` VARCHAR(64) NOT NULL COMMENT '服务名称',
  `instance_id` VARCHAR(64) NOT NULL COMMENT '实例ID',
  `status` VARCHAR(16) NOT NULL COMMENT '状态：success-成功 failure-失败 timeout-超时',
  `response_time` INT DEFAULT NULL COMMENT '响应时间(ms)',
  `response_code` VARCHAR(16) DEFAULT NULL COMMENT '响应码',
  `response_body` TEXT DEFAULT NULL COMMENT '响应内容',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `check_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '检查时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_check_id` (`check_id`),
  KEY `idx_service_name` (`service_name`),
  KEY `idx_instance_id` (`instance_id`),
  KEY `idx_status` (`status`),
  KEY `idx_check_time` (`check_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_health_check_logs_check_id` FOREIGN KEY (`check_id`) REFERENCES `service_health_checks` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='健康检查日志表';

-- 安全审计日志表
CREATE TABLE IF NOT EXISTS `security_audit_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `username` VARCHAR(64) DEFAULT NULL COMMENT '用户名',
  `client_ip` VARCHAR(64) NOT NULL COMMENT '客户端IP',
  `user_agent` VARCHAR(255) DEFAULT NULL COMMENT '用户代理',
  `event_type` VARCHAR(32) NOT NULL COMMENT '事件类型：login-登录 logout-登出 permission-权限 data-数据',
  `action` VARCHAR(32) NOT NULL COMMENT '动作：success-成功 failure-失败 deny-拒绝 allow-允许',
  `resource_type` VARCHAR(32) DEFAULT NULL COMMENT '资源类型',
  `resource_id` VARCHAR(64) DEFAULT NULL COMMENT '资源ID',
  `details` TEXT DEFAULT NULL COMMENT '详细信息',
  `session_id` VARCHAR(64) DEFAULT NULL COMMENT '会话ID',
  `request_id` VARCHAR(64) DEFAULT NULL COMMENT '请求ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_client_ip` (`client_ip`),
  KEY `idx_event_type` (`event_type`),
  KEY `idx_action` (`action`),
  KEY `idx_resource` (`resource_type`, `resource_id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_request_id` (`request_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='安全审计日志表';

-- 异常登录记录表
CREATE TABLE IF NOT EXISTS `abnormal_logins` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `username` VARCHAR(64) NOT NULL COMMENT '用户名',
  `login_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `client_ip` VARCHAR(64) NOT NULL COMMENT '客户端IP',
  `geo_location` VARCHAR(128) DEFAULT NULL COMMENT '地理位置',
  `user_agent` VARCHAR(255) DEFAULT NULL COMMENT '用户代理',
  `device_id` VARCHAR(64) DEFAULT NULL COMMENT '设备ID',
  `risk_level` VARCHAR(16) NOT NULL COMMENT '风险等级：low-低 medium-中 high-高',
  `risk_factors` JSON DEFAULT NULL COMMENT '风险因素，JSON格式',
  `detection_rules` JSON DEFAULT NULL COMMENT '检测规则，JSON格式',
  `action_taken` VARCHAR(32) NOT NULL COMMENT '采取的措施：none-无 notify-通知 block-阻止 verify-二次验证',
  `is_false_positive` TINYINT(1) DEFAULT NULL COMMENT '是否误报',
  `analyst_notes` TEXT DEFAULT NULL COMMENT '分析师备注',
  `resolution_time` TIMESTAMP NULL COMMENT '解决时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_login_time` (`login_time`),
  KEY `idx_client_ip` (`client_ip`),
  KEY `idx_risk_level` (`risk_level`),
  KEY `idx_action_taken` (`action_taken`),
  KEY `idx_is_false_positive` (`is_false_positive`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_abnormal_logins_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='异常登录记录表';

-- 系统安全配置表
CREATE TABLE IF NOT EXISTS `security_settings` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `setting_key` VARCHAR(64) NOT NULL COMMENT '配置键',
  `setting_value` TEXT NOT NULL COMMENT '配置值',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `setting_group` VARCHAR(32) NOT NULL DEFAULT 'default' COMMENT '配置组',
  `is_encrypted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否加密',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_by` BIGINT DEFAULT NULL COMMENT '更新人',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_key` (`tenant_id`, `setting_key`),
  KEY `idx_setting_group` (`setting_group`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_updated_by` (`updated_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统安全配置表';

-- API调用日志表
CREATE TABLE IF NOT EXISTS `api_request_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `trace_id` VARCHAR(64) NOT NULL COMMENT '跟踪ID',
  `service_name` VARCHAR(64) NOT NULL COMMENT '服务名称',
  `api_path` VARCHAR(255) NOT NULL COMMENT 'API路径',
  `method` VARCHAR(10) NOT NULL COMMENT '请求方法',
  `request_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '请求时间',
  `response_time` INT NOT NULL COMMENT '响应时间(ms)',
  `status_code` INT NOT NULL COMMENT '状态码',
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `client_ip` VARCHAR(64) DEFAULT NULL COMMENT '客户端IP',
  `user_agent` VARCHAR(255) DEFAULT NULL COMMENT '用户代理',
  `request_headers` TEXT DEFAULT NULL COMMENT '请求头',
  `request_body` TEXT DEFAULT NULL COMMENT '请求体',
  `response_body` TEXT DEFAULT NULL COMMENT '响应体',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_trace_id` (`trace_id`),
  KEY `idx_service_name` (`service_name`),
  KEY `idx_api_path` (`api_path`),
  KEY `idx_method` (`method`),
  KEY `idx_request_time` (`request_time`),
  KEY `idx_status_code` (`status_code`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_client_ip` (`client_ip`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API调用日志表';

-- 数据访问控制表
CREATE TABLE IF NOT EXISTS `data_access_controls` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `role_id` BIGINT NOT NULL COMMENT '角色ID',
  `resource_type` VARCHAR(32) NOT NULL COMMENT '资源类型',
  `resource_field` VARCHAR(64) DEFAULT NULL COMMENT '资源字段',
  `access_type` VARCHAR(16) NOT NULL COMMENT '访问类型：read-读 write-写 none-无',
  `condition_expression` TEXT DEFAULT NULL COMMENT '条件表达式',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_resource_field` (`role_id`, `resource_type`, `resource_field`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_resource_type` (`resource_type`),
  KEY `idx_access_type` (`access_type`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_data_access_controls_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据访问控制表';

-- 系统性能指标表
CREATE TABLE IF NOT EXISTS `performance_metrics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `service_name` VARCHAR(64) NOT NULL COMMENT '服务名称',
  `instance_id` VARCHAR(64) NOT NULL COMMENT '实例ID',
  `metric_type` VARCHAR(32) NOT NULL COMMENT '指标类型：cpu-CPU使用率 memory-内存使用率 disk-磁盘使用率 network-网络流量 response_time-响应时间',
  `value` DECIMAL(18,4) NOT NULL COMMENT '指标值',
  `unit` VARCHAR(16) DEFAULT NULL COMMENT '单位',
  `collect_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '采集时间',
  `hostname` VARCHAR(64) DEFAULT NULL COMMENT '主机名',
  `labels` JSON DEFAULT NULL COMMENT '标签，JSON格式',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_service_instance` (`service_name`, `instance_id`),
  KEY `idx_metric_type` (`metric_type`),
  KEY `idx_collect_time` (`collect_time`),
  KEY `idx_hostname` (`hostname`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统性能指标表';

-- 系统任务执行表
CREATE TABLE IF NOT EXISTS `task_executions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `task_id` VARCHAR(64) NOT NULL COMMENT '任务ID',
  `task_name` VARCHAR(64) NOT NULL COMMENT '任务名称',
  `task_type` VARCHAR(32) NOT NULL COMMENT '任务类型：cron-定时任务 once-一次性任务 stream-流处理任务',
  `status` VARCHAR(16) NOT NULL COMMENT '状态：pending-待执行 running-执行中 success-成功 failure-失败 cancelled-已取消',
  `start_time` TIMESTAMP NULL COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `duration` INT DEFAULT NULL COMMENT '执行时长(ms)',
  `result` TEXT DEFAULT NULL COMMENT '执行结果',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `worker_id` VARCHAR(64) DEFAULT NULL COMMENT '执行节点ID',
  `retry_count` INT NOT NULL DEFAULT 0 COMMENT '重试次数',
  `next_retry_time` TIMESTAMP NULL COMMENT '下次重试时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_task_id` (`task_id`),
  KEY `idx_task_name` (`task_name`),
  KEY `idx_task_type` (`task_type`),
  KEY `idx_status` (`status`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_worker_id` (`worker_id`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统任务执行表';
