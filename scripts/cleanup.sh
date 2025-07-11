echo "🧹 Cleaning up Docker resources..."

# Dừng tất cả containers
docker-compose -f docker-compose.dev.yml down --remove-orphans

# Xóa images không sử dụng
echo "🗑️ Removing unused images..."
docker image prune -f

# Xóa build cache
echo "🗑️ Removing build cache..."
docker builder prune -f

# Hiển thị disk usage
echo "💾 Current Docker disk usage:"
docker system df

echo "✅ Cleanup completed"