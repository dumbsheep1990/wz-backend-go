-- API网关和接口管理模块数据表结构

-- API服务表
CREATE TABLE IF NOT EXISTS `api_services` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '服务名称',
  `code` VARCHAR(64) NOT NULL COMMENT '服务编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '服务描述',
  `base_path` VARCHAR(128) NOT NULL COMMENT '基础路径',
  `protocol` VARCHAR(16) DEFAULT 'http' COMMENT '协议：http, https, grpc',
  `host` VARCHAR(128) DEFAULT NULL COMMENT '主机地址',
  `port` INT DEFAULT NULL COMMENT '端口号',
  `health_check_url` VARCHAR(255) DEFAULT NULL COMMENT '健康检查URL',
  `docs_url` VARCHAR(255) DEFAULT NULL COMMENT '文档URL',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `is_public` TINYINT(1) DEFAULT 0 COMMENT '是否公开：0-内部服务 1-公开服务',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示系统服务',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_is_public` (`is_public`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API服务表';

-- API接口表
CREATE TABLE IF NOT EXISTS `api_endpoints` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `service_id` BIGINT NOT NULL COMMENT '服务ID',
  `name` VARCHAR(128) NOT NULL COMMENT '接口名称',
  `path` VARCHAR(255) NOT NULL COMMENT '接口路径',
  `method` VARCHAR(20) NOT NULL COMMENT '请求方法：GET, POST, PUT, DELETE, etc.',
  `version` VARCHAR(16) DEFAULT 'v1' COMMENT 'API版本',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '接口描述',
  `request_schema` TEXT DEFAULT NULL COMMENT '请求参数结构(JSON格式)',
  `response_schema` TEXT DEFAULT NULL COMMENT '响应结构(JSON格式)',
  `example_request` TEXT DEFAULT NULL COMMENT '示例请求',
  `example_response` TEXT DEFAULT NULL COMMENT '示例响应',
  `timeout` INT DEFAULT 5000 COMMENT '超时时间(毫秒)',
  `rate_limit` INT DEFAULT 0 COMMENT '速率限制(每分钟请求数)，0表示不限制',
  `auth_required` TINYINT(1) DEFAULT 1 COMMENT '是否需要认证：0-不需要 1-需要',
  `permissions` VARCHAR(255) DEFAULT NULL COMMENT '所需权限，多个权限用逗号分隔',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `deprecated` TINYINT(1) DEFAULT 0 COMMENT '是否已废弃：0-否 1-是',
  `mock_enabled` TINYINT(1) DEFAULT 0 COMMENT '是否启用Mock：0-否 1-是',
  `mock_response` TEXT DEFAULT NULL COMMENT 'Mock响应数据',
  `tags` VARCHAR(255) DEFAULT NULL COMMENT '标签，多个标签用逗号分隔',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_service_path_method_version` (`service_id`, `path`, `method`, `version`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deprecated` (`deprecated`),
  KEY `idx_auth_required` (`auth_required`),
  CONSTRAINT `fk_api_endpoints_service_id` FOREIGN KEY (`service_id`) REFERENCES `api_services` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API接口表';

