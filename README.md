# taxi-finder

This project consists of two APIs: driver-location-api and matching-api, designed to facilitate the matching of drivers with nearby customers.

## Project Overview

The project consists of two APIs:

driver-location-api: Manages driver data and provides endpoints for uploading driver locations and retrieving driver information.

matching-api: Matches customers with the nearest available drivers based on their locations.

Both APIs are developed in Go and utilize MongoDB as the backend database.

## Setup 

Ensure you have the following dependencies installed on your system:

Docker
Docker Compose

## Build and run the application using Docker Compose:

```


docker-compose up --build

```


This command will build the Docker images for the APIs and spin up the containers.

Once the containers are up and running, you can access the APIs via the following endpoints:

Driver Location API: http://localhost:8080
Matching API: http://localhost:8081

## Environment Variables

The following environment variables can be configured in the .env file:

API_KEY: API key required for authenticating requests to both APIs.

JWT_SECRET: Secret key for JWT token generation (used in the matching-api).

Make sure to set these variables before running the application.

## Example cURL Requests for Driver API

### Get All Driver

```
curl --location 'http://localhost:8080/drivers/' \
--header 'Content-Type: application/json' \
--header 'X-API-Key: q7EJx3H8Lz9RdG5s1PfA2oKbF6cVgY4n' \
--data ''
```

### Create Driver Location and Insert DB

```
curl --location 'http://localhost:8080/drivers/create' \
--header 'Content-Type: application/json' \
--header 'X-API-Key: q7EJx3H8Lz9RdG5s1PfA2oKbF6cVgY4n' \
--data '{
    "count" : 3
}'
```

### Get Driver With an ID

```
curl --location 'http://localhost:8080/drivers/d5148549-a9a8-4f94-82b3-c8f683197911' \
--header 'Content-Type: application/json' \
--header 'X-API-Key: q7EJx3H8Lz9RdG5s1PfA2oKbF6cVgY4n' \
--data ''
```

## Example cURL Requests for Matching API

### Get Nearest Driver

```
curl --location 'http://localhost:8081/nearest-driver?latitude=41.0082&longitude=28.9784' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJhdXRoZW50aWNhdGVkIjp0cnVlfQ.bY0bTO6Krxex3CcM4VS3zOcfXffIpnML-FOPolhQ40U'

```






