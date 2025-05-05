DELETE FROM schema_migrations WHERE version = 1;
INSERT INTO schema_migrations (version, dirty) VALUES (1, false);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL
)