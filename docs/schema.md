# Database Schema Documentation
This project uses SQLite as a database, as i thought it was the based option for running on my pc and it requires no setup for reviewers.

## Database Engine
* **Engine** : SQLITE 3(`modernc.org/sqlite`)
* **Default Location**: `hr.db`(configurable via `DB_PATH`)

## Module: Employee Management

Stores employee profiles, department, and role access.

#### DDL:
```sql
CREATE TABLE IF NOT EXISTS employees(
id INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT NOT NULL,
email TEXT NOT NULL UNIQUE,
department TEXT NOT NULL,
role TEXT NOT NULL DEFAULT 'employee',
status TEXT NOT NULL DEFAULT 'ACTIVE',
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```
