-- 消息队列和事件总线模块数据表结构

-- 消息队列表
CREATE TABLE IF NOT EXISTS `message_queues` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '队列名称',
  `code` VARCHAR(64) NOT NULL COMMENT '队列编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '队列描述',
  `queue_type` VARCHAR(32) DEFAULT 'standard' COMMENT '队列类型：standard-标准队列 fifo-先进先出队列 priority-优先级队列 delay-延迟队列',
  `max_message_size` INT DEFAULT 262144 COMMENT '最大消息大小(字节)',
  `message_retention_period` INT DEFAULT 345600 COMMENT '消息保留时间(秒)',
  `visibility_timeout` INT DEFAULT 30 COMMENT '可见性超时(秒)',
  `max_receive_count` INT DEFAULT 10 COMMENT '最大接收次数',
  `dead_letter_queue` VARCHAR(64) DEFAULT NULL COMMENT '死信队列编码',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示系统队列',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_queue_type` (`queue_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息队列表';

-- 消息表
CREATE TABLE IF NOT EXISTS `messages` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `queue_id` BIGINT NOT NULL COMMENT '队列ID',
  `message_id` VARCHAR(64) NOT NULL COMMENT '消息ID',
  `producer_id` VARCHAR(128) DEFAULT NULL COMMENT '生产者ID',
  `group_id` VARCHAR(128) DEFAULT NULL COMMENT '消息组ID',
  `message_body` TEXT NOT NULL COMMENT '消息内容',
  `message_attributes` JSON DEFAULT NULL COMMENT '消息属性(JSON格式)',
  `delay_seconds` INT DEFAULT 0 COMMENT '延迟时间(秒)',
  `priority` INT DEFAULT 0 COMMENT '优先级，数字越大优先级越高',
  `deduplication_id` VARCHAR(128) DEFAULT NULL COMMENT '去重ID',
  `sequence_number` BIGINT DEFAULT NULL COMMENT '序列号',
  `enqueue_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '入队时间',
  `first_receive_time` TIMESTAMP NULL COMMENT '首次接收时间',
  `receive_count` INT DEFAULT 0 COMMENT '接收次数',
  `next_visible_time` TIMESTAMP NULL COMMENT '下次可见时间',
  `status` VARCHAR(32) DEFAULT 'available' COMMENT '状态：available-可用 invisible-不可见 delayed-延迟中 consumed-已消费 dead-死信',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_message_id` (`message_id`),
  KEY `idx_queue_id` (`queue_id`),
  KEY `idx_status` (`status`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_enqueue_time` (`enqueue_time`),
  KEY `idx_next_visible_time` (`next_visible_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_messages_queue_id` FOREIGN KEY (`queue_id`) REFERENCES `message_queues` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';

-- 消费者组表
CREATE TABLE IF NOT EXISTS `consumer_groups` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '消费者组名称',
  `code` VARCHAR(64) NOT NULL COMMENT '消费者组编码',
  `queue_id` BIGINT NOT NULL COMMENT '队列ID',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '消费者组描述',
  `max_consumers` INT DEFAULT 10 COMMENT '最大消费者数量',
  `active_consumers` INT DEFAULT 0 COMMENT '活跃消费者数量',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_queue` (`code`, `queue_id`),
  KEY `idx_queue_id` (`queue_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_consumer_groups_queue_id` FOREIGN KEY (`queue_id`) REFERENCES `message_queues` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消费者组表';

-- 消费者表
CREATE TABLE IF NOT EXISTS `consumers` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `consumer_id` VARCHAR(64) NOT NULL COMMENT '消费者ID',
  `group_id` BIGINT NOT NULL COMMENT '消费者组ID',
  `client_id` VARCHAR(128) DEFAULT NULL COMMENT '客户端ID',
  `client_type` VARCHAR(32) DEFAULT NULL COMMENT '客户端类型',
  `client_version` VARCHAR(32) DEFAULT NULL COMMENT '客户端版本',
  `client_ip` VARCHAR(64) DEFAULT NULL COMMENT '客户端IP',
  `status` VARCHAR(32) DEFAULT 'online' COMMENT '状态：online-在线 offline-离线 paused-暂停',
  `heartbeat_time` TIMESTAMP NULL COMMENT '最后心跳时间',
  `connect_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '连接时间',
  `disconnect_time` TIMESTAMP NULL COMMENT '断开时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_consumer_id` (`consumer_id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_status` (`status`),
  KEY `idx_heartbeat_time` (`heartbeat_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_consumers_group_id` FOREIGN KEY (`group_id`) REFERENCES `consumer_groups` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消费者表';

-- 消息消费记录表
CREATE TABLE IF NOT EXISTS `message_consumption_records` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `message_id` VARCHAR(64) NOT NULL COMMENT '消息ID',
  `consumer_id` VARCHAR(64) NOT NULL COMMENT '消费者ID',
  `group_id` BIGINT NOT NULL COMMENT '消费者组ID',
  `queue_id` BIGINT NOT NULL COMMENT '队列ID',
  `receive_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '接收时间',
  `process_time` TIMESTAMP NULL COMMENT '处理时间',
  `status` VARCHAR(32) DEFAULT 'received' COMMENT '状态：received-已接收 processing-处理中 completed-已完成 failed-失败',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `retry_count` INT DEFAULT 0 COMMENT '重试次数',
  `process_duration` INT DEFAULT NULL COMMENT '处理时长(毫秒)',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_message_consumer` (`message_id`, `consumer_id`),
  KEY `idx_consumer_id` (`consumer_id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_queue_id` (`queue_id`),
  KEY `idx_status` (`status`),
  KEY `idx_receive_time` (`receive_time`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息消费记录表';

-- 事件定义表
CREATE TABLE IF NOT EXISTS `event_definitions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '事件名称',
  `code` VARCHAR(64) NOT NULL COMMENT '事件编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '事件描述',
  `category` VARCHAR(64) DEFAULT 'default' COMMENT '事件分类',
  `schema` TEXT DEFAULT NULL COMMENT '事件数据结构(JSON Schema)',
  `example` TEXT DEFAULT NULL COMMENT '事件示例',
  `is_system` TINYINT(1) DEFAULT 0 COMMENT '是否系统事件：0-否 1-是',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示系统事件',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_category` (`category`),
  KEY `idx_is_system` (`is_system`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件定义表';

-- 事件总线表
CREATE TABLE IF NOT EXISTS `event_buses` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '事件总线名称',
  `code` VARCHAR(64) NOT NULL COMMENT '事件总线编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '事件总线描述',
  `type` VARCHAR(32) DEFAULT 'default' COMMENT '总线类型：default-默认 topic-主题 multicast-多播',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示系统总线',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件总线表';

-- 事件路由表
CREATE TABLE IF NOT EXISTS `event_routes` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `bus_id` BIGINT NOT NULL COMMENT '事件总线ID',
  `event_code` VARCHAR(64) NOT NULL COMMENT '事件编码，支持通配符*',
  `target_type` VARCHAR(32) NOT NULL COMMENT '目标类型：queue-消息队列 webhook-Webhook endpoint-服务接口 function-函数',
  `target_id` VARCHAR(128) NOT NULL COMMENT '目标ID',
  `target_params` JSON DEFAULT NULL COMMENT '目标参数(JSON格式)',
  `filter_expression` TEXT DEFAULT NULL COMMENT '过滤表达式',
  `transform_expression` TEXT DEFAULT NULL COMMENT '转换表达式',
  `priority` INT DEFAULT 0 COMMENT '优先级，数字越大优先级越高',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_bus_id` (`bus_id`),
  KEY `idx_event_code` (`event_code`),
  KEY `idx_target_type_id` (`target_type`, `target_id`),
  KEY `idx_priority` (`priority`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_event_routes_bus_id` FOREIGN KEY (`bus_id`) REFERENCES `event_buses` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件路由表';

-- 事件日志表
CREATE TABLE IF NOT EXISTS `event_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `event_id` VARCHAR(64) NOT NULL COMMENT '事件ID',
  `event_code` VARCHAR(64) NOT NULL COMMENT '事件编码',
  `event_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '事件时间',
  `source` VARCHAR(128) DEFAULT NULL COMMENT '事件源',
  `event_data` TEXT DEFAULT NULL COMMENT '事件数据(JSON格式)',
  `metadata` JSON DEFAULT NULL COMMENT '事件元数据',
  `bus_id` BIGINT DEFAULT NULL COMMENT '事件总线ID',
  `routing_count` INT DEFAULT 0 COMMENT '路由次数',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-待处理 processing-处理中 routed-已路由 failed-失败',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_event_id` (`event_id`),
  KEY `idx_event_code` (`event_code`),
  KEY `idx_event_time` (`event_time`),
  KEY `idx_source` (`source`),
  KEY `idx_bus_id` (`bus_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件日志表';

-- 事件路由日志表
CREATE TABLE IF NOT EXISTS `event_routing_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `event_id` VARCHAR(64) NOT NULL COMMENT '事件ID',
  `route_id` BIGINT NOT NULL COMMENT '路由ID',
  `target_type` VARCHAR(32) NOT NULL COMMENT '目标类型',
  `target_id` VARCHAR(128) NOT NULL COMMENT '目标ID',
  `process_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '处理时间',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-待处理 success-成功 failed-失败 rejected-拒绝',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `retry_count` INT DEFAULT 0 COMMENT '重试次数',
  `next_retry_time` TIMESTAMP NULL COMMENT '下次重试时间',
  `process_duration` INT DEFAULT NULL COMMENT '处理时长(毫秒)',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_event_route` (`event_id`, `route_id`),
  KEY `idx_event_id` (`event_id`),
  KEY `idx_route_id` (`route_id`),
  KEY `idx_status` (`status`),
  KEY `idx_process_time` (`process_time`),
  KEY `idx_next_retry_time` (`next_retry_time`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件路由日志表';

-- Webhook订阅表
CREATE TABLE IF NOT EXISTS `webhook_subscriptions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '订阅名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '订阅描述',
  `event_codes` TEXT NOT NULL COMMENT '订阅的事件编码，多个用逗号分隔，支持通配符*',
  `target_url` VARCHAR(255) NOT NULL COMMENT '目标URL',
  `http_method` VARCHAR(16) DEFAULT 'POST' COMMENT 'HTTP方法',
  `headers` TEXT DEFAULT NULL COMMENT '自定义HTTP头(JSON格式)',
  `content_type` VARCHAR(64) DEFAULT 'application/json' COMMENT '内容类型',
  `secret` VARCHAR(255) DEFAULT NULL COMMENT '签名密钥',
  `signature_header` VARCHAR(64) DEFAULT 'X-Webhook-Signature' COMMENT '签名HTTP头',
  `timeout` INT DEFAULT 5000 COMMENT '超时时间(毫秒)',
  `retry_count` INT DEFAULT 3 COMMENT '重试次数',
  `retry_interval` INT DEFAULT 60 COMMENT '重试间隔(秒)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建者ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Webhook订阅表';

-- Webhook调用日志表
CREATE TABLE IF NOT EXISTS `webhook_delivery_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `subscription_id` BIGINT NOT NULL COMMENT '订阅ID',
  `event_id` VARCHAR(64) NOT NULL COMMENT '事件ID',
  `request_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '请求时间',
  `response_time` TIMESTAMP NULL COMMENT '响应时间',
  `status_code` INT DEFAULT NULL COMMENT 'HTTP状态码',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-待处理 success-成功 failed-失败',
  `request_body` TEXT DEFAULT NULL COMMENT '请求体',
  `response_body` TEXT DEFAULT NULL COMMENT '响应体',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `retry_count` INT DEFAULT 0 COMMENT '重试次数',
  `next_retry_time` TIMESTAMP NULL COMMENT '下次重试时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_subscription_id` (`subscription_id`),
  KEY `idx_event_id` (`event_id`),
  KEY `idx_status` (`status`),
  KEY `idx_request_time` (`request_time`),
  KEY `idx_next_retry_time` (`next_retry_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_webhook_delivery_logs_subscription_id` FOREIGN KEY (`subscription_id`) REFERENCES `webhook_subscriptions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Webhook调用日志表';

-- 初始化基础消息队列
INSERT INTO `message_queues` (`name`, `code`, `description`, `queue_type`, `status`) VALUES
('默认队列', 'default-queue', '系统默认消息队列', 'standard', 1),
('用户事件队列', 'user-events', '用户相关事件队列', 'standard', 1),
('通知队列', 'notification-queue', '系统通知消息队列', 'standard', 1),
('交易队列', 'trade-queue', '交易相关消息队列', 'standard', 1),
('内容队列', 'content-queue', '内容相关消息队列', 'standard', 1),
('死信队列', 'dead-letter-queue', '无法处理的消息队列', 'standard', 1);

-- 初始化基础事件总线
INSERT INTO `event_buses` (`name`, `code`, `description`, `type`, `status`) VALUES
('系统事件总线', 'system-bus', '系统级事件总线', 'default', 1),
('用户事件总线', 'user-bus', '用户相关事件总线', 'topic', 1),
('业务事件总线', 'business-bus', '业务相关事件总线', 'topic', 1);

-- 初始化基础事件定义
INSERT INTO `event_definitions` (`name`, `code`, `description`, `category`, `is_system`, `status`) VALUES
('用户注册', 'user.registered', '新用户注册成功事件', 'user', 1, 1),
('用户登录', 'user.logged_in', '用户登录成功事件', 'user', 1, 1),
('用户更新', 'user.updated', '用户信息更新事件', 'user', 1, 1),
('内容创建', 'content.created', '新内容创建事件', 'content', 1, 1),
('内容更新', 'content.updated', '内容更新事件', 'content', 1, 1),
('内容发布', 'content.published', '内容发布事件', 'content', 1, 1),
('订单创建', 'order.created', '新订单创建事件', 'trade', 1, 1),
('订单支付', 'order.paid', '订单支付成功事件', 'trade', 1, 1),
('订单完成', 'order.completed', '订单完成事件', 'trade', 1, 1),
('系统通知', 'system.notification', '系统通知事件', 'system', 1, 1);
