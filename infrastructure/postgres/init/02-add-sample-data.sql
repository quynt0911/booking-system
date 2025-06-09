-- Add availability column to experts table if not exists
ALTER TABLE experts ADD COLUMN IF NOT EXISTS availability JSONB;

Insert sample users
INSERT INTO users (id, email, password_hash, fullname, phone, role, email_verified) VALUES
    ('11111111-1111-1111-1111-111111111111', 'user1@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'John Doe', '+1234567890', 'user', true),
    ('22222222-2222-2222-2222-222222222222', 'user2@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'Jane Smith', '+1234567891', 'user', true),
    ('33333333-3333-3333-3333-333333333333', 'expert1@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'Dr. Alice Johnson', '+1234567892', 'expert', true),
    ('44444444-4444-4444-4444-444444444444', 'expert2@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'Dr. Bob Wilson', '+1234567893', 'expert', true);

-- Insert sample experts
INSERT INTO experts (id, user_id, specialization, experience_years, hourly_rate, is_available, rating, total_reviews) VALUES
    ('55555555-5555-5555-5555-555555555555', '33333333-3333-3333-3333-333333333333', 'General Medicine', 10, 100.00, true, 4.8, 25),
    ('66666666-6666-6666-6666-666666666666', '44444444-4444-4444-4444-444444444444', 'Pediatrics', 8, 120.00, true, 4.9, 30);

-- Insert sample expert schedules
INSERT INTO expert_schedules (expert_id, day_of_week, start_time, end_time) VALUES
    ('55555555-5555-5555-5555-555555555555', 1, '09:00', '17:00'),
    ('55555555-5555-5555-5555-555555555555', 2, '09:00', '17:00'),
    ('55555555-5555-5555-5555-555555555555', 3, '09:00', '17:00'),
    ('55555555-5555-5555-5555-555555555555', 4, '09:00', '17:00'),
    ('55555555-5555-5555-5555-555555555555', 5, '09:00', '17:00'),
    ('66666666-6666-6666-6666-666666666666', 1, '10:00', '18:00'),
    ('66666666-6666-6666-6666-666666666666', 2, '10:00', '18:00'),
    ('66666666-6666-6666-6666-666666666666', 3, '10:00', '18:00'),
    ('66666666-6666-6666-6666-666666666666', 4, '10:00', '18:00'),
    ('66666666-6666-6666-6666-666666666666', 5, '10:00', '18:00');

-- Insert sample bookings
INSERT INTO bookings (id, user_id, expert_id, scheduled_datetime, duration_minutes, meeting_type, meeting_url, meeting_address, status, price, notes) VALUES
    ('77777777-7777-7777-7777-777777777777', '11111111-1111-1111-1111-111111111111', '55555555-5555-5555-5555-555555555555', 
    CURRENT_TIMESTAMP + INTERVAL '1 day', 60, 'online', 'https://meet.google.com/abc-defg-hij', NULL, 'pending', 100.00, 'First consultation'),
    
    ('88888888-8888-8888-8888-888888888888', '22222222-2222-2222-2222-222222222222', '66666666-6666-6666-6666-666666666666',
    CURRENT_TIMESTAMP + INTERVAL '2 days', 45, 'offline', NULL, '123 Medical Center, Room 456', 'confirmed', 90.00, 'Follow-up appointment'),
    
    ('99999999-9999-9999-9999-999999999999', '11111111-1111-1111-1111-111111111111', '55555555-5555-5555-5555-555555555555',
    CURRENT_TIMESTAMP + INTERVAL '3 days', 30, 'online', 'https://meet.google.com/xyz-uvw-rst', NULL, 'completed', 50.00, 'Quick check-up');

-- Insert sample booking status history
INSERT INTO booking_status_history (booking_id, old_status, new_status, changed_by, reason) VALUES
    ('77777777-7777-7777-7777-777777777777', NULL, 'pending', '11111111-1111-1111-1111-111111111111', 'Booking created'),
    ('88888888-8888-8888-8888-888888888888', NULL, 'pending', '22222222-2222-2222-2222-222222222222', 'Booking created'),
    ('88888888-8888-8888-8888-888888888888', 'pending', 'confirmed', '44444444-4444-4444-4444-444444444444', 'Booking confirmed by expert'),
    ('99999999-9999-9999-9999-999999999999', NULL, 'pending', '11111111-1111-1111-1111-111111111111', 'Booking created'),
    ('99999999-9999-9999-9999-999999999999', 'pending', 'confirmed', '33333333-3333-3333-3333-333333333333', 'Booking confirmed by expert'),
    ('99999999-9999-9999-9999-999999999999', 'confirmed', 'completed', '33333333-3333-3333-3333-333333333333', 'Consultation completed');

-- Insert sample notifications
INSERT INTO notifications (user_id, booking_id, title, message, type, is_read) VALUES
    ('11111111-1111-1111-1111-111111111111', '77777777-7777-7777-7777-777777777777', 'New Booking Created', 'Your booking has been created successfully', 'booking_created', false),
    ('22222222-2222-2222-2222-222222222222', '88888888-8888-8888-8888-888888888888', 'Booking Confirmed', 'Your booking has been confirmed by the expert', 'booking_confirmed', false);

-- Insert sample notification settings
INSERT INTO notification_settings (user_id, email_enabled, reminder_minutes) VALUES
    ('11111111-1111-1111-1111-111111111111', true, 60),
    ('22222222-2222-2222-2222-222222222222', true, 30),
    ('33333333-3333-3333-3333-333333333333', true, 60),
    ('44444444-4444-4444-4444-444444444444', true, 60);

-- Insert sample reviews
INSERT INTO reviews (booking_id, user_id, expert_id, rating, comment) VALUES
    ('99999999-9999-9999-9999-999999999999', '11111111-1111-1111-1111-111111111111', '55555555-5555-5555-5555-555555555555', 5, 'Great consultation, very professional');

-- Insert sample system settings
INSERT INTO system_settings (key, value, description, type, is_public) VALUES
    ('booking_min_duration', '30', 'Minimum booking duration in minutes', 'number', true),
    ('booking_max_duration', '120', 'Maximum booking duration in minutes', 'number', true),
    ('cancellation_policy_hours', '24', 'Hours before booking that cancellation is allowed', 'number', true); 