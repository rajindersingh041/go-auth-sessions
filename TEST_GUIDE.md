# Test Script for Go Auth Sessions

## Prerequisites
Ensure the server is running:
```powershell
go run .
```

## 1. Health Check
```bash
curl http://localhost:8080/health
```
Expected: `{"status":"ok","timestamp":"2025-11-10T22:16:37Z"}`

## 2. User Registration
```bash
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d "{\"username\":\"alice\",\"password\":\"password123\"}"
```
Expected: `{"message":"User created successfully"}`

## 3. User Login
```bash
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d "{\"username\":\"alice\",\"password\":\"password123\"}"
```
Expected: JWT token response

## 4. Create Order
```bash
curl -X POST http://localhost:8080/orders/alice -H "Content-Type: application/json" -d "{\"item\":\"Book\",\"quantity\":2}"
```
Expected: `{"message":"Order created successfully"}`

## 5. Get Orders
```bash
curl -X GET http://localhost:8080/orders/alice
```
Expected: List of orders for alice

## Database Testing

### Test with PostgreSQL (default)
```bash
# Set in .env file: DB_DRIVER=postgres
go run .
```

### Test with ClickHouse
```bash
# Set in .env file: DB_DRIVER=clickhouse
go run .
```

## Expected Server Output
```
2025/11/10 22:16:25 Using database driver: postgres
2025/11/10 22:16:25 Starting server on port 8080...
2025/11/10 22:16:25 Server started successfully. Press Ctrl+C to stop.
```