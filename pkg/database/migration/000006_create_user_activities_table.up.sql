CREATE TABLE user_activities (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INT NOT NULL,
    book_id INT NOT NULL,
    activity_type VARCHAR(50) CHECK (activity_type IN ('borrowed', 'returned', 'search')) NOT NULL DEFAULT 'search',
    activity_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);
