-- Create learning microservice tables following DDD principles

-- Create learn_teachers table
CREATE TABLE IF NOT EXISTS learn_teachers (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    avatar VARCHAR(255) NOT NULL,
    title VARCHAR(100) NOT NULL,
    introduction TEXT NOT NULL,
    specialties JSON NOT NULL COMMENT '专长领域',
    status VARCHAR(20) NOT NULL COMMENT '讲师状态: active, inactive, suspended',
    courses_count INT NOT NULL DEFAULT 0,
    students_count INT NOT NULL DEFAULT 0,
    rating DECIMAL(3, 2) NOT NULL DEFAULT 0,
    rating_count INT NOT NULL DEFAULT 0,
    contact_email VARCHAR(100) NOT NULL,
    contact_phone VARCHAR(20) NOT NULL,
    social_profiles JSON NOT NULL COMMENT '社交档案',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_teacher_user (user_id),
    INDEX idx_teacher_status (status),
    INDEX idx_teacher_rating (rating)
);

-- Create learn_categories table
CREATE TABLE IF NOT EXISTS learn_categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    icon VARCHAR(255) NOT NULL,
    parent_id VARCHAR(36) NULL,
    level INT NOT NULL DEFAULT 1 COMMENT '分类层级，从1开始',
    `order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
    courses_count INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_category_parent (parent_id),
    INDEX idx_category_level (level),
    INDEX idx_category_order (`order`),
    INDEX idx_category_active (is_active)
);

-- Create learn_courses table
CREATE TABLE IF NOT EXISTS learn_courses (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    teacher_id VARCHAR(36) NOT NULL,
    cover VARCHAR(255) NOT NULL,
    level VARCHAR(20) NOT NULL COMMENT '课程难度: beginner, intermediate, advanced',
    duration INT NOT NULL DEFAULT 0 COMMENT '总时长(分钟)',
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    discount_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL COMMENT '课程状态: draft, published, archived',
    category_ids JSON NOT NULL COMMENT '分类ID列表',
    tags JSON NOT NULL COMMENT '标签列表',
    chapters_count INT NOT NULL DEFAULT 0,
    lessons_count INT NOT NULL DEFAULT 0,
    enrollments_count INT NOT NULL DEFAULT 0,
    rating DECIMAL(3, 2) NOT NULL DEFAULT 0,
    rating_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL,
    INDEX idx_course_teacher (teacher_id),
    INDEX idx_course_status (status),
    INDEX idx_course_level (level),
    INDEX idx_course_price (price),
    INDEX idx_course_rating (rating),
    INDEX idx_course_created (created_at),
    INDEX idx_course_published (published_at),
    INDEX idx_course_popular (enrollments_count)
);

-- Create learn_chapters table
CREATE TABLE IF NOT EXISTS learn_chapters (
    id VARCHAR(36) PRIMARY KEY,
    course_id VARCHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    `order` INT NOT NULL DEFAULT 0,
    lesson_count INT NOT NULL DEFAULT 0,
    duration INT NOT NULL DEFAULT 0 COMMENT '章节总时长(分钟)',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_chapter_course (course_id),
    INDEX idx_chapter_order (`order`)
);

-- Create learn_lessons table
CREATE TABLE IF NOT EXISTS learn_lessons (
    id VARCHAR(36) PRIMARY KEY,
    course_id VARCHAR(36) NOT NULL,
    chapter_id VARCHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    type VARCHAR(20) NOT NULL COMMENT '课时类型: video, article, audio',
    status VARCHAR(20) NOT NULL COMMENT '课时状态: draft, published',
    `order` INT NOT NULL DEFAULT 0,
    duration INT NOT NULL DEFAULT 0 COMMENT '课时时长(分钟)',
    video_url VARCHAR(255) NOT NULL DEFAULT '',
    video_size BIGINT NOT NULL DEFAULT 0,
    article_content LONGTEXT NOT NULL DEFAULT '',
    audio_url VARCHAR(255) NOT NULL DEFAULT '',
    audio_size BIGINT NOT NULL DEFAULT 0,
    is_free BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL,
    INDEX idx_lesson_course (course_id),
    INDEX idx_lesson_chapter (chapter_id),
    INDEX idx_lesson_order (`order`),
    INDEX idx_lesson_type (type),
    INDEX idx_lesson_status (status),
    INDEX idx_lesson_free (is_free)
);

-- Create learn_enrollments table
CREATE TABLE IF NOT EXISTS learn_enrollments (
    id VARCHAR(36) PRIMARY KEY,
    course_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    order_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL COMMENT '报名状态: active, completed, expired, refunded',
    progress DECIMAL(5, 2) NOT NULL DEFAULT 0 COMMENT '学习进度(百分比)',
    completed_count INT NOT NULL DEFAULT 0 COMMENT '已完成课时数',
    total_count INT NOT NULL DEFAULT 0 COMMENT '总课时数',
    last_learn_time TIMESTAMP NULL,
    rating DECIMAL(3, 2) NULL,
    comment TEXT NULL,
    expires_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    INDEX idx_enrollment_course (course_id),
    INDEX idx_enrollment_user (user_id),
    INDEX idx_enrollment_order (order_id),
    INDEX idx_enrollment_status (status),
    INDEX idx_enrollment_created (created_at),
    INDEX idx_enrollment_progress (progress),
    UNIQUE KEY uk_enrollment_user_course (user_id, course_id)
);

-- Create learn_lesson_progress table for tracking individual lesson progress
CREATE TABLE IF NOT EXISTS learn_lesson_progress (
    id VARCHAR(36) PRIMARY KEY,
    enrollment_id VARCHAR(36) NOT NULL,
    lesson_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    progress DECIMAL(5, 2) NOT NULL DEFAULT 0 COMMENT '该课时学习进度(百分比)',
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    watch_duration INT NOT NULL DEFAULT 0 COMMENT '观看时长(秒)',
    last_position INT NOT NULL DEFAULT 0 COMMENT '上次观看位置(秒)',
    notes TEXT NULL COMMENT '个人笔记',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    INDEX idx_progress_enrollment (enrollment_id),
    INDEX idx_progress_lesson (lesson_id),
    INDEX idx_progress_user (user_id),
    INDEX idx_progress_completed (completed),
    UNIQUE KEY uk_progress_user_lesson (user_id, lesson_id)
);

-- Create learn_course_reviews table for course reviews and ratings
CREATE TABLE IF NOT EXISTS learn_course_reviews (
    id VARCHAR(36) PRIMARY KEY,
    course_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    enrollment_id VARCHAR(36) NOT NULL,
    rating DECIMAL(3, 2) NOT NULL,
    content TEXT NOT NULL,
    reply TEXT NULL COMMENT '讲师回复',
    replied_at TIMESTAMP NULL COMMENT '回复时间',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_review_course (course_id),
    INDEX idx_review_user (user_id),
    INDEX idx_review_enrollment (enrollment_id),
    INDEX idx_review_rating (rating),
    UNIQUE KEY uk_review_enrollment (enrollment_id)
);

-- Create learn_course_categories mapping table
CREATE TABLE IF NOT EXISTS learn_course_categories (
    id VARCHAR(36) PRIMARY KEY,
    course_id VARCHAR(36) NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_cc_course (course_id),
    INDEX idx_cc_category (category_id),
    UNIQUE KEY uk_course_category (course_id, category_id)
);

-- Create learn_certificates table for course completion certificates
CREATE TABLE IF NOT EXISTS learn_certificates (
    id VARCHAR(36) PRIMARY KEY,
    enrollment_id VARCHAR(36) NOT NULL,
    course_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    certificate_number VARCHAR(100) NOT NULL COMMENT '证书编号',
    issue_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    certificate_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_certificate_enrollment (enrollment_id),
    INDEX idx_certificate_course (course_id),
    INDEX idx_certificate_user (user_id),
    UNIQUE KEY uk_certificate_enrollment (enrollment_id),
    UNIQUE KEY uk_certificate_number (certificate_number)
);
