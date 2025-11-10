# Quick Start Guide for Postman Testing

## 🚀 Import Collections

### Step 1: Import Files
1. Open Postman
2. Click **Import** button (or Ctrl+O)
3. Drag and drop these files or browse to select them:
   - `Go-Auth-Sessions-API.postman_collection.json`
   - `Go-Auth-Sessions-Environment.postman_environment.json`
   - `Go-Auth-Sessions-Workflow.postman_collection.json`

### Step 2: Set Environment
1. In Postman, click the environment dropdown (top right)
2. Select **Go Auth Sessions Environment**
3. Make sure your Go server is running on `http://localhost:8080`

## 🧪 Two Ways to Test

### Option A: Complete Workflow Test (Recommended)
**Use the "Go Auth Sessions - Complete Workflow Test" collection**

1. Click on **Collections** in the left sidebar
2. Find **Go Auth Sessions - Complete Workflow Test**
3. Click the **Run** button (or use **Runner**)
4. Click **Run Go Auth Sessions - Complete Workflow Test**

**This will automatically:**
- ✅ Check server health
- ✅ Register a test user
- ✅ Login and save JWT token
- ✅ Get available products
- ✅ Create a test product
- ✅ Create an order with multiple products
- ✅ Generate an invoice from the order
- ✅ Update invoice status through the workflow
- ✅ Verify final state

### Option B: Manual Testing
**Use the "Go Auth Sessions - Complete API" collection**

1. Start with **Authentication > Login User** to get your JWT token
2. Test individual endpoints in any order
3. JWT token is automatically included in protected requests

## 🔧 Configuration

### Change Server URL
- Go to **Environments** tab
- Edit **Go Auth Sessions Environment**
- Change `base_url` value (default: `http://localhost:8080`)

### Change Test Credentials  
- Edit environment variables:
  - `username` (default: testuser)
  - `password` (default: password123)

## 📊 What Gets Tested

### **All 18 API Endpoints:**
- 6 Public endpoints (no authentication required)
- 12 Protected endpoints (JWT required)

### **Complete Business Workflow:**
1. User registration and authentication
2. Product catalog management
3. Multi-product order creation
4. Invoice generation with complete product details
5. Status tracking and updates

### **Automatic Validation:**
- Response status codes
- JWT token handling
- ID capture and reuse
- Data consistency across services
- Error handling

## 🎯 Success Indicators

### Console Output:
- ✅ Green checkmarks for successful operations
- 🎉 Final success message
- 📊 Summary statistics

### Test Results:
- All tests should pass (green)
- No failed assertions (red)
- Proper data flow between requests

## 🐛 Troubleshooting

### Common Issues:

1. **Connection Error**
   - Make sure Go server is running: `go run .`
   - Check base_url in environment

2. **401 Unauthorized**
   - Run Login request first to get JWT token
   - Check that token is saved in environment

3. **404 Not Found**
   - Verify endpoint URLs match your server routes
   - Check for typos in path parameters

4. **Empty Product Lists**
   - Server should auto-seed sample products
   - Check database connection and initialization

### Debug Steps:
1. Check Postman Console for detailed logs
2. Verify environment variables are populated
3. Test health endpoint first to confirm server connectivity
4. Run requests individually if workflow fails

## 📈 Advanced Usage

### Custom Test Data
Edit request bodies to test with different:
- Product names and prices
- Order quantities
- Invoice statuses
- User credentials

### Environment Variables
Available for use in requests:
- `{{base_url}}` - Server URL
- `{{jwt_token}}` - Authentication token
- `{{order_id}}` - Created order ID
- `{{invoice_id}}` - Created invoice ID
- `{{username}}` - Test username
- `{{password}}` - Test password

### Batch Testing
Use Postman Runner to:
- Run tests multiple times
- Test with different data sets
- Generate performance reports
- Export test results