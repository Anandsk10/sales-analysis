# Sales Data Analysis System

A comprehensive backend system for analyzing sales data from CSV files with RESTful APIs, automated data refresh, and performance optimization for large datasets.

## Features

- **Normalized Database Schema**: Efficient storage with proper relationships
- **CSV Data Loading**: Batch processing for large files (millions of rows)
- **Automated Data Refresh**: Daily scheduled refresh + on-demand triggers
- **RESTful Analytics APIs**: Comprehensive revenue and sales analysis
- **Performance Optimized**: Handles large datasets efficiently
- **Comprehensive Logging**: Detailed logs for monitoring and debugging
- **Docker Support**: Easy deployment with Docker Compose
- **Postman Collection**: Ready-to-use API collection for testing

## System Architecture

![System Design](docs/sales_system_design.png)

The system follows a layered architecture with:
- **Client Layer**: Postman collection, frontend apps, API consumers
- **API Layer**: RESTful endpoints with proper routing
- **Service Layer**: Business logic for analytics, data refresh, and customer management
- **Database Layer**: Normalized PostgreSQL schema with optimized indexes
- **Monitoring**: Health checks, structured logging, and error tracking

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+ (or Docker for containerized setup)
- Make (optional, for using Makefile commands)
- Postman (optional, for API testing with provided collection)

## Database Schema

![Entity Relationship Diagram](docs/er_diagram.png)

The system uses a normalized database schema with the following tables:

### Tables
- **customers**: Customer information (ID, name, email, address)
- **products**: Product catalog (ID, name, category)
- **orders**: Order details (ID, customer_id, region, date, payment_method, shipping_cost)
- **order_items**: Order line items (order_id, product_id, quantity, price, discount)
- **refresh_logs**: Data refresh activity logs

### Relationships
- One customer can have many orders
- One order can have many order items
- One product can be in many order items
- Foreign key constraints ensure data integrity

## Quick Start

### Option 1: Docker Compose (Recommended)

1. Clone the repository
```bash
git clone <repository-url>
cd sales-analysis
```

2. Start the services
```bash
make docker-up
```

3. Wait for services to start, then load sample data
```bash
make load-data
```

4. Import Postman Collection (Optional)
```bash
# The postman collection is available at: postman.json
# Import this file into Postman to get all API endpoints pre-configured
```

### Option 2: Local Development

1. Install PostgreSQL and create database
```bash
createdb salesdb
```

2. Set environment variables
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/salesdb?sslmode=disable"
export PORT=8080
export LOG_LEVEL=info
```

3. Install dependencies and run
```bash
go mod tidy
make run
```

4. Load sample data
```bash
curl -X POST "http://localhost:8080/api/v1/refresh?file_path=data/sales_data.csv"
```

## API Testing with Postman

### Postman Collection
The project includes a comprehensive Postman collection with all API endpoints pre-configured:

**File**: `postman.json`

### How to Use:
1. Open Postman
2. Click "Import" → "Upload Files"
3. Select `postman.json`
4. Configure environment variables:
   - `base_url`: `http://localhost:8080` (or your deployment URL)
   - `port`: `8080`

### Collection Structure:
```
Sales Analysis API/
├── Data Management/
│   ├── Trigger Data Refresh
│   └── Get Refresh Status
├── Revenue Analytics/
│   ├── Get Total Revenue
│   ├── Get Revenue by Product
│   ├── Get Revenue by Category
│   └── Get Revenue by Region
├── Product Analytics/
│   ├── Get Top Products
│   └── Get Top Products by Category
├── Customer Analytics/
│   ├── Get Customer Count
│   ├── Get Order Count
│   └── Get Average Order Value
└── Health Check/
    └── Service Health Status
```

### Pre-configured Variables:
- Date ranges for testing different periods
- Query parameters for pagination and filtering
- Authentication headers (if needed in future)

## API Endpoints

### Data Refresh
| Method | Endpoint | Description | Sample Response |
|--------|----------|-------------|-----------------|
| POST | `/api/v1/refresh` | Trigger data refresh | `{"message": "Data refresh triggered successfully", "status": "in_progress"}` |
| GET | `/api/v1/refresh/status` | Get refresh history | `{"data": [{"id": 1, "status": "success", "records_count": 6}]}` |

### Revenue Analytics
| Method | Endpoint | Query Params | Description | Sample Response |
|--------|----------|--------------|-------------|-----------------|
| GET | `/api/v1/analytics/revenue/total` | `start_date`, `end_date` | Total revenue | `{"data": {"revenue": 4337.94, "count": 6}}` |
| GET | `/api/v1/analytics/revenue/by-product` | `start_date`, `end_date` | Revenue by product | `{"data": [{"product_id": "P456", "product_name": "iPhone 15 Pro", "revenue": 2597.0}]}` |
| GET | `/api/v1/analytics/revenue/by-category` | `start_date`, `end_date` | Revenue by category | `{"data": [{"category": "Electronics", "revenue": 2946.99}]}` |
| GET | `/api/v1/analytics/revenue/by-region` | `start_date`, `end_date` | Revenue by region | `{"data": [{"region": "Asia", "revenue": 2776.95}]}` |

