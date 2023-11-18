# edge-receiver

This repository contains the code for the edge-receiver project. The edge-receiver is responsible for receiving and processing data from various sources.

## Development Instructions

To get started with development, follow the instructions below:

### Prerequisites

Before you begin, make sure you have the following installed on your system:

- Go programming language
- Reflex (install with `go install github.com/cespare/reflex@latest`)
- Docker and Docker Compose

### Setting up Kafka

The edge-receiver uses Kafka as a message broker. To run Kafka and the Kafka UI, follow these steps:

1. Open a terminal and navigate to the project directory.
2. Start Kafka and the Kafka UI using Docker Compose by running the following command:
    ```bash
    docker-compose -f kafka.yaml up
    ```

### Running the edge-receiver

To run the edge-receiver, follow these steps:

1. Open a terminal and navigate to the project directory.
2. Start the edge-receiver using the following command:
    ```bash
    make run
    ```

### Development Workflow

To streamline the development process, we recommend using Reflex. Reflex is a tool that automatically rebuilds and restarts your application whenever a file changes.

To use Reflex, open a terminal and navigate to the project directory. Then, run the following command:
```bash
make dev
```
