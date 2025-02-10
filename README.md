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
│   ├── test/                    # tests
│   │   ├── consumption_repository_test.go
│   │   ├── consumption_service_test.go
│   │   ├── integration_test.go
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

- Install **Go** (≥1.22)
- Install **PostgreSQL**
- Install **Gorilla Mux** (for routing)
- Install **GORM** (for database handling)
- Install **Testify** and **SQLMock** (for unit tests)
- **Place the file `test_bia.csv` in the root of the project before running the service.**

## **Installation (Local Setup)**

Clone the repository and install dependencies:

```sh
git clone https://github.com/JuanConde27/energy-microservice.git
cd energy-microservice
go mod tidy
```

### **Setting Up the Database**

Update your **PostgreSQL** connection details in `.env` if using a local database:

```sh
DATABASE_URL=YOUR_DATABASE_URL_HERE
```

### **Running the Microservice (Locally)**

Start the service:

```sh
go run main.go
```

The server will start at `http://localhost:3000`.

## **Running with Docker Compose**

If you don't have Docker installed, you can download it from the [Docker Desktop website](https://www.docker.com/products/docker-desktop).

To run the microservice with **Docker Compose**, execute:

```sh
docker-compose up --build
```

This will:

- Build and start the **PostgreSQL database** and the **Go application**.
- Automatically create and migrate the database.
- Expose the service on **port 3000**.

### **Stopping the Service**

```sh
docker-compose down
```

This will stop all running containers.

## **CSV File Loading**

The microservice loads energy consumption data from a **CSV file**. The path varies depending on whether you run the service **locally** or **inside Docker**:

### **Running Locally**

Place the CSV file in the **project root** and specify its path in `src/server/server.go`:

```go
csvPath := "test_bia.csv"
```

### **Running with Docker**

Since the application runs inside a **Docker container**, the path should be set to `/app/test_bia.csv` in `src/server/server.go`:

```go
csvPath := "/app/test_bia.csv"
```

Ensure the CSV file is included in the Docker context, or mount it as a volume in `docker-compose.yml`:

```yaml
volumes:
  - ./test_bia.csv:/app/test_bia.csv
```

## **Usage**

### **Energy Consumption Endpoint**

#### **GET **``

This endpoint retrieves energy consumption data based on parameters.

### **Query Parameters**

| Parameter     | Type   | Required | Description                                       |
| ------------- | ------ | -------- | ------------------------------------------------- |
| `meters_ids`  | string | ✅ Yes    | Comma-separated list of meter IDs (e.g., `1,2,3`) |
| `start_date`  | string | ✅ Yes    | Start date in **YYYY-MM-DD** format               |
| `end_date`    | string | ✅ Yes    | End date in **YYYY-MM-DD** format                 |
| `kind_period` | string | ✅ Yes    | `daily`, `weekly`, or `monthly`                   |

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
📌 Ejecutando GetConsumptionByPeriod
🔹 Meter IDs: [1]
🔹 Start Date: 2023-06-01 00:00:00 +0000 UTC
🔹 End Date: 2023-06-30 23:59:59 +0000 UTC
🔹 Period Type: monthly
📌 Respuesta HTTP: {"period":["JUN 2023"],"data_graph":[{"meter_id":1,"address":"Dirección mock","active":[600],"reactive_inductive":[0],"reactive_capacitive":[0],"exported":[0]}]}
--- PASS: TestGetConsumptionEndpoint (0.02s)
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

This **README** provides all necessary details for setting up and running the **Energy Microservice**. 

