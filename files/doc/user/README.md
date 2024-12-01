# User API Documentation

### 1. User Login API
API Endpoint:
```
POST http://localhost:8080/api/user/login
```
### Description:
This endpoint allows a user to log in to the system by providing their username and password.  
If the credentials are valid, a JWT token is generated and returned for authenticated access to protected endpoints.


### Request:
- Method: POST
- URL: /api/user/login
- Body (JSON):
```
{
    "username": "user123",
    "password": "password123"
}
```
- Fields:
    - username: The username of the user trying to log in. (string)
    - password: The password associated with the provided username. (string)

### Response:
- Success Response (HTTP 200):
```
{
    "code": 200,
    "status": "success",
    "message": "Access granted",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNThhYWU5NzgtNDhlOS00YTZhLWJmOWUtN2VkY2E3MDhiYmJkIiwidXNlcm5hbWUiOiJ1c2VyMTIzIiwicm9sZSI6InVzZXIiLCJleHAiOjE3MzMxMjc4Njd9.BcJt29ongfb3bUtObzZCbTtnxoJqDjOXL21aYK15ths"
    }
}
```
- Fields:
    - code: The HTTP status code. (integer)
    - status: The status of the request. (string)
    - message: A message indicating the result of the request. (string)
    - data.token: The JWT token generated for the user to access protected resources. (string)

- Failure Response (HTTP 401):
```
{
    "code": 401,
    "status": "failed",
    "message": "invalid credentials"
}
```
- Fields:
    - code: The HTTP status code. (integer)
    - status: The status of the request. (string)
    - message: An error message indicating that the credentials were invalid. (string)

