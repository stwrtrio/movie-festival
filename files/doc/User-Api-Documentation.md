# User API Documentation

## API List

| No  | API Name  | API Endpoint  | HTTP Method  |
|---|---|---|---|
|1.|User Register|/api/user/register|POST|
|2.|User Login|/api/user/login|POST|
|3.|User Logout|/api/user/logout|POST|

--- 

### 1. User Register API
#### API Endpoint:
```
http://localhost:8080/api/user/register
```
##### Description:
Allows a new user to register by providing a username and password.

##### Request:
- Method: `POST`
- URL: `/api/user/register`
- Body (JSON):
```
{
    "username": "user123",
    "password": "password123"
}
```
- Fields:
    - `username`: The username of the user trying to log in. (string)
        - Required
        - Must be a string
        - Maximum length: 50 characters

    - `password`: The password associated with the provided username. (string)
        - Required
        - Must be a string
        - Minimum length: 8 characters
        - Maximum length: 50 characters

#### Response:
##### Success Response (HTTP 200):
```
{
    "code": 201,
    "status": "success",
    "message": "User registered successfully"
}
```

- Fields:
    - code: The HTTP status code. (integer)
    - status: The status of the request. (string)
    - message: A message indicating the result of the request. (string)


##### Failure Response (HTTP 400):
```
{
    "code": 400,
    "status": "failed",
    "message": "username already exists"
}
```
- Fields:
    - code: The HTTP status code. (integer)
    - status: The status of the request. (string)
    - message: An error message indicating that the credentials were invalid. (string)

---

### 2. User Login API
#### API Endpoint:
```
http://localhost:8080/api/user/login
```
##### Description:
This endpoint allows a user to log in to the system by providing their username and password. If the credentials are valid, a JWT token is generated and returned for authenticated access to protected endpoints.

##### Request:
- Method: `POST`
- URL: `/api/user/login`
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

#### Response:
##### Success Response (HTTP 200):
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

##### Failure Response (HTTP 401):
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

##### Error Handling:
- 400 Bad Request: This error will be returned if the request body is malformed or missing required parameters.
- 401 Unauthorized: If the credentials (username or password) are incorrect, this error will be returned.

---

### 3. User Logout API
#### API Endpoint:
```
http://localhost:8080/api/user/logout
```
##### Description:
This endpoint allows a user to log out from the system by invalidating their JWT token. Upon logout, the token will be added to the blacklist in Redis to prevent further use.

##### Request:
- Method: `POST`
- URL: `/api/user/logout`
> There is no request body required for this endpoint. The authentication token should be passed via the Authorization header.

##### Request Header:
Authorization: Required. The value should be the Bearer token, which is the JWT token previously generated during login.
  Example:
```
Authorization: Bearer <your-jwt-token>
```
##### Success Response (HTTP 200):
```
{
    "code": 200,
    "status": "success",
    "message": "user logged out successfully"
}
```
##### Fields:
- code: HTTP status code of the response (200).
- status: The status of the request (success).
- message: A message confirming that the logout was successful.

##### Failure Response (HTTP 401):
```
{
    "code": 401,
    "status": "failed",
    "message": "Token has been revoked"
}
```
##### Fields:
- code: HTTP status code of the response (401).
- status: The status of the request (failed).
- message: A message indicating that the token is invalid or missing.











