CREATE TABLE rating_scores
(
    id         SERIAL PRIMARY KEY,
    athlete_id    integer   NOT NULL,
    score      numeric(10,2) NOT NULL,
    activity_id integer   NOT NULL,
    distance    numeric(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
CREATE INDEX idx_rating_scores_athlete_id ON rating_scores (athlete_id);
