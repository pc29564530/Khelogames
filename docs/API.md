# Khelogames API Documentation

## Overview
This document provides detailed information about all API endpoints in the Khelogames platform.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Most endpoints require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### Authentication Endpoints

#### User Registration
```http
POST /api/v1/auth/signup
```

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123",
  "phone_number": "+1234567890"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

#### User Login
```http
POST /api/v1/auth/signin
```

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "securepassword123"
}
```

#### Google OAuth
```http
POST /api/v1/auth/google
```

**Request Body:**
```json
{
  "access_token": "google-oauth-access-token"
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### User Management

#### Get User Profile
```http
GET /api/v1/users/:id
```

**Response:**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "profile_picture": "https://example.com/avatar.jpg",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Update User Profile
```http
PUT /api/v1/users/:id
```

**Request Body:**
```json
{
  "username": "newusername",
  "email": "newemail@example.com",
  "profile_picture": "https://example.com/new-avatar.jpg"
}
```

### Player Management

#### Create Player Profile
```http
POST /api/v1/players
```

**Request Body:**
```json
{
  "user_id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "date_of_birth": "1990-01-01",
  "position": "Batsman",
  "sport": "cricket"
}
```

#### Get Player Details
```http
GET /api/v1/players/:id
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "date_of_birth": "1990-01-01",
  "position": "Batsman",
  "sport": "cricket",
  "statistics": {
    "matches_played": 50,
    "runs_scored": 2500,
    "batting_average": 45.5
  }
}
```

#### Search Players
```http
GET /api/v1/players/search?q=john&sport=cricket&limit=10&offset=0
```

### Team Management

#### Create Team
```http
POST /api/v1/teams
```

**Request Body:**
```json
{
  "name": "Mumbai Indians",
  "sport": "cricket",
  "logo": "https://example.com/team-logo.jpg",
  "description": "Professional cricket team"
}
```

#### Get Team Details
```http
GET /api/v1/teams/:id
```

**Response:**
```json
{
  "id": 1,
  "name": "Mumbai Indians",
  "sport": "cricket",
  "logo": "https://example.com/team-logo.jpg",
  "description": "Professional cricket team",
  "players": [...],
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### List Teams
```http
GET /api/v1/teams?sport=cricket&limit=20&offset=0
```

### Tournament Management

#### Create Tournament
```http
POST /api/v1/tournaments
```

**Request Body:**
```json
{
  "name": "IPL 2024",
  "sport": "cricket",
  "start_date": "2024-03-22",
  "end_date": "2024-05-26",
  "format": "round_robin",
  "description": "Indian Premier League 2024"
}
```

#### Get Tournament Details
```http
GET /api/v1/tournaments/:id
```

#### List Tournaments
```http
GET /api/v1/tournaments?sport=cricket&status=active
```

#### Get Tournament Standings
```http
GET /api/v1/tournaments/:id/standing
```

**Response:**
```json
{
  "tournament_id": 1,
  "standings": [
    {
      "team_id": 1,
      "team_name": "Mumbai Indians",
      "played": 14,
      "won": 10,
      "lost": 4,
      "points": 20,
      "net_run_rate": 0.547
    }
  ]
}
```

### Cricket Match Management

#### Create Cricket Match
```http
POST /api/v1/cricket/matches
```

**Request Body:**
```json
{
  "tournament_id": 1,
  "team1_id": 1,
  "team2_id": 2,
  "match_date": "2024-03-22T19:30:00Z",
  "venue": "Wankhede Stadium",
  "format": "t20"
}
```

#### Update Match Score
```http
PUT /api/v1/cricket/matches/:id/score
```

**Request Body:**
```json
{
  "team1_score": 185,
  "team1_wickets": 5,
  "team1_overs": 20,
  "team2_score": 180,
  "team2_wickets": 8,
  "team2_overs": 20,
  "result": "team1_won"
}
```

#### Get Match Details
```http
GET /api/v1/cricket/matches/:id
```

#### Get Player Score
```http
GET /api/v1/cricket/players/:id/score?match_id=123
```

### Football Match Management

#### Create Football Match
```http
POST /api/v1/football/matches
```

**Request Body:**
```json
{
  "tournament_id": 1,
  "team1_id": 1,
  "team2_id": 2,
  "match_date": "2024-03-22T19:30:00Z",
  "venue": "Camp Nou",
  "format": "league"
}
```

#### Add Match Incident
```http
POST /api/v1/football/matches/:id/incidents
```

**Request Body:**
```json
{
  "type": "goal",
  "player_id": 1,
  "minute": 45,
  "description": "Penalty goal"
}
```

### Community Management

#### Create Community
```http
POST /api/v1/communities
```

**Request Body:**
```json
{
  "name": "Cricket Fans",
  "description": "A community for cricket enthusiasts",
  "sport": "cricket",
  "privacy": "public"
}
```

#### Join Community
```http
POST /api/v1/communities/:id/join
```

#### Get Community Threads
```http
GET /api/v1/communities/:id/threads
```

### Thread Management

#### Create Thread
```http
POST /api/v1/threads
```

**Request Body:**
```json
{
  "community_id": 1,
  "title": "Match Discussion: MI vs CSK",
  "content": "Let's discuss today's match..."
}
```

#### Get Thread Details
```http
GET /api/v1/threads/:id
```

#### Add Comment
```http
POST /api/v1/threads/:id/comments
```

**Request Body:**
```json
{
  "content": "Great match today!"
}
```

### File Upload

#### Upload Media
```http
POST /api/v1/upload
```

**Content-Type**: multipart/form-data

**Form Data:**
- `file`: The file to upload
- `type`: "profile_picture" | "match_media" | "community_media"

**Response:**
```json
{
  "url": "https://example.com/uploads/filename.jpg",
  "filename": "filename.jpg",
  "size": 1024000
}
```

## WebSocket Endpoints

### Real-time Messaging
Connect to WebSocket for real-time messaging:

```
ws://localhost:8080/ws?token=<jwt-token>
```

### Message Format
```json
{
  "type": "message",
  "data": {
    "community_id": 1,
    "content": "Hello everyone!",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

## Error Responses

### Standard Error Format
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "field": "email",
      "issue": "must be a valid email address"
    }
  }
}
```

### Common Error Codes
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `422` - Unprocessable Entity
- `500` - Internal Server Error

## Rate Limiting
- 100 requests per minute for authenticated users
- 20 requests per minute for unauthenticated users

## Pagination
List endpoints support pagination:
```
GET /api/v1/players?limit=20&offset=40
```

## Filtering
Most list endpoints support filtering:
```
GET /api/v1/tournaments?sport=cricket&status=active
```

## Sorting
Sort results using query parameters:
```
GET /api/v1/players?sort=created_at&order=desc
