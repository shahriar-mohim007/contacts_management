CREATE TABLE users (
id SERIAL PRIMARY KEY,                 -- Auto-incrementing ID for each user
name VARCHAR(100) NOT NULL,            -- User's full name
email VARCHAR(255) NOT NULL UNIQUE,    -- User's email (must be unique)
password VARCHAR(255) NOT NULL,   -- Hashed password for security
is_active BOOLEAN DEFAULT FALSE,       -- Indicates if the user is activated
activation_token VARCHAR(255),         -- Token for activating the account
created_at TIMESTAMP DEFAULT NOW(),    -- Timestamp of when the user was created
updated_at TIMESTAMP DEFAULT NOW()     -- Timestamp of the last update
);