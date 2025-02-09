# **Energy Microservice**

This is a **Golang-based microservice** for managing and retrieving energy consumption data. It allows querying consumption data by **daily, weekly, or monthly periods**.

## **Project Structure**

```
energy-microservice/
│── src/
│   ├── config/                  # Database and server configuration
│   │   ├── connection.go
│   │   ├── handler.go
│   ├── controllers/             # Controller handling request logic
│   │   ├── consumption.controller.go
│   ├── middlewares/             # Middleware functions
│   │   ├── recovery.middleware.go
│   ├── models/                  # Data models
│   │   ├── consumption.model.go
│   ├── repositories/            # Database queries
│   │   ├── consumption_repository.go
│   ├── routes/                  # API route definitions
│   │   ├── consumption.routes.go
│   ├── server/                  # Server setup and middleware
│   │   ├── middlewares.go
│   │   ├── routes.go
│   │   ├── server.go
│   ├── services/                # Business logic
│   │   ├── consumption.service.go
│   ├── test/                    # Unit tests
│   │   ├── consumption_repository_test.go
│   │   ├── consumption_service_test.go
│   ├── utils/                   # Utility functions
│   │   ├── loadCSVData.go
├── tmp/                         # Temporary storage (if needed)
├── .air.toml                    # Configuration for Air (if using live reload)
├── .env                         # Environment variables (ignored in Git)
├── .gitignore                   # Git ignored files
├── go.mod                        # Go modules dependencies
├── go.sum                        # Checksums for dependencies
├── main.go                       # Entry point
├── README.md                     # Project documentation
```

## **Prerequisites**

- Install **Go** (≥1.18)
- Install **PostgreSQL**
- Install **Gorilla Mux** (for routing)
- Install **GORM** (for database handling)
- Install **Testify** and **SQLMock** (for unit tests)

## **Installation**

Clone the repository and install dependencies:

```sh
git clone https://github.com/JuanConde27/energy-microservice.git
cd energy-microservice
go mod tidy
```

## **Setting Up the Database**

Update your **PostgreSQL** connection details in `.env`:

```sh
DATABASE_URL=YOUR_DATABASE_URL_HERE
```

## **Running the Microservice**

Start the service:

```sh
go run main.go
```

The server will start at `http://localhost:3000`.

## **CSV File Loading**

To load data from a CSV file, place the file in the project root and specify its path in `src/server/server.go`:

```go
csvPath := "test_bia.csv"
```

## **Usage**

### **Energy Consumption Endpoint**

#### **GET `/consumption`**

This endpoint retrieves energy consumption data based on parameters.

### **Query Parameters**

| Parameter      | Type     | Required | Description |
|---------------|---------|----------|-------------|
| `meters_ids`  | string  | ✅ Yes  | Comma-separated list of meter IDs (e.g., `1,2,3`) |
| `start_date`  | string  | ✅ Yes  | Start date in **YYYY-MM-DD** format |
| `end_date`    | string  | ✅ Yes  | End date in **YYYY-MM-DD** format |
| `kind_period` | string  | ✅ Yes  | `daily`, `weekly`, or `monthly` |

### **Examples**

#### **Daily Consumption**

```sh
curl "http://localhost:3000/consumption?meters_ids=1&start_date=2023-06-01&end_date=2023-06-10&kind_period=daily"
```

#### **Monthly Consumption**

```sh
curl "http://localhost:3000/consumption?meters_ids=1&start_date=2023-06-01&end_date=2023-07-10&kind_period=monthly"
```

#### **Weekly Consumption (Multiple Meters)**

```sh
curl "http://localhost:3000/consumption?meters_ids=1,2,3&start_date=2023-06-01&end_date=2023-06-26&kind_period=weekly"
```

### **How to Modify Query Parameters**

1. **Change meter IDs:** Update `meters_ids=1,2,3` with the desired meter IDs.
2. **Change date range:** Modify `start_date` and `end_date` (format: YYYY-MM-DD).
3. **Change period type:** Use `kind_period=daily`, `kind_period=weekly`, or `kind_period=monthly`.

## **Running Unit Tests**

Run all tests with:

```sh
go test -v ./src/test
```

Example output:

```
=== RUN   TestGetConsumptionByPeriod
✅ Passed: Daily consumption retrieval
--- PASS: TestGetConsumptionByPeriod (0.00s)
PASS
```

## **Git Flow Workflow**

### **Start a New Feature**

```sh
git flow feature start feature-name
```

### **Push the Feature Branch**

```sh
git push --set-upstream origin feature/feature-name
```

### **Finish a Feature**

```sh
git flow feature finish feature-name
git push origin develop
```

---

This **README** provides all necessary details for setting up and running the **Energy Microservice**. Let me know if you need any modifications! 🚀

