CREATE TABLE users (
    chat_id INTEGER primary key, -- telegram chatID.
    user_id INTEGER NOT NULL, -- userID инициатора взаимодействия.
    username VARCHAR(255), -- username инициатора взаимодействия.
    status VARCHAR(30) NOT NULL, -- статус диалога с пользователем.
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- created_at
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP -- updated_at
);

CREATE TABLE details (
        id INTEGER primary key AUTOINCREMENT, -- telegram chatID.
        chat_id INTEGER NOT NULL, -- chat_id
        name         VARCHAR(255), -- имя детали
        description TEXT, -- описание детали
        soft_deadline_at DATETIME, -- soft_deadline_at
        hard_deadline_at DATETIME, -- hard_deadline_at
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- created_at
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- updated_at

        FOREIGN KEY (chat_id)  REFERENCES users (chat_id)
);