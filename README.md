# Simple REST API Auth

**Simple REST API Auth** is a beginner project aimed at creating a basic authentication system using Go. This project serves as an introduction to building RESTful APIs, implementing JWT-based authentication, session-based authentication with Redis caching in Go, and generating API documentation using Swagger. It showcases fundamental concepts and practices in web development.

## Technologies Used
- **Go**: The primary programming language used for developing this API.
- **net/http**: The standard library used for routing and handling HTTP requests.
- **JWT-Based Auth**: Authentication system using JSON Web Tokens (JWT) for securely transmitting information between parties.
- **Session-Based Auth**: A session-based approach to manage user authentication, utilizing Redis for session caching.
- **Redis**: Used as a cache store to handle session data efficiently.
- **PostgreSQL**: The database used for storing user information and refresh tokens.
- **Slog**: Structured logging library to provide clear and actionable logs.
- **Validator**: For validating request payloads to ensure they meet expected formats.
- **Cleanenv**: A library for reading and processing environment variables.
- **Docker**: Containerization of the application, allowing for easy setup, scaling, and deployment.
- **Nginx**: Used as a reverse proxy to route incoming requests to the Go application.
- **Swagger**: API documentation generated using the Swaggo library.
- **Testify**: A testing library for writing unit tests in Go.
- **Mockery**: A library for generating mocks, used alongside Testify to test components in isolation.

## Installation

1. **Clone the Repository**:
    ```sh
    git clone https://github.com/Soloda1/go-rest-api-auth.git
    cd go-rest-api-auth
    ```

2. **Set Up Environment Variables**:
    Create a `.env` file in the root directory. You can use the `example.env` file as a template:
    ```sh
    cp example.env .env
    ```

3. **Build and Run the Docker Containers**:
    ```sh
    docker-compose up --build
    ```

## Utilities

1. **Run Unit Tests**:
    ```sh
    go test -v ./...
    ```

2. **API Documentation**:
    ```sh
    http://localhost/swagger/index.html
    ```


## Open Source
This project is open source, and you are welcome to explore, use, and contribute to its development.

## Project Status
This project is a work in progress and will continue to be improved.


