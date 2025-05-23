-- 搜索引擎和数据索引模块数据表结构

-- 索引配置表
CREATE TABLE IF NOT EXISTS `search_indices` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '索引名称',
  `code` VARCHAR(64) NOT NULL COMMENT '索引编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '索引描述',
  `entity_type` VARCHAR(64) NOT NULL COMMENT '实体类型：user, content, product, etc.',
  `index_type` VARCHAR(32) DEFAULT 'standard' COMMENT '索引类型：standard-标准 fulltext-全文 vector-向量',
  `analyzer` VARCHAR(32) DEFAULT 'standard' COMMENT '分析器：standard, keyword, simple, whitespace, etc.',
  `mapping_definition` TEXT DEFAULT NULL COMMENT '映射定义(JSON格式)',
  `settings` JSON DEFAULT NULL COMMENT '索引设置(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `rebuild_interval` INT DEFAULT 0 COMMENT '重建间隔(分钟)，0表示不自动重建',
  `last_rebuild_time` TIMESTAMP NULL COMMENT '最后重建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示系统索引',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_entity_type` (`entity_type`),
  KEY `idx_index_type` (`index_type`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='索引配置表';

-- 索引字段表
CREATE TABLE IF NOT EXISTS `search_index_fields` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `index_id` BIGINT NOT NULL COMMENT '索引ID',
  `field_name` VARCHAR(64) NOT NULL COMMENT '字段名称',
  `display_name` VARCHAR(64) DEFAULT NULL COMMENT '显示名称',
  `field_type` VARCHAR(32) NOT NULL COMMENT '字段类型：text, keyword, integer, float, date, boolean, geo_point, object, nested, etc.',
  `is_array` TINYINT(1) DEFAULT 0 COMMENT '是否数组：0-否 1-是',
  `is_searchable` TINYINT(1) DEFAULT 1 COMMENT '是否可搜索：0-否 1-是',
  `is_filterable` TINYINT(1) DEFAULT 0 COMMENT '是否可过滤：0-否 1-是',
  `is_sortable` TINYINT(1) DEFAULT 0 COMMENT '是否可排序：0-否 1-是',
  `is_primary` TINYINT(1) DEFAULT 0 COMMENT '是否主键：0-否 1-是',
  `is_required` TINYINT(1) DEFAULT 0 COMMENT '是否必需：0-否 1-是',
  `boost` FLOAT DEFAULT 1.0 COMMENT '权重提升',
  `analyzer` VARCHAR(32) DEFAULT NULL COMMENT '分析器',
  `search_analyzer` VARCHAR(32) DEFAULT NULL COMMENT '搜索分析器',
  `field_mapping` TEXT DEFAULT NULL COMMENT '字段映射(JSON格式)',
  `data_source` VARCHAR(255) DEFAULT NULL COMMENT '数据来源',
  `transform_expression` TEXT DEFAULT NULL COMMENT '转换表达式',
  `sort_order` INT DEFAULT 0 COMMENT '排序顺序',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_index_field` (`index_id`, `field_name`),
  KEY `idx_index_id` (`index_id`),
  KEY `idx_is_searchable` (`is_searchable`),
  KEY `idx_is_filterable` (`is_filterable`),
  KEY `idx_is_sortable` (`is_sortable`),
  KEY `idx_sort_order` (`sort_order`),
  CONSTRAINT `fk_search_index_fields_index_id` FOREIGN KEY (`index_id`) REFERENCES `search_indices` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='索引字段表';

-- 搜索同义词表
CREATE TABLE IF NOT EXISTS `search_synonyms` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `index_id` BIGINT DEFAULT NULL COMMENT '索引ID，为NULL表示全局同义词',
  `source_term` VARCHAR(255) NOT NULL COMMENT '源词',
  `target_terms` TEXT NOT NULL COMMENT '目标词，多个用逗号分隔',
  `synonym_type` VARCHAR(32) DEFAULT 'bidirectional' COMMENT '同义词类型：bidirectional-双向 directional-单向',
  `language` VARCHAR(16) DEFAULT 'zh_CN' COMMENT '语言',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_index_id` (`index_id`),
  KEY `idx_source_term` (`source_term`),
  KEY `idx_language` (`language`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_search_synonyms_index_id` FOREIGN KEY (`index_id`) REFERENCES `search_indices` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='搜索同义词表';

