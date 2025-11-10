# Test script for orders and invoices
Write-Host "Testing order and invoice functionality..." -ForegroundColor Green

# Test manually with curl commands - run each command separately
Write-Host "1. Start server: go run ." -ForegroundColor Yellow

Write-Host "2. Register: curl -X POST http://localhost:8080/register -H 'Content-Type: application/json' -d '{\"username\":\"testuser\",\"password\":\"testpass123\"}'" -ForegroundColor Yellow

Write-Host "3. Login: curl -X POST http://localhost:8080/login -H 'Content-Type: application/json' -d '{\"username\":\"testuser\",\"password\":\"testpass123\"}'" -ForegroundColor Yellow

Write-Host "`n--- SINGLE PRODUCT ORDER (Legacy) ---" -ForegroundColor Cyan
Write-Host "4a. Create Single Order: curl -X POST http://localhost:8080/orders/single -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_TOKEN' -d '{\"product_id\":1,\"quantity\":2}'" -ForegroundColor Yellow

Write-Host "`n--- MULTIPLE PRODUCTS IN ONE ORDER ---" -ForegroundColor Cyan
Write-Host "4b. Create Multi-Product Order: curl -X POST http://localhost:8080/orders -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_TOKEN' -d '{\"items\":[{\"product_id\":1,\"quantity\":2},{\"product_id\":2,\"quantity\":1},{\"product_id\":3,\"quantity\":3}]}'" -ForegroundColor Yellow

Write-Host "`n--- VIEW ORDERS ---" -ForegroundColor Cyan
Write-Host "5. Get All Orders: curl -X GET http://localhost:8080/orders -H 'Authorization: Bearer YOUR_TOKEN'" -ForegroundColor Yellow

Write-Host "`n--- CREATE INVOICES ---" -ForegroundColor Cyan
Write-Host "6. Create Invoice for Order: curl -X POST http://localhost:8080/invoices -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_TOKEN' -d '{\"order_id\":1}'" -ForegroundColor Yellow

Write-Host "`n--- VIEW INVOICES ---" -ForegroundColor Cyan
Write-Host "7. Get Invoice: curl -X GET http://localhost:8080/invoices/1 -H 'Authorization: Bearer YOUR_TOKEN'" -ForegroundColor Yellow

Write-Host "`n--- EXAMPLE MULTI-PRODUCT ORDER STRUCTURE ---" -ForegroundColor Magenta
Write-Host "The new order structure allows multiple products in one order:" -ForegroundColor White
Write-Host "{" -ForegroundColor Gray
Write-Host "  \`"items\`": [" -ForegroundColor Gray
Write-Host "    {\`"product_id\`": 1, \`"quantity\`": 2},  # 2 Laptops" -ForegroundColor Gray
Write-Host "    {\`"product_id\`": 2, \`"quantity\`": 1},  # 1 Phone" -ForegroundColor Gray
Write-Host "    {\`"product_id\`": 3, \`"quantity\`": 3}   # 3 Headphones" -ForegroundColor Gray
Write-Host "  ]" -ForegroundColor Gray
Write-Host "}" -ForegroundColor Gray
Write-Host ""
Write-Host "This creates ONE order with multiple products, automatic tax calculation," -ForegroundColor White
Write-Host "and a single invoice that itemizes all products!" -ForegroundColor White

Write-Host "`nNote: Replace YOUR_TOKEN with the actual JWT token from step 3" -ForegroundColor Red