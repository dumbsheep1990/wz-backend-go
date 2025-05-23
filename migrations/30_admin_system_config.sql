-- 管理系统配置和控制面板表结构

-- 系统配置表
CREATE TABLE IF NOT EXISTS `admin_system_configs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `config_group` VARCHAR(64) NOT NULL COMMENT '配置分组',
  `config_key` VARCHAR(128) NOT NULL COMMENT '配置键',
  `config_value` TEXT NOT NULL COMMENT '配置值',
  `value_type` VARCHAR(32) DEFAULT 'string' COMMENT '值类型：string, number, boolean, json, array',
  `name` VARCHAR(64) NOT NULL COMMENT '配置名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '配置描述',
  `sort_order` INT DEFAULT 0 COMMENT '排序顺序',
  `is_system` TINYINT(1) DEFAULT 0 COMMENT '是否系统配置：0-否 1-是',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `updated_by` BIGINT DEFAULT NULL COMMENT '更新人',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID，NULL表示系统级配置',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_group_key_tenant` (`config_group`, `config_key`, `tenant_id`),
  KEY `idx_config_group` (`config_group`),
  KEY `idx_status` (`status`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 控制面板布局表
CREATE TABLE IF NOT EXISTS `admin_dashboard_layouts` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `admin_id` BIGINT NOT NULL COMMENT '管理员ID',
  `layout_name` VARCHAR(64) NOT NULL COMMENT '布局名称',
  `is_default` TINYINT(1) DEFAULT 0 COMMENT '是否默认布局：0-否 1-是',
  `layout_data` JSON NOT NULL COMMENT '布局数据(JSON格式)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_is_default` (`is_default`),
  KEY `idx_tenant_id` (`tenant_id`),
  CONSTRAINT `fk_admin_dashboard_layouts_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `admin_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='控制面板布局表';

-- 控制面板组件表
CREATE TABLE IF NOT EXISTS `admin_dashboard_widgets` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `widget_code` VARCHAR(64) NOT NULL COMMENT '组件编码',
  `widget_name` VARCHAR(64) NOT NULL COMMENT '组件名称',
  `widget_type` VARCHAR(32) NOT NULL COMMENT '组件类型：chart-图表 list-列表 statistic-统计 custom-自定义',
  `component` VARCHAR(128) NOT NULL COMMENT '前端组件',
  `icon` VARCHAR(64) DEFAULT NULL COMMENT '图标',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `config_schema` JSON DEFAULT NULL COMMENT '配置模式(JSON格式)',
  `default_config` JSON DEFAULT NULL COMMENT '默认配置(JSON格式)',
  `data_api` VARCHAR(255) DEFAULT NULL COMMENT '数据API',
  `permissions` VARCHAR(255) DEFAULT NULL COMMENT '所需权限，多个用逗号分隔',
  `is_system` TINYINT(1) DEFAULT 0 COMMENT '是否系统组件：0-否 1-是',
  `status` INT DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_widget_code` (`widget_code`),
  KEY `idx_widget_type` (`widget_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='控制面板组件表';

-- 用户组件实例表
CREATE TABLE IF NOT EXISTS `admin_widget_instances` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `layout_id` BIGINT NOT NULL COMMENT '布局ID',
  `widget_id` BIGINT NOT NULL COMMENT '组件ID',
  `instance_id` VARCHAR(64) NOT NULL COMMENT '实例ID',
  `title` VARCHAR(64) DEFAULT NULL COMMENT '实例标题',
  `position_x` INT NOT NULL COMMENT 'X位置',
  `position_y` INT NOT NULL COMMENT 'Y位置',
  `width` INT NOT NULL COMMENT '宽度',
  `height` INT NOT NULL COMMENT '高度',
  `config` JSON DEFAULT NULL COMMENT '配置(JSON格式)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_layout_instance` (`layout_id`, `instance_id`),
  KEY `idx_layout_id` (`layout_id`),
  KEY `idx_widget_id` (`widget_id`),
  CONSTRAINT `fk_admin_widget_instances_layout_id` FOREIGN KEY (`layout_id`) REFERENCES `admin_dashboard_layouts` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_admin_widget_instances_widget_id` FOREIGN KEY (`widget_id`) REFERENCES `admin_dashboard_widgets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户组件实例表';

-- 初始化系统配置
INSERT INTO `admin_system_configs` (`config_group`, `config_key`, `config_value`, `value_type`, `name`, `description`, `is_system`, `status`)
VALUES 
('system', 'site_name', '网站管理系统', 'string', '站点名称', '网站名称，显示在浏览器标题栏', 1, 1),
('system', 'site_logo', '/assets/logo.png', 'string', '站点Logo', '网站Logo图片路径', 1, 1),
('system', 'site_favicon', '/favicon.ico', 'string', '站点图标', '浏览器标签页显示的小图标', 1, 1),
('system', 'copyright', '© 2025 管理系统 版权所有', 'string', '版权信息', '网站底部显示的版权信息', 1, 1),
('system', 'allow_register', 'true', 'boolean', '允许注册', '是否允许新用户注册', 1, 1),
('system', 'default_theme', 'light', 'string', '默认主题', '系统默认主题：light-浅色 dark-深色', 1, 1),
('system', 'system_notice', '欢迎使用管理系统', 'string', '系统公告', '登录页和首页显示的系统公告', 1, 1),
('email', 'smtp_host', 'smtp.example.com', 'string', 'SMTP服务器', '邮件服务器地址', 1, 1),
('email', 'smtp_port', '465', 'number', 'SMTP端口', '邮件服务器端口', 1, 1),
('email', 'smtp_user', 'admin@example.com', 'string', 'SMTP用户名', '邮件服务器用户名', 1, 1),
('email', 'smtp_password', '******', 'string', 'SMTP密码', '邮件服务器密码', 1, 1),
('email', 'mail_from', 'admin@example.com', 'string', '发件人地址', '系统发送邮件的发件人地址', 1, 1),
('email', 'mail_from_name', '系统管理员', 'string', '发件人名称', '系统发送邮件的发件人名称', 1, 1),
('upload', 'upload_max_size', '10', 'number', '上传大小限制', '上传文件大小限制(MB)', 1, 1),
('upload', 'allowed_extensions', 'jpg,jpeg,png,gif,doc,docx,xls,xlsx,pdf,zip,rar', 'string', '允许的扩展名', '允许上传的文件扩展名，用逗号分隔', 1, 1),
('upload', 'storage_type', 'local', 'string', '存储类型', '文件存储类型：local-本地 oss-对象存储', 1, 1),
('security', 'password_min_length', '6', 'number', '密码最小长度', '用户密码最小长度限制', 1, 1),
('security', 'password_strength', 'medium', 'string', '密码强度要求', '密码强度要求：low-低 medium-中 high-高', 1, 1),
('security', 'login_fail_limit', '5', 'number', '登录失败限制', '连续登录失败次数限制', 1, 1),
('security', 'lock_time', '30', 'number', '锁定时间', '账户锁定时间(分钟)', 1, 1);

-- 初始化仪表盘组件
INSERT INTO `admin_dashboard_widgets` (`widget_code`, `widget_name`, `widget_type`, `component`, `icon`, `description`, `is_system`, `status`)
VALUES 
('system_info', '系统信息', 'statistic', 'SystemInfoWidget', 'el-icon-monitor', '显示系统基本信息', 1, 1),
('quick_nav', '快捷导航', 'custom', 'QuickNavWidget', 'el-icon-link', '常用功能快捷导航', 1, 1),
('user_statistic', '用户统计', 'chart', 'UserStatisticWidget', 'el-icon-user', '用户相关统计数据', 1, 1),
('content_statistic', '内容统计', 'chart', 'ContentStatisticWidget', 'el-icon-document', '内容相关统计数据', 1, 1),
('recent_orders', '最近订单', 'list', 'RecentOrdersWidget', 'el-icon-shopping-cart-full', '显示最近订单列表', 1, 1),
('task_list', '待办任务', 'list', 'TaskListWidget', 'el-icon-s-claim', '显示待办任务列表', 1, 1),
('server_monitor', '服务器监控', 'chart', 'ServerMonitorWidget', 'el-icon-cpu', '服务器资源监控', 1, 1),
('notification_center', '通知中心', 'list', 'NotificationCenterWidget', 'el-icon-bell', '系统通知和消息', 1, 1);
