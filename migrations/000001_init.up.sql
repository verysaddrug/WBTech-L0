CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) UNIQUE NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(255) NOT NULL,
    delivery_name VARCHAR(255) NOT NULL,
    delivery_phone VARCHAR(20) NOT NULL,
    delivery_zip VARCHAR(20) NOT NULL,
    delivery_city VARCHAR(255) NOT NULL,
    delivery_address VARCHAR(255) NOT NULL,
    delivery_region VARCHAR(255) NOT NULL,
    delivery_email VARCHAR(255) NOT NULL,
    payment_transaction VARCHAR(255) NOT NULL,
    payment_request_id VARCHAR(255) NOT NULL,
    payment_currency VARCHAR(5) NOT NULL,
    payment_provider VARCHAR(255) NOT NULL,
    payment_amount INT NOT NULL,
    payment_dt TIMESTAMP NOT NULL,
    payment_bank VARCHAR(255) NOT NULL,
    payment_delivery_cost INT NOT NULL,
    payment_goods_total INT NOT NULL,
    payment_custom_fee INT NOT NULL,
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255),
    shardkey VARCHAR(255) NOT NULL,
    sm_id INT,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(255) NOT NULL
);
CREATE TABLE order_items (
    item_id SERIAL PRIMARY KEY,
    order_id INT,
    chrt_id INT NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INT NOT NULL,
    size VARCHAR(255) NOT NULL,
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE
);