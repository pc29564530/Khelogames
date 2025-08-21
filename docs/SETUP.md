# Khelogames Setup Guide

This guide provides step-by-step instructions for setting up the Khelogames backend API on your local machine and production environment.

## 📋 Prerequisites

### System Requirements
- **Go**: 1.21.2 or higher
- **PostgreSQL**: 12 or higher
- **RabbitMQ**: 3.8 or higher
- **Redis**: 6.0 or higher (optional, for caching)
- **Git**: Latest version

### Operating System Support
- macOS (Intel/Apple Silicon)
- Linux (Ubuntu 20.04+, CentOS 8+)
- Windows 10/11 (with WSL2 recommended)

## 🛠️ Local Development Setup

### 1. Install Dependencies

#### macOS
```bash
# Install Go
brew install go

# Install PostgreSQL
brew install postgresql
brew services start postgresql

# Install RabbitMQ
brew install rabbitmq
brew services start rabbitmq

# Install Redis (optional)
brew install redis
brew services start redis
```

#### Ubuntu/Debian
```bash
# Install Go
wget https://go.dev/dl/go1.21.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql

# Install RabbitMQ
sudo apt install rabbitmq-server
sudo systemctl start rabbitmq-server

# Install Redis (optional)
sudo apt install redis-server
sudo systemctl start redis-server
```

#### Windows (WSL2)
```bash
# Follow Ubuntu/Debian instructions above
```

### 2. Clone and Setup Project

```bash
# Clone repository
git clone https://github.com/your-org/khelogames.git
cd khelogames

# Install Go dependencies
go mod download
go mod tidy
```

### 3. Database Setup

#### Create Database
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE khelogames;

# Create user (optional)
CREATE USER khelogames WITH PASSWORD 'yourpassword';
GRANT ALL PRIVILEGES ON DATABASE khelogames TO khelogames;

# Exit
\q
```

#### Run Database Migrations
```bash
# If you have migration files
go run cmd/migrate/main.go up

# Or use the provided SQL files
psql -U postgres -d khelogames -f database/schema.sql
```

### 4. Environment Configuration

Create `app.env` file in the project root:
```bash
# Database Configuration
DB_DRIVER=postgres
DB_SOURCE=postgresql://username:password@localhost:5432/khelogames?sslmode=disable

# Server Configuration
SERVER_ADDRESS=0.0.0.0:8080
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=168h

# JWT Configuration
TOKEN_SYMMETRIC_KEY=your-32-character-secret-key-here

# Twilio Configuration (for SMS)
ACCOUNT_SID=your-twilio-account-sid
AUTH_TOKEN=your-twilio-auth-token
VERIFY_SERVICE_SID=your-verify-service-sid
YOUR_TWILIO_PHONE_NUMBER=+1234567890

# Google OAuth Configuration
CLIENT_ID=your-google-client-id
CLIENT_SECRET=your-google-client-secret

# RabbitMQ Configuration
RABBIT_SOURCE=amqp://guest:guest@localhost:5672/

# Redis Configuration (optional)
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 5. Start Services

#### Start PostgreSQL
```bash
# macOS
brew services start postgresql

# Ubuntu/Debian
sudo systemctl start postgresql

# Windows
net start postgresql
```

#### Start RabbitMQ
```bash
# macOS
brew services start rabbitmq

# Ubuntu/Debian
sudo systemctl start rabbitmq-server

# Enable management plugin
sudo rabbitmq-plugins enable rabbitmq_management
# Access at http://localhost:15672 (guest/guest)
```

#### Start Redis (optional)
```bash
# macOS
brew services start redis

# Ubuntu/Debian
sudo systemctl start redis-server
```

### 6. Run the Application

#### Development Mode
```bash
# Run with hot reload
go run main.go

# Or with air for hot reload
air
```

#### Production Mode
```bash
# Build the application
go build -o khelogames main.go

# Run the binary
./khelogames
```

## 🐳 Docker Setup

### 1. Using Docker Compose (Recommended)

Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: khelogames
      POSTGRES_USER: khelogames
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  rabbitmq:
    image: rabbitmq:3-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  khelogames:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - rabbitmq
      - redis
    environment:
      DB_DRIVER: postgres
      DB_SOURCE: postgresql://khelogames:password@postgres:5432/khelogames?sslmode=disable
      RABBIT_SOURCE: amqp://guest:guest@rabbitmq:5672/
      REDIS_ADDRESS: redis:6379
    volumes:
      - ./uploads:/app/uploads

volumes:
  postgres_data:
```

### 2. Run with Docker Compose
```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 🚀 Production Deployment

### 1. Environment Setup

