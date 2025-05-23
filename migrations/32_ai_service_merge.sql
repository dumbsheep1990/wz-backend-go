-- AI服务表结构合并

-- 客服对话记录表（对19_ai_service.sql的补充）
CREATE TABLE IF NOT EXISTS `ai_chat_records` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `session_id` VARCHAR(64) NOT NULL COMMENT '会话ID',
  `message` TEXT NOT NULL COMMENT '用户消息内容',
  `response` TEXT NOT NULL COMMENT 'AI回复内容',
  `message_type` VARCHAR(32) DEFAULT 'text' COMMENT '消息类型：text-文本 image-图片 voice-语音 video-视频',
  `intent` VARCHAR(64) DEFAULT NULL COMMENT '意图识别结果',
  `entity` JSON DEFAULT NULL COMMENT '实体识别结果',
  `sentiment` VARCHAR(32) DEFAULT NULL COMMENT '情感分析结果',
  `feedback` VARCHAR(32) DEFAULT NULL COMMENT '用户反馈：positive-正面 negative-负面',
  `feedback_content` VARCHAR(512) DEFAULT NULL COMMENT '反馈内容',
  `query_duration` INT DEFAULT NULL COMMENT '查询处理时间(毫秒)',
  `source` VARCHAR(32) DEFAULT 'web' COMMENT '来源：web-网页 app-应用 api-接口',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI客服对话记录表';

-- 内容审核记录表（对19_ai_service.sql的补充）
CREATE TABLE IF NOT EXISTS `ai_content_reviews` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `content_id` BIGINT NOT NULL COMMENT '内容ID',
  `content_type` VARCHAR(32) NOT NULL COMMENT '内容类型：text-文本 image-图片 video-视频 audio-音频',
  `review_type` VARCHAR(32) NOT NULL COMMENT '审核类型：auto-自动 manual-人工',
  `content_hash` VARCHAR(64) DEFAULT NULL COMMENT '内容哈希值',
  `is_passed` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否通过：0-未通过 1-通过',
  `reject_reason` VARCHAR(255) DEFAULT NULL COMMENT '拒绝原因',
  `risk_level` VARCHAR(32) DEFAULT NULL COMMENT '风险级别：low-低 medium-中 high-高',
  `risk_detail` JSON DEFAULT NULL COMMENT '风险详情',
  `review_model` VARCHAR(64) DEFAULT NULL COMMENT '审核模型',
  `confidence` DECIMAL(5,4) DEFAULT NULL COMMENT '置信度',
  `review_time` INT DEFAULT NULL COMMENT '审核耗时(毫秒)',
  `reviewer_id` BIGINT DEFAULT NULL COMMENT '审核员ID(人工审核时)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_content` (`content_id`, `content_type`),
  KEY `idx_content_hash` (`content_hash`),
  KEY `idx_is_passed` (`is_passed`),
  KEY `idx_risk_level` (`risk_level`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI内容审核记录表';

-- AI服务调用日志表
CREATE TABLE IF NOT EXISTS `ai_service_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `service_type` VARCHAR(32) NOT NULL COMMENT '服务类型：chat-聊天 moderation-审核 embedding-嵌入 generation-生成',
  `model_id` BIGINT DEFAULT NULL COMMENT '模型ID',
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `request_id` VARCHAR(64) NOT NULL COMMENT '请求ID',
  `request_data` JSON DEFAULT NULL COMMENT '请求数据',
  `response_data` JSON DEFAULT NULL COMMENT '响应数据',
  `token_count` INT DEFAULT 0 COMMENT 'token数量',
  `duration` INT DEFAULT 0 COMMENT '处理时间(毫秒)',
  `status` VARCHAR(32) DEFAULT 'success' COMMENT '状态：success-成功 failed-失败',
  `error_message` VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
  `cost` DECIMAL(10,6) DEFAULT 0 COMMENT '成本',
  `ip` VARCHAR(64) DEFAULT NULL COMMENT 'IP地址',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_service_type` (`service_type`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_request_id` (`request_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI服务调用日志表';

-- AI推荐策略表（对21_recommendation.sql的补充）
CREATE TABLE IF NOT EXISTS `ai_recommendation_strategies` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '策略名称',
  `code` VARCHAR(32) NOT NULL COMMENT '策略编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '策略描述',
  `algorithm` VARCHAR(64) NOT NULL COMMENT '算法：collaborative_filtering-协同过滤 content_based-基于内容 hybrid-混合',
  `parameters` JSON DEFAULT NULL COMMENT '参数配置',
  `weights` JSON DEFAULT NULL COMMENT '权重配置',
  `filters` JSON DEFAULT NULL COMMENT '过滤条件',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI推荐策略表';

-- AI训练任务表
CREATE TABLE IF NOT EXISTS `ai_training_tasks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `task_id` VARCHAR(64) NOT NULL COMMENT '任务ID',
  `name` VARCHAR(64) NOT NULL COMMENT '任务名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '任务描述',
  `model_id` BIGINT NOT NULL COMMENT '模型ID',
  `task_type` VARCHAR(32) NOT NULL COMMENT '任务类型：fine_tuning-微调 training-训练',
  `dataset_ids` JSON DEFAULT NULL COMMENT '数据集ID列表',
  `parameters` JSON DEFAULT NULL COMMENT '训练参数',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-等待 running-运行中 completed-完成 failed-失败',
  `progress` DECIMAL(5,2) DEFAULT 0 COMMENT '进度',
  `start_time` TIMESTAMP NULL COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `result_model_id` BIGINT DEFAULT NULL COMMENT '结果模型ID',
  `metrics` JSON DEFAULT NULL COMMENT '训练指标',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_task_id` (`task_id`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_ai_training_tasks_model_id` FOREIGN KEY (`model_id`) REFERENCES `ai_models` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI训练任务表';
