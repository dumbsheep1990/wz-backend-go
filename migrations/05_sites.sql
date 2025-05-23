-- 站点模块数据表结构

-- 站点表
CREATE TABLE IF NOT EXISTS `sites` (
  `id` VARCHAR(64) NOT NULL COMMENT '站点ID',
  `name` VARCHAR(128) NOT NULL COMMENT '站点名称',
  `description` TEXT DEFAULT NULL COMMENT '站点描述',
  `domain` VARCHAR(128) DEFAULT NULL COMMENT '站点域名',
  `logo` VARCHAR(255) DEFAULT NULL COMMENT '站点Logo',
  `favicon` VARCHAR(255) DEFAULT NULL COMMENT '站点图标',
  `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
  `primary_color` VARCHAR(32) DEFAULT NULL COMMENT '主色调',
  `secondary_color` VARCHAR(32) DEFAULT NULL COMMENT '辅助色',
  `accent_color` VARCHAR(32) DEFAULT NULL COMMENT '强调色',
  `text_color` VARCHAR(32) DEFAULT NULL COMMENT '文字颜色',
  `background_color` VARCHAR(32) DEFAULT NULL COMMENT '背景颜色',
  `font_family` VARCHAR(64) DEFAULT NULL COMMENT '字体',
  `header_style` VARCHAR(32) DEFAULT NULL COMMENT '页头样式：standard, centered, minimal',
  `border_radius` VARCHAR(32) DEFAULT NULL COMMENT '边框圆角：none, small, medium, large',
  `custom_css` TEXT DEFAULT NULL COMMENT '自定义CSS',
  `navigation` JSON DEFAULT NULL COMMENT '导航配置',
  `footer` JSON DEFAULT NULL COMMENT '页脚配置',
  `thumbnail` VARCHAR(255) DEFAULT NULL COMMENT '缩略图',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `published_at` TIMESTAMP NULL DEFAULT NULL COMMENT '发布时间',
  `status` VARCHAR(32) NOT NULL DEFAULT 'draft' COMMENT '状态：draft, published, archived',
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_domain` (`domain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='站点表';

-- 页面表
CREATE TABLE IF NOT EXISTS `pages` (
  `id` VARCHAR(64) NOT NULL COMMENT '页面ID',
  `site_id` VARCHAR(64) NOT NULL COMMENT '站点ID',
  `name` VARCHAR(128) NOT NULL COMMENT '页面名称',
  `slug` VARCHAR(128) NOT NULL COMMENT '页面Slug',
  `title` VARCHAR(255) NOT NULL COMMENT '页面标题',
  `description` TEXT DEFAULT NULL COMMENT '页面描述',
  `keywords` JSON DEFAULT NULL COMMENT '关键词',
  `is_homepage` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否是首页',
  `layout` VARCHAR(32) DEFAULT 'default' COMMENT '布局类型：default, full-width, sidebar',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  PRIMARY KEY (`id`),
  KEY `idx_site_id` (`site_id`),
  KEY `idx_slug` (`slug`),
  KEY `idx_is_homepage` (`is_homepage`),
  KEY `idx_sort_order` (`sort_order`),
  CONSTRAINT `fk_pages_site_id` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='页面表';

-- 区块表
CREATE TABLE IF NOT EXISTS `sections` (
  `id` VARCHAR(64) NOT NULL COMMENT '区块ID',
  `page_id` VARCHAR(64) NOT NULL COMMENT '页面ID',
  `type` VARCHAR(64) NOT NULL COMMENT '区块类型',
  `title` VARCHAR(255) DEFAULT NULL COMMENT '区块标题',
  `settings` JSON DEFAULT NULL COMMENT '区块设置',
  `style` JSON DEFAULT NULL COMMENT '区块样式',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  PRIMARY KEY (`id`),
  KEY `idx_page_id` (`page_id`),
  KEY `idx_sort_order` (`sort_order`),
  CONSTRAINT `fk_sections_page_id` FOREIGN KEY (`page_id`) REFERENCES `pages` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='区块表';

-- 组件表
CREATE TABLE IF NOT EXISTS `components` (
  `id` VARCHAR(64) NOT NULL COMMENT '组件ID',
  `section_id` VARCHAR(64) NOT NULL COMMENT '区块ID',
  `type` VARCHAR(64) NOT NULL COMMENT '组件类型',
  `name` VARCHAR(128) NOT NULL COMMENT '组件名称',
  `settings` JSON DEFAULT NULL COMMENT '组件设置',
  `content` JSON DEFAULT NULL COMMENT '组件内容',
  `style` JSON DEFAULT NULL COMMENT '组件样式',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  PRIMARY KEY (`id`),
  KEY `idx_section_id` (`section_id`),
  KEY `idx_sort_order` (`sort_order`),
  CONSTRAINT `fk_components_section_id` FOREIGN KEY (`section_id`) REFERENCES `sections` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件表';

-- 站点模板表
CREATE TABLE IF NOT EXISTS `site_templates` (
  `id` VARCHAR(64) NOT NULL COMMENT '模板ID',
  `name` VARCHAR(128) NOT NULL COMMENT '模板名称',
  `thumbnail` VARCHAR(255) DEFAULT NULL COMMENT '缩略图',
  `description` TEXT DEFAULT NULL COMMENT '模板描述',
  `config` JSON DEFAULT NULL COMMENT '模板配置',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='站点模板表';

-- 组件分类表
CREATE TABLE IF NOT EXISTS `component_categories` (
  `id` VARCHAR(64) NOT NULL COMMENT '分类ID',
  `name` VARCHAR(64) NOT NULL COMMENT '分类名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件分类表';

-- 组件定义表
CREATE TABLE IF NOT EXISTS `component_definitions` (
  `id` VARCHAR(64) NOT NULL COMMENT '定义ID',
  `category_id` VARCHAR(64) NOT NULL COMMENT '分类ID',
  `type` VARCHAR(64) NOT NULL COMMENT '组件类型',
  `name` VARCHAR(128) NOT NULL COMMENT '组件名称',
  `icon` VARCHAR(64) DEFAULT NULL COMMENT '图标',
  `description` TEXT DEFAULT NULL COMMENT '描述',
  `default_settings` JSON DEFAULT NULL COMMENT '默认设置',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  CONSTRAINT `fk_component_definitions_category_id` FOREIGN KEY (`category_id`) REFERENCES `component_categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件定义表';

-- 站点访问统计表
CREATE TABLE IF NOT EXISTS `site_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `site_id` VARCHAR(64) NOT NULL COMMENT '站点ID',
  `date` DATE NOT NULL COMMENT '日期',
  `page_views` INT NOT NULL DEFAULT 0 COMMENT '页面浏览量',
  `unique_visitors` INT NOT NULL DEFAULT 0 COMMENT '唯一访客数',
  `bounce_rate` DECIMAL(5,2) DEFAULT NULL COMMENT '跳出率',
  `avg_session_duration` INT DEFAULT NULL COMMENT '平均会话时长(秒)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_site_date` (`site_id`, `date`),
  KEY `idx_site_id` (`site_id`),
  CONSTRAINT `fk_site_statistics_site_id` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='站点访问统计表';

-- 页面访问统计表
CREATE TABLE IF NOT EXISTS `page_statistics` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `page_id` VARCHAR(64) NOT NULL COMMENT '页面ID',
  `site_id` VARCHAR(64) NOT NULL COMMENT '站点ID',
  `date` DATE NOT NULL COMMENT '日期',
  `page_views` INT NOT NULL DEFAULT 0 COMMENT '页面浏览量',
  `unique_views` INT NOT NULL DEFAULT 0 COMMENT '唯一浏览量',
  `avg_time_on_page` INT DEFAULT NULL COMMENT '平均页面停留时间(秒)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_page_date` (`page_id`, `date`),
  KEY `idx_page_id` (`page_id`),
  KEY `idx_site_id` (`site_id`),
  CONSTRAINT `fk_page_statistics_page_id` FOREIGN KEY (`page_id`) REFERENCES `pages` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_page_statistics_site_id` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='页面访问统计表';
