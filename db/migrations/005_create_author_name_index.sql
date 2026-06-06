
-- +goose Up
CREATE INDEX idx_author_name_hash ON author USING hash (name);

-- +goose Down
DROP INDEX idx_author_name_hash;