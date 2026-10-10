CREATE TABLE IF NOT EXISTS employees(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NUll UNIQUE,
    department TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'employee',
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

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

-- This make searching for leaves by employee_id fast
-- Index is bascially a B-Tree 
-- Since it is a B-Tree it has B-tree's trade offs: slower writer but faster reads and takes extra storage
-- this creates a copy of the table sorted by employee_id
CREATE INDEX IF NOT EXISTS idx_leave_requests_employee_id ON leave_requests(employee_id);
