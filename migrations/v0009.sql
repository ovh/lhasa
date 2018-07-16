-- +migrate Down

UPDATE releases SET badge_ratings = NULL;

-- +migrate Up
