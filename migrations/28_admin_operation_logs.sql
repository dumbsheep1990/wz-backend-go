-- 管理员操作日志表结构

-- 管理员操作日志表
CREATE TABLE IF NOT EXISTS `admin_operation_logs` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `admin_id` BIGINT NOT NULL COMMENT '管理员ID',
  `admin_name` VARCHAR(64) NOT NULL COMMENT '管理员名称',
  `module` VARCHAR(64) NOT NULL COMMENT '操作模块',
  `action` VARCHAR(32) NOT NULL COMMENT '操作类型：create-创建 update-更新 delete-删除 query-查询 login-登录 logout-登出',
  `resource_type` VARCHAR(64) NOT NULL COMMENT '资源类型',
  `resource_id` VARCHAR(64) DEFAULT NULL COMMENT '资源ID',
  `ip` VARCHAR(64) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` VARCHAR(512) DEFAULT NULL COMMENT '用户代理',
  `request_url` VARCHAR(512) DEFAULT NULL COMMENT '请求URL',
  `request_method` VARCHAR(16) DEFAULT NULL COMMENT '请求方法',
  `request_data` TEXT DEFAULT NULL COMMENT '请求数据',
  `response_code` INT DEFAULT NULL COMMENT '响应状态码',
  `response_data` TEXT DEFAULT NULL COMMENT '响应数据',
  `execution_time` INT DEFAULT NULL COMMENT '执行时间(毫秒)',
  `status` INT DEFAULT 1 COMMENT '状态：0-失败 1-成功',
  `remark` VARCHAR(512) DEFAULT NULL COMMENT '备注',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `tenant_id` BIGINT DEFAULT NULL COMMENT '租户ID',
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_module` (`module`),
  KEY `idx_action` (`action`),
  KEY `idx_resource_type` (`resource_type`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员操作日志表';
