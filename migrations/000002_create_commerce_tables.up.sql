-- Create commerce_products table
CREATE TABLE IF NOT EXISTS commerce_products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    original_price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    featured BOOLEAN NOT NULL DEFAULT FALSE,
    store_id VARCHAR(36) NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    thumbnail VARCHAR(255) NOT NULL,
    images JSON NOT NULL,
    specifications JSON NOT NULL,
    tags JSON NOT NULL,
    view_count INT NOT NULL DEFAULT 0,
    sales_count INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_product_store (store_id),
    INDEX idx_product_category (category_id),
    INDEX idx_product_featured (featured),
    INDEX idx_product_active (is_active),
    INDEX idx_product_created (created_at),
    INDEX idx_product_views (view_count),
    INDEX idx_product_sales (sales_count)
);

-- Create commerce_stores table
CREATE TABLE IF NOT EXISTS commerce_stores (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    owner_id VARCHAR(36) NOT NULL,
    logo_url VARCHAR(255) NOT NULL,
    province VARCHAR(50) NOT NULL,
    city VARCHAR(50) NOT NULL,
    district VARCHAR(50) NOT NULL,
    address VARCHAR(255) NOT NULL,
    contact_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    rating DECIMAL(2, 1) NOT NULL DEFAULT 5.0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_store_owner (owner_id),
    INDEX idx_store_location (province, city, district),
    INDEX idx_store_rating (rating),
    INDEX idx_store_active (is_active)
);

-- Create commerce_categories table
CREATE TABLE IF NOT EXISTS commerce_categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    parent_id VARCHAR(36),
    icon_url VARCHAR(255) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    level INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_category_parent (parent_id),
    INDEX idx_category_level (level),
    INDEX idx_category_sort (sort_order),
    INDEX idx_category_active (is_active)
);
