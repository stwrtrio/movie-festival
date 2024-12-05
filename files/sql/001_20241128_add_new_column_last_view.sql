ALTER TABLE movie_festival.movie_views
ADD COLUMN last_viewed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER view_count;
