-- +goose Up
CREATE TABLE gophers (
    id   INTEGER NOT NULL PRIMARY KEY /*!40101 AUTO_INCREMENT */,
    name VARCHAR(255) NOT NULL,
    job  VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS gophers;
