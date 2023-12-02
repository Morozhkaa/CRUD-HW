SELECT 'CREATE DATABASE orders' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'orders')\gexec
\c orders;

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price INTEGER NOT NULL,
    total_cost INTEGER NOT NULL,
    status VARCHAR(255) NOT NULL
);