-- 停用词表
CREATE TABLE IF NOT EXISTS `search_stop_words` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `index_id` BIGINT DEFAULT NULL COMMENT '索引ID，为NULL表示全局停用词',
  `word` VARCHAR(64) NOT NULL COMMENT '停用词',
  `language` VARCHAR(16) DEFAULT 'zh_CN' COMMENT '语言',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_index_word_language` (`index_id`, `word`, `language`),
  KEY `idx_index_id` (`index_id`),
  KEY `idx_word` (`word`),
  KEY `idx_language` (`language`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='停用词表';

-- 搜索模板表
CREATE TABLE IF NOT EXISTS `search_templates` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '模板名称',
  `code` VARCHAR(64) NOT NULL COMMENT '模板编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '模板描述',
  `template_type` VARCHAR(32) DEFAULT 'query' COMMENT '模板类型：query-查询 filter-过滤 aggregation-聚合',
  `template_content` TEXT NOT NULL COMMENT '模板内容(JSON格式)',
  `params` TEXT DEFAULT NULL COMMENT '参数定义(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_template_type` (`template_type`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='搜索模板表';

-- 搜索记录表
CREATE TABLE IF NOT EXISTS `search_records` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT DEFAULT NULL COMMENT '用户ID，为NULL表示匿名用户',
  `query` VARCHAR(255) NOT NULL COMMENT '搜索关键词',
  `index_codes` VARCHAR(255) DEFAULT NULL COMMENT '索引编码，多个用逗号分隔',
  `filters` JSON DEFAULT NULL COMMENT '过滤条件(JSON格式)',
  `sorts` VARCHAR(255) DEFAULT NULL COMMENT '排序条件，格式为field:direction，多个用逗号分隔',
  `page` INT DEFAULT 1 COMMENT '页码',
  `page_size` INT DEFAULT 10 COMMENT '每页条数',
  `result_count` INT DEFAULT NULL COMMENT '结果数量',
  `search_time` INT DEFAULT NULL COMMENT '搜索耗时(毫秒)',
  `client_ip` VARCHAR(64) DEFAULT NULL COMMENT '客户端IP',
  `user_agent` VARCHAR(255) DEFAULT NULL COMMENT '用户代理',
  `search_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '搜索时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_query` (`query`),
  KEY `idx_search_time` (`search_time`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='搜索记录表';

-- 热门搜索表
CREATE TABLE IF NOT EXISTS `search_hot_words` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `word` VARCHAR(255) NOT NULL COMMENT '搜索词',
  `count` INT NOT NULL DEFAULT 1 COMMENT '搜索次数',
  `category` VARCHAR(64) DEFAULT 'default' COMMENT '分类',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `is_recommend` TINYINT(1) DEFAULT 0 COMMENT '是否推荐：0-否 1-是',
  `recommend_sort` INT DEFAULT 0 COMMENT '推荐排序',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `date` DATE NOT NULL COMMENT '统计日期',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_word_date_tenant` (`word`, `date`, `tenant_id`),
  KEY `idx_count` (`count`),
  KEY `idx_category` (`category`),
  KEY `idx_status` (`status`),
  KEY `idx_is_recommend` (`is_recommend`),
  KEY `idx_date` (`date`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='热门搜索表';

-- 搜索建议表
CREATE TABLE IF NOT EXISTS `search_suggestions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `word` VARCHAR(255) NOT NULL COMMENT '建议词',
  `source_type` VARCHAR(32) DEFAULT 'auto' COMMENT '来源类型：auto-自动生成 manual-手动添加',
  `weight` INT DEFAULT 0 COMMENT '权重',
  `category` VARCHAR(64) DEFAULT 'default' COMMENT '分类',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_word_tenant` (`word`, `tenant_id`),
  KEY `idx_weight` (`weight`),
  KEY `idx_category` (`category`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='搜索建议表';

-- 索引任务表
CREATE TABLE IF NOT EXISTS `search_index_tasks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `index_id` BIGINT NOT NULL COMMENT '索引ID',
  `task_type` VARCHAR(32) NOT NULL COMMENT '任务类型：create-创建索引 rebuild-重建索引 update-更新索引 delete-删除索引',
  `status` VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending-待处理 processing-处理中 completed-已完成 failed-失败',
  `progress` INT DEFAULT 0 COMMENT '进度(0-100)',
  `total_records` INT DEFAULT 0 COMMENT '总记录数',
  `processed_records` INT DEFAULT 0 COMMENT '已处理记录数',
  `start_time` TIMESTAMP NULL COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `params` JSON DEFAULT NULL COMMENT '任务参数(JSON格式)',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建者ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_index_id` (`index_id`),
  KEY `idx_task_type` (`task_type`),
  KEY `idx_status` (`status`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_search_index_tasks_index_id` FOREIGN KEY (`index_id`) REFERENCES `search_indices` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='索引任务表';

-- 搜索过滤器表
CREATE TABLE IF NOT EXISTS `search_filters` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '过滤器名称',
  `code` VARCHAR(64) NOT NULL COMMENT '过滤器编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '过滤器描述',
  `filter_type` VARCHAR(32) NOT NULL COMMENT '过滤器类型：term-精确匹配 range-范围 geo-地理位置 custom-自定义',
  `field_name` VARCHAR(64) NOT NULL COMMENT '字段名称',
  `display_name` VARCHAR(64) DEFAULT NULL COMMENT '显示名称',
  `options` JSON DEFAULT NULL COMMENT '选项配置(JSON格式)',
  `default_value` TEXT DEFAULT NULL COMMENT '默认值',
  `sort_order` INT DEFAULT 0 COMMENT '排序顺序',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `index_id` BIGINT NOT NULL COMMENT '索引ID',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_index` (`code`, `index_id`),
  KEY `idx_index_id` (`index_id`),
  KEY `idx_filter_type` (`filter_type`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_search_filters_index_id` FOREIGN KEY (`index_id`) REFERENCES `search_indices` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='搜索过滤器表';

-- 向量存储配置表
CREATE TABLE IF NOT EXISTS `vector_store_configs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '向量存储名称',
  `code` VARCHAR(64) NOT NULL COMMENT '向量存储编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '向量存储描述',
  `store_type` VARCHAR(32) DEFAULT 'faiss' COMMENT '存储类型：faiss, milvus, pinecone, etc.',
  `dimension` INT NOT NULL COMMENT '向量维度',
  `metric_type` VARCHAR(32) DEFAULT 'cosine' COMMENT '相似度度量类型：cosine, euclidean, dot, etc.',
  `index_type` VARCHAR(32) DEFAULT 'hnsw' COMMENT '索引类型：hnsw, flat, ivf, etc.',
  `index_params` JSON DEFAULT NULL COMMENT '索引参数(JSON格式)',
  `search_params` JSON DEFAULT NULL COMMENT '搜索参数(JSON格式)',
  `model_name` VARCHAR(128) DEFAULT NULL COMMENT '嵌入模型名称',
  `connection_params` JSON DEFAULT NULL COMMENT '连接参数(JSON格式)',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_store_type` (`store_type`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='向量存储配置表';

-- 向量数据表
CREATE TABLE IF NOT EXISTS `vector_data` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `store_id` BIGINT NOT NULL COMMENT '向量存储ID',
  `vector_id` VARCHAR(64) NOT NULL COMMENT '向量ID',
  `object_id` VARCHAR(128) DEFAULT NULL COMMENT '对象ID',
  `object_type` VARCHAR(64) DEFAULT NULL COMMENT '对象类型',
  `metadata` JSON DEFAULT NULL COMMENT '元数据(JSON格式)',
  `vector_data` LONGBLOB DEFAULT NULL COMMENT '向量数据',
  `text_data` TEXT DEFAULT NULL COMMENT '文本数据',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_store_vector` (`store_id`, `vector_id`),
  KEY `idx_store_id` (`store_id`),
  KEY `idx_object_id_type` (`object_id`, `object_type`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_vector_data_store_id` FOREIGN KEY (`store_id`) REFERENCES `vector_store_configs` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='向量数据表';

-- 初始化基础索引配置
INSERT INTO `search_indices` (`name`, `code`, `description`, `entity_type`, `index_type`, `status`) VALUES
('用户索引', 'user-index', '用户数据搜索索引', 'user', 'standard', 1),
('内容索引', 'content-index', '内容数据搜索索引', 'content', 'fulltext', 1),
('产品索引', 'product-index', '产品数据搜索索引', 'product', 'standard', 1),
('全站搜索', 'global-search', '全站统一搜索索引', 'global', 'fulltext', 1);

-- 初始化基础用户索引字段
INSERT INTO `search_index_fields` (`index_id`, `field_name`, `display_name`, `field_type`, `is_searchable`, `is_filterable`, `is_sortable`, `is_primary`, `sort_order`) VALUES
(1, 'id', 'ID', 'keyword', 0, 1, 1, 1, 1),
(1, 'username', '用户名', 'text', 1, 0, 0, 0, 2),
(1, 'nickname', '昵称', 'text', 1, 0, 0, 0, 3),
(1, 'email', '邮箱', 'keyword', 1, 1, 0, 0, 4),
(1, 'phone', '手机号', 'keyword', 1, 1, 0, 0, 5),
(1, 'status', '状态', 'integer', 0, 1, 1, 0, 6),
(1, 'created_at', '创建时间', 'date', 0, 1, 1, 0, 7);

-- 初始化基础内容索引字段
INSERT INTO `search_index_fields` (`index_id`, `field_name`, `display_name`, `field_type`, `is_searchable`, `is_filterable`, `is_sortable`, `is_primary`, `sort_order`) VALUES
(2, 'id', 'ID', 'keyword', 0, 1, 1, 1, 1),
(2, 'title', '标题', 'text', 1, 0, 0, 0, 2),
(2, 'content', '内容', 'text', 1, 0, 0, 0, 3),
(2, 'summary', '摘要', 'text', 1, 0, 0, 0, 4),
(2, 'category_id', '分类ID', 'keyword', 0, 1, 0, 0, 5),
(2, 'tags', '标签', 'keyword', 1, 1, 0, 0, 6),
(2, 'status', '状态', 'integer', 0, 1, 1, 0, 7),
(2, 'published_at', '发布时间', 'date', 0, 1, 1, 0, 8),
(2, 'view_count', '浏览数', 'integer', 0, 1, 1, 0, 9);

-- 初始化基础产品索引字段
INSERT INTO `search_index_fields` (`index_id`, `field_name`, `display_name`, `field_type`, `is_searchable`, `is_filterable`, `is_sortable`, `is_primary`, `sort_order`) VALUES
(3, 'id', 'ID', 'keyword', 0, 1, 1, 1, 1),
(3, 'name', '名称', 'text', 1, 0, 0, 0, 2),
(3, 'description', '描述', 'text', 1, 0, 0, 0, 3),
(3, 'price', '价格', 'float', 0, 1, 1, 0, 4),
(3, 'category_id', '分类ID', 'keyword', 0, 1, 0, 0, 5),
(3, 'tags', '标签', 'keyword', 1, 1, 0, 0, 6),
(3, 'status', '状态', 'integer', 0, 1, 1, 0, 7),
(3, 'created_at', '创建时间', 'date', 0, 1, 1, 0, 8),
(3, 'sales_count', '销量', 'integer', 0, 1, 1, 0, 9);
