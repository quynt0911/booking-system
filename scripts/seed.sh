set -e

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-booking_system}
DB_USER=${DB_USER:-directus}
DB_PASSWORD=${DB_PASSWORD:-directus}

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Starting database seeding...${NC}"

# Run seed data
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" << 'EOF'

-- Seed Users
INSERT INTO users (id, email, password_hash, fullname, role, phone, email_verified, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'admin@teknix.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Admin System', 'admin', '+84901234567', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002', 'expert1@teknix.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Nguyễn Văn A', 'expert', '+84901234568', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003', 'expert2@teknix.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Trần Thị B', 'expert', '+84901234569', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440004', 'user1@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Lê Văn C', 'user', '+84901234570', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440005', 'user2@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Phạm Thị D', 'user', '+84901234571', true, NOW(), NOW());

-- Seed Experts
INSERT INTO experts (id, user_id, specialization, bio, hourly_rate, years_experience, rating, total_reviews, is_available, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440002', 'Technology Consulting', 'Chuyên gia công nghệ với 10 năm kinh nghiệm', 500000, 10, 4.8, 45, true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440003', 'Business Strategy', 'Chuyên gia chiến lược kinh doanh', 750000, 8, 4.9, 32, true, NOW(), NOW());

-- Seed Schedules
INSERT INTO schedules (id, expert_id, day_of_week, start_time, end_time, is_recurring, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440010', 1, '09:00:00', '17:00:00', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440021', '550e8400-e29b-41d4-a716-446655440010', 2, '09:00:00', '17:00:00', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440022', '550e8400-e29b-41d4-a716-446655440011', 1, '10:00:00', '18:00:00', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440023', '550e8400-e29b-41d4-a716-446655440011', 3, '10:00:00', '18:00:00', true, NOW(), NOW());

-- Seed some bookings
INSERT INTO bookings (id, user_id, expert_id, booking_date, start_time, end_time, consultation_type, notes, status, total_amount, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440010', CURRENT_DATE + INTERVAL '1 day', '10:00:00', '11:00:00', 'online', 'Tư vấn về hệ thống ERP', 'confirmed', 500000, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440031', '550e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440011', CURRENT_DATE + INTERVAL '2 days', '14:00:00', '15:00:00', 'offline', 'Tư vấn chiến lược marketing', 'pending', 750000, NOW(), NOW());

EOF

echo -e "${GREEN}✓ Seed data inserted successfully!${NC}"  