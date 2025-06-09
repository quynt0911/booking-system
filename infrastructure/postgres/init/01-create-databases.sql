-- Tạo databases cho từng service
CREATE DATABASE user_service_db;
CREATE DATABASE expert_service_db; 
CREATE DATABASE booking_service_db;
CREATE DATABASE notification_service_db;
CREATE DATABASE review_service_db;
CREATE DATABASE worker_service_db;

-- Tạo users cho từng service
CREATE USER user_service WITH PASSWORD 'user_pass';
CREATE USER expert_service WITH PASSWORD 'expert_pass';
CREATE USER booking_service WITH PASSWORD 'booking_pass';
CREATE USER notification_service WITH PASSWORD 'notification_pass';
CREATE USER review_service WITH PASSWORD 'review_pass';
CREATE USER worker_service WITH PASSWORD 'worker_pass';

-- Phân quyền
GRANT ALL PRIVILEGES ON DATABASE user_service_db TO user_service;
GRANT ALL PRIVILEGES ON DATABASE expert_service_db TO expert_service;
GRANT ALL PRIVILEGES ON DATABASE booking_service_db TO booking_service;
GRANT ALL PRIVILEGES ON DATABASE notification_service_db TO notification_service;
GRANT ALL PRIVILEGES ON DATABASE review_service_db TO review_service;
GRANT ALL PRIVILEGES ON DATABASE worker_service_db TO worker_service;