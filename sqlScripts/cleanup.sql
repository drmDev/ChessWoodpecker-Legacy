DROP TABLE IF EXISTS user_progress;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS puzzle_progress (
    progress_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    puzzle_id INTEGER REFERENCES puzzles(puzzle_id),
    status VARCHAR(10) NOT NULL CHECK (status IN ('success', 'failed')),
    attempted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, puzzle_id)
);