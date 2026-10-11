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

## Module: Leave Management

Stores employee leave requests,date intervals, and approval status.
```sql
CREATE TABLE IF NOT EXISTS leave_requests(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER NOT NULL,
    leave_type TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    reason TEXT,
    status TEXT NOT NULL DEFAULT 'PENDING',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```
### Indexing
- This make searching for leaves by employee_id faster
- Uses SQLite's B-Tree index.
- Since it is a B-Tree it has B-tree's trade offs: slower writer but faster reads and takes extra storage
- Creates a copy of the table sorted by employee_id

```sql
CREATE INDEX IF NOT EXISTS idx_leave_requests_employee_id ON leave_requests(employee_id);
```
