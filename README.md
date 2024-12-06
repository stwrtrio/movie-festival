# Movie-Festival

The Movie Festival API is a backend service for managing movies, genres, artists, and viewing statistics for a movie festival application. It provides endpoints to manage movies, retrieve the most viewed movie and genre, and perform various CRUD operations.

### Features
- Movies: Create, update, delete, and retrieve movie details including title, description, genres, artists, and viewing statistics.
- Genres: Manage movie genres and associate them with movies.
- Artists: Manage artists involved in movies and associate them with movies.
- Most Viewed: Retrieve the most viewed movie and genre based on view statistics.

### Technologies Used
- Golang: The API is built using the Go programming language.
- Echo: A fast and flexible web framework for building APIs.
- MySQL: The API uses MySQL to store movie, genre, artist, and view data.
- SQL Mock: Used for unit testing repository functions with mocked database calls.
- Testify: Testing framework to assert expected behavior during tests.
- GoMock: For mocking the repository in unit tests.
- Redis: Used for caching frequently accessed data.
- JWT: Tokens are securely signed using a secret key to prevent tampering.
- Swagger: Documentation and testing API

## API Endpoints
#### Api List Documentation
- [Admin API Endpoints](./files/doc/Admin-Api-Documentation.md#admin-api-documentation)
- [User API Endpoints](./files/doc/User-Api-Documentation.md#user-api-documentation)

#### Testing
To test the endpoint, you can open swagger with this link: http://localhost:8080/swagger/index.html
Follong the instruction in the API List Documentation.

## Database Schema
### Tables
- movies: Stores movie details such as title, description, duration, watch URL, and view statistics.
- genres: Stores movie genres.
- artists: Stores movie artists (e.g., actors, directors).
- movie_genres: Junction table to associate movies with genres.
- movie_artists: Junction table to associate movies with artists.
- movie_views: Stores the view count for each movie.
- votes: Stores the movie voted by user.

For table structures files is included in directory ``files/sql``

# Getting Started
### Prerequisites
- Go: Install Go (version 1.16+).
- MySQL: Set up MySQL and create a database for the project.
- Postman or cURL: Use Postman or cURL to test the API endpoints.

### installation
1. Clone the repository:
```
git clone https://github.com/stwrtrio/movie-festival.git
cd movie-festival
```

2. Install Dependencies:
```
go mod tidy
```

3. Set up the MySQL database. Create a .env file and configure your database connection:
```
# Database Connection
DB_USER=root
DB_PASSWORD=
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=movie_festival

# Redis Connection
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=

# App 
SERVER_PORT=8080
CACHE_DEFAULT_EXPIRATION=3600s

#JWT
JWT_SECRET=replace_this
JWT_EXPIRY=24h
```
4. Run the application:
```
make run
```

### Steps Performed by make run
1. Linting:
    - Runs syntax and formatting checks using go vet and go fmt.
    - Ensures the code adheres to Go best practices and standards.
2. Testing:
    - Executes all unit and integration tests in the project using go test.
    - Confirms that all tests pass without errors.
3. Swagger:
    - Generate Documentation and testing for API
4. Running the Application:
    - If both linting and testing succeed, the application is started using go run.

### Expected Behavior
- On Success: The application is started after the linting and testing steps pass without errors.
- On Failure:
    - If linting fails, the process stops, and an error message is displayed.
    - If any tests fail, the process stops, and test failure details are displayed.


## Testing
To run tests for the repository and controllers, you can use go test.
1. Run tests:
```
make test
```
2. Or you can run tests for a specific file:
```
go test ./tests/controllers/movie_controller_test.go
```

I've attached the Postman collection file to the ```files/doc``` folder for reference. This file contains the API endpoints that can be using for the test. Just import that file to your postman