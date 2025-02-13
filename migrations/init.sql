-- Создание таблицы пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    balance INT NOT NULL DEFAULT 1000
);

-- Создание таблицы истории переводов монет
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    sender_id INT REFERENCES users(id),
    receiver_id INT REFERENCES users(id),
    amount INT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы покупок
CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    item_name TEXT NOT NULL, 
    quantity INT NOT NULL DEFAULT 1,
    price INT NOT NULL,
);
