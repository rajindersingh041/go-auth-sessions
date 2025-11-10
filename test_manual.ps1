# Simple test script
Write-Host "Testing invoice functionality..." -ForegroundColor Green

# Test manually with curl commands - run each command separately
Write-Host "1. Start server: go run ." -ForegroundColor Yellow
Write-Host "2. Register: curl -X POST http://localhost:8080/register -H 'Content-Type: application/json' -d '{\"username\":\"testuser\",\"password\":\"testpass123\"}'" -ForegroundColor Yellow
Write-Host "3. Login: curl -X POST http://localhost:8080/login -H 'Content-Type: application/json' -d '{\"username\":\"testuser\",\"password\":\"testpass123\"}'" -ForegroundColor Yellow
Write-Host "4. Create Order: curl -X POST http://localhost:8080/orders -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_TOKEN' -d '{\"product_id\":1,\"quantity\":2}'" -ForegroundColor Yellow
Write-Host "5. Create Invoice: curl -X POST http://localhost:8080/invoices -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_TOKEN' -d '{\"order_id\":1}'" -ForegroundColor Yellow
Write-Host "6. Get Invoice: curl -X GET http://localhost:8080/invoices/1 -H 'Authorization: Bearer YOUR_TOKEN'" -ForegroundColor Yellow