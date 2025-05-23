-- 组件模块数据表结构

-- 组件分类表
CREATE TABLE IF NOT EXISTS `component_categories` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '分类名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '分类描述',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件分类表';

-- 组件表
CREATE TABLE IF NOT EXISTS `components_library` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL COMMENT '组件名称',
  `category_id` BIGINT NOT NULL COMMENT '分类ID',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '组件描述',
  `content` TEXT DEFAULT NULL COMMENT '组件内容模板',
  `status` VARCHAR(32) NOT NULL DEFAULT 'active' COMMENT '状态：active, deprecated',
  `version` VARCHAR(32) DEFAULT '1.0.0' COMMENT '版本号',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `ext` TEXT DEFAULT NULL COMMENT '扩展字段，JSON格式',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_name` (`name`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_components_category_id` FOREIGN KEY (`category_id`) REFERENCES `component_categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件库表';

-- 组件配置项表
CREATE TABLE IF NOT EXISTS `component_configs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `component_id` BIGINT NOT NULL COMMENT '组件ID',
  `key` VARCHAR(64) NOT NULL COMMENT '配置键',
  `type` VARCHAR(32) NOT NULL COMMENT '配置类型：string, number, boolean, object, array',
  `default_value` TEXT DEFAULT NULL COMMENT '默认值',
  `label` VARCHAR(64) NOT NULL COMMENT '标签名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `required` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否必填',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_component_id` (`component_id`),
  KEY `idx_sort_order` (`sort_order`),
  CONSTRAINT `fk_component_configs_component_id` FOREIGN KEY (`component_id`) REFERENCES `components_library` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件配置项表';

-- 组件样式表
CREATE TABLE IF NOT EXISTS `component_styles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `component_id` BIGINT NOT NULL COMMENT '组件ID',
  `key` VARCHAR(64) NOT NULL COMMENT '样式键',
  `type` VARCHAR(32) NOT NULL COMMENT '样式类型：color, size, spacing, typography, border',
  `default_value` VARCHAR(255) DEFAULT NULL COMMENT '默认值',
  `label` VARCHAR(64) NOT NULL COMMENT '标签名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_component_id` (`component_id`),
  KEY `idx_sort_order` (`sort_order`),
  CONSTRAINT `fk_component_styles_component_id` FOREIGN KEY (`component_id`) REFERENCES `components_library` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件样式表';

-- 组件使用记录表
CREATE TABLE IF NOT EXISTS `component_usages` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `component_id` BIGINT NOT NULL COMMENT '组件ID',
  `site_id` VARCHAR(64) NOT NULL COMMENT '站点ID',
  `page_id` VARCHAR(64) NOT NULL COMMENT '页面ID',
  `section_id` VARCHAR(64) NOT NULL COMMENT '区块ID',
  `instance_id` VARCHAR(64) NOT NULL COMMENT '实例ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_component_id` (`component_id`),
  KEY `idx_site_id` (`site_id`),
  KEY `idx_page_id` (`page_id`),
  KEY `idx_section_id` (`section_id`),
  KEY `idx_instance_id` (`instance_id`),
  CONSTRAINT `fk_component_usages_component_id` FOREIGN KEY (`component_id`) REFERENCES `components_library` (`id`),
  CONSTRAINT `fk_component_usages_site_id` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件使用记录表';

-- 组件版本历史表
CREATE TABLE IF NOT EXISTS `component_versions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `component_id` BIGINT NOT NULL COMMENT '组件ID',
  `version` VARCHAR(32) NOT NULL COMMENT '版本号',
  `content` TEXT NOT NULL COMMENT '组件内容',
  `changelog` TEXT DEFAULT NULL COMMENT '变更日志',
  `created_by` BIGINT DEFAULT NULL COMMENT '创建人',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_component_version` (`component_id`, `version`),
  KEY `idx_component_id` (`component_id`),
  CONSTRAINT `fk_component_versions_component_id` FOREIGN KEY (`component_id`) REFERENCES `components_library` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件版本历史表';

-- 初始化基础组件分类数据
INSERT INTO `component_categories` (`name`, `description`) VALUES 
('基础组件', '基础UI组件如文本、按钮、分隔线等'),
('布局组件', '布局相关组件如容器、行、列等'),
('媒体组件', '媒体相关组件如图片、视频、轮播图等'),
('表单组件', '表单相关组件如输入框、选择框、单选框等'),
('数据组件', '数据展示组件如表格、列表、图表等'),
('导航组件', '导航相关组件如导航栏、菜单、面包屑等'),
('高级组件', '复杂功能组件如地图、日历、评论区等');
