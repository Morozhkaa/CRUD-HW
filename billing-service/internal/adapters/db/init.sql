SELECT 'CREATE DATABASE balances' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'balances')\gexec
\c balances;

CREATE TABLE balances (
    username VARCHAR(255) PRIMARY KEY,
    balance INTEGER NOT NULL
);
