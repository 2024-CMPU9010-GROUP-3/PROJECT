# Backend

## Usage

There are two options to run the backend, either locally or using Docker.

### Run in Docker

Requirements:

- Docker or Docker-compliant container runtime

Navigate to the `PROJECT/Backend` and build the Docker image:

```
docker build -t cmpu9010-backend .
```

Then run the container with the following command:

```
docker run cmpu9010-backend -p 8080:8080
```

The server will then be reachable on `localhost:8080`

### Run Locally

Requirements:

- Golang >= 1.23.1

Navigate to `PROJECT/Backend` and run this command:

```
go run cmd/main/main.go
```

The server will then be reachable on `localhost:8080`

## Routes

#### `GET` `/hello`

**Parameters:** None
**Accepts:** None
**Response:** Plaintext
This route returns _"World!"_ when queried to show the server is working and accepting requests.
