CREATE TABLE users
(
user_id serial PRIMARY KEY,
password text NOT NULL,
email text UNIQUE NOT NULL
);

