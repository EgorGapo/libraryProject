-- +goose Up
CREATE INDEX idx_name_book_hash ON book USING hash (name);


-- +goose Down
DROP INDEX idx_name_book_hash;