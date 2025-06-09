-- +migrate Up
-- 用户积分表
CREATE TABLE IF NOT EXISTS user_points (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    points INT NOT NULL COMMENT '积分值',
    total_points INT NOT NULL COMMENT '用户总积分（冗余字段）',
    type TINYINT NOT NULL COMMENT '积分类型：1-增加，2-减少',
    source VARCHAR(50) NOT NULL COMMENT '积分来源：sign-签到，purchase-购买，activity-活动等',
    description VARCHAR(200) NOT NULL COMMENT '积分描述',
    related_id BIGINT DEFAULT 0 COMMENT '关联ID',
    related_type VARCHAR(50) DEFAULT '' COMMENT '关联类型',
    tenant_id BIGINT DEFAULT 0 COMMENT '租户ID',
    operator_id BIGINT DEFAULT 0 COMMENT '操作者ID',
    is_revoked BOOLEAN DEFAULT FALSE COMMENT '是否已撤销',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_source (source),
    INDEX idx_created_at (created_at),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分记录表';

-- 积分规则表
CREATE TABLE IF NOT EXISTS points_rules (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sign_in_points INT NOT NULL DEFAULT 5 COMMENT '签到积分',
    comment_points INT NOT NULL DEFAULT 2 COMMENT '评论积分',
    share_points INT NOT NULL DEFAULT 3 COMMENT '分享积分',
    article_points INT NOT NULL DEFAULT 10 COMMENT '发布文章积分',
    invite_points INT NOT NULL DEFAULT 20 COMMENT '邀请积分',
    purchase_rate INT NOT NULL DEFAULT 10 COMMENT '购买积分比例（元:积分=1:10）',
    max_daily_points INT NOT NULL DEFAULT 100 COMMENT '每日最大获取积分',
    enable_exchange BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否可兑换商品',
    exchange_rate INT NOT NULL DEFAULT 100 COMMENT '兑换比例（积分:元=100:1）',
    min_exchange_points INT NOT NULL DEFAULT 1000 COMMENT '最小兑换积分',
    tenant_id BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分规则表';

-- 创建默认的积分规则
INSERT INTO points_rules (tenant_id) VALUES (0);

-- +migrate Down
DROP TABLE IF EXISTS user_points;
DROP TABLE IF EXISTS points_rules; 