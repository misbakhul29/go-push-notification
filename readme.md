# Go Push Notification Service

A high-performance, distributed push notification server built with **Go**, **RabbitMQ**, **Redis**, and **WebSockets**. This service consumes notification jobs from a message queue and delivers them in real-time to connected browser clients.

## 🏗 Architecture

The system follows a **Pub/Sub** pattern to ensure scalability across multiple server instances.

1.  **Ingestion:** External services publish notification jobs to a **RabbitMQ** queue.
2.  **Consumption:** The Go server consumes messages from RabbitMQ.
3.  **Routing:** The server publishes the message to a specific **Redis** channel (`user:{id}:notify`).
4.  **Delivery:** The specific Go instance holding the user's WebSocket connection subscribes to that Redis channel and pushes the payload to the browser.

## 🚀 Prerequisites

* **Go** 1.22+
* **Docker** & **Docker Compose**
* **Make** (optional, for build commands)

## 🛠 Installation & Setup

1.  **Clone the repository**
    ```bash
    git clone [https://github.com/your-username/go-push-service.git](https://github.com/your-username/go-push-service.git)
    cd go-push-service
    ```

2.  **Environment Configuration**
    Copy the example environment file:
    ```bash
    cp .env.example .env
    ```
    *Modify `.env` if you need to change ports or credentials (defaults work out-of-the-box).*

3.  **Start Infrastructure (RabbitMQ & Redis)**
    ```bash
    docker-compose up -d
    ```

4.  **Install Go Dependencies**
    ```bash
    go mod tidy
    ```

## 🏃 running the Server

You can run the server directly using Go or via the Makefile.

```bash
# OR Standard Go Command
go run cmd/server/main.go