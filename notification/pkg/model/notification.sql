CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    chat_id VARCHAR(255),
    sender VARCHAR(255),
    receiver VARCHAR(255),
    msg TEXT,
    reference TEXT
);

CREATE TABLE user_chats (
    user_id VARCHAR(255) NOT NULL,
    chat_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (user_id, chat_id)
);
