BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    username      VARCHAR(40)  NOT NULL PRIMARY KEY,
    email         VARCHAR(64)  NOT NULL UNIQUE,
    is_active     BOOLEAN      NOT NULL DEFAULT FALSE,
    password_hash VARCHAR(512) NOT NULL,

    first_name    VARCHAR(32)  NOT NULL DEFAULT '',
    last_name     VARCHAR(32)  NOT NULL DEFAULT '',
    avatar_id     VARCHAR(32)  NULL     DEFAULT NULL
);

COMMIT;