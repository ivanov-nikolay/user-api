CREATE TABLE IF NOT EXISTS "users"
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    gender VARCHAR(10) NOT NULL,
    status VARCHAR(50) NOT NULL,
    birthday TIMESTAMP,
    join_date TIMESTAMP NOT NULL
) WITH (oids = false);
