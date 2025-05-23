-- 系统设置模块数据表结构

-- 系统设置表
CREATE TABLE IF NOT EXISTS `system_settings` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `key` VARCHAR(128) NOT NULL COMMENT '设置键',
  `value` TEXT DEFAULT NULL COMMENT '设置值',
  `type` VARCHAR(32) DEFAULT 'string' COMMENT '设置类型：string, number, boolean, json, array',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '设置描述',
  `group` VARCHAR(64) DEFAULT 'default' COMMENT '设置分组',
  `is_public` TINYINT(1) DEFAULT 1 COMMENT '是否公开：0-私有（仅管理员可见） 1-公开',
  `is_system` TINYINT(1) DEFAULT 0 COMMENT '是否系统内置：0-否 1-是',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示全局设置',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_key_tenant` (`key`, `tenant_id`),
  KEY `idx_group` (`group`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_is_public` (`is_public`),
  KEY `idx_is_system` (`is_system`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统设置表';

-- 系统参数表
CREATE TABLE IF NOT EXISTS `sys_params` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `param_name` VARCHAR(128) NOT NULL COMMENT '参数名称',
  `param_key` VARCHAR(128) NOT NULL COMMENT '参数键',
  `param_value` TEXT DEFAULT NULL COMMENT '参数值',
  `param_type` VARCHAR(32) DEFAULT 'string' COMMENT '参数类型：string, number, boolean, json, array',
  `param_group` VARCHAR(64) DEFAULT 'default' COMMENT '参数分组',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '参数描述',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `is_system` TINYINT(1) DEFAULT 0 COMMENT '是否系统内置：0-否 1-是',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示全局参数',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_param_key_tenant` (`param_key`, `tenant_id`),
  KEY `idx_param_group` (`param_group`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_system` (`is_system`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统参数表';

-- 系统配置表
CREATE TABLE IF NOT EXISTS `system_configs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `site_name` VARCHAR(128) DEFAULT NULL COMMENT '站点名称',
  `site_logo` VARCHAR(255) DEFAULT NULL COMMENT '站点Logo',
  `site_favicon` VARCHAR(255) DEFAULT NULL COMMENT '站点图标',
  `site_description` VARCHAR(255) DEFAULT NULL COMMENT '站点描述',
  `site_keywords` VARCHAR(255) DEFAULT NULL COMMENT '站点关键词',
  `site_copyright` VARCHAR(255) DEFAULT NULL COMMENT '站点版权信息',
  `site_icp` VARCHAR(64) DEFAULT NULL COMMENT 'ICP备案号',
  `site_psr` VARCHAR(64) DEFAULT NULL COMMENT '公安备案号',
  `site_domain` VARCHAR(128) DEFAULT NULL COMMENT '站点域名',
  `site_status` INT DEFAULT 1 COMMENT '站点状态：0-关闭 1-开启',
  `site_close_message` TEXT DEFAULT NULL COMMENT '站点关闭提示信息',
  `logo_config` JSON DEFAULT NULL COMMENT 'Logo配置(JSON)',
  `login_config` JSON DEFAULT NULL COMMENT '登录配置(JSON)',
  `register_config` JSON DEFAULT NULL COMMENT '注册配置(JSON)',
  `upload_config` JSON DEFAULT NULL COMMENT '上传配置(JSON)',
  `payment_config` JSON DEFAULT NULL COMMENT '支付配置(JSON)',
  `mail_config` JSON DEFAULT NULL COMMENT '邮件配置(JSON)',
  `sms_config` JSON DEFAULT NULL COMMENT '短信配置(JSON)',
  `oss_config` JSON DEFAULT NULL COMMENT '对象存储配置(JSON)',
  `cdn_config` JSON DEFAULT NULL COMMENT 'CDN配置(JSON)',
  `cache_config` JSON DEFAULT NULL COMMENT '缓存配置(JSON)',
  `security_config` JSON DEFAULT NULL COMMENT '安全配置(JSON)',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示全局配置',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 字典表
CREATE TABLE IF NOT EXISTS `sys_dictionaries` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL COMMENT '字典名称',
  `type` VARCHAR(128) NOT NULL COMMENT '字典类型',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '字典描述',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，为NULL表示全局字典',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_type_tenant` (`type`, `tenant_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字典表';

-- 字典详情表
CREATE TABLE IF NOT EXISTS `sys_dictionary_details` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `dictionary_id` BIGINT NOT NULL COMMENT '字典ID',
  `label` VARCHAR(128) NOT NULL COMMENT '字典标签',
  `value` VARCHAR(255) NOT NULL COMMENT '字典值',
  `color` VARCHAR(32) DEFAULT NULL COMMENT '颜色',
  `sort` INT DEFAULT 0 COMMENT '排序',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_dictionary_id` (`dictionary_id`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_dictionary_details_dictionary_id` FOREIGN KEY (`dictionary_id`) REFERENCES `sys_dictionaries` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字典详情表';

-- 初始化基本系统设置
INSERT INTO `system_settings` (`key`, `value`, `type`, `description`, `group`, `is_public`, `is_system`) VALUES
('site_name', 'WZ Platform', 'string', '站点名称', 'site', 1, 1),
('site_logo', '/static/logo.png', 'string', '站点Logo', 'site', 1, 1),
('site_favicon', '/static/favicon.ico', 'string', '站点图标', 'site', 1, 1),
('site_description', 'WZ开发平台', 'string', '站点描述', 'site', 1, 1),
('site_keywords', 'WZ,平台,开发', 'string', '站点关键词', 'site', 1, 1),
('site_status', '1', 'number', '站点状态：0-关闭 1-开启', 'site', 1, 1),
('register_enabled', '1', 'boolean', '是否允许注册', 'user', 1, 1),
('login_captcha', '1', 'boolean', '登录是否需要验证码', 'security', 1, 1),
('default_storage', 'local', 'string', '默认存储方式', 'storage', 0, 1),
('max_upload_size', '10', 'number', '最大上传大小(MB)', 'storage', 1, 1);

-- 初始化基本字典
INSERT INTO `sys_dictionaries` (`name`, `type`, `description`, `status`) VALUES
('性别', 'gender', '用户性别', 1),
('状态', 'status', '通用状态', 1),
('是否', 'yes_no', '是否选项', 1);

-- 初始化字典详情
INSERT INTO `sys_dictionary_details` (`dictionary_id`, `label`, `value`, `color`, `sort`, `status`) VALUES
(1, '男', '1', 'blue', 1, 1),
(1, '女', '2', 'pink', 2, 1),
(1, '未知', '0', 'gray', 3, 1),
(2, '启用', '1', 'green', 1, 1),
(2, '禁用', '0', 'red', 2, 1),
(3, '是', '1', 'green', 1, 1),
(3, '否', '0', 'red', 2, 1);
