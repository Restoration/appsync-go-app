CREATE DATABASE chatdb;

USE chatdb;

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL
);
