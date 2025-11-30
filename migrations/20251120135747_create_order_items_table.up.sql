CREATE TABLE IF NOT EXISTS order_items
(
    id          VARCHAR(50) PRIMARY KEY,
    order_id        VARCHAR(50)   NOT NULL,
    product_id VARCHAR(50) NOT NULL,
    price       DECIMAL(10, 2) NOT NULL CHECK (price > 0),
    quantity       INTEGER        NOT NULL,
    created_at  TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP
    );