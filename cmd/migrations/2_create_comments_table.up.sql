CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          content TEXT NOT NULL,
                          post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
                          parent_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
                          comment_author TEXT NOT NULL
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
