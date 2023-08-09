CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    password VARCHAR(100) NOT NULL,
    blocked BOOLEAN NOT NULL DEFAULT false,
    registration_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
