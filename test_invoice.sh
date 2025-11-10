#!/bin/bash

# Test script for invoice functionality

echo "Starting server..."
go run . &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo "1. Testing user registration..."
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}' \
  && echo

echo "2. Testing user login..."
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}' | jq -r '.token')

echo "Token: $TOKEN"

echo "3. Testing order creation..."
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"product_id":1,"quantity":2}' \
  && echo

echo "4. Testing invoice creation..."
curl -X POST http://localhost:8080/invoices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"order_id":1}' \
  && echo

echo "5. Testing invoice retrieval..."
curl -X GET http://localhost:8080/invoices/1 \
  -H "Authorization: Bearer $TOKEN" \
  && echo

# Kill the server
kill $SERVER_PID
echo "Test completed."