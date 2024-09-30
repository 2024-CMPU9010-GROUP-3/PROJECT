# Backend

This is the private backend of our application, which will only be used by internal services such as our machine learning containers. It is written in Go using the net/http package of the standard library. The server is equipped with a middleware that will log the URL, status code and response time of any request to std out.

## Usage

There are two options to run the backend, either locally or using Docker.

### Run in Docker

Requirements:

- Docker or Docker-compliant container runtime

Navigate to the `PROJECT/Backend` and build the Docker image:

```
docker build -t cmpu9010-backend-private .
```

Then run the container with the following command:

```
docker run -p 8081:8081 cmpu9010-backend-private
```

The server will then be reachable on `localhost:8081`

### Run Locally

Requirements:

- Golang >= 1.23.1

Navigate to `PROJECT/Backend/private` and run this command:

```
go run cmd/main/main.go
```

The server will then be reachable on `localhost:8081`

## Routes

### Points

#### `POST` `/v1/private/points`

- **Access:** private
- **Path Parameters:** None
- **Query Parameters:** None
- **Accepts:** Point details
- **Response:** JSON (PointId)
- **Description:** This route creates a new point in the database and returns its id

#### `PUT` `/v1/private/points/{id}`

- **Access:** private
- **Path Parameters:** PointId
- **Query Parameters:** None
- **Accepts:** Point details
- **Response:** None
- **Description:** This route updates an existing point's details

#### `DELETE` `/v1/private/points/{id}`

- **Access:** private
- **Path Parameters:** PointId
- **Query Parameters:** None
- **Accepts:** None
- **Response:** None
- **Description:** This route deletes a point from the database
