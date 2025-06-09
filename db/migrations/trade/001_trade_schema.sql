-- 交易服务数据库表结构
USE wz_backend;

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id VARCHAR(64) NOT NULL COMMENT '订单ID，业务唯一标识',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    product_id BIGINT NOT NULL COMMENT '产品ID',
    product_type VARCHAR(50) NOT NULL COMMENT '产品类型',
    quantity INT NOT NULL DEFAULT 1 COMMENT '数量',
    amount DECIMAL(12, 2) NOT NULL COMMENT '金额',
    currency VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
    status VARCHAR(20) NOT NULL COMMENT '订单状态：pending(待支付), paid(已支付), canceled(已取消), refunded(已退款), expired(已过期)',
    payment_id VARCHAR(64) COMMENT '支付ID',
    payment_type VARCHAR(20) COMMENT '支付类型：alipay, wechat, bank_transfer',
    payment_time TIMESTAMP NULL COMMENT '支付时间',
    description VARCHAR(500) COMMENT '描述',
    metadata TEXT COMMENT '元数据，JSON格式',
    client_ip VARCHAR(64) COMMENT '客户端IP',
    device_id VARCHAR(100) COMMENT '设备ID',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    expire_time TIMESTAMP NULL COMMENT '过期时间',
    UNIQUE KEY uk_order_id (order_id),
    INDEX idx_orders_user_id (user_id),
    INDEX idx_orders_status (status),
    INDEX idx_orders_payment_id (payment_id),
    INDEX idx_orders_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 订单明细表
CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id VARCHAR(64) NOT NULL COMMENT '订单ID',
    product_id BIGINT NOT NULL COMMENT '产品ID',
    product_type VARCHAR(50) NOT NULL COMMENT '产品类型',
    product_name VARCHAR(200) NOT NULL COMMENT '产品名称',
    quantity INT NOT NULL DEFAULT 1 COMMENT '数量',
    unit_price DECIMAL(12, 2) NOT NULL COMMENT '单价',
    total_price DECIMAL(12, 2) NOT NULL COMMENT '总价',
    discount DECIMAL(12, 2) DEFAULT 0 COMMENT '折扣',
    metadata TEXT COMMENT '元数据，JSON格式',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_order_items_order_id (order_id),
    INDEX idx_order_items_product_id (product_id),
    CONSTRAINT fk_order_items_order_id FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 支付记录表
CREATE TABLE IF NOT EXISTS payments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    payment_id VARCHAR(64) NOT NULL COMMENT '支付ID，业务唯一标识',
    order_id VARCHAR(64) NOT NULL COMMENT '订单ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    amount DECIMAL(12, 2) NOT NULL COMMENT '支付金额',
    currency VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
    payment_type VARCHAR(20) NOT NULL COMMENT '支付类型：alipay, wechat, bank_transfer',
    status VARCHAR(20) NOT NULL COMMENT '支付状态：pending(处理中), success(成功), failed(失败)',
    transaction_id VARCHAR(100) COMMENT '第三方交易ID',
    payment_time TIMESTAMP NULL COMMENT '支付时间',
    callback_time TIMESTAMP NULL COMMENT '回调时间',
    callback_data TEXT COMMENT '回调原始数据',
    client_ip VARCHAR(64) COMMENT '客户端IP',
    metadata TEXT COMMENT '元数据，JSON格式',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_payment_id (payment_id),
    INDEX idx_payments_order_id (order_id),
    INDEX idx_payments_user_id (user_id),
    INDEX idx_payments_status (status),
    INDEX idx_payments_transaction_id (transaction_id),
    INDEX idx_payments_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 退款表
