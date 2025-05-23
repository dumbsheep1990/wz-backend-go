-- 网关服务路由关联表结构

-- 服务API映射表
CREATE TABLE IF NOT EXISTS `service_api_mappings` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `service_id` BIGINT NOT NULL COMMENT '服务ID',
  `api_id` BIGINT NOT NULL COMMENT 'API ID',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_service_api` (`service_id`, `api_id`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_api_id` (`api_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_service_api_mappings_service_id` FOREIGN KEY (`service_id`) REFERENCES `api_services` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_service_api_mappings_api_id` FOREIGN KEY (`api_id`) REFERENCES `api_routes` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务API映射表';

-- 客户端订阅服务表
CREATE TABLE IF NOT EXISTS `client_service_subscriptions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `client_id` BIGINT NOT NULL COMMENT '客户端ID',
  `service_id` BIGINT NOT NULL COMMENT '服务ID',
  `access_level` VARCHAR(32) NOT NULL DEFAULT 'read' COMMENT '访问级别：read-只读 write-读写 admin-管理',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `start_time` TIMESTAMP NULL COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_client_service` (`client_id`, `service_id`),
  KEY `idx_client_id` (`client_id`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_status` (`status`),
  KEY `idx_end_time` (`end_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_client_service_subscriptions_client_id` FOREIGN KEY (`client_id`) REFERENCES `api_clients` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_client_service_subscriptions_service_id` FOREIGN KEY (`service_id`) REFERENCES `api_services` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户端订阅服务表';

-- 服务依赖表
CREATE TABLE IF NOT EXISTS `service_dependencies` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `service_id` BIGINT NOT NULL COMMENT '服务ID',
  `dependency_service_id` BIGINT NOT NULL COMMENT '依赖的服务ID',
  `dependency_type` VARCHAR(32) NOT NULL DEFAULT 'required' COMMENT '依赖类型：required-必需 optional-可选',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_service_dependency` (`service_id`, `dependency_service_id`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_dependency_service_id` (`dependency_service_id`),
  KEY `idx_dependency_type` (`dependency_type`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_service_dependencies_service_id` FOREIGN KEY (`service_id`) REFERENCES `api_services` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_service_dependencies_dependency_id` FOREIGN KEY (`dependency_service_id`) REFERENCES `api_services` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务依赖表';
