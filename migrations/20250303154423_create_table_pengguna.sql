-- +goose Up
-- +goose StatementBegin
CREATE TABLE pengguna
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    username   VARCHAR(255) NOT NULL,
    password   TEXT         NOT NULL,
    first_name VARCHAR(100) NULL,
    last_name  VARCHAR(100) NULL,
    email      VARCHAR(255) NOT NULL,
    status     VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT drop table pengguna;
-- +goose StatementEnd
