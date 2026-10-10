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
 [] Add example request and response
- 500 Internal Server Error: Database write failure.
 [] Add example request and response

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







## Planned / Future Endpoints

- `GET /api/v1/employees`: List employees with pagination
- `DELETE /api/v1/employees/{id}`: Delete employee by ID, it is a soft delete, `status = 'INACTIVE'`
- `PUT /api/v1/employees/{id}`: Update employee details(department, role)


