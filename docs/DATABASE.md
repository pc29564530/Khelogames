# Khelogames Database Documentation

## Overview
This document provides comprehensive documentation for the Khelogames database schema, including table structures, relationships, and query examples.

## Database Schema

### Core Tables

#### Users
Stores user account information.

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    profile_picture VARCHAR(255),
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

#### Players
Stores player profile information.

```sql
CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    date_of_birth DATE,
    position VARCHAR(50),
    sport VARCHAR(20) NOT NULL CHECK (sport IN ('cricket', 'football')),
    profile_picture VARCHAR(255),
    biography TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_players_user_id ON players(user_id);
CREATE INDEX idx_players_sport ON players(sport);
```

#### Teams
Stores team information.

```sql
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    sport VARCHAR(20) NOT NULL CHECK (sport IN ('cricket', 'football')),
    logo VARCHAR(255),
    description TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_teams_sport ON teams(sport);
CREATE INDEX idx_teams_created_by ON teams(created_by);
```

#### Tournaments
Stores tournament information.

```sql
CREATE TABLE tournaments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    sport VARCHAR(20) NOT NULL CHECK (sport IN ('cricket', 'football')),
    format VARCHAR(20) NOT NULL CHECK (format IN ('round_robin', 'knockout', 'league')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'active', 'completed', 'cancelled')),
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_tournaments_sport ON tournaments(sport);
CREATE INDEX idx_tournaments_status ON tournaments(status);
CREATE INDEX idx_tournaments_dates ON tournaments(start_date, end_date);
```

### Match Management Tables

#### Matches
Stores match information.

```sql
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    tournament_id INTEGER REFERENCES tournaments(id),
    team1_id INTEGER REFERENCES teams(id),
    team2_id INTEGER REFERENCES teams(id),
    match_date TIMESTAMP NOT NULL,
    venue VARCHAR(255),
    format VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'live', 'completed', 'cancelled')),
    result VARCHAR(20) CHECK (result IN ('team1_won', 'team2_won', 'draw', 'no_result')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_matches_tournament ON matches(tournament_id);
CREATE INDEX idx_matches_teams ON matches(team1_id, team2_id);
CREATE INDEX idx_matches_date ON matches(match_date);
CREATE INDEX idx_matches_status ON matches(status);
```

#### Cricket Scores
Stores cricket match scores.

```sql
CREATE TABLE cricket_scores (
    id SERIAL PRIMARY KEY,
    match_id INTEGER REFERENCES matches(id) ON DELETE CASCADE,
    team_id INTEGER REFERENCES teams(id),
    innings INTEGER NOT NULL,
    score INTEGER DEFAULT 0,
    wickets INTEGER DEFAULT 0,
    overs DECIMAL(3,1) DEFAULT 0,
    is_declared BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(match_id, team_id, innings)
);

-- Indexes
CREATE INDEX idx_cricket_scores_match ON cricket_scores(match_id);
CREATE INDEX idx_cricket_scores_team ON cricket_scores(team_id);
```

#### Cricket Player Scores
Stores individual player scores for cricket matches.

```sql
CREATE TABLE cricket_player_scores (
    id SERIAL PRIMARY KEY,
    match_id INTEGER REFERENCES matches(id) ON DELETE CASCADE,
    player_id INTEGER REFERENCES players(id),
    team_id INTEGER REFERENCES teams(id),
    runs INTEGER DEFAULT 0,
    balls INTEGER DEFAULT 0,
    fours INTEGER DEFAULT 0,
    sixes INTEGER DEFAULT 0,
    strike_rate DECIMAL(5,2),
    is_out BOOLEAN DEFAULT FALSE,
    how_out VARCHAR(50),
    bowler_id INTEGER REFERENCES players(id),
    fielder_id INTEGER REFERENCES players(id),
    overs_bowled DECIMAL(3,1) DEFAULT 0,
    runs_conceded INTEGER DEFAULT 0,
    wickets_taken INTEGER DEFAULT 0,
    economy_rate DECIMAL(4,2),
    catches INTEGER DEFAULT 0,
    run_outs INTEGER DEFAULT 0,
    stumping INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_cricket_player_scores_match ON cricket_player_scores(match_id);
CREATE INDEX idx_cricket_player_scores_player ON cricket_player_scores(player_id);
CREATE INDEX idx_cricket_player_scores_team ON cricket_player_scores(team_id);
```

#### Football Scores
Stores football match scores.

```sql
CREATE TABLE football_scores (
    id SERIAL PRIMARY KEY,
    match_id INTEGER REFERENCES matches(id) ON DELETE CASCADE,
    team_id INTEGER REFERENCES teams(id),
    goals INTEGER DEFAULT 0,
    own_goals INTEGER DEFAULT 0,
    penalties INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(match_id, team_id)
);

-- Indexes
CREATE INDEX idx_football_scores_match ON football_scores(match_id);
CREATE INDEX idx_football_scores_team ON football_scores(team_id);
```

#### Football Incidents
Stores football match incidents.

```sql
CREATE TABLE football_incidents (
    id SERIAL PRIMARY KEY,
    match_id INTEGER REFERENCES matches(id) ON DELETE CASCADE,
    player_id INTEGER REFERENCES players(id),
    team_id INTEGER REFERENCES teams(id),
    type VARCHAR(20) NOT NULL CHECK (type IN ('goal', 'own_goal', 'penalty', 'yellow_card', 'red_card', 'substitution')),
    minute INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_football_incidents_match ON football_incidents(match_id);
CREATE INDEX idx_football_incidents_player ON football_incidents(player_id);
CREATE INDEX idx_football_incidents_type ON football_incidents(type);
```

### Community Tables

#### Communities
Stores community information.

