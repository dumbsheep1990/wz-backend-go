-- Drop all commerce tables in reverse order to avoid foreign key constraints
DROP TABLE IF EXISTS commerce_products;
DROP TABLE IF EXISTS commerce_stores;
DROP TABLE IF EXISTS commerce_categories;
