-- 缓存管理和服务监控模块数据表结构

-- 缓存配置表
CREATE TABLE IF NOT EXISTS `cache_configs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '缓存名称',
  `code` VARCHAR(64) NOT NULL COMMENT '缓存编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '缓存描述',
  `cache_type` VARCHAR(32) DEFAULT 'memory' COMMENT '缓存类型：memory-内存 redis-Redis memcached-Memcached',
  `connection_params` JSON DEFAULT NULL COMMENT '连接参数(JSON格式)',
  `ttl` INT DEFAULT 3600 COMMENT '默认过期时间(秒)',
  `max_size` INT DEFAULT 10000 COMMENT '最大缓存项数，仅适用于内存缓存',
  `max_memory` INT DEFAULT 0 COMMENT '最大内存(MB)，0表示不限制',
  `cleanup_interval` INT DEFAULT 60 COMMENT '清理间隔(秒)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示系统缓存',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_cache_type` (`cache_type`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='缓存配置表';

-- 缓存键规则表
CREATE TABLE IF NOT EXISTS `cache_key_rules` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `cache_id` BIGINT NOT NULL COMMENT '缓存ID',
  `name` VARCHAR(128) NOT NULL COMMENT '规则名称',
  `key_pattern` VARCHAR(255) NOT NULL COMMENT '键模式，支持通配符*',
  `ttl` INT DEFAULT NULL COMMENT '过期时间(秒)，为NULL使用缓存默认值',
  `max_size` INT DEFAULT NULL COMMENT '最大项数，为NULL使用缓存默认值',
  `priority` INT DEFAULT 0 COMMENT '优先级，数字越大优先级越高',
  `invalidation_strategy` VARCHAR(32) DEFAULT 'expire' COMMENT '失效策略：expire-自动过期 manual-手动失效 update-更新时失效',
  `invalidation_tags` VARCHAR(255) DEFAULT NULL COMMENT '失效标签，多个用逗号分隔',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_cache_id` (`cache_id`),
  KEY `idx_priority` (`priority`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_cache_key_rules_cache_id` FOREIGN KEY (`cache_id`) REFERENCES `cache_configs` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='缓存键规则表';

-- 缓存统计表
CREATE TABLE IF NOT EXISTS `cache_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `cache_id` BIGINT NOT NULL COMMENT '缓存ID',
  `key_pattern` VARCHAR(255) DEFAULT NULL COMMENT '键模式，为NULL表示整个缓存',
  `total_keys` INT DEFAULT 0 COMMENT '键总数',
  `hit_count` BIGINT DEFAULT 0 COMMENT '命中次数',
  `miss_count` BIGINT DEFAULT 0 COMMENT '未命中次数',
  `hit_rate` FLOAT DEFAULT 0 COMMENT '命中率',
  `avg_access_time` FLOAT DEFAULT 0 COMMENT '平均访问时间(毫秒)',
  `eviction_count` INT DEFAULT 0 COMMENT '驱逐次数',
  `memory_usage` BIGINT DEFAULT 0 COMMENT '内存使用量(字节)',
  `stat_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '统计时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_cache_id` (`cache_id`),
  KEY `idx_stat_time` (`stat_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_cache_statistics_cache_id` FOREIGN KEY (`cache_id`) REFERENCES `cache_configs` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='缓存统计表';

-- 缓存操作日志表
CREATE TABLE IF NOT EXISTS `cache_operation_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `cache_id` BIGINT NOT NULL COMMENT '缓存ID',
  `operation_type` VARCHAR(32) NOT NULL COMMENT '操作类型：get-获取 set-设置 delete-删除 clear-清空 flush-刷新',
  `key` VARCHAR(255) DEFAULT NULL COMMENT '缓存键',
  `key_pattern` VARCHAR(255) DEFAULT NULL COMMENT '键模式，用于批量操作',
  `value_size` INT DEFAULT NULL COMMENT '值大小(字节)',
  `status` VARCHAR(32) DEFAULT 'success' COMMENT '状态：success-成功 failed-失败',
  `error_message` VARCHAR(255) DEFAULT NULL COMMENT '错误信息',
  `operation_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  `duration` INT DEFAULT NULL COMMENT '操作耗时(毫秒)',
  `operator_id` BIGINT DEFAULT NULL COMMENT '操作者ID',
  `operator_type` VARCHAR(32) DEFAULT NULL COMMENT '操作者类型：user-用户 system-系统 service-服务',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_cache_id` (`cache_id`),
  KEY `idx_operation_type` (`operation_type`),
  KEY `idx_operation_time` (`operation_time`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_cache_operation_logs_cache_id` FOREIGN KEY (`cache_id`) REFERENCES `cache_configs` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='缓存操作日志表';

-- 服务节点表
CREATE TABLE IF NOT EXISTS `service_nodes` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `node_id` VARCHAR(64) NOT NULL COMMENT '节点ID',
  `service_name` VARCHAR(128) NOT NULL COMMENT '服务名称',
  `instance_id` VARCHAR(64) NOT NULL COMMENT '实例ID',
  `host` VARCHAR(128) NOT NULL COMMENT '主机地址',
  `port` INT NOT NULL COMMENT '端口号',
  `scheme` VARCHAR(16) DEFAULT 'http' COMMENT '协议：http, https, etc.',
  `metadata` JSON DEFAULT NULL COMMENT '元数据(JSON格式)',
  `weight` INT DEFAULT 100 COMMENT '权重',
  `status` VARCHAR(32) DEFAULT 'online' COMMENT '状态：online-在线 offline-离线 unhealthy-不健康',
  `last_heartbeat` TIMESTAMP NULL COMMENT '最后心跳时间',
  `registration_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  `version` VARCHAR(32) DEFAULT NULL COMMENT '版本',
  `environment` VARCHAR(32) DEFAULT 'production' COMMENT '环境：production-生产 testing-测试 development-开发',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_node_id` (`node_id`),
  UNIQUE KEY `idx_service_instance` (`service_name`, `instance_id`),
  KEY `idx_service_name` (`service_name`),
  KEY `idx_status` (`status`),
  KEY `idx_last_heartbeat` (`last_heartbeat`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务节点表';

-- 服务健康检查表
CREATE TABLE IF NOT EXISTS `service_health_checks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `node_id` VARCHAR(64) NOT NULL COMMENT '节点ID',
  `check_id` VARCHAR(64) NOT NULL COMMENT '检查ID',
  `name` VARCHAR(128) NOT NULL COMMENT '检查名称',
  `type` VARCHAR(32) NOT NULL COMMENT '检查类型：http-HTTP请求 tcp-TCP连接 script-脚本 ttl-TTL grpc-gRPC',
  `target` VARCHAR(255) NOT NULL COMMENT '检查目标：URL、TCP地址或脚本',
  `interval` INT DEFAULT 30 COMMENT '检查间隔(秒)',
  `timeout` INT DEFAULT 5 COMMENT '超时时间(秒)',
  `deregister_after` INT DEFAULT 300 COMMENT '注销时间(秒)，超过该时间未恢复则注销服务',
  `status` VARCHAR(32) DEFAULT 'passing' COMMENT '状态：passing-通过 warning-警告 critical-严重 unknown-未知',
  `output` TEXT DEFAULT NULL COMMENT '检查输出',
  `last_check_time` TIMESTAMP NULL COMMENT '最后检查时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_node_check` (`node_id`, `check_id`),
  KEY `idx_node_id` (`node_id`),
  KEY `idx_status` (`status`),
  KEY `idx_last_check_time` (`last_check_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务健康检查表';

-- 服务指标表
CREATE TABLE IF NOT EXISTS `service_metrics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `node_id` VARCHAR(64) NOT NULL COMMENT '节点ID',
  `metric_name` VARCHAR(64) NOT NULL COMMENT '指标名称',
  `metric_value` DOUBLE NOT NULL COMMENT '指标值',
  `metric_type` VARCHAR(32) DEFAULT 'gauge' COMMENT '指标类型：gauge-仪表 counter-计数器 histogram-直方图 summary-摘要 timer-计时器',
  `unit` VARCHAR(32) DEFAULT NULL COMMENT '单位',
  `dimensions` JSON DEFAULT NULL COMMENT '维度信息(JSON格式)',
  `timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间戳',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_node_id` (`node_id`),
  KEY `idx_metric_name` (`metric_name`),
  KEY `idx_timestamp` (`timestamp`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务指标表';

-- 告警规则表
CREATE TABLE IF NOT EXISTS `alert_rules` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '规则名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '规则描述',
  `service_name` VARCHAR(128) DEFAULT NULL COMMENT '服务名称，为NULL表示适用于所有服务',
  `metric_name` VARCHAR(64) NOT NULL COMMENT '指标名称',
  `condition` VARCHAR(255) NOT NULL COMMENT '告警条件，例如：>90, <10, ==0',
  `duration` INT DEFAULT 60 COMMENT '持续时间(秒)，满足条件持续该时间后触发告警',
  `severity` VARCHAR(32) DEFAULT 'warning' COMMENT '严重性：info-信息 warning-警告 error-错误 critical-严重',
  `notification_channels` VARCHAR(255) DEFAULT NULL COMMENT '通知渠道，多个用逗号分隔',
  `message_template` TEXT DEFAULT NULL COMMENT '消息模板',
  `recovery_condition` VARCHAR(255) DEFAULT NULL COMMENT '恢复条件，为NULL时使用默认恢复逻辑',
  `silenced` TINYINT(1) DEFAULT 0 COMMENT '是否静默：0-否 1-是',
  `silence_until` TIMESTAMP NULL COMMENT '静默截止时间',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_service_name` (`service_name`),
  KEY `idx_metric_name` (`metric_name`),
  KEY `idx_severity` (`severity`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警规则表';

-- 告警记录表
CREATE TABLE IF NOT EXISTS `alert_records` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `alert_id` VARCHAR(64) NOT NULL COMMENT '告警ID',
  `rule_id` BIGINT NOT NULL COMMENT '规则ID',
  `service_name` VARCHAR(128) DEFAULT NULL COMMENT '服务名称',
  `node_id` VARCHAR(64) DEFAULT NULL COMMENT '节点ID',
  `metric_name` VARCHAR(64) NOT NULL COMMENT '指标名称',
  `metric_value` DOUBLE NOT NULL COMMENT '指标值',
  `condition` VARCHAR(255) NOT NULL COMMENT '告警条件',
  `status` VARCHAR(32) DEFAULT 'firing' COMMENT '状态：firing-触发中 resolved-已解决 acknowledged-已确认',
  `severity` VARCHAR(32) DEFAULT 'warning' COMMENT '严重性：info-信息 warning-警告 error-错误 critical-严重',
  `summary` VARCHAR(255) NOT NULL COMMENT '告警摘要',
  `description` TEXT DEFAULT NULL COMMENT '告警描述',
  `start_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `duration` INT DEFAULT NULL COMMENT '持续时间(秒)',
  `notification_sent` TINYINT(1) DEFAULT 0 COMMENT '是否已发送通知：0-否 1-是',
  `acknowledged_by` BIGINT DEFAULT NULL COMMENT '确认者ID',
  `acknowledged_time` TIMESTAMP NULL COMMENT '确认时间',
  `acknowledge_comment` TEXT DEFAULT NULL COMMENT '确认评论',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_alert_id` (`alert_id`),
  KEY `idx_rule_id` (`rule_id`),
  KEY `idx_service_name` (`service_name`),
  KEY `idx_node_id` (`node_id`),
  KEY `idx_status` (`status`),
  KEY `idx_severity` (`severity`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_alert_records_rule_id` FOREIGN KEY (`rule_id`) REFERENCES `alert_rules` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警记录表';

-- 监控仪表板表
CREATE TABLE IF NOT EXISTS `monitoring_dashboards` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '仪表板名称',
  `code` VARCHAR(64) NOT NULL COMMENT '仪表板编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '仪表板描述',
  `layout` TEXT DEFAULT NULL COMMENT '布局配置(JSON格式)',
  `is_system` TINYINT(1) DEFAULT 0 COMMENT '是否系统仪表板：0-否 1-是',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建者ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_is_system` (`is_system`),
  KEY `idx_status` (`status`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='监控仪表板表';

-- 监控面板表
CREATE TABLE IF NOT EXISTS `monitoring_panels` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `dashboard_id` BIGINT NOT NULL COMMENT '仪表板ID',
  `name` VARCHAR(128) NOT NULL COMMENT '面板名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '面板描述',
  `panel_type` VARCHAR(32) NOT NULL COMMENT '面板类型：chart-图表 table-表格 gauge-仪表盘 stat-统计 text-文本 etc.',
  `metrics` TEXT NOT NULL COMMENT '指标配置(JSON格式)',
  `visualization` JSON DEFAULT NULL COMMENT '可视化配置(JSON格式)',
  `position` JSON NOT NULL COMMENT '位置配置(JSON格式)',
  `datasource` VARCHAR(128) DEFAULT NULL COMMENT '数据源',
  `refresh_interval` INT DEFAULT 60 COMMENT '刷新间隔(秒)',
  `time_range` VARCHAR(64) DEFAULT 'last_1_hour' COMMENT '时间范围：last_5_minutes, last_1_hour, last_1_day, etc.',
  `sort_order` INT DEFAULT 0 COMMENT '排序顺序',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_dashboard_id` (`dashboard_id`),
  KEY `idx_panel_type` (`panel_type`),
  KEY `idx_sort_order` (`sort_order`),
  CONSTRAINT `fk_monitoring_panels_dashboard_id` FOREIGN KEY (`dashboard_id`) REFERENCES `monitoring_dashboards` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='监控面板表';

-- 日志配置表
CREATE TABLE IF NOT EXISTS `log_configs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '配置名称',
  `code` VARCHAR(64) NOT NULL COMMENT '配置编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '配置描述',
  `logger_name` VARCHAR(128) DEFAULT NULL COMMENT '日志记录器名称',
  `log_level` VARCHAR(16) DEFAULT 'INFO' COMMENT '日志级别：TRACE, DEBUG, INFO, WARN, ERROR',
  `format` VARCHAR(32) DEFAULT 'json' COMMENT '日志格式：json, text, etc.',
  `output_targets` JSON DEFAULT NULL COMMENT '输出目标配置(JSON格式)',
  `retention_days` INT DEFAULT 30 COMMENT '保留天数',
  `max_size` INT DEFAULT 100 COMMENT '最大大小(MB)',
  `max_files` INT DEFAULT 10 COMMENT '最大文件数',
  `compress` TINYINT(1) DEFAULT 1 COMMENT '是否压缩：0-否 1-是',
  `sampling_rate` FLOAT DEFAULT 1.0 COMMENT '采样率(0-1)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_logger_name` (`logger_name`),
  KEY `idx_log_level` (`log_level`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日志配置表';

-- 初始化基础缓存配置
INSERT INTO `cache_configs` (`name`, `code`, `description`, `cache_type`, `ttl`, `status`) VALUES
('系统缓存', 'system-cache', '系统级缓存', 'memory', 3600, 1),
('用户缓存', 'user-cache', '用户数据缓存', 'memory', 7200, 1),
('内容缓存', 'content-cache', '内容数据缓存', 'memory', 1800, 1),
('配置缓存', 'config-cache', '系统配置缓存', 'memory', 86400, 1);

-- 初始化基础告警规则
INSERT INTO `alert_rules` (`name`, `description`, `metric_name`, `condition`, `duration`, `severity`, `status`) VALUES
('CPU使用率过高', '服务器CPU使用率超过90%', 'cpu_usage_percent', '>90', 300, 'warning', 1),
('内存使用率过高', '服务器内存使用率超过85%', 'memory_usage_percent', '>85', 300, 'warning', 1),
('磁盘空间不足', '磁盘剩余空间低于10%', 'disk_free_percent', '<10', 300, 'error', 1),
('服务不可用', '服务健康检查失败', 'service_health_status', '==0', 180, 'critical', 1),
('API响应时间过长', 'API响应时间超过5秒', 'api_response_time_ms', '>5000', 180, 'warning', 1);

-- 初始化基础监控仪表板
INSERT INTO `monitoring_dashboards` (`name`, `code`, `description`, `is_system`, `status`) VALUES
('系统概览', 'system-overview', '系统级监控指标概览', 1, 1),
('服务监控', 'service-monitoring', '微服务监控仪表板', 1, 1),
('API性能', 'api-performance', 'API性能监控仪表板', 1, 1),
('资源使用', 'resource-usage', '系统资源使用监控', 1, 1);
