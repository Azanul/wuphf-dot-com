CREATE TABLE user_chats (
    user_id VARCHAR(255) NOT NULL,
    chat_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (user_id, chat_id)
);
