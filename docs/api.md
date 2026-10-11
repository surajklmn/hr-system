# HR System API documentation

Base URL: `http://localhost:8080/api/v1`

## 1. Health Check
### `GET /health`
Return the status of the service

#### Example:
```bash
curl -i http://localhost:8080/health
```
#### Response( 200 OK):

    {
        "status": "ok"
    }

## 2. Employees Module

### POST /api/v1/employees

note: Requires Header : X-User-Role: admin.
`eg, -H "X-User-Role:admin"`

Create a new employee record

#### Request Headers:
- Content-Type: application/json

#### Request Body Parameters:

Field | Type | Required | Default | Description
name    string  Yes         -       Full name
email   string  Yes         -       Unique email
department string Yes       -       Assigned department
role    string  No      "employee"  Allowed: admin, manager, employee

#### Example request:

```bash
curl -i -X POST http://localhost:8080/api/v1/employees \
-H "Content-Type: application/json" \
-H "X-User-Role:admin" \
-d '{
        "name":"Ram Bahadur",
        "email":"ram@honda.com",
        "department":"Sales",
        "role":"manager"
    }'
```

#### Responses:
- 201 Created:
```json
{
        "id":1,
        "name":"Ram Bahadur",
        "email":"ram@honda.com",
        "department":"Sales",
        "role":"manager",
        "created_at":"2026-10-10T12:19:23.818904766Z"
        "updated_at":"2026-10-10T12:19:23.818904766Z"
    }
```
- 400 Bad Request: Valiation error(e.g. missing name, invalid email format, or unauthorized role.
- 401 Unauthorized: missing `X-User-Role` header.
- 403 Forbidden: insufficient permissions.
- 409 Conflict: email already exists.
- 500 Internal Server Error: Database write failure.

### Get /api/v1/employees/{id}
Fetches an employee profile by their id

#### URL Parameters

Paramater | Type                 | Required           |Description
id        |integer(int64)        | Yes                | empoyee's PK ID



### Example Request:

`curl -i http://localhost:8080/api/v1/employees/1`

### Responses:
- 200 OK
```json
{
        "id": 1,
        "name":"Ram Bahadur",
        "email":"ram@honda.com",
        "department":"Sales",
        "role":"manager",
        "status":"ACTIVE",
        "created_at":"2026-10-10T12:19:23.818904766Z",
        "updated_at":"2026-10-10T12:19:23.818904766Z"
    }
```

- 400 Bad Request: Malformed ID (e.g /api/v1/employees/asdf).
- 404 Not Found: No employee found with the given ID.

## 3. Leave Module

### POST /api/v1/leaves
Submit a new leave request

#### Request Headers:
- Content-Type: application/json

#### Request Body Parameters:

Field | Type               |Required  | Descriptong
employee_id| integer(int64)|Yes
leave_type |string         |Yes        |Allowed: ANNUAL,SICK, UNPAID
start_date |string         |Yes        |Format: YYYY-MM-DD
end_date   |string         |Yes        |Format: YYYY-MM-DD
reason     |string         |No         |Optional

#### Example Request:
```bash
curl -i -X POST http://localhost:8080/api/v1/leaves \
-H "Content-Type: application/json" \
-d '{
        "employee_id":1,
        "leave_type":"ANNUAL",
        "start_date":"2026-11-01",
        "end_date": "2026-11-05",
        "reason" : "Vacation"
    }'
```
#### Responses:
- 201 Created:
```json
{"id": 1, 
"employee_id": 1,
"leave_type": "ANNUAL",
"start_date": "2026-11-01",
"end_date": "2026-11-05",
"reason": "Autumn vacation",
"created_at": "2026-10-11T04:12:18.471415449Z",                     
"updated_at": "2026-10-11T04:12:18.471415449Z"}
```

- 400 Bad Request: Validation error(e.g invalid date format, end_date before start_date, invalid leave_type).
- 404 Not Found: Employee doesn't exist.

### GET /api/v1/leaves/{id}

Fetches a leave request by its ID

#### URL Parameters:

Paramter | Type          | Required | Description
id       |integer(int64) | Yes      | Leave requests's PK ID


#### Example Request:
`curl -i http://localhost:8080/api/v1/leaves/1`

- 200 OK:

```json

{"id":1,
"employee_id":1,
"leave_type":"ANNUAL",
"start_date":"2026-11-01",
"end_date":"2026-11-05",
"reason":"Vacation",
"status":"PENDING",
"created_at":"2026-10-11T04:12:18.471415449Z",
"updated_at":"2026-10-11T04:45:38.183375704Z"}
```

- 400 Bad Request: invalid ID(e.g /api/v1/leaves/abc)
- 404 Not Found: Leave request not found.


### POST /api/v1/leaves/{id}/approve
note: Requires Header: X-User-Role: admin or manager
eg: `-H "X-User-Role:admin"`
eg: `-H "X-User-Role:manager"`
Approves a pending leave request.

### URL Parameters:
Paramter | Type          | Required | Description
id       |integer(int64) | Yes      | Leave requests's PK ID

#### Examples Request:
`curl -i -X POST http://localhost:8080/api/v1/leaves/1/approve -H "X-User-Role:manager"`

#### Responses:
- 200 OK: Returns leave with "status":"APPROVED".
- 401 Unauthorized: missing `X-User-Role` header.
- 403 Forbidden: insufficient permissions.
- 404 Not Found: Leave request not found.
- 409 Conflict: Leave request is not in PENDING status.

### POST /api/v1/leaves/{id}/reject
note: Requires Header : `X-User-Role:admin` or `X-User-Role:manager`
Rejects a pending leave request.

#### Example Request:
`curl -i -X POST http://localhost:8080/api/v1/leaves/1/reject -H "X-User-Role:manager"`

#### Responses:
- 200 OK: Returns leave with `"status":"REJECTED"`.
- 401 Unauthorized: missing `X-User-Role` header.
- 403 Forbidden: insufficient permissions.
- 404 Not Found: leave request not found.
- 409 Conflict: leave request is not in PENDING status.


## Planned / Future Endpoints

- `GET /api/v1/employees`: List employees with pagination
- `DELETE /api/v1/employees/{id}`: Delete employee by ID, it is a soft delete, `status = 'INACTIVE'`
- `PUT /api/v1/employees/{id}`: Update employee details(department, role)
- `GET /api/v1/leaves?employee_id={id}`: List allleaves for a specific employee.
- Leave Balance Deduction: Track employee balance and deduct days on approval.
- Detect overlapping leaves and return 409 Conflict


