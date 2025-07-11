echo "🚀 Starting development environment..."

# Dừng và xóa containers cũ nếu có
echo "🧹 Cleaning up old containers..."
docker-compose -f docker-compose.dev.yml down --remove-orphans

# Build và start services
echo "🔨 Building and starting services..."
docker-compose -f docker-compose.dev.yml up --build -d

# Chờ services khởi động
echo "⏳ Waiting for services to start..."
sleep 10

# Kiểm tra trạng thái services
echo "📊 Service status:"
docker-compose -f docker-compose.dev.yml ps

# Hiển thị logs
echo "📜 Starting to follow logs (Ctrl+C to stop following)..."
docker-compose -f docker-compose.dev.yml logs -f