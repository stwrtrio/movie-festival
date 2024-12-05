CREATE TABLE IF NOT EXISTS movie_festival.movie_artists (
    movie_id VARCHAR(50) NOT NULL,
    artist_id VARCHAR(50) NOT NULL,
    PRIMARY KEY (movie_id, artist_id),
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE
);

ALTER TABLE movie_festival.movie_artists ADD UNIQUE KEY unique_movie_artist (movie_id, artist_id);