-- 搜索服务数据库表结构
USE wz_backend;

-- 搜索记录表
CREATE TABLE IF NOT EXISTS search_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT,
    keyword VARCHAR(200) NOT NULL,
    category VARCHAR(50),
    tags VARCHAR(500),  -- 逗号分隔的标签
    location VARCHAR(100),
    longitude DOUBLE,
    latitude DOUBLE,
    result_num INT NOT NULL DEFAULT 0,  -- 搜索结果数量
    device_id VARCHAR(100),
    ip VARCHAR(50),
    user_agent VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_search_logs_user_id (user_id),
    INDEX idx_search_logs_keyword (keyword),
    INDEX idx_search_logs_created_at (created_at),
    INDEX idx_search_logs_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 搜索建议表 - 存储常用搜索词，用于自动提示
CREATE TABLE IF NOT EXISTS search_suggestions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    keyword VARCHAR(200) NOT NULL,
    category VARCHAR(50),
    frequency INT NOT NULL DEFAULT 1,  -- 搜索频率
    is_enabled TINYINT NOT NULL DEFAULT 1,  -- 是否启用
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_suggestion_keyword_category (keyword, category),
    INDEX idx_suggestions_frequency (frequency)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 热搜词表
CREATE TABLE IF NOT EXISTS hot_searches (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    keyword VARCHAR(200) NOT NULL,
    category VARCHAR(50),
    count INT NOT NULL DEFAULT 0,  -- 搜索次数
    trend TINYINT NOT NULL DEFAULT 0,  -- 1:上升，0:持平，-1:下降
    is_promoted TINYINT NOT NULL DEFAULT 0,  -- 是否推广
    sort_order INT NOT NULL DEFAULT 0,  -- 排序权重，越大越靠前
    operator_id BIGINT,  -- 操作员ID
    start_time TIMESTAMP NULL,  -- 热搜开始时间
    end_time TIMESTAMP NULL,  -- 热搜结束时间
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_hot_keyword_category (keyword, category),
    INDEX idx_hot_searches_count (count),
    INDEX idx_hot_searches_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 搜索统计表 - 按天统计
CREATE TABLE IF NOT EXISTS search_statistics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    stat_date DATE NOT NULL,  -- 统计日期
    keyword VARCHAR(200) NOT NULL,
    category VARCHAR(50),
    search_count INT NOT NULL DEFAULT 0,  -- 搜索次数
    user_count INT NOT NULL DEFAULT 0,  -- 用户数
    result_count INT NOT NULL DEFAULT 0,  -- 结果总数
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_stat_date_keyword_category (stat_date, keyword, category),
    INDEX idx_statistics_date (stat_date),
    INDEX idx_statistics_count (search_count)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 搜索同义词表 - 用于同义词扩展
CREATE TABLE IF NOT EXISTS search_synonyms (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    word VARCHAR(100) NOT NULL,
    synonyms VARCHAR(1000) NOT NULL,  -- 逗号分隔的同义词列表
    is_enabled TINYINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_synonym_word (word),
    INDEX idx_synonyms_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 搜索停用词表 - 用于过滤无用词
CREATE TABLE IF NOT EXISTS search_stopwords (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    word VARCHAR(100) NOT NULL,
    is_enabled TINYINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_stopword (word),
    INDEX idx_stopwords_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 搜索索引配置表 - 用于管理各类内容的索引配置
CREATE TABLE IF NOT EXISTS search_index_configs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource_type VARCHAR(50) NOT NULL,  -- 资源类型：post, review, etc.
    index_name VARCHAR(100) NOT NULL,  -- 索引名称
    fields VARCHAR(1000) NOT NULL,  -- 索引字段，JSON格式
    boost_fields VARCHAR(1000),  -- 权重字段，JSON格式
    filter_fields VARCHAR(1000),  -- 过滤字段，JSON格式
    is_enabled TINYINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_index_config_type (resource_type),
    INDEX idx_index_configs_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 搜索排序规则表 - 用于自定义搜索结果排序规则
CREATE TABLE IF NOT EXISTS search_sort_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    sort_fields VARCHAR(1000) NOT NULL,  -- 排序字段，JSON格式
    is_default TINYINT NOT NULL DEFAULT 0,  -- 是否默认
    is_enabled TINYINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_sort_rule_name (name),
    INDEX idx_sort_rules_default (is_default),
    INDEX idx_sort_rules_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
