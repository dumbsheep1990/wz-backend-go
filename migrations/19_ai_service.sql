-- AI服务数据表结构

-- AI模型配置表
CREATE TABLE IF NOT EXISTS `ai_models` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '模型名称',
  `code` VARCHAR(64) NOT NULL COMMENT '模型编码',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '模型描述',
  `model_type` VARCHAR(32) NOT NULL COMMENT '模型类型：llm-自然语言处理 ocr-光学字符识别 mome-多模态 audio-语音处理 recommendation-推荐系统',
  `provider` VARCHAR(64) DEFAULT NULL COMMENT '提供商：openai, azure, huggingface, etc.',
  `version` VARCHAR(32) DEFAULT NULL COMMENT '模型版本',
  `endpoint` VARCHAR(255) DEFAULT NULL COMMENT 'API端点',
  `config` JSON DEFAULT NULL COMMENT '模型配置(JSON格式)',
  `credentials` TEXT DEFAULT NULL COMMENT '认证凭证（加密存储）',
  `parameters` JSON DEFAULT NULL COMMENT '默认参数(JSON格式)',
  `max_tokens` INT DEFAULT 2048 COMMENT '最大token数',
  `price_per_token` DECIMAL(10, 8) DEFAULT 0 COMMENT '每token价格',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `is_default` TINYINT(1) DEFAULT 0 COMMENT '是否默认模型：0-否 1-是',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示系统模型',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_model_type` (`model_type`),
  KEY `idx_provider` (`provider`),
  KEY `idx_status` (`status`),
  KEY `idx_is_default` (`is_default`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI模型配置表';

-- AI应用表
CREATE TABLE IF NOT EXISTS `ai_applications` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '应用名称',
  `code` VARCHAR(64) NOT NULL COMMENT '应用编码',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '应用描述',
  `application_type` VARCHAR(32) NOT NULL COMMENT '应用类型：chatbot-聊天机器人 content-内容生成 analysis-数据分析 image-图像处理',
  `model_id` BIGINT NOT NULL COMMENT '模型ID',
  `prompt_template` TEXT DEFAULT NULL COMMENT '提示词模板',
  `system_message` TEXT DEFAULT NULL COMMENT '系统消息',
  `parameters` JSON DEFAULT NULL COMMENT '应用参数(JSON格式)',
  `input_schema` JSON DEFAULT NULL COMMENT '输入结构(JSON格式)',
  `output_schema` JSON DEFAULT NULL COMMENT '输出结构(JSON格式)',
  `rate_limit` INT DEFAULT 0 COMMENT '速率限制(每分钟请求数)，0表示不限制',
  `daily_quota` INT DEFAULT 0 COMMENT '每日配额，0表示不限制',
  `access_control` JSON DEFAULT NULL COMMENT '访问控制(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `is_public` TINYINT(1) DEFAULT 0 COMMENT '是否公开：0-私有 1-公开',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建者ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_application_type` (`application_type`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_public` (`is_public`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_ai_applications_model_id` FOREIGN KEY (`model_id`) REFERENCES `ai_models` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI应用表';

-- AI对话会话表
CREATE TABLE IF NOT EXISTS `ai_conversations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `conversation_id` VARCHAR(64) NOT NULL COMMENT '会话ID',
  `title` VARCHAR(255) DEFAULT NULL COMMENT '会话标题',
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `application_id` BIGINT NOT NULL COMMENT '应用ID',
  `model_id` BIGINT NOT NULL COMMENT '模型ID',
  `system_message` TEXT DEFAULT NULL COMMENT '系统消息',
  `parameters` JSON DEFAULT NULL COMMENT '会话参数(JSON格式)',
  `metadata` JSON DEFAULT NULL COMMENT '元数据(JSON格式)',
  `status` VARCHAR(32) DEFAULT 'active' COMMENT '状态：active-活跃 archived-已归档 deleted-已删除',
  `last_message_time` TIMESTAMP NULL COMMENT '最后消息时间',
  `message_count` INT DEFAULT 0 COMMENT '消息数量',
  `token_count` INT DEFAULT 0 COMMENT 'token数量',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_conversation_id` (`conversation_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_status` (`status`),
  KEY `idx_last_message_time` (`last_message_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_ai_conversations_application_id` FOREIGN KEY (`application_id`) REFERENCES `ai_applications` (`id`),
  CONSTRAINT `fk_ai_conversations_model_id` FOREIGN KEY (`model_id`) REFERENCES `ai_models` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI对话会话表';

-- AI对话消息表
CREATE TABLE IF NOT EXISTS `ai_messages` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `message_id` VARCHAR(64) NOT NULL COMMENT '消息ID',
  `conversation_id` VARCHAR(64) NOT NULL COMMENT '会话ID',
  `role` VARCHAR(32) NOT NULL COMMENT '角色：user-用户 assistant-助手 system-系统',
  `content` TEXT NOT NULL COMMENT '消息内容',
  `attachments` JSON DEFAULT NULL COMMENT '附件(JSON格式)',
  `tokens` INT DEFAULT 0 COMMENT 'token数量',
  `parent_message_id` VARCHAR(64) DEFAULT NULL COMMENT '父消息ID',
  `metadata` JSON DEFAULT NULL COMMENT '元数据(JSON格式)',
  `feedback` VARCHAR(32) DEFAULT NULL COMMENT '反馈：like-点赞 dislike-点踩',
  `feedback_note` VARCHAR(512) DEFAULT NULL COMMENT '反馈备注',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_message_id` (`message_id`),
  KEY `idx_conversation_id` (`conversation_id`),
  KEY `idx_role` (`role`),
  KEY `idx_parent_message_id` (`parent_message_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI对话消息表';

-- AI内容生成表
CREATE TABLE IF NOT EXISTS `ai_generated_contents` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `content_id` VARCHAR(64) NOT NULL COMMENT '内容ID',
  `title` VARCHAR(255) DEFAULT NULL COMMENT '内容标题',
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `application_id` BIGINT NOT NULL COMMENT '应用ID',
  `model_id` BIGINT NOT NULL COMMENT '模型ID',
  `prompt` TEXT NOT NULL COMMENT '提示词',
  `content` TEXT NOT NULL COMMENT '生成内容',
  `content_type` VARCHAR(32) DEFAULT 'text' COMMENT '内容类型：text-文本 image-图像 code-代码 audio-音频',
  `parameters` JSON DEFAULT NULL COMMENT '生成参数(JSON格式)',
  `tokens` INT DEFAULT 0 COMMENT 'token数量',
  `duration` INT DEFAULT 0 COMMENT '生成时长(毫秒)',
  `status` VARCHAR(32) DEFAULT 'completed' COMMENT '状态：pending-待处理 processing-处理中 completed-已完成 failed-失败',
  `error_message` VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
  `feedback` VARCHAR(32) DEFAULT NULL COMMENT '反馈：like-点赞 dislike-点踩',
  `feedback_note` VARCHAR(512) DEFAULT NULL COMMENT '反馈备注',
  `metadata` JSON DEFAULT NULL COMMENT '元数据(JSON格式)',
  `is_public` TINYINT(1) DEFAULT 0 COMMENT '是否公开：0-私有 1-公开',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_content_id` (`content_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_content_type` (`content_type`),
  KEY `idx_status` (`status`),
  KEY `idx_is_public` (`is_public`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_ai_generated_contents_application_id` FOREIGN KEY (`application_id`) REFERENCES `ai_applications` (`id`),
  CONSTRAINT `fk_ai_generated_contents_model_id` FOREIGN KEY (`model_id`) REFERENCES `ai_models` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI内容生成表';

-- AI资源使用记录表
CREATE TABLE IF NOT EXISTS `ai_usage_records` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `application_id` BIGINT NOT NULL COMMENT '应用ID',
  `model_id` BIGINT NOT NULL COMMENT '模型ID',
  `resource_type` VARCHAR(32) NOT NULL COMMENT '资源类型：conversation-对话 generation-生成 embedding-嵌入 analysis-分析',
  `resource_id` VARCHAR(64) DEFAULT NULL COMMENT '资源ID',
  `input_tokens` INT DEFAULT 0 COMMENT '输入token数量',
  `output_tokens` INT DEFAULT 0 COMMENT '输出token数量',
  `total_tokens` INT DEFAULT 0 COMMENT '总token数量',
  `duration` INT DEFAULT 0 COMMENT '处理时长(毫秒)',
  `cost` DECIMAL(10, 6) DEFAULT 0 COMMENT '成本',
  `status` VARCHAR(32) DEFAULT 'success' COMMENT '状态：success-成功 failed-失败',
  `usage_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '使用时间',
  `ip` VARCHAR(64) DEFAULT NULL COMMENT 'IP地址',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_resource_type` (`resource_type`),
  KEY `idx_resource_id` (`resource_id`),
  KEY `idx_usage_time` (`usage_time`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_ai_usage_records_application_id` FOREIGN KEY (`application_id`) REFERENCES `ai_applications` (`id`),
  CONSTRAINT `fk_ai_usage_records_model_id` FOREIGN KEY (`model_id`) REFERENCES `ai_models` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI资源使用记录表';

-- AI知识库表
CREATE TABLE IF NOT EXISTS `ai_knowledge_bases` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '知识库名称',
  `code` VARCHAR(64) NOT NULL COMMENT '知识库编码',
  `description` VARCHAR(512) DEFAULT NULL COMMENT '知识库描述',
  `embeddings_model` VARCHAR(128) DEFAULT NULL COMMENT '嵌入模型',
  `retrieval_model` VARCHAR(32) DEFAULT 'similarity' COMMENT '检索模型：similarity-相似度 semantic-语义',
  `chunk_size` INT DEFAULT 1000 COMMENT '分块大小',
  `chunk_overlap` INT DEFAULT 200 COMMENT '分块重叠',
  `vector_dimension` INT DEFAULT 1536 COMMENT '向量维度',
  `metadata_fields` JSON DEFAULT NULL COMMENT '元数据字段(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建者ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_tenant` (`code`, `tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI知识库表';

-- AI知识库文档表
CREATE TABLE IF NOT EXISTS `ai_knowledge_documents` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `document_id` VARCHAR(64) NOT NULL COMMENT '文档ID',
  `knowledge_base_id` BIGINT NOT NULL COMMENT '知识库ID',
  `title` VARCHAR(255) NOT NULL COMMENT '文档标题',
  `content` TEXT DEFAULT NULL COMMENT '文档内容',
  `content_type` VARCHAR(32) DEFAULT 'text' COMMENT '内容类型：text-文本 html-HTML markdown-Markdown pdf-PDF etc.',
  `file_url` VARCHAR(255) DEFAULT NULL COMMENT '文件URL',
  `source_url` VARCHAR(255) DEFAULT NULL COMMENT '来源URL',
  `metadata` JSON DEFAULT NULL COMMENT '元数据(JSON格式)',
  `chunk_count` INT DEFAULT 0 COMMENT '分块数量',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-待处理 processing-处理中 completed-已完成 failed-失败',
  `error_message` VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建者ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_document_id` (`document_id`),
  KEY `idx_knowledge_base_id` (`knowledge_base_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_ai_knowledge_documents_kb_id` FOREIGN KEY (`knowledge_base_id`) REFERENCES `ai_knowledge_bases` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI知识库文档表';

-- AI知识库分块表
CREATE TABLE IF NOT EXISTS `ai_knowledge_chunks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `chunk_id` VARCHAR(64) NOT NULL COMMENT '分块ID',
  `document_id` VARCHAR(64) NOT NULL COMMENT '文档ID',
  `knowledge_base_id` BIGINT NOT NULL COMMENT '知识库ID',
  `content` TEXT NOT NULL COMMENT '分块内容',
  `embedding` LONGBLOB DEFAULT NULL COMMENT '向量嵌入',
  `embedding_model` VARCHAR(128) DEFAULT NULL COMMENT '嵌入模型',
  `metadata` JSON DEFAULT NULL COMMENT '元数据(JSON格式)',
  `token_count` INT DEFAULT 0 COMMENT 'token数量',
  `position` INT DEFAULT 0 COMMENT '位置',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_chunk_id` (`chunk_id`),
  KEY `idx_document_id` (`document_id`),
  KEY `idx_knowledge_base_id` (`knowledge_base_id`),
  KEY `idx_position` (`position`),
  CONSTRAINT `fk_ai_knowledge_chunks_kb_id` FOREIGN KEY (`knowledge_base_id`) REFERENCES `ai_knowledge_bases` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI知识库分块表';

-- AI图像生成表
CREATE TABLE IF NOT EXISTS `ai_generated_images` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `image_id` VARCHAR(64) NOT NULL COMMENT '图像ID',
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
  `application_id` BIGINT NOT NULL COMMENT '应用ID',
  `model_id` BIGINT NOT NULL COMMENT '模型ID',
  `prompt` TEXT NOT NULL COMMENT '提示词',
  `negative_prompt` TEXT DEFAULT NULL COMMENT '负面提示词',
  `image_url` VARCHAR(255) NOT NULL COMMENT '图像URL',
  `thumbnail_url` VARCHAR(255) DEFAULT NULL COMMENT '缩略图URL',
  `width` INT DEFAULT 512 COMMENT '宽度',
  `height` INT DEFAULT 512 COMMENT '高度',
  `parameters` JSON DEFAULT NULL COMMENT '生成参数(JSON格式)',
  `duration` INT DEFAULT 0 COMMENT '生成时长(毫秒)',
  `status` VARCHAR(32) DEFAULT 'completed' COMMENT '状态：pending-待处理 processing-处理中 completed-已完成 failed-失败',
  `error_message` VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
  `feedback` VARCHAR(32) DEFAULT NULL COMMENT '反馈：like-点赞 dislike-点踩',
  `is_public` TINYINT(1) DEFAULT 0 COMMENT '是否公开：0-私有 1-公开',
  `cost` DECIMAL(10, 6) DEFAULT 0 COMMENT '成本',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_image_id` (`image_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_public` (`is_public`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_ai_generated_images_application_id` FOREIGN KEY (`application_id`) REFERENCES `ai_applications` (`id`),
  CONSTRAINT `fk_ai_generated_images_model_id` FOREIGN KEY (`model_id`) REFERENCES `ai_models` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI图像生成表';

-- 初始化基础AI模型
INSERT INTO `ai_models` (`name`, `code`, `description`, `model_type`, `provider`, `version`, `status`, `is_default`) VALUES
('GPT-3.5 Turbo', 'gpt-3.5-turbo', 'OpenAI GPT-3.5 Turbo模型', 'nlp', 'openai', '3.5', 1, 1),
('GPT-4', 'gpt-4', 'OpenAI GPT-4模型', 'nlp', 'openai', '4.0', 1, 0),
('DALL-E 3', 'dall-e-3', 'OpenAI DALL-E 3图像生成模型', 'cv', 'openai', '3.0', 1, 0),
('Embedding Model', 'text-embedding-ada-002', '文本嵌入模型', 'nlp', 'openai', '002', 1, 1);

-- 初始化基础AI应用
INSERT INTO `ai_applications` (`name`, `code`, `description`, `application_type`, `model_id`, `prompt_template`, `status`, `is_public`) VALUES
('智能客服', 'smart-assistant', '智能客服聊天机器人', 'chatbot', 1, '你是一个专业、友好的客服助手，请根据用户的问题提供准确的回答。', 1, 1),
('内容创作助手', 'content-creator', '帮助用户创作各类内容', 'content', 1, '你是一个专业的内容创作助手，请根据用户的需求生成高质量的内容。', 1, 1),
('图像生成器', 'image-generator', '根据文本描述生成图像', 'image', 3, '根据以下描述生成图像：', 1, 1);

-- 初始化基础AI知识库
INSERT INTO `ai_knowledge_bases` (`name`, `code`, `description`, `embeddings_model`, `status`) VALUES
('产品知识库', 'product-kb', '包含产品信息、使用说明和常见问题', 'text-embedding-ada-002', 1),
('帮助文档', 'help-docs', '系统帮助文档和教程', 'text-embedding-ada-002', 1);
