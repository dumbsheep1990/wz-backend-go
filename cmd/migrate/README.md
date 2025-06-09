# Database Migration Tool

A command-line tool for managing database migrations in the WanZhi platform, following DDD principles.

## Usage

### Creating a New Migration

To create a new migration for a specific domain:

```bash
go run cmd/migrate/main.go -create -domain=learning -name="add_course_rating_table"
```

This will create two files:
- `db/migrations/learning/00X_add_course_rating_table.up.sql`
- `db/migrations/learning/00X_add_course_rating_table.down.sql`

Where `X` is the next sequence number for the domain.

### Running Migrations

To run migrations for all domains:

```bash
go run cmd/migrate/main.go -migrate -db="user:password@tcp(localhost:3306)/wz_backend"
```

To run migrations for specific domains:

```bash
go run cmd/migrate/main.go -migrate -domain=learning,commerce -db="user:password@tcp(localhost:3306)/wz_backend"
```

## Environment Variables

You can set the following environment variables instead of using command-line flags:

- `DATABASE_URL`: Database connection string

## Migration Structure

Migrations are organized by domain:

- `/db/migrations/core/` - Core system tables
- `/db/migrations/user/` - User-related tables
- `/db/migrations/learning/` - Learning microservice tables
- `/db/migrations/commerce/` - Commerce microservice tables
- etc.

Each domain directory contains numbered SQL migration files in sequence.