```sql
CREATE TABLE communities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    sport VARCHAR(20),
    privacy VARCHAR(20) DEFAULT 'public' CHECK (privacy IN ('public', 'private')),
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_communities_sport ON communities(sport);
CREATE INDEX idx_communities_privacy ON communities(privacy);
```

#### Community Members
Stores community membership information.

```sql
CREATE TABLE community_members (
    id SERIAL PRIMARY KEY,
    community_id INTEGER REFERENCES communities(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member' CHECK (role IN ('admin', 'moderator', 'member')),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(community_id, user_id)
);

-- Indexes
CREATE INDEX idx_community_members_community ON community_members(community_id);
CREATE INDEX idx_community_members_user ON community_members(user_id);
```

#### Threads
Stores discussion threads.

```sql
CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    community_id INTEGER REFERENCES communities(id) ON DELETE CASCADE,
    author_id INTEGER REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT,
    is_pinned BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_threads_community ON threads(community_id);
CREATE INDEX idx_threads_author ON threads(author_id);
CREATE INDEX idx_threads_created ON threads(created_at);
```

#### Comments
Stores thread comments.

```sql
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    thread_id INTEGER REFERENCES threads(id) ON DELETE CASCADE,
    author_id INTEGER REFERENCES users(id),
    content TEXT NOT NULL,
    parent_id INTEGER REFERENCES comments(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_comments_thread ON comments(thread_id);
CREATE INDEX idx_comments_author ON comments(author_id);
CREATE INDEX idx_comments_parent ON comments(parent_id);
```

### Tournament Tables

#### Tournament Teams
Stores tournament team registrations.

```sql
CREATE TABLE tournament_teams (
    id SERIAL PRIMARY KEY,
    tournament_id INTEGER REFERENCES tournaments(id) ON DELETE CASCADE,
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    group_name VARCHAR(10),
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tournament_id, team_id)
);

-- Indexes
CREATE INDEX idx_tournament_teams_tournament ON tournament_teams(tournament_id);
CREATE INDEX idx_tournament_teams_team ON tournament_teams(team_id);
```

#### Tournament Standings
Stores tournament standings.

```sql
CREATE TABLE tournament_standings (
    id SERIAL PRIMARY KEY,
    tournament_id INTEGER REFERENCES tournaments(id) ON DELETE CASCADE,
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    matches_played INTEGER DEFAULT 0,
    won INTEGER DEFAULT 0,
    lost INTEGER DEFAULT 0,
    drawn INTEGER DEFAULT 0,
    points INTEGER DEFAULT 0,
    net_run_rate DECIMAL(8,5) DEFAULT 0,
    goals_for INTEGER DEFAULT 0,
    goals_against INTEGER DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tournament_id, team_id)
);

-- Indexes
CREATE INDEX idx_tournament_standings_tournament ON tournament_standings(tournament_id);
CREATE INDEX idx_tournament_standings_team ON tournament_standings(team_id);
```

### Utility Tables

#### Sessions
Stores user sessions.

```sql
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255),
    client_ip VARCHAR(45),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
```

#### Follows
Stores user follow relationships.

```sql
CREATE TABLE follows (
    id SERIAL PRIMARY KEY,
    follower_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    following_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(follower_id, following_id)
);

-- Indexes
CREATE INDEX idx_follows_follower ON follows(follower_id);
CREATE INDEX idx_follows_following ON follows(following_id);
```

## Query Examples

### Complex Queries

#### Get Tournament Standings with Team Details
```sql
SELECT 
    ts.tournament_id,
    t.name as tournament_name,
    te.name as team_name,
    te.logo as team_logo,
    ts.matches_played,
    ts.won,
    ts.lost,
    ts.drawn,
    ts.points,
    ts.net_run_rate
FROM tournament_standings ts
JOIN tournaments t ON ts.tournament_id = t.id
JOIN teams te ON ts.team_id = te.id
WHERE ts.tournament_id = $1
ORDER BY ts.points DESC, ts.net_run_rate DESC;
```

#### Get Player Statistics
```sql
SELECT 
    p.id,
    p.first_name,
    p.last_name,
    COUNT(DISTINCT m.id) as matches_played,
    SUM(cps.runs) as total_runs,
    AVG(cps.runs) as batting_average,
    SUM(cps.wickets_taken) as total_wickets,
    AVG(cps.economy_rate) as bowling_average
FROM players p
JOIN cricket_player_scores cps ON p.id = cps.player_id
JOIN matches m ON cps.match_id = m.id
WHERE p.sport = 'cricket'
GROUP BY p.id, p.first_name, p.last_name;
```

#### Get Community Activity
```sql
SELECT 
    c.id,
    c.name,
    COUNT(DISTINCT cm.user_id) as member_count,
    COUNT(DISTINCT t.id) as thread_count,
    COUNT(DISTINCT co.id) as comment_count,
    MAX(t.created_at) as last_activity
FROM communities c
LEFT JOIN community_members cm ON c.id = cm.community_id
LEFT JOIN threads t ON c.id = t.community_id
LEFT JOIN comments co ON t.id = co.thread_id
WHERE c.id = $1
GROUP BY c.id, c.name;
```

## Database Maintenance

### Backup Strategy
```bash
# Daily backup
pg_dump -h localhost -U username khelogames > backup_$(date +%Y%m%d).sql

# Restore backup
psql -h localhost -U username khelogames < backup_20240101.sql
```

### Performance Optimization
```sql
-- Analyze tables
ANALYZE users;
ANALYZE matches;

-- Vacuum and analyze
VACUUM ANALYZE;

-- Update statistics
UPDATE statistics;
```

### Monitoring Queries
```sql
-- Check table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
WHERE schemaname NOT IN ('information_schema', 'pg_catalog')
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Check slow queries
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY total_time DESC
LIMIT 10;
