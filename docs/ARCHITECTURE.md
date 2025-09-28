# Khelogames Architecture Documentation

## Overview

This document describes the architecture of the Khelogames backend API, including system design, component interactions, and data flow.

## System Architecture

### High-Level Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                        Load Balancer                        │
│                    (Nginx/CloudFlare)                      │
└─────────────────────┬───────────────────────────────────────┘
                    │
┌─────────────────────┴───────────────────────────────────────┐
│                    API Gateway                              │
│                  (Gin-gonic)                                 │
└─────────────────────┬───────────────────────────────────────┘
                    │
┌─────────────────────┴───────────────────────────────────────┐
│                   Application Layer                         │
├─────────────────────────────────────────────────────────────┤
│  Auth Service  │  Sports Service  │  Tournament Service    │
│  User Service  │  Team Service   │  Community Service   │
└─────────────────────┬───────────────────────────────────────┘
                    │
┌─────────────────────┴───────────────────────────────────────┐
│                    Data Layer                               │
├─────────────────────────────────────────────────────────────┤
│   PostgreSQL   │    Redis    │    RabbitMQ    │  MinIO   │
└─────────────────────────────────────────────────────────────┘
```

## Component Architecture

### 1. API Gateway (Gin-gonic)
- **Purpose**: HTTP request routing and middleware
- **Features**:
  - Request validation
  - Rate limiting
  - CORS handling
  - Request logging
  - Error handling

### 2. Authentication Service
- **Purpose**: User authentication and authorization
- **Components**:
  - JWT token generation/validation
  - Google OAuth integration
  - Password hashing (bcrypt)
  - Session management

### 3. Sports Service
- **Purpose**: Manage sports-specific functionality
- **Sub-services**:
  - **Cricket Service**: Match scoring, player statistics
  - **Football Service**: Match incidents, lineups, statistics

### 4. Tournament Service
- **Purpose**: Tournament lifecycle management
- **Features**:
  - Tournament creation
  - Group stage management
  - Knockout stage automation
  - Standings calculation

### 5. Community Service
- **Purpose**: Social features
- **Features**:
  - Community creation/management
  - Thread discussions
  - Real-time messaging
  - File sharing

## Database Design

### Entity Relationship Diagram
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    Users    │    │  Players    │    │   Teams     │
├─────────────┤    ├─────────────┤    ├─────────────┤
│ id (PK)     │    │ id (PK)     │    │ id (PK)     │
│ username    │    │ user_id     │    │ name        │
│ email       │    │ first_name  │    │ sport       │
│ password    │    │ last_name   │    │ logo        │
│ created_at  │    │ position    │    │ created_at  │
└─────────────┘    └─────────────┘    └─────────────┘
        │                   │                   │
        │                   │                   │
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Tournaments │    │   Matches   │    │ Communities │
├─────────────┤    ├─────────────┤    ├─────────────┤
│ id (PK)     │    │ id (PK)     │    │ id (PK)     │
│ name        │    │ tournament  │    │ name        │
│ sport       │    │ team1_id    │    │ description │
│ start_date  │    │ team2_id    │    │ sport       │
│ end_date    │    │ match_date  │    │ created_at  │
└─────────────┘    └─────────────┘    └─────────────┘
```

### Database Schema Details

#### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    profile_picture VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Service Communication

### Synchronous Communication
- **HTTP/REST**: Between microservices
- **gRPC**: For internal service communication (future enhancement)

### Asynchronous Communication
- **RabbitMQ**: For event-driven architecture
- **WebSocket**: For real-time features

### Event Flow Example
```
User Creates Tournament
    ↓
Tournament Service
    ↓
RabbitMQ Event: tournament.created
    ↓
Notification Service → Email/SMS
    ↓
Community Service → Create discussion thread
```

## Data Flow

### 1. User Registration Flow
```
Client → API Gateway → Auth Service → Database
   ↓
Token Generation → JWT Token → Client
```

### 2. Match Score Update Flow
```
Client → API Gateway → Sports Service → Database
   ↓
Event: score.updated → RabbitMQ → WebSocket → Clients
```

### 3. File Upload Flow
```
Client → API Gateway → Upload Service → MinIO/S3
   ↓
Database Update → Event: file.uploaded → Notification
```

## Security Architecture

### Authentication & Authorization
- **JWT Tokens**: Stateless authentication
- **Refresh Tokens**: Long-lived tokens for renewal
- **OAuth 2.0**: Google OAuth integration
- **Role-based Access Control**: User roles and permissions

### Security Measures
- **HTTPS**: TLS encryption
- **Rate Limiting**: API rate limiting per user/IP
- **Input Validation**: Request validation middleware
- **SQL Injection Prevention**: Parameterized queries
- **XSS Protection**: Input sanitization
- **CORS**: Cross-origin resource sharing configuration

## Scalability Design

### Horizontal Scaling
- **Load Balancer**: Nginx for distributing requests
- **Database Sharding**: Partition data by sport/tournament
- **Caching**: Redis for frequently accessed data
- **CDN**: CloudFront for static assets

### Vertical Scaling
- **Database Optimization**: Indexes, query optimization
- **Connection Pooling**: Database connection management
- **Resource Optimization**: Efficient memory usage

## Monitoring & Observability

### Logging
- **Structured Logging**: JSON format logs
- **Log Levels**: Debug, Info, Warn, Error
- **Centralized Logging**: ELK Stack (future enhancement)

### Metrics
- **Application Metrics**: Request count, response time
- **Business Metrics**: Active users, matches created
- **System Metrics**: CPU, memory, disk usage

### Tracing
- **Distributed Tracing**: Jaeger/OpenTelemetry
- **Request Tracing**: Trace ID propagation

## Deployment Architecture

### Development Environment
- **Local**: Docker Compose
- **Database**: Local PostgreSQL
- **Message Queue**: Local RabbitMQ

### Staging Environment
- **Cloud**: AWS/GCP/Azure
- **Database**: Managed PostgreSQL
- **Message Queue**: Managed RabbitMQ
- **Storage**: Cloud storage (S3/GCS)

### Production Environment
- **Load Balancer**: Cloud Load Balancer
- **Auto-scaling**: Kubernetes/ECS
- **Database**: Multi-AZ PostgreSQL
- **Caching**: Redis Cluster
- **CDN**: CloudFront/CloudFlare

## Technology Stack Details

### Core Technologies
- **Language**: Go 1.21.2
- **Framework**: Gin-gonic
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Message Queue**: RabbitMQ 3.11

### Supporting Technologies
- **Logging**: Logrus
- **Configuration**: Viper
- **Database ORM**: SQLC (SQL compiler)
- **Testing**: Testify, Ginkgo
- **Monitoring**: Prometheus, Grafana

### Infrastructure
- **Container**: Docker
- **Orchestration**: Kubernetes (future)
- **CI/CD**: GitHub Actions
- **Cloud**: AWS/GCP/Azure
