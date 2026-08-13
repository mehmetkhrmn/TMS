# Ticket Management API

A RESTful Ticket Management API built with **Go (Golang)**, **Gin Framework**, and **PostgreSQL**. This project provides a simple ticket tracking system that enables communication between customers and support representatives through tickets and answers.

---

## Features

### Ticket Operations

- Create tickets
- Retrieve all tickets
- Retrieve a ticket by ID
- Filter tickets by status
- Update ticket information
- Change ticket status

### Answer Operations

- Create answers for tickets
- Retrieve all answers of a ticket
- Retrieve a specific answer
- Update answers

### Customer Operations

- Create customers
- Retrieve all customers
- Retrieve a customer by ID
- Update customer information

### Representative Operations

- Create representatives
- Retrieve all representatives
- Retrieve a representative by ID
- Update representative information

---

## Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin
- **Database:** PostgreSQL
- **Database Driver:** `database/sql` + `lib/pq`

---

# Project Structure

```text
TMS/
├── internal/
│   ├── database/
│   │   └── database.sql       # Database configuration
│   ├── models/
│   │   └── models.go          # Data models
│   ├── repository/
│   │   └── repository.go      # Database operations
│   └── routers/
│       └── routers.go         # API routes
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
└── README.md                  # Project documentation
```

---

## Installation

### Clone the repository

```bash
git clone https://github.com/mehmetkhrmn/TMS.git
cd TMS
```

### Install dependencies

```bash
go mod tidy
```

### Run the application

```bash
go run main.go
```

The server will start at:

```text
http://localhost:8080
```

---

# Data Models

## Ticket

```go
type Ticket struct {
	ID            int       `json:"id"`
	Subject       string    `json:"subject"`
	Description   string    `json:"description"`
	CustomerID    int       `json:"customer_id"`
	CustomerEmail string    `json:"customer_email"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Status        string    `json:"status"`
}
```

| Field | Type | Description |
|------|------|-------------|
| id | int | Ticket identifier |
| subject | string | Ticket subject |
| description | string | Detailed issue description |
| customer_id | int | Customer ID |
| customer_email | string | Customer email |
| created_at | time.Time | Creation timestamp |
| updated_at | time.Time | Last update timestamp |
| status | string | Ticket status (`open`, `in_progress`, `resolved`, `closed`) |

---

## Customer

```go
type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

| Field | Type | Description |
|------|------|-------------|
| id | int | Customer identifier |
| name | string | Customer name |
| email | string | Customer email |
| created_at | time.Time | Creation timestamp |
| updated_at | time.Time | Last update timestamp |

---

## Representative

```go
type Representative struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

| Field | Type | Description |
|------|------|-------------|
| id | int | Representative identifier |
| name | string | Representative name |
| created_at | time.Time | Creation timestamp |
| updated_at | time.Time | Last update timestamp |

---

## Answer

```go
type Answer struct {
	ID         int       `json:"id"`
	AnswerText string    `json:"answer"`
	RepID      int       `json:"representative_id"`
	TicketID   int       `json:"ticket_id"`
	AnsweredAt time.Time `json:"answered_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
```

| Field | Type | Description |
|------|------|-------------|
| id | int | Answer identifier |
| answer | string | Representative's response |
| representative_id | int | Representative ID |
| ticket_id | int | Related ticket ID |
| answered_at | time.Time | Creation timestamp |
| updated_at | time.Time | Last update timestamp |

---

# API Endpoints

## Tickets

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tickets` | Returns all tickets |
| GET | `/tickets/:ticket_id` | Return a ticket by ID |
| GET | `/tickets?ticket_status=in_progress` | Filter tickets by status |
| POST | `/tickets` | Create a ticket |
| PUT | `/tickets/:ticket_id` | Update a ticket |
| PATCH | `/tickets/:ticket_id?status=open` | Change ticket status |

---

## Answers

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tickets/:ticket_id/answers` | Returns all answers |
| GET | `/tickets/:ticket_id/answers/:answer_id` | Return a specific answer |
| POST | `/tickets/:ticket_id/answers` | Create an answer |
| PUT | `/answers/:answer_id` | Update an answer |

---

## Customers

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/customers` | Returns all customers |
| GET | `/customers/:customer_id` | Return customer by ID |
| POST | `/customers` | Create a customer |
| PUT | `/customers/:customer_id` | Update customer |

---

## Representatives

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/representatives` | Returns all representatives |
| GET | `/representatives/:representative_id` | Return representative by ID |
| POST | `/representatives` | Create a representative |
| PUT | `/representatives/:representative_id` | Update representative |

---

# Example Requests

## Create a Ticket

```http
POST http://localhost:8080/tickets
Content-Type: application/json
```

```json
{
  "subject": "Unable to access my cart",
  "description": "Products disappear after being added to the cart.",
  "customer_id": 1,
  "customer_email": "test@mail.com"
}
```

---

## Update a Ticket

```http
PUT http://localhost:8080/tickets/10
Content-Type: application/json
```

```json
{
  "description": "The cart button is not working."
}
```

> Only include the fields you want to update.

---

## Change Ticket Status

```http
PATCH http://localhost:8080/tickets/10?status=closed
```

---

## Create an Answer

```http
POST http://localhost:8080/tickets/12/answers
Content-Type: application/json
```

```json
{
  "answer": "Please contact your bank because your payment appears to be blocked.",
  "representative_id": 3
}
```

---

# Design Details

- The API validates that `customer_id` and `customer_email` belong to the same customer.
- If they do not match, the API returns **400 Bad Request**.
- Customers must exist before tickets can be created.
- Representatives must exist before answers can be created.
- Database-generated fields (`id`, `created_at`, `updated_at`) cannot be modified by clients.
- During update operations, only the related fields included in the request body are updated.
- Supported ticket statuses:
  - `open`
  - `in_progress`
  - `resolved`
  - `closed`
## Running the Project

### 1. Clone the repository

```bash
git clone https://github.com/mehmetkhrmn/TMS.git
cd TMS
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Create the PostgreSQL database

Open PostgreSQL and create a new database.

```sql
CREATE DATABASE tms;
```
or 
in bash
```
psql -U postgres
```
login and create by using command above.
### 4. Import the database schema

Import the schema located in the `database` directory.

Using `psql`:

```bash
psql -U postgres -d tms -f internal/database/schema.sql
```

Or, if you are using GoLand:

1. Open the **Database** tool window.
2. Connect to your PostgreSQL server.
3. Right-click the **tms** database.
4. Select **Run SQL Script...**
5. Choose `ticket-service/internal/database/schema.sql`.
6. Execute the script.

### 5. Configure the database connection

Create and update the PostgreSQL connection settings in the .env file if needed.

Example:

```text
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres (or your username)
DB_PASSWORD=your_password
DB_NAME=tms
DB_SSLMODE=disable
```

### 6. Run the application

```bash
go run main.go
```

The API will start at:

```text
http://localhost:8080
```

---



Example request:

```http
GET http://localhost:8080/tickets
```

If the API returns a response, the installation was successful.

> **Note:** Before creating tickets or answers, you must first create at least one customer and one representative.

.
