# PowerShell test script for invoice functionality

Write-Host "Starting server..." -ForegroundColor Green
$serverProcess = Start-Process -FilePath "go" -ArgumentList "run", "." -PassThru -NoNewWindow
Start-Sleep -Seconds 5

try {
    Write-Host "1. Testing user registration..." -ForegroundColor Yellow
    $registerResponse = Invoke-RestMethod -Uri "http://localhost:8080/register" -Method POST -ContentType "application/json" -Body '{"username":"testuser","password":"testpass123"}'
    Write-Host "Registration: $registerResponse" -ForegroundColor Cyan

    Write-Host "2. Testing user login..." -ForegroundColor Yellow
    $loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method POST -ContentType "application/json" -Body '{"username":"testuser","password":"testpass123"}'
    $token = $loginResponse.token
    Write-Host "Token: $token" -ForegroundColor Cyan

    Write-Host "3. Testing order creation..." -ForegroundColor Yellow
    $headers = @{ Authorization = "Bearer $token" }
    $orderResponse = Invoke-RestMethod -Uri "http://localhost:8080/orders" -Method POST -ContentType "application/json" -Headers $headers -Body '{"product_id":1,"quantity":2}'
    Write-Host "Order: $orderResponse" -ForegroundColor Cyan

    Write-Host "4. Testing invoice creation..." -ForegroundColor Yellow
    $invoiceResponse = Invoke-RestMethod -Uri "http://localhost:8080/invoices" -Method POST -ContentType "application/json" -Headers $headers -Body '{"order_id":1}'
    Write-Host "Invoice: $invoiceResponse" -ForegroundColor Cyan

    Write-Host "5. Testing invoice retrieval..." -ForegroundColor Yellow
    $getInvoiceResponse = Invoke-RestMethod -Uri "http://localhost:8080/invoices/1" -Method GET -Headers $headers
    Write-Host "Retrieved Invoice: $($getInvoiceResponse | ConvertTo-Json -Depth 10)" -ForegroundColor Cyan

} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
} finally {
    Write-Host "Stopping server..." -ForegroundColor Green
    Stop-Process -Id $serverProcess.Id -Force
}

Write-Host "Test completed." -ForegroundColor Green