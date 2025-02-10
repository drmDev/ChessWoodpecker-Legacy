CREATE TABLE user_sessions (
    session_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    current_puzzle_index INT NOT NULL,
    total_time BIGINT NOT NULL,
    puzzle_times JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);