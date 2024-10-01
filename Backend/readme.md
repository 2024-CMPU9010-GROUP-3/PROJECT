# Backend

This is the backend of our application. It has two configurations (private and public) which are defined using build tags at compile time. The private backend will only be used by internal services such as our machine learning containers. The public backend provides the endpoints used by the frontend when deployed.

The backend is written in Go using the net/http package of the standard library. The server is equipped with a middleware that will log the URL, status code and response time of any request to std out.

## Routes

- [Public Routes](./routes-public.md)
- [Private Routes](./routes-private.md)

## Usage

There are two options to run the backend, either locally or using Docker.

### Run in Docker

Requirements:

- Docker or Docker-compliant container runtime

Navigate to the `PROJECT/Backend` and build the Docker image:

#### Build Private Backend

```
docker build -f Dockerfile.private -t cmpu9010-backend-private .
```

Then run the container with the following command:

```
docker run -p 8081:8080 cmpu9010-backend-private
```

The server will then be reachable on `localhost:8081`

#### Build Public Backend

```
docker build -f Dockerfile.public -t cmpu9010-backend-public .
```

Then run the container with the following command:

```
docker run -p 8080:8080 cmpu9010-backend-public
```

The server will then be reachable on `localhost:8080`

### Run Locally

#### Run Private Backend

Requirements:

- Golang >= 1.23.1

Navigate to `PROJECT/Backend` and run this command:

```
go run -tags private cmd/main/main.go
```

The server will then be reachable on `localhost:8080`

#### Run Public Backend

Requirements:

- Golang >= 1.23.1

Navigate to `PROJECT/Backend` and run this command:

```
go run -tags public cmd/main/main.go
```

The server will then be reachable on `localhost:8080`

### Command Line Arguments

- `-p <Port Number>` or `-port <Port Number>` sets the port number the server will listen on - default: `8080`
