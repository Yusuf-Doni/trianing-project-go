-- Create users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create sessions table for session management
CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR(64) PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);

-- Insert default admin user (password: admin123)
INSERT INTO users (username, password, email, role) VALUES 
('admin', 'admin123', 'admin@book.com', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Insert sample users for testing
INSERT INTO users (username, password, email, role) VALUES 
('user1', 'password123', 'user1@example.com', 'user'),
('manager', 'manager123', 'manager@book.com', 'manager')
ON CONFLICT (username) DO NOTHING;
