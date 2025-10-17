# MKV Mender - Project Structure

## Overview

MKV Mender is a complete full-stack application for community-driven media file renaming.

## Directory Structure

```
mkvmender/
├── cmd/                          # Application entry points
│   ├── cli/                      # CLI application
│   │   ├── main.go              # CLI entry point
│   │   ├── hash.go              # Hash command
│   │   ├── lookup.go            # Lookup command
│   │   ├── rename.go            # Rename command
│   │   ├── upload.go            # Upload command
│   │   ├── vote.go              # Vote command
│   │   ├── batch.go             # Batch processing
│   │   ├── login.go             # Login/config command
│   │   └── register.go          # User registration
│   └── server/                   # API server
│       └── main.go              # Server entry point
│
├── internal/                     # Private application code
│   ├── api/                      # API client for CLI
│   │   ├── client.go            # HTTP client
│   │   └── config.go            # Configuration management
│   ├── database/                 # Database layer
│   │   ├── database.go          # Database connection
│   │   ├── users.go             # User operations
│   │   ├── file_hashes.go       # File hash operations
│   │   ├── submissions.go       # Submission operations
│   │   └── votes.go             # Voting operations
│   ├── handlers/                 # HTTP handlers
│   │   ├── handlers.go          # API endpoints
│   │   ├── middleware.go        # Authentication & CORS
│   │   └── response.go          # Response helpers
│   ├── hasher/                   # File hashing
│   │   └── hasher.go            # SHA-256 hashing
│   └── models/                   # Data models
│       └── models.go            # Shared data structures
│
├── migrations/                   # Database migrations
│   └── 001_initial_schema.sql  # Initial database schema
│
├── bin/                          # Compiled binaries (gitignored)
│   ├── mkvmender               # CLI binary
│   └── mkvmender-server        # Server binary
│
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore rules
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Makefile                     # Build automation
├── README.md                    # Main documentation
├── QUICKSTART.md               # Quick start guide
└── PROJECT_STRUCTURE.md        # This file
```

## Key Components

### 1. CLI Application (cmd/cli/)
- **Purpose**: User-facing command-line interface
- **Commands**: hash, lookup, rename, upload, vote, batch, login, register
- **Framework**: Cobra CLI framework
- **Configuration**: Stored in ~/.mkvmender/config.yaml

### 2. API Server (cmd/server/)
- **Purpose**: HTTP API for data storage and retrieval
- **Framework**: Standard Go net/http
- **Endpoints**: 
  - Public: /api/health, /api/register, /api/lookup
  - Protected: /api/upload, /api/vote, /api/vote/delete
- **Authentication**: Bearer token (API key)

### 3. Database Layer (internal/database/)
- **Technology**: Turso (distributed SQLite)
- **Tables**: users, file_hashes, naming_submissions, votes, naming_metadata
- **Features**: Full CRUD operations, vote counting, submission ranking

### 4. File Hasher (internal/hasher/)
- **Algorithm**: SHA-256
- **Features**: Full file hashing, partial hashing for speed
- **Output**: Hash string + file size

### 5. API Client (internal/api/)
- **Purpose**: HTTP client for CLI to communicate with server
- **Features**: Configuration management, automatic authentication
- **Transport**: HTTP/REST with JSON

## Data Flow

### Upload Flow
```
User → CLI → Hash File → API Client → Server → Database
```

### Lookup Flow
```
User → CLI → Hash File → API Client → Server → Database → Return Submissions → Display
```

### Rename Flow
```
User → CLI → Hash → Lookup → Select → Confirm → Rename Local File
```

### Voting Flow
```
User → CLI → Vote Command → API Client → Server → Database → Update Vote Count
```

## Technology Stack

- **Language**: Go 1.24+
- **CLI Framework**: Cobra
- **Database**: Turso (libSQL/SQLite)
- **HTTP**: Standard library (net/http)
- **Configuration**: YAML (gopkg.in/yaml.v3)
- **Database Driver**: github.com/tursodatabase/libsql-client-go

## Security Features

1. **API Key Authentication**: Each user has unique API key
2. **Bearer Token**: Passed in Authorization header
3. **Secure Storage**: API keys stored with restricted permissions (0600)
4. **Vote Constraints**: One vote per user per submission
5. **Input Validation**: All inputs validated before database operations

## Scalability

- **Edge Distribution**: Turso provides global edge database
- **Stateless API**: Server can be horizontally scaled
- **Efficient Hashing**: Partial hashing option for large files
- **Indexed Queries**: All lookups use database indexes

## Future Enhancements

- [ ] Web UI for browsing submissions
- [ ] Automatic file recognition using metadata
- [ ] Integration with media servers (Plex, Jellyfin)
- [ ] Bulk voting system
- [ ] User reputation system
- [ ] Rate limiting and abuse prevention
- [ ] Submission editing/deletion
- [ ] Advanced search and filtering
- [ ] Mobile app

## Build Process

1. **Dependencies**: `go mod tidy`
2. **Build**: `make build`
3. **Test**: `make test`
4. **Install**: `make install`

## Configuration Files

- **Server**: Environment variables (.env)
- **CLI**: ~/.mkvmender/config.yaml
- **Database**: migrations/001_initial_schema.sql

## API Documentation

See README.md for full API endpoint documentation.
