CREATE TABLE notifications (
    id VARCHAR(255) SERIAL PRIMARY KEY,
    chat_id VARCHAR(255),
    sender VARCHAR(255),
    receiver VARCHAR(255),
    msg TEXT
);