-- API访问密钥表
CREATE TABLE IF NOT EXISTS `api_keys` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '密钥名称',
  `key_id` VARCHAR(64) NOT NULL COMMENT '密钥ID',
  `key_secret` VARCHAR(255) NOT NULL COMMENT '密钥秘钥',
  `type` VARCHAR(32) DEFAULT 'user' COMMENT '密钥类型：user-用户密钥 service-服务密钥 tenant-租户密钥',
  `owner_id` BIGINT DEFAULT NULL COMMENT '所有者ID',
  `owner_type` VARCHAR(32) DEFAULT NULL COMMENT '所有者类型：user, tenant, service',
  `expire_time` TIMESTAMP NULL COMMENT '过期时间，为NULL表示永不过期',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `ip_whitelist` TEXT DEFAULT NULL COMMENT 'IP白名单，多个IP用逗号分隔',
  `allowed_services` TEXT DEFAULT NULL COMMENT '允许访问的服务，多个服务编码用逗号分隔，为空表示允许所有',
  `allowed_endpoints` TEXT DEFAULT NULL COMMENT '允许访问的接口，格式为serviceCode:path:method，多个用逗号分隔，为空表示允许所有',
  `rate_limit` INT DEFAULT 0 COMMENT '速率限制(每分钟请求数)，0表示不限制',
  `quota_limit` INT DEFAULT 0 COMMENT '配额限制(每天请求数)，0表示不限制',
  `last_used_time` TIMESTAMP NULL COMMENT '最后使用时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_key_id` (`key_id`),
  KEY `idx_owner_id_type` (`owner_id`, `owner_type`),
  KEY `idx_status` (`status`),
  KEY `idx_expire_time` (`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API访问密钥表';

-- API访问日志表
CREATE TABLE IF NOT EXISTS `api_access_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `trace_id` VARCHAR(64) NOT NULL COMMENT '追踪ID',
  `service_id` BIGINT DEFAULT NULL COMMENT '服务ID',
  `endpoint_id` BIGINT DEFAULT NULL COMMENT '接口ID',
  `path` VARCHAR(255) NOT NULL COMMENT '请求路径',
  `method` VARCHAR(20) NOT NULL COMMENT '请求方法',
  `client_ip` VARCHAR(64) DEFAULT NULL COMMENT '客户端IP',
  `user_agent` VARCHAR(255) DEFAULT NULL COMMENT '用户代理',
  `request_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '请求时间',
  `response_time` TIMESTAMP NULL COMMENT '响应时间',
  `latency` INT DEFAULT NULL COMMENT '延迟时间(毫秒)',
  `status_code` INT DEFAULT NULL COMMENT '状态码',
  `request_size` INT DEFAULT NULL COMMENT '请求大小(字节)',
  `response_size` INT DEFAULT NULL COMMENT '响应大小(字节)',
  `request_body` TEXT DEFAULT NULL COMMENT '请求体(可选存储)',
  `response_body` TEXT DEFAULT NULL COMMENT '响应体(可选存储)',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `api_key_id` BIGINT DEFAULT NULL COMMENT 'API密钥ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_trace_id` (`trace_id`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_endpoint_id` (`endpoint_id`),
  KEY `idx_request_time` (`request_time`),
  KEY `idx_status_code` (`status_code`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_api_key_id` (`api_key_id`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API访问日志表';

-- API监控指标表
CREATE TABLE IF NOT EXISTS `api_metrics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `service_id` BIGINT DEFAULT NULL COMMENT '服务ID',
  `endpoint_id` BIGINT DEFAULT NULL COMMENT '接口ID',
  `metric_name` VARCHAR(64) NOT NULL COMMENT '指标名称：requests, errors, latency, etc.',
  `metric_value` DOUBLE NOT NULL COMMENT '指标值',
  `timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间戳',
  `dimensions` JSON DEFAULT NULL COMMENT '维度信息(JSON格式)',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_endpoint_id` (`endpoint_id`),
  KEY `idx_metric_name` (`metric_name`),
  KEY `idx_timestamp` (`timestamp`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API监控指标表';

-- API限流规则表
CREATE TABLE IF NOT EXISTS `api_rate_limits` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `target_type` VARCHAR(32) NOT NULL COMMENT '目标类型：service-服务 endpoint-接口 ip-IP地址 user-用户 tenant-租户',
  `target_id` VARCHAR(128) NOT NULL COMMENT '目标ID',
  `limit_type` VARCHAR(32) NOT NULL COMMENT '限制类型：rate-速率限制 quota-配额限制',
  `time_unit` VARCHAR(16) NOT NULL COMMENT '时间单位：second, minute, hour, day',
  `limit_value` INT NOT NULL COMMENT '限制值',
  `action` VARCHAR(32) DEFAULT 'block' COMMENT '超限动作：block-阻止 throttle-节流 log-仅记录',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_target_type_id_limit_type` (`target_type`, `target_id`, `limit_type`, `time_unit`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API限流规则表';

-- API路由表
CREATE TABLE IF NOT EXISTS `api_routes` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '路由名称',
  `path_pattern` VARCHAR(255) NOT NULL COMMENT '路径模式',
  `service_id` BIGINT NOT NULL COMMENT '目标服务ID',
  `target_path` VARCHAR(255) DEFAULT NULL COMMENT '目标路径，为空则使用原路径',
  `methods` VARCHAR(128) DEFAULT NULL COMMENT '允许的HTTP方法，多个用逗号分隔，为空表示全部',
  `priority` INT DEFAULT 0 COMMENT '优先级，数字越大优先级越高',
  `strip_prefix` TINYINT(1) DEFAULT 0 COMMENT '是否去除前缀：0-否 1-是',
  `preserve_host` TINYINT(1) DEFAULT 0 COMMENT '是否保留主机头：0-否 1-是',
  `timeout` INT DEFAULT 5000 COMMENT '超时时间(毫秒)',
  `retry_count` INT DEFAULT 0 COMMENT '重试次数',
  `auth_required` TINYINT(1) DEFAULT 1 COMMENT '是否需要认证：0-不需要 1-需要',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_priority` (`priority`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_api_routes_service_id` FOREIGN KEY (`service_id`) REFERENCES `api_services` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API路由表';

-- API插件表
CREATE TABLE IF NOT EXISTS `api_plugins` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '插件名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '插件描述',
  `type` VARCHAR(32) NOT NULL COMMENT '插件类型：auth, security, traffic, transform, log',
  `handler` VARCHAR(128) NOT NULL COMMENT '处理器名称',
  `config` JSON DEFAULT NULL COMMENT '插件配置(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API插件表';

-- API插件实例表
CREATE TABLE IF NOT EXISTS `api_plugin_instances` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `plugin_id` BIGINT NOT NULL COMMENT '插件ID',
  `target_type` VARCHAR(32) NOT NULL COMMENT '目标类型：global-全局 service-服务 route-路由 endpoint-接口',
  `target_id` BIGINT NOT NULL COMMENT '目标ID',
  `config` JSON DEFAULT NULL COMMENT '插件实例配置(JSON格式)',
  `priority` INT DEFAULT 0 COMMENT '优先级，数字越大优先级越高',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_plugin_id` (`plugin_id`),
  KEY `idx_target_type_id` (`target_type`, `target_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_api_plugin_instances_plugin_id` FOREIGN KEY (`plugin_id`) REFERENCES `api_plugins` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API插件实例表';

-- 初始化基础API服务
INSERT INTO `api_services` (`name`, `code`, `description`, `base_path`, `protocol`, `host`, `port`, `status`, `is_public`) VALUES
('用户服务', 'user-service', '用户管理相关接口', '/api/users', 'http', 'localhost', 8001, 1, 1),
('租户服务', 'tenant-service', '租户管理相关接口', '/api/tenants', 'http', 'localhost', 8002, 1, 0),
('内容服务', 'content-service', '内容管理相关接口', '/api/contents', 'http', 'localhost', 8003, 1, 1),
('交易服务', 'trade-service', '交易管理相关接口', '/api/trades', 'http', 'localhost', 8004, 1, 1),
('文件服务', 'file-service', '文件存储相关接口', '/api/files', 'http', 'localhost', 8005, 1, 1);

-- 初始化基础API插件
INSERT INTO `api_plugins` (`name`, `description`, `type`, `handler`, `status`) VALUES
('jwt-auth', 'JWT认证插件', 'auth', 'jwt_auth_handler', 1),
('cors', '跨域资源共享插件', 'security', 'cors_handler', 1),
('rate-limit', '请求速率限制插件', 'traffic', 'rate_limit_handler', 1),
('ip-restriction', 'IP限制插件', 'security', 'ip_restriction_handler', 1),
('request-transform', '请求转换插件', 'transform', 'request_transform_handler', 1),
('response-transform', '响应转换插件', 'transform', 'response_transform_handler', 1),
('access-log', '访问日志插件', 'log', 'access_log_handler', 1);
