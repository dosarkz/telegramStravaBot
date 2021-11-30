CREATE TABLE workouts
(
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     text NULL,
    status_id       int default 0,
    created_at      TIMESTAMP           NOT NULL,
    updated_at      TIMESTAMP           NOT NULL,
    deleted_at      TIMESTAMP NULL
);

