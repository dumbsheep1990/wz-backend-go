-- 后端数据库架构
CREATE DATABASE IF NOT EXISTS wz_backend DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE wz_backend;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1: 活跃, 0: 未激活',
    is_verified TINYINT NOT NULL DEFAULT 0 COMMENT '1: 已验证, 0: 未验证',
    is_company_verified TINYINT NOT NULL DEFAULT 0 COMMENT '1: 已验证, 0: 未验证',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_username (username),
    INDEX idx_users_email (email),
    INDEX idx_users_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户详细信息表
CREATE TABLE IF NOT EXISTS user_details (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    real_name VARCHAR(50),
    id_card VARCHAR(30),
    avatar VARCHAR(255),
    gender TINYINT COMMENT '1: 男, 2: 女, 0: 未知',
    birthday DATE,
    address VARCHAR(200),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_id (user_id),
    CONSTRAINT fk_user_details_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 企业认证表
CREATE TABLE IF NOT EXISTS company_verifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    company_type TINYINT NOT NULL COMMENT '1: 企业, 2: 集团, 3: 政府机构, 4: 上市公司',
    company_name VARCHAR(100) NOT NULL,
    -- 通用字段
    business_license VARCHAR(255) COMMENT '营业执照，企业、集团、上市公司需要',
    committee_letter VARCHAR(255) COMMENT '委托书，所有类型可选',
    -- 企业特定字段
    org_code_cert VARCHAR(255) COMMENT '组织机构代码证，企业类型需要',
    agency_cert VARCHAR(255) COMMENT '代理机构证明，企业类型可选',
    -- 集团特定字段
    org_structure VARCHAR(255) COMMENT '组织架构说明，集团类型需要',
    -- 政府机构特定字段
    unified_social_credit VARCHAR(255) COMMENT '统一社会信用代码证，政府机构类型需要',
    -- 上市公司特定字段
    listing_cert VARCHAR(255) COMMENT '上市证明，上市公司类型需要',
    -- 其他通用字段
    contact_person VARCHAR(50) NOT NULL,
    contact_phone VARCHAR(20) NOT NULL,
    uploaded_document TEXT COMMENT '上传的机构设立文本，适用于所有类型',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '0: 待审核, 1: 已通过, 2: 已拒绝',
    remark VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_id (user_id),
    CONSTRAINT fk_company_verifications_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户登录日志表
CREATE TABLE IF NOT EXISTS user_login_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    login_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    login_ip VARCHAR(50),
    user_agent VARCHAR(255),
    device_type TINYINT COMMENT '1: 网页, 2: 移动端, 3: 平板, 4: 其他',
    login_status TINYINT NOT NULL COMMENT '1: 成功, 0: 失败',
    INDEX idx_user_login_logs_user_id (user_id),
    CONSTRAINT fk_user_login_logs_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 企业入驻表
CREATE TABLE IF NOT EXISTS enterprise_registrations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    company_name VARCHAR(100) NOT NULL,
    company_type TINYINT NOT NULL COMMENT '1: 企业, 2: 集团, 3: 政府机构/NGO/协会, 4: 科研所',
    contact_person VARCHAR(50) NOT NULL,
    job_position VARCHAR(50) NOT NULL,
    region VARCHAR(100) NOT NULL,
    verification_method VARCHAR(50) NOT NULL,
    detailed_address VARCHAR(255) NOT NULL,
    location_latitude DECIMAL(10, 7),
    location_longitude DECIMAL(10, 7),
    status TINYINT NOT NULL DEFAULT 0 COMMENT '0: 待审核, 1: 已通过, 2: 已拒绝',
    remark VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_id (user_id),
    CONSTRAINT fk_enterprise_registrations_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分类表
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(200),
    parent_id BIGINT NOT NULL DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1: 启用, 0: 禁用',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_categories_parent_id (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 帖子表
CREATE TABLE IF NOT EXISTS posts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1: 正常, 0: 已删除, 2: 审核中, 3: 已拒绝',
    view_count INT NOT NULL DEFAULT 0,
    like_count INT NOT NULL DEFAULT 0,
    comment_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_posts_user_id (user_id),
    INDEX idx_posts_category_id (category_id),
    INDEX idx_posts_status (status),
    CONSTRAINT fk_posts_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_posts_category_id FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 评论表
CREATE TABLE IF NOT EXISTS reviews (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1: 正常, 0: 已删除, 2: 审核中, 3: 已拒绝',
    like_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_reviews_post_id (post_id),
    INDEX idx_reviews_user_id (user_id),
    INDEX idx_reviews_status (status),
    CONSTRAINT fk_reviews_post_id FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_reviews_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 内容状态日志表
CREATE TABLE IF NOT EXISTS content_status_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource_type VARCHAR(20) NOT NULL COMMENT '帖子, 评论等',
    resource_id BIGINT NOT NULL,
    status TINYINT NOT NULL COMMENT '1: 正常, 0: 已删除, 2: 审核中, 3: 已拒绝',
    reason VARCHAR(255),
    operator_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_content_status_logs_resource (resource_type, resource_id),
    INDEX idx_content_status_logs_operator_id (operator_id),
    CONSTRAINT fk_content_status_logs_operator_id FOREIGN KEY (operator_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 热门内容表
CREATE TABLE IF NOT EXISTS hot_contents (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource_type VARCHAR(20) NOT NULL COMMENT '帖子, 评论等',
    resource_id BIGINT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    operator_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_hot_contents_resource (resource_type, resource_id),
    INDEX idx_hot_contents_type_order (resource_type, sort_order),
    CONSTRAINT fk_hot_contents_operator_id FOREIGN KEY (operator_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户行为日志表
CREATE TABLE IF NOT EXISTS user_behavior_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    action VARCHAR(20) NOT NULL COMMENT '查看, 点赞, 评论等',
    resource_type VARCHAR(20) NOT NULL COMMENT '帖子, 评论等',
    resource_id BIGINT NOT NULL,
    ip VARCHAR(50),
    user_agent VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_behavior_logs_user_id (user_id),
    INDEX idx_user_behavior_logs_resource (resource_type, resource_id),
    CONSTRAINT fk_user_behavior_logs_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