CREATE TABLE IF NOT EXISTS refunds (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    refund_id VARCHAR(64) NOT NULL COMMENT '退款ID，业务唯一标识',
    order_id VARCHAR(64) NOT NULL COMMENT '订单ID',
    payment_id VARCHAR(64) COMMENT '支付ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    amount DECIMAL(12, 2) NOT NULL COMMENT '退款金额',
    currency VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
    status VARCHAR(20) NOT NULL COMMENT '退款状态：pending(待处理), approved(已批准), processing(处理中), success(成功), rejected(已拒绝), failed(失败)',
    reason VARCHAR(200) NOT NULL COMMENT '退款原因',
    description VARCHAR(500) COMMENT '描述',
    processed_by VARCHAR(100) COMMENT '处理人',
    process_time TIMESTAMP NULL COMMENT '处理时间',
    refund_transaction_id VARCHAR(100) COMMENT '退款交易ID',
    metadata TEXT COMMENT '元数据，JSON格式',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_refund_id (refund_id),
    INDEX idx_refunds_order_id (order_id),
    INDEX idx_refunds_payment_id (payment_id),
    INDEX idx_refunds_user_id (user_id),
    INDEX idx_refunds_status (status),
    INDEX idx_refunds_created_at (created_at),
    CONSTRAINT fk_refunds_order_id FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 账户余额表
CREATE TABLE IF NOT EXISTS account_balances (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    currency VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
    available DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '可用余额',
    pending DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '待结算余额',
    frozen DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '冻结余额',
    total DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '总余额',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_account_user_currency (user_id, currency),
    INDEX idx_account_balances_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 交易记录表
CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    transaction_id VARCHAR(64) NOT NULL COMMENT '交易ID，业务唯一标识',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    related_id VARCHAR(64) COMMENT '关联ID（订单ID或退款ID）',
    related_type VARCHAR(20) COMMENT '关联类型：order, refund',
    type VARCHAR(20) NOT NULL COMMENT '交易类型：payment(支付), refund(退款), withdraw(提现), adjustment(调整)',
    amount DECIMAL(12, 2) NOT NULL COMMENT '金额',
    currency VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
    balance_before DECIMAL(12, 2) NOT NULL COMMENT '交易前余额',
    balance_after DECIMAL(12, 2) NOT NULL COMMENT '交易后余额',
    status VARCHAR(20) NOT NULL COMMENT '交易状态：pending(处理中), success(成功), failed(失败)',
    description VARCHAR(500) COMMENT '描述',
    metadata TEXT COMMENT '元数据，JSON格式',
    operator_id VARCHAR(100) COMMENT '操作员ID',
    client_ip VARCHAR(64) COMMENT '客户端IP',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_transaction_id (transaction_id),
    INDEX idx_transactions_user_id (user_id),
    INDEX idx_transactions_related_id (related_id),
    INDEX idx_transactions_type (type),
    INDEX idx_transactions_status (status),
    INDEX idx_transactions_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 财务日报表
CREATE TABLE IF NOT EXISTS financial_daily_reports (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    report_date DATE NOT NULL COMMENT '报表日期',
    currency VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
    income DECIMAL(15, 2) NOT NULL DEFAULT 0 COMMENT '收入',
    refund DECIMAL(15, 2) NOT NULL DEFAULT 0 COMMENT '退款',
    net DECIMAL(15, 2) NOT NULL DEFAULT 0 COMMENT '净收入',
    order_count INT NOT NULL DEFAULT 0 COMMENT '订单数',
    payment_count INT NOT NULL DEFAULT 0 COMMENT '支付数',
    refund_count INT NOT NULL DEFAULT 0 COMMENT '退款数',
    user_count INT NOT NULL DEFAULT 0 COMMENT '用户数',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_daily_report_date_currency (report_date, currency),
    INDEX idx_daily_reports_date (report_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 财务月报表
CREATE TABLE IF NOT EXISTS financial_monthly_reports (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    report_year INT NOT NULL COMMENT '报表年份',
    report_month INT NOT NULL COMMENT '报表月份',
    currency VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
    income DECIMAL(15, 2) NOT NULL DEFAULT 0 COMMENT '收入',
    refund DECIMAL(15, 2) NOT NULL DEFAULT 0 COMMENT '退款',
    net DECIMAL(15, 2) NOT NULL DEFAULT 0 COMMENT '净收入',
    order_count INT NOT NULL DEFAULT 0 COMMENT '订单数',
    payment_count INT NOT NULL DEFAULT 0 COMMENT '支付数',
    refund_count INT NOT NULL DEFAULT 0 COMMENT '退款数',
    user_count INT NOT NULL DEFAULT 0 COMMENT '用户数',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_monthly_report_year_month_currency (report_year, report_month, currency),
    INDEX idx_monthly_reports_year_month (report_year, report_month)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 付款方式配置表
CREATE TABLE IF NOT EXISTS payment_methods (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    method_code VARCHAR(50) NOT NULL COMMENT '方式代码',
    method_name VARCHAR(100) NOT NULL COMMENT '方式名称',
    method_type VARCHAR(20) NOT NULL COMMENT '方式类型：alipay, wechat, bank_transfer',
    config TEXT NOT NULL COMMENT '配置，JSON格式',
    is_enabled TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_payment_method_code (method_code),
    INDEX idx_payment_methods_enabled (is_enabled),
    INDEX idx_payment_methods_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
