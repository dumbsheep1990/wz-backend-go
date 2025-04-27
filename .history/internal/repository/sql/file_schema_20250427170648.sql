-- 文件服务数据库表结构

-- 文件表
CREATE TABLE IF NOT EXISTS `files` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `file_id` varchar(64) NOT NULL COMMENT '文件ID',
    `user_id` bigint NOT NULL COMMENT '用户ID',
    `file_type` varchar(20) NOT NULL COMMENT '文件类型：image/video',
    `file_name` varchar(255) NOT NULL COMMENT '文件名',
    `file_size` bigint NOT NULL COMMENT '文件大小',
    `file_url` varchar(255) NOT NULL COMMENT '文件URL',
    `storage_type` varchar(20) NOT NULL COMMENT '存储类型：local/object',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_file_id` (`file_id`),
    KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件表'; 