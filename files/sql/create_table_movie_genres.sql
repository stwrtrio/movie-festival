CREATE TABLE IF NOT EXISTS movie_festival.movie_genres (
    movie_id VARCHAR(50) NOT NULL,
    genre_id VARCHAR(50) NOT NULL,
    PRIMARY KEY (movie_id, genre_id),
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
);

ALTER TABLE movie_festival.movie_genres ADD UNIQUE KEY unique_movie_genre (movie_id, genre_id);