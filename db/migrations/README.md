# WanZhi Backend Database Migrations

This directory contains all database migration files organized by domain/microservice, following the Domain-Driven Design (DDD) principles used throughout the project.

## Directory Structure

Each subdirectory represents a specific domain/bounded context in our system:

- `admin/` - Admin panel tables
- `ai/` - AI-related tables
- `commerce/` - Commerce and store management
- `core/` - Core system tables
- `file/` - File storage and management
- `interaction/` - User interactions
- `learning/` - Learning platform (courses, teachers, enrollments)
- `navigation/` - Navigation menus and structures
- `notification/` - Notification system
- `search/` - Search functionality
- `statistics/` - Analytics and statistics
- `tenant/` - Multi-tenancy support
- `trade/` - Orders, payments, and transactions
- `user/` - User accounts and profiles

## Migration Order

The migrations should be run in the following order to ensure proper dependency relationships:

1. Core tables (users, configurations)
2. Domain-specific tables (learning, commerce, etc.)

Each migration file follows the naming convention: `<sequence_number>_<description>.[up|down].sql`
