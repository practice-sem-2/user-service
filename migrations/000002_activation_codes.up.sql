BEGIN;

CREATE TABLE users_activation_codes
(
    username VARCHAR(40) REFERENCES users ON DELETE CASCADE,
    code     VARCHAR(64) NOT NULL
);

END;