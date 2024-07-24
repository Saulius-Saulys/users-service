CREATE SCHEMA IF NOT EXISTS users AUTHORIZATION test;

CREATE TABLE users.users
(
    id         UUID,
    first_name VARCHAR(75)              NOT NULL,
    last_name  VARCHAR(75)              NOT NULL,
    nickname   VARCHAR(75)              NOT NULL,
    password   TEXT                     NOT NULL,
    email      VARCHAR(100)             NOT NULL,
    country    CHAR(3)                  NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    PRIMARY KEY (id)
);

-- Index for country created because it is a common filter in the application
CREATE INDEX idx_country
    ON users.users (country);