#### Server Requirements
- **OS**: Ubuntu 20.04 LTS or CentOS 8
- **CPU**: 4+ cores
- **RAM**: 8GB+
- **Storage**: 100GB+ SSD
- **Network**: 1Gbps+

### 2. Security Setup

#### Firewall Configuration
```bash
# Ubuntu/Debian
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 5432/tcp  # PostgreSQL (restrict to specific IPs)
sudo ufw allow 5672/tcp  # RabbitMQ (restrict to specific IPs)
sudo ufw enable
```

#### SSL Certificate Setup
```bash
# Install certbot
sudo apt install certbot nginx

# Get SSL certificate
sudo certbot --nginx -d yourdomain.com
```

### 3. Database Production Setup

#### PostgreSQL Configuration
```bash
# Install PostgreSQL
sudo apt install postgresql postgresql-contrib

# Configure for production
sudo nano /etc/postgresql/12/main/postgresql.conf

# Optimize settings
max_connections = 200
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 4MB
min_wal_size = 1GB
max_wal_size = 4GB
```

### 4. Application Deployment

#### Using Systemd Service
Create `/etc/systemd/system/khelogames.service`:
```ini
[Unit]
Description=Khelogames API
After=network.target

[Service]
Type=simple
User=khelogames
WorkingDirectory=/opt/khelogames
ExecStart=/opt/khelogames/khelogames
Restart=always
RestartSec=5
Environment=APP_ENV=production

[Install]
WantedBy=multi-user.target
```

#### Deploy Application
```bash
# Create application user
sudo useradd -r -s /bin/false khelogames

# Create directories
sudo mkdir -p /opt/khelogames
sudo chown khelogames:khelogames /opt/khelogames

# Copy application files
sudo cp khelogames /opt/khelogames/
sudo cp -r uploads /opt/khelogames/
sudo cp app.env /opt/khelogames/

# Set permissions
sudo chown -R khelogames:khelogames /opt/khelogames

# Start service
sudo systemctl enable khelogames
sudo systemctl start khelogames
```

### 5. Monitoring Setup

#### Install Prometheus
```bash
# Download and install Prometheus
wget https://github.com/prometheus/prometheus/releases/download/v2.40.0/prometheus-2.40.0.linux-amd64.tar.gz
tar xvfz prometheus-2.40.0.linux-amd64.tar.gz
sudo cp prometheus-2.40.0.linux-amd64/prometheus /usr/local/bin/
```

#### Install Grafana
```bash
sudo apt-get install -y software-properties-common
sudo add-apt-repository "deb https://packages.grafana.com/oss/deb stable main"
sudo apt-get update
sudo apt-get install grafana
sudo systemctl enable grafana-server
sudo systemctl start grafana-server
```

## 🔧 Configuration Management

### Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_DRIVER` | Database driver | postgres | Yes |
| `DB_SOURCE` | Database connection string | - | Yes |
| `SERVER_ADDRESS` | Server bind address | 0.0.0.0:8080 | No |
| `ACCESS_TOKEN_DURATION` | JWT access token duration | 15m | No |
| `REFRESH_TOKEN_DURATION` | JWT refresh token duration | 168h | No |
| `TOKEN_SYMMETRIC_KEY` | JWT signing key | - | Yes |
| `RABBIT_SOURCE` | RabbitMQ connection string | - | Yes |
| `REDIS_ADDRESS` | Redis connection string | - | No |
| `CLIENT_ID` | Google OAuth client ID | - | No |
| `CLIENT_SECRET` | Google OAuth client secret | - | No |

## 🧪 Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./api/auth -v

# Run integration tests
go test -tags=integration ./...
```

### Load Testing
```bash
# Install Apache Bench
sudo apt install apache2-utils

# Test API endpoints
ab -n 1000 -c 10 http://localhost:8080/api/v1/health
```

## 🆘 Troubleshooting

### Common Issues

#### Database Connection Issues
```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check connection
psql -h localhost -U username -d khelogames -c "SELECT 1"
```

#### RabbitMQ Connection Issues
```bash
# Check RabbitMQ status
sudo systemctl status rabbitmq-server

# Check management interface
curl http://localhost:15672
```

#### Port Already in Use
```bash
# Find process using port 8080
sudo lsof -i :8080

# Kill process
sudo kill -9 <PID>
```

#### Permission Issues
```bash
# Fix file permissions
sudo chown -R $USER:$USER /path/to/project

# Fix database permissions
sudo -u postgres psql
GRANT ALL PRIVILEGES ON DATABASE khelogames TO username;
```

### Debug Mode
```bash
# Enable debug logging
export LOG_LEVEL=debug
go run main.go
```

## 📞 Support

For additional support:
- Check the [troubleshooting guide](docs/TROUBLESHOOTING.md)
- Open an issue on GitHub
- Contact support@khelogames.com
