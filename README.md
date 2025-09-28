# Khelogames - Sports Platform Backend API

A comprehensive backend API for managing sports tournaments, matches, and communities with support for cricket and football.

## 🏆 Overview

Khelogames is a Go-based backend API that provides a complete sports management platform. It supports tournament management, live match scoring, player profiles, team management, and community features for both cricket and football sports.

## 🚀 Features

### Sports Management
- **Cricket**: Complete match management with batting/bowling scorecards, player statistics, and live scoring
- **Football**: Match management with incidents, lineups, substitutions, and statistics

### Tournament System
- Create and manage tournaments
- Group stages and knockout formats
- Automatic standings calculation
- Team and player statistics tracking

### User Management
- JWT-based authentication
- Google OAuth integration
- User profiles and player management
- Community membership system

### Real-time Features
- WebSocket-based messaging system
- Real-time match updates
- Live chat for communities
- RabbitMQ message queue integration

### Media & Files
- File upload system with Tus protocol
- Profile picture management
- Match media uploads
- Cloud storage integration

## 🛠️ Technology Stack

- **Language**: Go 1.21.2
- **Framework**: Gin-gonic (HTTP router)
- **Database**: PostgreSQL with pgx driver
- **Authentication**: JWT tokens
- **Real-time**: WebSocket with Gorilla
- **Message Queue**: RabbitMQ
- **File Upload**: Tus protocol
- **SMS**: Twilio
- **OAuth**: Google OAuth2

## 📋 Prerequisites

- Go 1.21.2 or higher
- PostgreSQL 12+
- RabbitMQ 3.8+
- Redis (optional, for caching)

## 🔧 Installation & Setup

### 1. Clone the Repository
```bash
git clone https://github.com/your-org/khelogames.git
cd khelogames
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Database Setup
```bash
# Create database
createdb khelogames

# Run migrations (if available)
# Or use the provided SQL files in database/
```

### 4. Environment Configuration
Create `app.env` file:
```bash
DB_DRIVER=postgres
DB_SOURCE=postgresql://username:password@localhost:5432/khelogames?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=168h
TOKEN_SYMMETRIC_KEY=your-32-character-secret-key
RABBIT_SOURCE=amqp://guest:guest@localhost:5672/
```

### 5. Start Services
```bash
# Start PostgreSQL
pg_ctl -D /usr/local/var/postgres start

# Start RabbitMQ
rabbitmq-server

# Run the application
go run main.go
```

## 🏗️ Architecture

### Project Structure
```
khelogames/
├── api/                    # API endpoints
│   ├── auth/              # Authentication endpoints
│   ├── handlers/          # General handlers
│   ├── sports/            # Sports-specific APIs
│   ├── tournaments/       # Tournament management
│   ├── teams/             # Team management
│   ├── players/           # Player profiles
│   └── messenger/         # Real-time messaging
├── database/              # Database models and queries
├── util/                  # Utility functions
├── logger/                # Logging configuration
├── token/                 # JWT token management
└── uploads/               # File upload handling
```

### Key Components

1. **Authentication System**
   - JWT token-based authentication
   - Google OAuth2 integration
   - Session management
   - Password reset via SMS

2. **Sports Engine**
   - Cricket: Batting/bowling scorecards, player stats
   - Football: Incidents, lineups, substitutions, statistics
   - Live match updates
   - Match result calculation

3. **Tournament System**
   - Tournament creation and management
   - Group stages and knockout formats
   - Automatic standings calculation
   - Team/player statistics

4. **Community Features**
   - Community creation and management
   - Thread-based discussions
   - Real-time messaging
   - File sharing

## 📡 API Endpoints

### Authentication
- `POST /api/v1/auth/signup` - User registration
- `POST /api/v1/auth/signin` - User login
- `POST /api/v1/auth/google` - Google OAuth
- `POST /api/v1/auth/refresh` - Token refresh

### Sports
- **Cricket**
  - `GET /api/v1/cricket/matches` - List cricket matches
  - `POST /api/v1/cricket/matches` - Create cricket match
  - `PUT /api/v1/cricket/matches/:id/score` - Update match score
  - `GET /api/v1/cricket/players/:id/stats` - Player statistics

- **Football**
  - `GET /api/v1/football/matches` - List football matches
  - `POST /api/v1/football/matches` - Create football match
  - `PUT /api/v1/football/matches/:id/incidents` - Add match incidents

### Tournaments
- `GET /api/v1/tournaments` - List tournaments
- `POST /api/v1/tournaments` - Create tournament
- `GET /api/v1/tournaments/:id/standing` - Tournament standings
- `GET /api/v1/tournaments/:id/matches` - Tournament matches

### Teams & Players
- `GET /api/v1/teams` - List teams
- `POST /api/v1/teams` - Create team
- `GET /api/v1/players` - List players
- `POST /api/v1/players` - Create player profile

### Communities
- `GET /api/v1/communities` - List communities
- `POST /api/v1/communities` - Create community
- `POST /api/v1/communities/:id/join` - Join community
- `GET /api/v1/communities/:id/threads` - Community threads

## 🧪 Testing

### Unit Tests
```bash
go test ./...
```

### Integration Tests
```bash
go test -tags=integration ./...
```

## 📊 Database Schema

### Key Tables
- `users` - User accounts
- `players` - Player profiles
- `teams` - Team information
- `tournaments` - Tournament details
- `matches` - Match records
- `cricket_scores` - Cricket match scores
- `football_scores` - Football match scores
- `communities` - Community information
- `threads` - Discussion threads

## 🔐 Security Features

- JWT token-based authentication
- Password hashing with bcrypt
- Rate limiting
- CORS configuration
- Input validation and sanitization
- SQL injection prevention

## 🚀 Deployment

### Docker Deployment
```bash
# Build the Docker image
docker build -t khelogames .

# Run with Docker Compose
docker-compose up -d
```

### Production Setup
1. Set up PostgreSQL with proper configuration
2. Configure environment variables
3. Set up SSL certificates
4. Configure reverse proxy (nginx)
5. Set up monitoring and logging

## 📈 Monitoring & Logging

- Structured logging with Logrus
- Request/response logging
- Error tracking and alerting
- Performance monitoring
- Health check endpoints

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support, email support@khelogames.com or join our community Discord server.
