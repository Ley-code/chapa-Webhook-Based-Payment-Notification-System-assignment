# Chapa - Webhook-Based Payment Notification System Assignment

This repository contains the solution for the take-home assignment from Chapa Financial Technologies S.C.

The project simulates a simple, webhook-based payment notification system built in Go. It consists of two main components:
1.  A **Payment Processor** server that accepts payment requests, simulates processing, and dispatches a webhook upon completion.
2.  A mock **Client Server** that listens for and logs these incoming webhook notifications.

## Features

-   **Simulated Payment Processing**: Accepts payment details via a RESTful API endpoint.
-   **Asynchronous Webhook Notifications**: Uses goroutines to process payments in the background without blocking the initial API response.
-   **Clean Architecture**: The code is structured following Clean Architecture principles for better maintainability, scalability, and testability.
-   **Robust Unit Tests**: The core business logic is verified with a unit test suite that covers key success and failure scenarios.
-   **Concurrent & Safe**: The in-memory data store is protected with a mutex to safely handle concurrent requests.

## Architectural Approach: Clean Architecture

This project is built using the principles of **Clean Architecture** to ensure a robust and maintainable codebase. This approach decouples the business logic from external concerns like databases, frameworks, and other third services.

The core layers are:

-   **`domain`**: Contains the core business entities and objects.
-   **`repository`**: Defines the interfaces for data persistence and provides the in-memory implementation.
-   **`usecase`**: Orchestrates the business logic and data flow. It contains the application's specific rules.
-   **`handler`**: The outermost layer, responsible for handling HTTP requests and interfacing with the usecase.

This separation of concerns makes the system highly testable, as demonstrated by the unit tests for the usecase layer.

## Project Structure

```plaintext
.
├── server/
│   ├── domain/
|   |   ├── payment_entity.go
│   ├── handler/
|   |   ├── payment_handler.go
│   ├── repository/
|   |   ├── payment_repository.go
│   ├── usecase/
│   │   ├── payment_usecase.go
│   │   └── payment_usecase_test.go 
│   └── main.go
└── client/
    └── main.go
```

## Getting Started

### Prerequisites

-   Go (version 1.18 or higher)
-   Git

### How to Run

Follow these steps to run the processor and client servers locally. You will need three separate terminal windows.

#### Step 1: Clone the Repository

```bash
git clone https://github.com/{your-github-username}/chapa-webhook-notification-system-assignment.git
cd chapa-webhook-notification-system-assignment
```

#### Step 2: Run the Client Server (Terminal 1)

This server will listen for incoming webhook notifications.

```bash
cd client
go run main.go
```
You should see the output: `Mock Client listening for webhooks on http://localhost:8081/webhook`

#### Step 3: Run the Payment Processor Server (Terminal 2)

This server will accept payment requests.

```bash
cd server
go run main.go
```
You should see the output: `Payment Processor server starting on http://localhost:8080`

#### Step 4: Send a Test Payment Request (Terminal 3)

Use `curl` to send a POST request to the payment processor.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 150.75,
    "currency": "ETB",
    "webhookUrl": "http://localhost:8081/webhook"
  }' \
  http://localhost:8080/api/v1/payment
```

## Testing

This project includes a unit test suite for the core business logic located in the `usecase` layer. The tests verify the system's reliability by ensuring the synchronous part of the payment creation process works as expected.

The test suite uses a flexible mock repository (a "test double") that acts as a "spy" to confirm that:
-   The repository's `Create` method is called correctly.
-   A new payment is created with the initial `PENDING` status.
-   A unique ID is generated for the new payment.

This approach validates the correctness of the business logic in isolation, which is a key benefit of the Clean Architecture pattern.

### Running the Tests

To run the unit tests for the server application, navigate to the `server/` directory and execute the following command:

```bash
go test ./... -v
```

The `-v` flag enables verbose output, which will list the test cases as they run and confirm they have passed.

## Future Enhancements

This implementation provides a solid, well-tested foundation. For a production-ready system, the following features would be the next logical steps:

-   **Persistent Storage**: Replace the in-memory repository with a persistent database (e.g., PostgreSQL or MongoDB).
-   **Enhanced Security**: Implement HMAC webhook signature verification to allow the client to confirm that a webhook is authentic.
-   **Configuration Management**: Externalize configuration like server ports using environment variables.
-   **Robust Error Handling & Retries**: Add a retry mechanism for failed webhook deliveries to handle temporary client unavailability.
-   **Containerization**: Add `Dockerfile`s and a `docker-compose.yml` to streamline the setup and deployment process.