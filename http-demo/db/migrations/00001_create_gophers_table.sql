-- +goose Up
CREATE TABLE gophers (
    id   INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    job  VARCHAR(255)
);

-- +goose Down
DROP TABLE IF EXISTS gophers;
