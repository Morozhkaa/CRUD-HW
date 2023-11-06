SELECT 'CREATE DATABASE users' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'users')\gexec
\c users;

CREATE TABLE users (
    username VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(255)
);
