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

## Features
- **JWT Authentication**: Secure user authentication with JSON Web Tokens.
- **Session Management**: Session-based authentication using Redis caching.
- **API Documentation**: Automatically generated API documentation using Swagger and Swaggo.
- **Dockerized Application**: The project is fully containerized using Docker, making it easy to set up and deploy.
- **Reverse Proxy with Nginx**: Nginx is configured as a reverse proxy to route requests to the Go API server.

## Open Source
This project is open source, and you are welcome to explore, use, and contribute to its development.

## Project Status
This project is a work in progress and will continue to be improved.


