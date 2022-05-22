createdb -p5433 wb

CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(50) UNIQUE,
    request_id VARCHAR(20),
    currency CHAR(3) NOT NULL,
    provider VARCHAR(20) NOT NULL,
    amount INT NOT NULL,
    payment_dt INT NOT NULL,
    bank VARCHAR(50) NOT NULL,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE item (
    id SERIAL PRIMARY KEY,
    chrt_id INT,
    track_number VARCHAR(20),
    price INT NOT NULL,
    rid VARCHAR(50),
    name VARCHAR(50) NOT NULL,
    sale INT,
    size VARCHAR(10),
    total_price INT,
    nm_id INT,
    brand VARCHAR(50),
    status INT
);

CREATE TABLE delivery (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(12) NOT NULL,
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(50) NOT NULL,
    address VARCHAR(50) NOT NULL,
    region VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL
);

CREATE TABLE order (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) UNIQUE NOT NULL,
    track_number VARCHAR(20) NOT NULL,
    entry VARCHAR(10),
    delivery_id INT NOT NULL,
    payment_id INT NOT NULL,
    items INT[],
    CONSTRAINT fk_delivery
        FOREIGN KEY(delivery_id) 
	    REFERENCES delivery(id),
    CONSTRAINT fk_payment,
        FOREIGN KEY(payment_id),
        REFERENCES payment(id),
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(20),
    customer_id VARCHAR(50) NOT NULL,
    delivery_service VARCHAR(20) NOT NULL,
    shardkey VARCHAR(10),
    sm_id INT,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard VARCHAR(10)
);

CREATE TABLE timetest (
    date TIMESTAMPTZ NOT NULL
);

INSERT INTO timetest (date) VALUES ('2021-11-26T06:22:19Z');

SELECT to_char(date, 'YYYY-MM-DD"T"HH:MI:SS TZ') FROM timetest;
SELECT to_char(date AT TIME ZONE 'z', 'YYYY-MM-DD"T"HH:MI:SSZ') FROM timetest;
