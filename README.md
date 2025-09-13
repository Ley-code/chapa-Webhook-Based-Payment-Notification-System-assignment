
# Chapa - Webhook-Based Payment Notification System Assignment

This repository contains the solution for the take-home assignment from Chapa Financial Technologies S.C.

The project is a simulation of a webhook-based payment system, built in Go. It has been enhanced with production-ready features like containerization with Docker, secure webhook verification using HMAC signatures, and environment-based configuration.

## Features

-   **Containerized with Docker**: Both the server and client services are fully containerized for a simple, one-command setup using Docker Compose.
-   **Secure Webhooks (HMAC)**: Webhook payloads are cryptographically signed using HMAC-SHA256 to ensure their integrity and authenticity.
-   **Environment-Based Configuration**: All configuration (ports, secrets) is managed via environment variables, following 12-Factor App principles.
-   **Clean Architecture**: The code is structured into distinct layers for better maintainability, scalability, and testability.
-   **Asynchronous Processing**: Uses goroutines for non-blocking payment processing, ensuring a responsive API.
-   **Robust Unit Tests**: The core business logic is verified with a comprehensive unit test suite.

## Architectural Approach: Clean Architecture

This project is built using the principles of **Clean Architecture** to ensure a robust and maintainable codebase. This approach decouples the business logic from external concerns like databases, frameworks, and other third-party services.

The core layers are:
-   **`domain`**: Contains the core business entities.
-   **`repository`**: Defines the interfaces for data persistence.
-   **`usecase`**: Orchestrates the business logic and data flow.
-   **`handler`**: The outermost layer, responsible for handling HTTP requests.

## Project Structure

```plaintext
.
├── docker-compose.yml
├── .gitignore
├── server/
│   ├── .env.example       
│   ├── Dockerfile
│   └── ... (source code)
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
    ├── .env.example 
    ├── Dockerfile
    └── main.go
```

## Getting Started (Docker Recommended)

### Prerequisites

-   Docker
-   Docker Compose
-   Git
-   Go (only required for running tests locally)

### How to Run

This project is designed to be run easily using Docker Compose.

#### Step 1: Clone the Repository

```bash
git clone https://github.com/{your-github-username}/chapa-webhook-notification-system-assignment.git
cd chapa-webhook-notification-system-assignment
```

#### Step 2: Set Up Environment Variables

This project uses `.env` files to manage secrets and configuration, which are kept out of Git for security. Template files (`.env.example`) are provided as a blueprint.

**For the server:**
```bash
cp server/.env.example server/.env
```

**For the client:**
```bash
cp client/.env.example client/.env
```
> **Important:** The `WEBHOOK_SECRET_KEY` in both `.env` files must match for the webhook signature verification to work. The default value is already set correctly in the example files.

#### Step 3: Build and Run with Docker Compose

From the root of the project, run this single command:
```bash
docker-compose up --build
```
This command will build the Docker images, create a private network for them, and start the containers. You will see the logs from both the server and client in your terminal.

#### Step 4: Send a Test Payment Request

Open a **new terminal** and use `curl` to send a request to the server.

> **Note:** Inside the Docker network, services communicate using their service names. Therefore, the `webhookUrl` must be `http://client:8081`, not `localhost`.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 250.00,
    "currency": "USD",
    "webhookUrl": "http://client:8081/webhook"
  }' \
  http://localhost:8080/api/v1/payment
```

### Expected Outcome

1.  Your `curl` command will receive an immediate **`202 Accepted`** response.
2.  In the `docker-compose` logs, the **server** will log that it is processing the payment.
3.  After 3 seconds, the server will log that it sent the webhook with an `X-Chapa-Signature` header:).
4.  The **client** will log that it has received the webhook and that the **signature was successfully verified**.

## Testing

To run the unit tests for the server's business logic, navigate to the `server/` directory and run:
```bash
go test ./... -v
```

## Future Enhancements

This implementation provides a solid, secure, and well-tested foundation. The next logical steps to make it fully production-grade would be:

-   **Persistent Storage**: Replace the in-memory repository with a persistent database (e.g., PostgreSQL or MongoDB) to ensure data durability.
-   **Robust Error Handling & Retries**: Implement a durable queue and retry mechanism (e.g., with exponential backoff) for failed webhook deliveries to handle cases where the client server is temporarily unavailable.
```