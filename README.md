# Weather Radar

## Description

- Weather Radar is a Go-based application designed to aggregate weather data from multiple sources.
- It provides a REST API for retrieving weather information and submitting feedback.
- The application is optimized for efficient performance with in-memory caching and includes rate-limiting middleware to
  prevent excessive requests.

## Features

* Retrieve weather data from multiple sources
* Submit feedback on weather data
* Rate-limiting middleware to prevent excessive requests
* Feedback submission with Basic Auth
* In-memory caching for improved performance
* Dockerized application for easy deployment

## Installation

### Used tools

- Go 1.23.2
- Docker

### Setup

1. Clone the repository:
    ```bash
    git clone
    ```
2. Navigate to the project directory:
    ```bash
    cd weather-radar
   ```

### Run locally

1. Create a `.env` file in the root directory and add the following environment variables:
    ```text
    WEATHER_API_KEY=your_api_key
    OPEN_WEATHER_API_KEY=your_api_key
    ```
2. Run the application:
    ```bash
    go run cmd/main.go
    ```
### Run with Docker
1. Change following env vars in docker-compose.yml file:
     ```text
     WEATHER_API_KEY=your_api_key
     OPEN_WEATHER_API_KEY=your_api_key
     ```
2. Run docker-compose:
    ```bash
    docker compose up -d
    ```
### Access the application
3. The application will be available at `http://localhost:1312`.

## API Docs

- Directory: `docs` contains the API documentation in Swagger format. Also, it contains http file with test requests.

### HTTP file

- Import the `weather_api.http` file into your IDE to test the API endpoints.

### Swagger

- Swag CLI for generating Swagger documentation:
    ```bash
    go get -u github.com/swaggo/swag/cmd/swag
    ```
- Generate Swagger documentation:
    ```bash
    swag init -g cmd/main.go
    ```
- The Swagger documentation will be available at `http://localhost:1312/swagger-ui/index.html`.

