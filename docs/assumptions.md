# This document documents the assumptions i have made when designing the API

## 1. Timezone Standards
- All timestamps are generated and stored in UTC using `time.Now().UTC()`
- This avoids server side locality bugs and Daylight Saving Time anomalies. Timezone localization is the responsibility of frontend client applications.

## 2. Leave lifecycle
- All newly created leave requests begin in `PENDING` state.
- Once a leave is marked `APPROVED` or `REJECTED`, its state is final and immuteable. Any attempt to change change the state returns `409 Conflict`
- State transitions are executed atomically in SQL(`WHERE ID= ? AND status ="PENDING"`) to eliminate race conditions.

## 3. Modular Boundary
- `leave` package is forbiddine from quering the `employees` table direclty.
- As there's no foreign key, `leave_requests` table can exist on its own.
- `leave_requests` and `employees` table talk through Go interfaces, as these modules might be split into separate services.

## 4. Role-Based Access
- Authorization uses the `X-User-Role` HTTP header (`admin, manager, employee`)
- Request should contain the `X-User-Role` header.
- This keeps the project minimal without adding third-party auth dependencies

## 5. Security and Resource Protecion
- Request payload sizes are strictly bounded to 1MB to protect the server from DoS attacks.
- SQL queries use parameter binding(`?`) across all repositories to prevent SQL injection.
