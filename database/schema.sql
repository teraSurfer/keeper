CREATE TABLE IF NOT EXISTS todos (
    id INTEGER PRIMARY KEY,
    title text NOT NULL,
    description text,
    completed boolean
);