-- +goose Up
CREATE TABLE friends (
    user_id UUID,
    friend_id UUID,

    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL,

    CONSTRAINT user_not_friend CHECK(user_id <> friend_id),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE,

    PRIMARY KEY (user_id, friend_id)
);

-- +goose Down
DROP TABLE friends;