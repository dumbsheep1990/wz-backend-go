-- 工作流和数据集成关系表结构

-- 工作流节点关联表
CREATE TABLE IF NOT EXISTS `workflow_node_relations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `workflow_id` BIGINT NOT NULL COMMENT '工作流ID',
  `source_node_id` BIGINT NOT NULL COMMENT '源节点ID',
  `target_node_id` BIGINT NOT NULL COMMENT '目标节点ID',
  `condition` TEXT DEFAULT NULL COMMENT '条件表达式',
  `priority` INT DEFAULT 0 COMMENT '优先级',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_workflow_id` (`workflow_id`),
  KEY `idx_source_node_id` (`source_node_id`),
  KEY `idx_target_node_id` (`target_node_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_workflow_node_relations_workflow_id` FOREIGN KEY (`workflow_id`) REFERENCES `workflows` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_workflow_node_relations_source_node_id` FOREIGN KEY (`source_node_id`) REFERENCES `workflow_nodes` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_workflow_node_relations_target_node_id` FOREIGN KEY (`target_node_id`) REFERENCES `workflow_nodes` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流节点关联表';

-- 数据集成映射表
CREATE TABLE IF NOT EXISTS `data_integration_mappings` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `integration_id` BIGINT NOT NULL COMMENT '集成ID',
  `source_field` VARCHAR(128) NOT NULL COMMENT '源字段',
  `target_field` VARCHAR(128) NOT NULL COMMENT '目标字段',
  `transform_rule` TEXT DEFAULT NULL COMMENT '转换规则',
  `is_required` TINYINT(1) DEFAULT 0 COMMENT '是否必需：0-否 1-是',
  `default_value` VARCHAR(255) DEFAULT NULL COMMENT '默认值',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_integration_id` (`integration_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_data_integration_mappings_integration_id` FOREIGN KEY (`integration_id`) REFERENCES `data_integrations` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据集成映射表';

-- 服务事件订阅表
CREATE TABLE IF NOT EXISTS `service_event_subscriptions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `service_id` BIGINT NOT NULL COMMENT '服务ID',
  `event_type` VARCHAR(128) NOT NULL COMMENT '事件类型',
  `subscriber_service_id` BIGINT NOT NULL COMMENT '订阅服务ID',
  `callback_url` VARCHAR(255) DEFAULT NULL COMMENT '回调URL',
  `filter_condition` TEXT DEFAULT NULL COMMENT '过滤条件',
  `retry_strategy` JSON DEFAULT NULL COMMENT '重试策略(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_service_event_subscriber` (`service_id`, `event_type`, `subscriber_service_id`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_event_type` (`event_type`),
  KEY `idx_subscriber_service_id` (`subscriber_service_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务事件订阅表';

-- 服务同步状态表
CREATE TABLE IF NOT EXISTS `service_sync_states` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `source_service_id` BIGINT NOT NULL COMMENT '源服务ID',
  `target_service_id` BIGINT NOT NULL COMMENT '目标服务ID',
  `resource_type` VARCHAR(64) NOT NULL COMMENT '资源类型',
  `last_sync_time` TIMESTAMP NULL COMMENT '最后同步时间',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-待同步 syncing-同步中 completed-已完成 failed-失败',
  `error_message` VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
  `sync_count` INT DEFAULT 0 COMMENT '同步记录数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_source_target_resource` (`source_service_id`, `target_service_id`, `resource_type`),
  KEY `idx_source_service_id` (`source_service_id`),
  KEY `idx_target_service_id` (`target_service_id`),
  KEY `idx_resource_type` (`resource_type`),
  KEY `idx_status` (`status`),
  KEY `idx_last_sync_time` (`last_sync_time`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务同步状态表';
