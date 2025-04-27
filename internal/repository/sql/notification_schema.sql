-- 通知服务数据库表结构

-- 通知表
CREATE TABLE IF NOT EXISTS `notifications` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_id` bigint NOT NULL COMMENT '用户ID',
    `type` varchar(20) NOT NULL COMMENT '通知类型',
    `content` text NOT NULL COMMENT '通知内容',
    `read` tinyint NOT NULL DEFAULT '0' COMMENT '是否已读',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user` (`user_id`),
    KEY `idx_read` (`read`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知表';

-- 邮件通知记录表
CREATE TABLE IF NOT EXISTS `email_notifications` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `to` varchar(255) NOT NULL COMMENT '收件人',
    `subject` varchar(255) NOT NULL COMMENT '主题',
    `content` text NOT NULL COMMENT '内容',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-待发送，1-已发送，2-发送失败',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='邮件通知记录表';

-- 短信通知记录表
CREATE TABLE IF NOT EXISTS `sms_notifications` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `phone` varchar(20) NOT NULL COMMENT '手机号',
    `content` text NOT NULL COMMENT '内容',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-待发送，1-已发送，2-发送失败',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信通知记录表'; 