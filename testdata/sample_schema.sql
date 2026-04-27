-- Sample SQL file for testing clipper database capabilities

-- Create tables
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT TRUE
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    content TEXT,
    published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts(id),
    user_id INTEGER REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample data
INSERT INTO users (username, email, password_hash) VALUES
('alice', 'alice@example.com', 'hashed_password_1'),
('bob', 'bob@example.com', 'hashed_password_2'),
('charlie', 'charlie@example.com', 'hashed_password_3');

INSERT INTO posts (user_id, title, content, published) VALUES
(1, 'First Post', 'This is the content of the first post.', TRUE),
(1, 'Second Post', 'Content of the second post.', FALSE),
(2, 'Bob''s Post', 'Bob here with some content.', TRUE);

INSERT INTO comments (post_id, user_id, content) VALUES
(1, 2, 'Great post!'),
(1, 3, 'I agree with Bob.'),
(3, 1, 'Thanks for sharing!');

-- Sample queries
-- Get all published posts with author info
SELECT p.id, p.title, p.content, u.username, p.created_at
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.published = TRUE
ORDER BY p.created_at DESC;

-- Get comment count per post
SELECT p.title, COUNT(c.id) as comment_count
FROM posts p
LEFT JOIN comments c ON p.id = c.post_id
GROUP BY p.id, p.title
ORDER BY comment_count DESC;

-- Get users with their post count
SELECT u.username, COUNT(p.id) as post_count
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
GROUP BY u.id, u.username
ORDER BY post_count DESC;