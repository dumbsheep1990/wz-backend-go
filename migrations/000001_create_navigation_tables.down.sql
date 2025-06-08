-- Drop indexes
DROP INDEX IF EXISTS idx_websites_category_id;
DROP INDEX IF EXISTS idx_websites_is_featured;
DROP INDEX IF EXISTS idx_websites_view_count;
DROP INDEX IF EXISTS idx_categories_sort_order;

-- Drop tables
DROP TABLE IF EXISTS navigation_websites;
DROP TABLE IF EXISTS navigation_categories;
