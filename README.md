# taxi-finder

This project consists of two APIs: driver-location-api and matching-api, designed to facilitate the matching of drivers with nearby customers.

## Project Overview

The project consists of two APIs:

driver-location-api: Manages driver data and provides endpoints for uploading driver locations and retrieving driver information.

matching-api: Matches customers with the nearest available drivers based on their locations.

Both APIs are developed in Go and utilize MongoDB as the backend database.

##Setup 

Ensure you have the following dependencies installed on your system:

Docker
Docker Compose

## Build and run the application using Docker Compose:

docker-compose up --build

This command will build the Docker images for the APIs and spin up the containers.

Once the containers are up and running, you can access the APIs via the following endpoints:

Driver Location API: http://localhost:8080
Matching API: http://localhost:8081

#Environment Variables

The following environment variables can be configured in the .env file:

API_KEY: API key required for authenticating requests to both APIs.

JWT_SECRET: Secret key for JWT token generation (used in the matching-api).

Make sure to set these variables before running the application.