### Product Analytics
| Method | Endpoint | Query Params | Description | Sample Response |
|--------|----------|--------------|-------------|-----------------|
| GET | `/api/v1/analytics/products/top` | `start_date`, `end_date`, `limit` | Top products by quantity | `{"data": [{"product_id": "P789", "product_name": "Levi's 501 Jeans", "total_sold": 3}]}` |
| GET | `/api/v1/analytics/products/top/by-category` | `start_date`, `end_date`, `category`, `limit` | Top products in category | `{"data": [{"product_id": "P456", "product_name": "iPhone 15 Pro", "total_sold": 3}]}` |

### Customer Analytics
| Method | Endpoint | Query Params | Description | Sample Response |
|--------|----------|--------------|-------------|-----------------|
| GET | `/api/v1/analytics/customers/count` | `start_date`, `end_date` | Unique customer count | `{"data": {"customer_count": 3}}` |
| GET | `/api/v1/analytics/orders/count` | `start_date`, `end_date` | Total order count | `{"data": {"order_count": 6}}` |
| GET | `/api/v1/analytics/orders/average-value` | `start_date`, `end_date` | Average order value | `{"data": {"average_order_value": 722.99}}` |

### Health Check
| Method | Endpoint | Description | Sample Response |
|--------|----------|-------------|-----------------|
| GET | `/health` | Service health status | `{"status": "healthy", "timestamp": "2024-01-01T12:00:00Z", "database": "connected"}` |

## Usage Examples

### 1. Using cURL

#### Load Data
```bash
curl -X POST "http://localhost:8080/api/v1/refresh?file_path=data/sales_data.csv"
```

#### Get Total Revenue for Date Range
```bash
curl "http://localhost:8080/api/v1/analytics/revenue/total?start_date=2024-01-01&end_date=2024-12-31"
```

#### Get Top 5 Products
```bash
curl "http://localhost:8080/api/v1/analytics/products/top?limit=5&start_date=2024-01-01&end_date=2024-12-31"
```

#### Get Revenue by Category
```bash
curl "http://localhost:8080/api/v1/analytics/revenue/by-category?start_date=2024-01-01&end_date=2024-12-31"
```

### 2. Using Postman Collection

1. Import the collection: `postman.json`
2. Set environment variables:
   - `base_url`: Your server URL (default: `http://localhost:8080`)
3. Run requests individually or use Postman's Collection Runner for batch testing
4. Use pre-configured test data and date ranges
5. Monitor response times and validate JSON schemas

## Configuration

Environment variables:
- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

## Data Format

Expected CSV format:
```csv
Order ID,Product ID,Customer ID,Product Name,Category,Region,Date of Sale,Quantity Sold,Unit Price,Discount,Shipping Cost,Payment Method,Customer Name,Customer Email,Customer Address
1001,P123,C456,UltraBoost Running Shoes,Shoes,North America,2023-12-15,2,180.00,0.1,10.00,Credit Card,John Smith,johnsmith@email.com,"123 Main St, Anytown, CA 12345"
```

## Performance Optimizations

1. **Batch Processing**: CSV data loaded in batches of 1000 records
2. **Database Indexing**: Indexes on frequently queried columns
3. **Connection Pooling**: Configured connection limits
4. **Transaction Management**: Bulk operations wrapped in transactions
5. **Query Optimization**: Efficient SQL queries with proper joins
6. **Memory Management**: Streaming CSV processing for large files
7. **Caching Strategy**: Query result caching for frequently accessed data

## Logging

- **Structured JSON logging** with configurable levels
- **Request/response logging** middleware
- **Data refresh activity logs** with detailed metrics
- **Error tracking** and debugging information
- **Performance metrics** including query execution times
- **Log rotation** and retention policies

### Log Locations
- **Docker**: `docker-compose logs app`
- **Local**: `./logs/` directory

## Testing

### Unit Tests
```bash
make test
```

### Integration Tests
```bash
go test -v ./internal/services/...
```

### API Testing with Postman
1. Import collection: `postman.json`
2. Run collection tests
3. Use Newman for CI/CD integration:
```bash
newman run postman.json -e postman/environments/local.json
```

## Monitoring

- **Health check endpoint**: `/health` with database connectivity status
- **Refresh status tracking**: `/api/v1/refresh/status` with detailed logs
- **Structured logs** for monitoring system integration
- **Database connection health** monitoring with alerts
- **Performance metrics** collection and reporting
- **Error rate tracking** and alerting
