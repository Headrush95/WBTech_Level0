CREATE TABLE IF NOT EXISTS orders
(
    order_uid  VARCHAR(19) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    delivery INT REFERENCES delivery(id) ON DELETE CASCADE NOT NULL,
    payment INT REFERENCES transactions(id) ON DELETE CASCADE NOT NULL,
    items INT[] REFERENCES items(id) ON DELETE CASCADE NOT NULL, --?????
    locale VARCHAR(2) NOT NULL ,
    internal_signature VARCHAR(255),
    sharedkey VARCHAR(5) NOT NULL,
    sm_id INT,
    date_created TIMESTAMP NOT NULL,
    off_shard VARCHAR(5) NOT NULL
);

CREATE TABLE IF NOT EXISTS items
(
    id SERIAL PRIMARY KEY,
    chrt_id INT NOT NULL,
    track_number VARCHAR(255) NOT NULL ,
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

CREATE TABLE IF NOT EXISTS delivery
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255),
    zip VARCHAR(6),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS transactions
(
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(5),
    provider varchar(255),
    amount INT,
    payment_dt TIMESTAMP,
    bank VARCHAR(255),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);