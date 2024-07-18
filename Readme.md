Microservices Architecture with Kafka, PostgreSQL, and NATS

This project implements a microservices architecture using Golang for backend services, Kafka for messaging, PostgreSQL for data storage, and NATS for event-driven communication.
Technologies Used

    Golang: Used for building microservices.
    Gin: HTTP web framework for Golang.
    Kafka: Distributed event streaming platform.
    PostgreSQL: Relational database management system.
    NATS: Lightweight and fast messaging system.
    Docker: Containerization for easy deployment.
    Docker Compose: Tool for defining and running multi-container Docker applications.

Project Structure

The project is structured into two microservices:

    Transactions Microservice
        Handles transactions and interacts with the transactions_microservice PostgreSQL database.
        Consumes messages from Kafka topics.
        Exposes endpoints to manage transactions.

    Users Microservice
        Manages user data and interactions.
        Uses the users_microservice PostgreSQL database.
        Consumes messages from Kafka topics.
        Exposes endpoints for user management.

Setup
Prerequisites

    Docker and Docker Compose installed on your system.

Running the Application

    Clone the repository:

    bash

git clone <repository-url>
cd <repository-directory>

Start the services using Docker Compose:

bash

    docker-compose up

    This command will build the Docker images and start the containers for Kafka, PostgreSQL databases, NATS, and the two microservices (transactions-microservice and users-microservice).

    Access the services:
        Transactions Microservice: http://localhost:8081
        Users Microservice: http://localhost:8080

Environment Variables

    PostgreSQL:
        POSTGRES_USER: PostgreSQL username.
        POSTGRES_PASSWORD: PostgreSQL password.
        POSTGRES_DB: PostgreSQL database name.

    Kafka:
        KAFKA_BROKER: Kafka broker address.
        KAFKA_ZOOKEEPER_CONNECT: Zookeeper address for Kafka.

    NATS:
        NATS_SERVER: NATS server address.

    Microservices:
        DB_HOST: Hostname for PostgreSQL.
        DB_PORT: PostgreSQL port.
        DB_USER: PostgreSQL username.
        DB_PASSWORD: PostgreSQL password.
        DB_NAME: PostgreSQL database name.
        PORT: Port on which the microservice will run.

Development

For local development:

    Ensure Go is installed.
    Use go run to run individual services.
    Use go test to run tests.

Troubleshooting

    If any service fails to start, ensure all dependencies (like Kafka, PostgreSQL) are up and running.
    Check logs (docker-compose logs <service-name>) for detailed error messages.

Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your improvements.
