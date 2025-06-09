-- 插入初始角色
INSERT INTO roles (id, name, description, parent_id) VALUES
('role_1', '超级管理员', '系统超级管理员，拥有所有权限', NULL),
('role_2', '内容管理员', '负责内容管理，包括文章、评论等', NULL),
('role_3', '用户管理员', '负责用户管理，包括用户审核、封禁等', NULL),
('role_4', '运营管理员', '负责网站运营，包括活动、推广等', NULL)
ON DUPLICATE KEY UPDATE 
    description = VALUES(description);

-- 插入角色权限
INSERT INTO role_permissions (role_id, permission) VALUES
-- 超级管理员权限
('role_1', 'admin:*'),
('role_1', 'site:*'),
('role_1', 'user:*'),
('role_1', 'content:*'),
('role_1', 'system:*'),

-- 内容管理员权限
('role_2', 'content:view'),
('role_2', 'content:create'),
('role_2', 'content:edit'),
('role_2', 'content:delete'),
('role_2', 'content:approve'),
('role_2', 'content:reject'),

-- 用户管理员权限
('role_3', 'user:view'),
('role_3', 'user:create'),
('role_3', 'user:edit'),
('role_3', 'user:delete'),
('role_3', 'user:ban'),
('role_3', 'user:unban'),

-- 运营管理员权限
('role_4', 'content:view'),
('role_4', 'user:view'),
('role_4', 'operation:*')
ON DUPLICATE KEY UPDATE 
    role_id = VALUES(role_id);

-- 插入默认超级管理员账号 (密码为 'admin123', 使用bcrypt加密)
INSERT INTO admins (id, username, password_hash, role_id, status)
VALUES (1, 'admin', '$2a$10$Oz6YME9/oeOQGp9EFu9e5OQA1enVVOZKB5xoerGIl0/bEIxB2sX0e', 'role_1', 1)
ON DUPLICATE KEY UPDATE 
    role_id = VALUES(role_id),
    status = VALUES(status);

-- 针对万知网站的特定"入同"分类权限
INSERT INTO role_permissions (role_id, permission) VALUES
-- 超级管理员拥有所有分类权限
('role_1', 'category:同用'),
('role_1', 'category:同好'),
('role_1', 'category:同购'),
('role_1', 'category:同年'),
('role_1', 'category:同游'),
('role_1', 'category:同在'),
('role_1', 'category:同市'),
('role_1', 'category:同企'),
('role_1', 'category:同亲'),
('role_1', 'category:同班'),
('role_1', 'category:同师'),
('role_1', 'category:同业'),
('role_1', 'category:同网'),
('role_1', 'category:同工'),
('role_1', 'category:同务'),
('role_1', 'category:同艺'),
('role_1', 'category:同玩'),
('role_1', 'category:同闲'),
('role_1', 'category:同拍'),
('role_1', 'category:同乡'),
('role_1', 'category:同学'),

-- 内容管理员只拥有部分分类权限
('role_2', 'category:同用'),
('role_2', 'category:同好'),
('role_2', 'category:同购'),
('role_2', 'category:同年'),
('role_2', 'category:同游')
ON DUPLICATE KEY UPDATE 
    role_id = VALUES(role_id);
