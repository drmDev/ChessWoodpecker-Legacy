CREATE TABLE IF NOT EXISTS puzzles (
    puzzle_id SERIAL PRIMARY KEY,
    category TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL
);