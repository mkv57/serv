CREATE TABLE session
(
    id_session serial primary key,
    user_id INTEGER references users not null,
    token TEXT not null,
    ip TEXT not null,
    user_agent TEXT not null,
    created_at timestamp not null default now()
);