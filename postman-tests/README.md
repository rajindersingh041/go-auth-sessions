# Postman Test Collection for Go Auth Sessions

This folder contains comprehensive Postman collections for testing all API endpoints with automated JWT token management.

## 📁 Files Included

### 1. `Go-Auth-Sessions-API.postman_collection.json`
Complete API collection with all endpoints organized by domain:

#### **🔐 Authentication**
- **Register User** - Creates a new user account
- **Login User** - Authenticates and automatically saves JWT token

#### **⚡ Health Check**
- **Health Status** - System health endpoint

#### **📦 Products (Public)**
- **Get All Products** - List all available products
- **Get Product by ID** - Retrieve specific product details
- **Get Products by Category** - Filter products by category

#### **📦 Products (Protected)**
- **Create Product** - Add new product (requires JWT)
- **Update Product Stock** - Modify product availability (requires JWT)

#### **🛒 Orders**
- **Create Order (Multiple Products)** - Create order with multiple items
- **Create Single Product Order** - Create order with single item
- **Get User Orders** - Retrieve user's orders
- **Get Orders by Username (Legacy)** - Legacy endpoint for username-based queries
- **Create Order by Username (Legacy)** - Legacy endpoint for creating orders

#### **🧾 Invoices**
- **Create Invoice from Order** - Generate invoice from existing order
- **Get Invoice by ID** - Retrieve invoice details
- **Get Invoice by Order ID** - Get invoice using order ID
- **Get Invoices by User ID** - List all user invoices
- **Update Invoice Status** - Change invoice status (draft, sent, paid, cancelled)

### 2. `Go-Auth-Sessions-Environment.postman_environment.json`
Environment variables for easy configuration:
- `base_url` - API server URL (default: http://localhost:8080)
- `jwt_token` - Automatically managed JWT token
- `order_id` - Automatically captured from order creation
- `invoice_id` - Automatically captured from invoice creation
- `username` - Test username (default: testuser)
- `password` - Test password (default: password123)

## 🚀 How to Use

### Step 1: Import into Postman
1. Open Postman
2. Click **Import** button
3. Select both JSON files:
   - `Go-Auth-Sessions-API.postman_collection.json`
   - `Go-Auth-Sessions-Environment.postman_environment.json`

### Step 2: Set Environment
1. In Postman, select **Go Auth Sessions Environment** from the environment dropdown
2. Make sure your Go server is running on `http://localhost:8080`

### Step 3: Run Tests

#### **Option A: Manual Testing**
1. Start with **Authentication > Register User** to create an account
2. Run **Authentication > Login User** to get JWT token (automatically saved)
3. Test any protected endpoints - JWT token is automatically included

#### **Option B: Automated Testing**
1. Use **Runner** in Postman to run the entire collection
2. The collection will automatically:
   - Register a user
   - Login and capture JWT token
   - Test all endpoints in logical order
   - Create orders and invoices with proper dependencies

## ✨ Key Features

### 🔄 **Automatic Token Management**
- JWT token is automatically captured from login response
- All protected endpoints automatically use the saved token
- No manual token copying required

### 📋 **Automatic ID Capture**
- Order IDs are automatically captured and used for invoice creation
- Invoice IDs are automatically captured for status updates
- Seamless workflow testing

### 🧪 **Comprehensive Coverage**
- Tests all 18 API endpoints
- Covers both public and protected routes
- Includes legacy endpoints for backward compatibility
- Tests complete business workflow: Registration → Products → Orders → Invoices

### 📊 **Smart Test Scripts**
- Automatic response validation
- Console logging for debugging
- Environment variable management
- Error handling and reporting

## 🔧 Configuration

### Change Server URL
Update the `base_url` environment variable if your server runs on a different port or host.

### Change Test Credentials
Modify the `username` and `password` environment variables to use different test credentials.

### Custom Test Data
Edit the request bodies in the collection to test with different product names, quantities, etc.

## 📝 Testing Workflow

### **Recommended Testing Order:**
1. **Health Check** - Verify server is running
2. **Register User** - Create test account
3. **Login User** - Get authentication token
4. **Get All Products** - View available products
5. **Create Product** - Add new product (optional)
6. **Create Order** - Place an order
7. **Create Invoice** - Generate invoice from order
8. **Update Invoice Status** - Change status to 'paid'

### **Business Logic Testing:**
- Products can be viewed without authentication
- Orders require valid JWT tokens
- Invoices aggregate data from orders, products, and users
- Status tracking works across orders and invoices

## 🐛 Troubleshooting

### Common Issues:
1. **401 Unauthorized**: Make sure to run Login first to get JWT token
2. **404 Not Found**: Verify server is running and base_url is correct
3. **Order/Invoice not found**: Ensure you've created orders before creating invoices

### Debug Tips:
- Check Postman Console for detailed logs
- Verify environment variables are set correctly
- Make sure server is running with correct database configuration