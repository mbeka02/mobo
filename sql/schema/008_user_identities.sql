-- +goose Up
CREATE TABLE IF NOT EXISTS user_identities (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider varchar(50) NOT NULL,
    provider_user_id varchar(256) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_user_identities_user_id ON user_identities(user_id);

-- +goose Down
DROP TABLE user_identities;
