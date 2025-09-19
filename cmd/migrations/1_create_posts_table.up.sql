CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       content TEXT NOT NULL,
                       comments_are_allowed BOOLEAN NOT NULL DEFAULT true,
                       post_author TEXT NOT NULL
);
