set -e

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-consultation_db}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting database migration...${NC}"

# Wait for database to be ready
echo "Waiting for database to be ready..."
until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
  >&2 echo "Database is unavailable - sleeping"
  sleep 1
done

echo -e "${GREEN}Database is ready!${NC}"

# Run migrations for each service
services=("user-service" "expert-service" "booking-service" "notification-service")

for service in "${services[@]}"; do
    echo -e "${YELLOW}Running migrations for $service...${NC}"
    
    migration_dir="services/$service/migrations"
    if [ -d "$migration_dir" ]; then
        for migration_file in "$migration_dir"/*.sql; do
            if [ -f "$migration_file" ]; then
                echo "Applying migration: $(basename "$migration_file")"
                PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$migration_file"
            fi
        done
        echo -e "${GREEN}✓ $service migrations completed${NC}"
    else
        echo -e "${RED}⚠ Migration directory not found for $service${NC}"
    fi
done

echo -e "${GREEN}All migrations completed successfully!${NC}"