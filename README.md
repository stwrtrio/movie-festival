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
- Redis: Used for caching frequently accessed data.


## API Endpoints
### Admin APIs
- GET /api/admin/movies/most-viewed Retrieve most viewed movie.
- GET /api/admin/movies/most-viewed-genre Retrieve most viewed movie genre.
- POST /api/movies/:id/view To incrementing the view count for the movie
- POST /api/admin/movies Create a new movie.
- POST /api/admin/movies/:id Update an existing movie.

## Database Schema
### Tables
- movies: Stores movie details such as title, description, duration, watch URL, and view statistics.
- genres: Stores movie genres.
- artists: Stores movie artists (e.g., actors, directors).
- movie_genres: Junction table to associate movies with genres.
- movie_artists: Junction table to associate movies with artists.
- movie_views: Stores the view count for each movie.

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
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=movie_festival
```
4. Run the application:
```
go run main.go
```

## Testing
To run tests for the repository and controllers, you can use go test.
1. Run tests:
```
go test ./tests/*
```
2. Or you can run tests for a specific file:
```
go test ./tests/controllers/movie_controller_test.go
```

I've attached the Postman collection file to the ```files/doc``` folder for reference. This file contains the API endpoints that can be using for the test.