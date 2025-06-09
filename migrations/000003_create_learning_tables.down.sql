-- Drop learning microservice tables in reverse order to avoid foreign key constraint issues

DROP TABLE IF EXISTS learn_certificates;
DROP TABLE IF EXISTS learn_course_categories;
DROP TABLE IF EXISTS learn_course_reviews;
DROP TABLE IF EXISTS learn_lesson_progress;
DROP TABLE IF EXISTS learn_enrollments;
DROP TABLE IF EXISTS learn_lessons;
DROP TABLE IF EXISTS learn_chapters;
DROP TABLE IF EXISTS learn_courses;
DROP TABLE IF EXISTS learn_categories;
DROP TABLE IF EXISTS learn_teachers;
