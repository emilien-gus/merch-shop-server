CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    balance INT NOT NULL DEFAULT 1000,  -- Начальный баланс 1000 монет
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount INT NOT NULL CHECK (amount > 0),  -- Сумма перевода должна быть положительной
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id INT NOT NULL REFERENCES merch_items(id) ON DELETE CASCADE,
    amount INT NOT NULL CHECK (amount > 0), -- Количество купленных товаров
    total_price INT NOT NULL, -- Итоговая сумма
    created_at TIMESTAMP DEFAULT NOW()
);
