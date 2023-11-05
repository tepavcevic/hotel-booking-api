Hotel Booking is a simple REST API that allows users to make bookings for a given hotel.

## :wrench: Technologies

- Golang
- Docker
- PostgreSQL

# :checkered_flag: How to run (needs Docker :whale:)

- Copy the `.env.template` file content into a new file called `.env`
- After using `make docker`, the API will be available in `http://localhost:8080/api/v1`

### Running the application locally using Go and Docker

```zsh
  make run
```

### Running with Docker (Go not required)
```zsh
  make docker
```

### Seeding the database
```zsh
  make seed
```

## :vertical_traffic_light: Testing

### Run tests
```zsh
  make test
```


# :triangular_flag_on_post: Endpoints

## `/users`

### `/api/v1/users - POST` Creates a new user. Example of request body:

```json
{
  "firstName": "Maria",
  "lastName": "Curie",
  "email": "maria@curie.com",
  "password": "testpassword"
}
```

### `/api/v1/users - GET` Fetch all users

### `/api/v1/users/:id - GET` Get user by ID

### `/api/v1/users/:id - PUT` Update user by ID

```json
{
  "firstName": "Maria",
  "lastName": "Curie",
}
```

### `/api/v1/users/:id - DELETE` Delete user by ID

### Details

- Users endpoints requires authentication
- New users get isAdmin flag as false by default
- Users can update only first and last name
- New users must be within these validating constraints
    - min firstName length = 2
	- min lastName length  = 2
	- min password length  = 8


## `/auth`

### `/api/auth - POST` Creates a JWT token. Example:
```json
// Request body
{
  "email": "admin@admin.com",
  "password": "admin_admin"
}
// Response
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1hcmlvQGRyZXpuamFrLmNvbSIsImV4cGlyZXMiOjE2OTg0MTI3NTUsInVzZXJJRCI6IjY1M2I4MDdjNjQ3OWE1MjIwM2UzZDAxMiJ9.lpJXrWSotbClr1692xoI1L_FovaVrkeOC627BV6P9zc"
}
```

### Details

- Auth endpoint requires no authentication
