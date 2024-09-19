CREATE TABLE refresh_tokens (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    refresh_hash VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);