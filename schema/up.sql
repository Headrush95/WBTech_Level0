CREATE TABLE IF NOT EXISTS delivery
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255),
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(255),
    email VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS orders
(
    order_uid  VARCHAR(19) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULl UNIQUE,
    entry VARCHAR(50) NOT NULL,
    delivery INT REFERENCES delivery(id) ON DELETE CASCADE NOT NULL,
--     payment INT REFERENCES transactions(transaction) ON DELETE CASCADE NOT NULL,--уйти от id и связать таблицы через order_uid=transaction
    locale VARCHAR(2) NOT NULL,
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(5) NOT NULL,
    sm_id INT,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(5) NOT NULL
);

CREATE TABLE IF NOT EXISTS items
(
    chrt_id INT NOT NULL,
    track_number VARCHAR(255) REFERENCES orders(track_number) ON DELETE CASCADE NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(25) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INT,
    size VARCHAR(5) NOT NULL,
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INT NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions
(
--     id SERIAL PRIMARY KEY, --уйти от id и связать таблицы через order_uid=transaction
    transaction VARCHAR(19) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    request_id VARCHAR(255) NOT NULL,
    currency VARCHAR(5) NOT NULL,
    provider varchar(255) NOT NULL,
    amount INT NOT NULL,
    payment_dt TIMESTAMP NOT NULL,
    bank VARCHAR(255) NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT
);

CREATE INDEX track_number_key ON items (track_number);



