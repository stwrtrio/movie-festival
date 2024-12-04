# Admin API Documentation

## API List

| No  | API Name  | API Endpoint  | HTTP Method  |
|---|---|---|---|
|1.|Create a new movie|/api/admin/movie|POST|
|2.|Update an existing movie|/api/admin/movies/:id|POST|
|3.|Retrieve most viewed movie|/api/admin/movies/most-viewed|GET|
|4.|Retrieve most viewed movie genre|/api/admin/movies/most-viewed-genres|GET|

--- 

### 1. Create Movie
#### API Endpoint:
```
http://localhost:8080/api/admin/movies
```
##### Description:
To create movie

##### Request:
- Method: `POST`
- URL: `/api/admin/movie`
- Body (JSON):
```
{
    "title": "Inception",
    "description": "A mind-bending thriller",
    "duration": 148,
    "genres": [
        "Sci-Fi",
        "Thriller", 
        "Action"
    ],
    "watch_url": "http://example.com/inception.mp4",
    "artists": [
        "Leonardo DiCaprio",
        "Joseph Gordon-Levitt",
        "Elliot Page"
    ]
}
```
- Fields:
    - `title`: The title of the movie. (string)
        - Required
        - Must be a string
        - Maximum length: 150 characters
    - `description`: The description of the movie. (string)
        - Required
        - Must be a string
    - `duration`: Duration of the movie. (integer)
        - Required
        - Must be a string
    - `genres`: Genre of the movie. (array,string)
        - Required
        - Must be a array of string
    - `watch_url`: The URL of the movie. (string)
        - Required
        - Must be a string
    - `artists`: The artist of the movie. (array,string)
        - Required
        - Must be a array of string

#### Response:
##### Success Response (HTTP 201):
```
{
    "code": 201,
    "status": "success",
    "message": "Movie created successfully"
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
    "message": "Key: 'CreateMovieRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"
}
```
- Fields:
    - code: The HTTP status code. (integer)
    - status: The status of the request. (string)
    - message: An error message indicating that the request were invalid. (string)

---

### 2. Update Existing Movie API
#### API Endpoint:
```
http://localhost:8080/api/admin/movies/:id
```
##### Description:
This endpoint allows admin to update existing movie details.

##### Request:
- Method: `POST`
- URL: `/api/admin/movies/:id`
- Body (JSON):
```
{
    "title": "Inception 2",
    "description": "A mind-bending thriller update",
    "duration": 150,
    "genres": [
        "Sci-Fi",
        "Thriller"
    ],
    "watch_url": "http://example.com/inception.mp4",
    "artists": [
        "Leonardo DiCaprio",
        "Joseph Gordon-Levitt",
        "Elliot Page"
    ]
}
```
- Fields:
    - `title`: The title of the movie. (string)
        - Required
        - Must be a string
        - Maximum length: 150 characters
    - `description`: The description of the movie. (string)
        - Required
        - Must be a string
    - `duration`: Duration of the movie. (integer)
        - Required
        - Must be a string
    - `genres`: Genre of the movie. (array,string)
        - Required
        - Must be a array of string
    - `watch_url`: The URL of the movie. (string)
        - Required
        - Must be a string
    - `artists`: The artist of the movie. (array,string)
        - Required
        - Must be a array of string

#### Response:
##### Success Response (HTTP 201):
```
{
    "code": 201,
    "status": "success",
    "message": "Movie updated successfully"
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
    "message": "movie is not exists"
}
```
- Fields:
    - code: The HTTP status code. (integer)
    - status: The status of the request. (string)
    - message: An error message indicating that the movie is not exists. (string)

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











