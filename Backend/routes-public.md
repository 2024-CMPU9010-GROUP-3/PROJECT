# Public Routes

## Auth

### `GET` `/v1/public/auth/User/{id}`

- **Access:** public
- **Path Parameters:** UserId
- **Query Parameters:** None
- **Accepts:** None
- **Response:** JSON (User details)
- **Description:** This route returns the user details object for the given UserId

### `POST` `/v1/public/auth/User/`

- **Access:** public
- **Path Parameters:** None
- **Query Parameters:** None
- **Accepts:** User details
- **Response:** JSON (UserId)
- **Description:** This route creates a new user in the database and returns their id

### `DELETE` `/v1/public/auth/User/{id}`

- **Access:** public
- **Path Parameters:** UserId
- **Query Parameters:** None
- **Accepts:** None
- **Response:** None
- **Description:** This route deletes a user from the database

### `PUT` `/v1/public/auth/User/{id}`

- **Access:** public
- **Path Parameters:** UserId
- **Query Parameters:** None
- **Accepts:** JSON (User details)
- **Response:** None
- **Description:** This route updates an existing user's details

### `POST` `/v1/public/auth/User/login`

- **Access:** public
- **Path Parameters:** None
- **Query Parameters:** None
- **Accepts:** JSON (User login)
- **Response:** JSON (Bearer Token)
- **Description:** This route allows a user to login with a username and password. If authentication succeeds, returns a bearer token for authentication.

## Points

### `GET` `/v1/public/points/byRadius?long={}&lat={}&radius={}`

- **Access:** public
- **Path Parameters:** None
- **Query Parameters:** Longitude, Latitude, Radius
- **Accepts:** None
- **Response:** JSON (List of points)
- **Description:** This route returns all points in a square with a given radius around the given latitude and longitude

### `GET` `/v1/public/points/{id}`

- **Access:** public
- **Path Parameters:** PointId
- **Query Parameters:** None
- **Accepts:** None
- **Response:** JSON (Point details)
- **Description:** This route retrieves details about the given point from the database
