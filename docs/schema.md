# Database Schema Documentation
This project uses SQLite as a database, as i thought it was the based option for running on my pc and it requires no setup for reviewers.

## Database Engine
* **Engine** : SQLITE 3(`modernc.org/sqlite`)
* **Default Location**: `hr.db`(configurable via `DB_PATH`)
* **Concurrency Settings**: Write-Ahead Logging(`WAL` mode), 5000ms `busy_timeout`, single connection pool(`MaxOpenConns=1`)

## Data Ownership and Boundaries

- Employee Management module owns table `employees`.
- Leave Management module owns table `leave_requests`.

> To preserve modular independence, the `leave_requests` table doesn't define a database-level foreign key to `employees`. Instead, employee existence is verified at the application boundary via `EmployeeCheckeer` interface.

## Module 1: Employee Management

Stores employee profiles, department, and role access.
#### DDL:
```sql
CREATE TABLE IF NOT EXISTS employees(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE COLLATE NOCASE,
        department TEXT NOT NULL,
        role TEXT NOT NULL DEFAULT 'employee' CHECK (role IN('admin','manager','employee')),
        status TEXT NOT NULL DEFAULT 'ACTIVE'CHECK (status In('ACTIVE','INACTIVE')),
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );

```

## Module 2: Leave Management

Stores employee leave requests,date intervals, and approval status.
```sql
CREATE TABLE IF NOT EXISTS leave_requests(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        employee_id INTEGER NOT NULL,
        leave_type TEXT NOT NULL CHECK(leave_type IN ('ANNUAL','SICK','UNPAID')),
        start_date TEXT NOT NULL,
        end_date TEXT NOT NULL,
        reason TEXT,
        status TEXT NOT NULL DEFAULT 'PENDING' CHECK(status IN('PENDING','APPROVED','REJECTED')),
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        CHECK (end_date >= start_date)
        );

```
### Indexing
- This makes searching for leaves by employee_id faster
- Uses SQLite's B-Tree index.
- Since it is a B-Tree it has B-tree's trade offs: slower writer but faster reads and takes extra storage
- Creates a copy of the table sorted by employee_id

```sql
CREATE INDEX IF NOT EXISTS idx_leave_requests_employee_id ON leave_requests(employee_id);
```
