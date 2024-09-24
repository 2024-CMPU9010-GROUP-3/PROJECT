# Backend

This is a backend written in Go using the net/http package of the standard library. The server is equipped with a middleware that will log the URL, status code and response time of any request to std out.

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
docker run -p 8080:8080 cmpu9010-backend
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

- **Parameters:** None  
- **Accepts:** None  
- **Response:** Plaintext  
- **Description:** This route returns _"World!"_ when queried to show the server is working and accepting requests.  
