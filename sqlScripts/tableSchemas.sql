CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL
);

CREATE TABLE puzzles (
    puzzle_id SERIAL PRIMARY KEY,
    category TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL
);

CREATE TABLE user_progress (
    progress_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    puzzle_id INT REFERENCES puzzles(puzzle_id) ON DELETE CASCADE,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, puzzle_id)  -- Ensures a user can complete a puzzle only once
);