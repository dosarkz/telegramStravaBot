CREATE TABLE workout_users
(
    id         SERIAL PRIMARY KEY,
    user_id    integer   NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
CREATE INDEX idx_workout_users_user_id ON workout_users (user_id);
