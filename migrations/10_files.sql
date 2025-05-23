-- 文件存储模块数据表结构

-- 文件信息表
CREATE TABLE IF NOT EXISTS `files` (
  `id` VARCHAR(64) NOT NULL COMMENT '文件ID',
  `file_url` VARCHAR(255) NOT NULL COMMENT '文件URL',
  `file_type` VARCHAR(32) NOT NULL COMMENT '文件类型',
  `file_name` VARCHAR(128) DEFAULT NULL COMMENT '文件名称',
  `file_ext` VARCHAR(32) DEFAULT NULL COMMENT '文件扩展名',
  `file_size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小(字节)',
  `mime_type` VARCHAR(128) DEFAULT NULL COMMENT 'MIME类型',
  `storage_type` VARCHAR(32) NOT NULL DEFAULT 'local' COMMENT '存储类型：local, oss, cos, s3',
  `storage_path` VARCHAR(255) DEFAULT NULL COMMENT '存储路径',
  `user_id` VARCHAR(64) DEFAULT NULL COMMENT '上传用户ID',
  `tenant_id` VARCHAR(64) DEFAULT NULL COMMENT '租户ID',
  `resource_type` VARCHAR(32) DEFAULT NULL COMMENT '关联资源类型',
  `resource_id` VARCHAR(64) DEFAULT NULL COMMENT '关联资源ID',
  `is_public` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否公开：0-私有 1-公开',
  `is_temp` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否临时文件：0-否 1-是',
  `status` VARCHAR(32) NOT NULL DEFAULT 'active' COMMENT '状态：active, deleted',
  `md5` VARCHAR(32) DEFAULT NULL COMMENT '文件MD5',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_file_type` (`file_type`),
  KEY `idx_resource` (`resource_type`, `resource_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_md5` (`md5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件信息表';

-- 文件分类表
CREATE TABLE IF NOT EXISTS `file_categories` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '分类名称',
  `code` VARCHAR(32) NOT NULL COMMENT '分类编码',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '分类描述',
  `parent_id` BIGINT DEFAULT 0 COMMENT '父分类ID',
  `allowed_types` VARCHAR(255) DEFAULT NULL COMMENT '允许的文件类型(逗号分隔)',
  `max_size` BIGINT DEFAULT 0 COMMENT '最大文件大小(字节)',
  `sort_order` INT DEFAULT 0 COMMENT '排序顺序',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件分类表';

-- 文件访问日志表
CREATE TABLE IF NOT EXISTS `file_access_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `file_id` VARCHAR(64) NOT NULL COMMENT '文件ID',
  `user_id` VARCHAR(64) DEFAULT NULL COMMENT '访问用户ID',
  `ip` VARCHAR(64) DEFAULT NULL COMMENT '访问IP',
  `user_agent` VARCHAR(255) DEFAULT NULL COMMENT '用户代理',
  `referer` VARCHAR(255) DEFAULT NULL COMMENT '来源URL',
  `access_type` VARCHAR(32) NOT NULL DEFAULT 'view' COMMENT '访问类型：view, download',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '访问时间',
  PRIMARY KEY (`id`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_access_type` (`access_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件访问日志表';

-- 文件标签关联表
CREATE TABLE IF NOT EXISTS `file_tags` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `file_id` VARCHAR(64) NOT NULL COMMENT '文件ID',
  `tag_name` VARCHAR(64) NOT NULL COMMENT '标签名称',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_tag` (`file_id`, `tag_name`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_tag_name` (`tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件标签关联表';

-- 文件存储配置表
CREATE TABLE IF NOT EXISTS `file_storage_configs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `storage_type` VARCHAR(32) NOT NULL COMMENT '存储类型：local, oss, cos, s3',
  `provider` VARCHAR(64) NOT NULL COMMENT '提供商：aliyun, tencent, aws, etc',
  `name` VARCHAR(64) NOT NULL COMMENT '配置名称',
  `is_default` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否默认：0-否 1-是',
  `access_key` VARCHAR(128) DEFAULT NULL COMMENT 'AccessKey',
  `secret_key` VARCHAR(255) DEFAULT NULL COMMENT 'SecretKey',
  `endpoint` VARCHAR(255) DEFAULT NULL COMMENT '访问端点',
  `bucket` VARCHAR(64) DEFAULT NULL COMMENT '存储桶/容器名称',
  `region` VARCHAR(64) DEFAULT NULL COMMENT '区域',
  `base_url` VARCHAR(255) DEFAULT NULL COMMENT '基础URL',
  `directory` VARCHAR(255) DEFAULT NULL COMMENT '目录',
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_storage_type` (`storage_type`),
  KEY `idx_is_default` (`is_default`),
  KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件存储配置表';

-- 初始化文件分类数据
INSERT INTO `file_categories` (`name`, `code`, `description`, `parent_id`, `allowed_types`, `max_size`, `sort_order`) VALUES
('图片', 'image', '图片文件', 0, 'jpg,jpeg,png,gif,bmp,webp', 10485760, 1),
('文档', 'document', '文档文件', 0, 'doc,docx,xls,xlsx,ppt,pptx,pdf,txt', 20971520, 2),
('视频', 'video', '视频文件', 0, 'mp4,avi,mov,wmv,flv', 104857600, 3),
('音频', 'audio', '音频文件', 0, 'mp3,wav,ogg,flac,aac', 31457280, 4),
('压缩包', 'archive', '压缩文件', 0, 'zip,rar,7z,tar,gz', 52428800, 5),
('其他', 'other', '其他文件', 0, NULL, 10485760, 6);

-- 初始化文件存储配置数据
INSERT INTO `file_storage_configs` (`storage_type`, `provider`, `name`, `is_default`, `endpoint`, `bucket`, `base_url`, `directory`, `is_enabled`) VALUES
('local', 'local', '本地存储', 1, NULL, NULL, '/uploads', 'uploads', 1);
