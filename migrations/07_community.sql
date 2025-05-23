-- 社区模块数据表结构

-- 同X申请表
CREATE TABLE IF NOT EXISTS `similar_applications` (
  `id` VARCHAR(36) NOT NULL COMMENT '申请ID',
  `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
  `application_type` VARCHAR(20) NOT NULL COMMENT '申请类型(同用/同好/同购等)',
  `name` VARCHAR(50) NOT NULL COMMENT '姓名',
  `gender` VARCHAR(10) DEFAULT NULL COMMENT '性别',
  `birthplace` VARCHAR(100) DEFAULT NULL COMMENT '出生地点',
  `occupation` VARCHAR(50) DEFAULT NULL COMMENT '职业',
  `education` VARCHAR(50) DEFAULT NULL COMMENT '学历',
  `work_position` VARCHAR(50) DEFAULT NULL COMMENT '工作岗位',
  `work_place` VARCHAR(100) DEFAULT NULL COMMENT '工作地点',
  `hobby` TEXT DEFAULT NULL COMMENT '爱好',
  `address` VARCHAR(200) DEFAULT NULL COMMENT '地址',
  `contact_type` VARCHAR(20) DEFAULT NULL COMMENT '联系方式类型',
  `contact_value` VARCHAR(100) DEFAULT NULL COMMENT '联系方式值',
  `description` TEXT DEFAULT NULL COMMENT '个人简介',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态(pending/approved/rejected)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_application_type` (`application_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同X申请表';

-- 社区群组表
CREATE TABLE IF NOT EXISTS `community_groups` (
  `id` VARCHAR(36) NOT NULL COMMENT '群组ID',
  `name` VARCHAR(100) NOT NULL COMMENT '群组名称',
  `category` VARCHAR(20) NOT NULL COMMENT '群组分类(同X分类)',
  `description` TEXT DEFAULT NULL COMMENT '群组描述',
  `avatar` VARCHAR(255) DEFAULT NULL COMMENT '群组头像',
  `cover_image` VARCHAR(255) DEFAULT NULL COMMENT '群组封面图',
  `creator_id` VARCHAR(36) NOT NULL COMMENT '创建者ID',
  `member_count` INT NOT NULL DEFAULT 0 COMMENT '成员数量',
  `post_count` INT NOT NULL DEFAULT 0 COMMENT '帖子数量',
  `is_public` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否公开',
  `join_mode` VARCHAR(20) NOT NULL DEFAULT 'free' COMMENT '加入方式：free(自由加入), apply(申请加入), invite(仅邀请)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active, locked, deleted',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_category` (`category`),
  KEY `idx_creator_id` (`creator_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='社区群组表';

-- 群组成员表
CREATE TABLE IF NOT EXISTS `group_members` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `group_id` VARCHAR(36) NOT NULL COMMENT '群组ID',
  `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
  `role` VARCHAR(20) NOT NULL DEFAULT 'member' COMMENT '角色：admin(管理员), moderator(版主), member(普通成员)',
  `join_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active, muted, banned',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_group_user` (`group_id`, `user_id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role` (`role`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群组成员表';

-- 群组帖子表
CREATE TABLE IF NOT EXISTS `group_posts` (
  `id` VARCHAR(36) NOT NULL COMMENT '帖子ID',
  `group_id` VARCHAR(36) NOT NULL COMMENT '群组ID',
  `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
  `title` VARCHAR(255) DEFAULT NULL COMMENT '标题',
  `content` TEXT NOT NULL COMMENT '内容',
  `media_urls` TEXT DEFAULT NULL COMMENT '媒体URL列表，JSON格式',
  `view_count` INT NOT NULL DEFAULT 0 COMMENT '查看次数',
  `like_count` INT NOT NULL DEFAULT 0 COMMENT '点赞次数',
  `comment_count` INT NOT NULL DEFAULT 0 COMMENT '评论次数',
  `is_pinned` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否置顶',
  `is_essence` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否精华',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active, hidden, deleted',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群组帖子表';

-- 群组帖子评论表
CREATE TABLE IF NOT EXISTS `group_post_comments` (
  `id` VARCHAR(36) NOT NULL COMMENT '评论ID',
  `post_id` VARCHAR(36) NOT NULL COMMENT '帖子ID',
  `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
  `content` TEXT NOT NULL COMMENT '内容',
  `parent_id` VARCHAR(36) DEFAULT NULL COMMENT '父评论ID',
  `like_count` INT NOT NULL DEFAULT 0 COMMENT '点赞次数',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active, hidden, deleted',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群组帖子评论表';

-- 群组活动表
CREATE TABLE IF NOT EXISTS `group_activities` (
  `id` VARCHAR(36) NOT NULL COMMENT '活动ID',
  `group_id` VARCHAR(36) NOT NULL COMMENT '群组ID',
  `creator_id` VARCHAR(36) NOT NULL COMMENT '创建者ID',
  `title` VARCHAR(255) NOT NULL COMMENT '活动标题',
  `description` TEXT DEFAULT NULL COMMENT '活动描述',
  `cover_image` VARCHAR(255) DEFAULT NULL COMMENT '封面图',
  `location` VARCHAR(255) DEFAULT NULL COMMENT '活动地点',
  `start_time` TIMESTAMP NOT NULL COMMENT '开始时间',
  `end_time` TIMESTAMP DEFAULT NULL COMMENT '结束时间',
  `max_participants` INT DEFAULT NULL COMMENT '最大参与人数',
  `current_participants` INT NOT NULL DEFAULT 0 COMMENT '当前参与人数',
  `status` VARCHAR(20) NOT NULL DEFAULT 'planned' COMMENT '状态：planned, ongoing, completed, canceled',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_creator_id` (`creator_id`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群组活动表';

-- 活动参与者表
CREATE TABLE IF NOT EXISTS `activity_participants` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `activity_id` VARCHAR(36) NOT NULL COMMENT '活动ID',
  `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
  `status` VARCHAR(20) NOT NULL DEFAULT 'joined' COMMENT '状态：joined, canceled',
  `join_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
  `cancel_time` TIMESTAMP NULL DEFAULT NULL COMMENT '取消时间',
  `attendance` VARCHAR(20) DEFAULT NULL COMMENT '出席情况：attended, absent',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_activity_user` (`activity_id`, `user_id`),
  KEY `idx_activity_id` (`activity_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='活动参与者表';

-- 社区通知表
CREATE TABLE IF NOT EXISTS `community_notifications` (
  `id` VARCHAR(36) NOT NULL COMMENT '通知ID',
  `user_id` VARCHAR(36) NOT NULL COMMENT '接收用户ID',
  `sender_id` VARCHAR(36) DEFAULT NULL COMMENT '发送者ID',
  `type` VARCHAR(50) NOT NULL COMMENT '通知类型',
  `content` TEXT NOT NULL COMMENT '通知内容',
  `resource_type` VARCHAR(50) DEFAULT NULL COMMENT '资源类型',
  `resource_id` VARCHAR(36) DEFAULT NULL COMMENT '资源ID',
  `is_read` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已读',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `read_at` TIMESTAMP NULL DEFAULT NULL COMMENT '阅读时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_resource` (`resource_type`, `resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='社区通知表';

-- 同X分类表
CREATE TABLE IF NOT EXISTS `similar_categories` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `type` VARCHAR(20) NOT NULL COMMENT '分类类型(同用/同好/同购等)',
  `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
  `description` TEXT DEFAULT NULL COMMENT '分类描述',
  `icon` VARCHAR(255) DEFAULT NULL COMMENT '图标',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active, inactive',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_type` (`type`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同X分类表';

-- 初始化同X分类数据
INSERT INTO `similar_categories` (`type`, `name`, `description`, `sort_order`, `status`) VALUES
('同用', '同用', '相同用户', 1, 'active'),
('同好', '同好', '相同爱好', 2, 'active'),
('同购', '同购', '相同购物习惯', 3, 'active'),
('同年', '同年', '相同年份', 4, 'active'),
('同游', '同游', '相同游戏', 5, 'active'),
('同在', '同在', '相同位置', 6, 'active'),
('同市', '同市', '相同城市', 7, 'active'),
('同企', '同企', '相同企业', 8, 'active'),
('同亲', '同亲', '相同亲属关系', 9, 'active'),
('同班', '同班', '相同班级', 10, 'active'),
('同师', '同师', '相同老师', 11, 'active'),
('同业', '同业', '相同行业', 12, 'active'),
('同网', '同网', '相同网络', 13, 'active'),
('同工', '同工', '相同工作', 14, 'active'),
('同务', '同务', '相同事务', 15, 'active'),
('同艺', '同艺', '相同艺术爱好', 16, 'active'),
('同玩', '同玩', '相同玩乐方式', 17, 'active'),
('同闲', '同闲', '相同闲暇活动', 18, 'active'),
('同拍', '同拍', '相同摄影爱好', 19, 'active'),
('同乡', '同乡', '相同家乡', 20, 'active'),
('同学', '同学', '相同学校', 21, 'active');
