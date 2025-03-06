CREATE TABLE books
(
id serial PRIMARY key,
title VARCHAR (50) UNIQUE NOT NULL,
year INTEGER NOT NULL,
created_at time,
updated_at time,
user_id INTEGER 
);