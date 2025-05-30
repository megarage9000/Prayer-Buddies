-- +goose Up
CREATE TABLE prayers (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    sender UUID NOT NULL,
    receiver UUID NOT NULL,
    prayer TEXT NOT NULL,

    -- To ensure that sender and receiver are not equal
    CONSTRAINT sender_not_receiver CHECK (sender <> receiver),

    FOREIGN KEY (sender) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE prayers;