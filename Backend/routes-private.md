# Private Routes

## Points

### `POST` `/v1/private/points`

- **Access:** private
- **Path Parameters:** None
- **Query Parameters:** None
- **Accepts:** Point details
- **Response:** JSON (PointId)
- **Description:** This route creates a new point in the database and returns its id

### `PUT` `/v1/private/points/{id}`

- **Access:** private
- **Path Parameters:** PointId
- **Query Parameters:** None
- **Accepts:** Point details
- **Response:** None
- **Description:** This route updates an existing point's details

### `DELETE` `/v1/private/points/{id}`

- **Access:** private
- **Path Parameters:** PointId
- **Query Parameters:** None
- **Accepts:** None
- **Response:** None
- **Description:** This route deletes a point from the database