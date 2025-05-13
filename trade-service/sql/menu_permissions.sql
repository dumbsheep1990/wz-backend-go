-- 添加交易管理菜单
INSERT INTO `system_menu` (`id`, `parent_id`, `name`, `path`, `component`, `redirect`, `perms`, `type`, `icon`, `order_num`, `hidden`, `affix`, `cache_on`, `created_at`, `updated_at`, `deleted_at`) VALUES
(100, 0, '交易管理', '/trade', 'Layout', NULL, '', 1, 'el-icon-shopping-cart-full', 5, 0, 0, 0, now(), now(), NULL);

-- 添加订单管理菜单
INSERT INTO `system_menu` (`id`, `parent_id`, `name`, `path`, `component`, `redirect`, `perms`, `type`, `icon`, `order_num`, `hidden`, `affix`, `cache_on`, `created_at`, `updated_at`, `deleted_at`) VALUES
(101, 100, '订单管理', 'orders', 'trade/orders/index', NULL, 'orders:view', 2, 'el-icon-tickets', 1, 0, 0, 1, now(), now(), NULL);

-- 添加订单管理按钮权限
INSERT INTO `system_menu` (`id`, `parent_id`, `name`, `path`, `component`, `redirect`, `perms`, `type`, `icon`, `order_num`, `hidden`, `affix`, `cache_on`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1011, 101, '订单查看', '', '', NULL, 'orders:view', 3, '', 1, 0, 0, 0, now(), now(), NULL),
(1012, 101, '订单编辑', '', '', NULL, 'orders:edit', 3, '', 2, 0, 0, 0, now(), now(), NULL),
(1013, 101, '订单删除', '', '', NULL, 'orders:delete', 3, '', 3, 0, 0, 0, now(), now(), NULL),
(1014, 101, '订单导出', '', '', NULL, 'orders:export', 3, '', 4, 0, 0, 0, now(), now(), NULL),
(1015, 101, '订单统计', '', '', NULL, 'orders:stats', 3, '', 5, 0, 0, 0, now(), now(), NULL);

-- 添加支付配置菜单
INSERT INTO `system_menu` (`id`, `parent_id`, `name`, `path`, `component`, `redirect`, `perms`, `type`, `icon`, `order_num`, `hidden`, `affix`, `cache_on`, `created_at`, `updated_at`, `deleted_at`) VALUES
(102, 100, '支付配置', 'payment-config', 'trade/payment/config', NULL, 'payment:config', 2, 'el-icon-money', 2, 0, 0, 1, now(), now(), NULL);

-- 添加支付配置按钮权限
INSERT INTO `system_menu` (`id`, `parent_id`, `name`, `path`, `component`, `redirect`, `perms`, `type`, `icon`, `order_num`, `hidden`, `affix`, `cache_on`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1021, 102, '支付配置查看', '', '', NULL, 'payment:view', 3, '', 1, 0, 0, 0, now(), now(), NULL),
(1022, 102, '支付配置编辑', '', '', NULL, 'payment:edit', 3, '', 2, 0, 0, 0, now(), now(), NULL); 