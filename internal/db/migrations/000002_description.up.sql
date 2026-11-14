CREATE TABLE authors (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    biography TEXT NOT NULL,
    birth INTEGER NOT NULL,
    nationality TEXT NOT NULL
);

CREATE UNIQUE INDEX uk_authors_name ON authors (name);