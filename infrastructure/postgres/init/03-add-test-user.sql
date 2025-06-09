-- -- Delete existing test user if exists
-- DELETE FROM users WHERE email = 'user1@example.com';

-- Insert test user with properly hashed password (password123)
INSERT INTO users (id, email, password_hash, fullname, phone, role, email_verified) VALUES
    ('11111111-1111-1111-1111-111111111111', 'user1@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'John Doe', '+1234567890', 'user', true);

-- Insert another test user with the same password
INSERT INTO users (id, email, password_hash, fullname, phone, role, email_verified) VALUES
    ('22222222-2222-2222-2222-222222222222', 'user2@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Jane Doe', '+1234567891', 'user', true); 