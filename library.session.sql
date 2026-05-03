CREATE TABLE authors (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO authors (id, name)
VALUES ('11111111-1111-1111-1111-111111111111', 'Ivan Ivanov');

SELECT * FROM authors; 

