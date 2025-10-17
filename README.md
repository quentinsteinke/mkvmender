# MKV Mender

A community-driven CLI tool for renaming movie and TV show rips from Blu-ray and DVDs. MKV Mender uses file hashing to match your files with naming data submitted by other users, making it easy to maintain a clean media library.

## Features

- **Hash-based matching**: Uses SHA-256 to uniquely identify files
- **Community-driven**: Users can upload and vote on naming submissions
- **Interactive CLI**: Easy-to-use command-line interface
- **Batch processing**: Process entire directories at once
- **Voting system**: Upvote/downvote naming submissions to surface the best options
- **Metadata support**: Include title, year, season, episode, quality, and source information

## Architecture

- **CLI**: Go-based command-line tool
- **API Server**: RESTful HTTP API
- **Database**: Turso (distributed SQLite) for edge performance

## Installation

### Prerequisites

- Go 1.24 or higher
- Turso account (for running the server)

### Build from source

```bash
# Clone the repository
git clone https://github.com/quentinsteinke/mkvmender.git
cd mkvmender

# Install dependencies
go mod tidy

# Build CLI and server
go build -o bin/mkvmender ./cmd/cli
go build -o bin/mkvmender-server ./cmd/server

# Optionally, install to your PATH
cp bin/mkvmender /usr/local/bin/
```

## Usage

### Server Setup

1. Create a Turso database:
```bash
turso db create mkvmender
turso db show mkvmender
```

2. Set environment variables:
```bash
export TURSO_DATABASE_URL="libsql://[your-database-url]"
export TURSO_AUTH_TOKEN="[your-auth-token]"
export PORT=8080
```

3. Run the server:
```bash
./bin/mkvmender-server
```

### CLI Usage

#### Register a new account

```bash
mkvmender register
```

This will create a new user account and provide you with an API key.

#### Configure authentication

```bash
mkvmender login --api-key YOUR_API_KEY
```

Or set it interactively:
```bash
mkvmender login
```

#### Hash a file

```bash
mkvmender hash movie.mkv
```

#### Look up naming options

```bash
mkvmender lookup movie.mkv
```

#### Rename a file interactively

```bash
mkvmender rename movie.mkv
```

#### Upload a naming submission

For a movie:
```bash
mkvmender upload movie.mkv \
  --type movie \
  --name "The Matrix (1999).mkv" \
  --title "The Matrix" \
  --year 1999 \
  --quality 1080p \
  --source Blu-ray
```

For a TV show:
```bash
mkvmender upload episode.mkv \
  --type tv \
  --name "Breaking Bad - S01E01.mkv" \
  --title "Breaking Bad" \
  --season 1 \
  --episode 1 \
  --quality 1080p \
  --source Blu-ray
```

#### Vote on submissions

```bash
mkvmender vote movie.mkv
```

This opens an interactive menu where you can:
1. See all naming submissions for the file
2. Select which submission to vote on
3. Choose to upvote or downvote
4. See updated rankings immediately

#### Batch process a directory

```bash
mkvmender batch /path/to/movies
```

## API Endpoints

### Public Endpoints

- `GET /api/health` - Health check
- `POST /api/register` - Register new user
- `GET /api/lookup?hash=<hash>` - Look up naming submissions

### Protected Endpoints (require authentication)

- `POST /api/upload` - Upload naming submission
- `POST /api/vote` - Vote on submission
- `DELETE /api/vote/delete?submission_id=<id>` - Remove vote

## Database Schema

- **users**: User accounts and API keys
- **file_hashes**: Unique file hashes and metadata
- **naming_submissions**: User-submitted file names
- **votes**: User votes on submissions
- **naming_metadata**: Extended metadata for submissions

## Configuration

CLI configuration is stored in `~/.mkvmender/config.yaml`:

```yaml
api_key: your-api-key-here
base_url: http://localhost:8080
```

## Development

### Project Structure

```
mkvmender/
├── cmd/
│   ├── cli/          # CLI application
│   └── server/       # API server
├── internal/
│   ├── hasher/       # File hashing
│   ├── api/          # API client
│   ├── database/     # Database layer
│   ├── models/       # Data models
│   └── handlers/     # HTTP handlers
├── migrations/       # Database migrations
└── pkg/             # Public libraries
```

### Running Tests

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - See LICENSE file for details

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Turso](https://turso.tech) - Distributed SQLite database
- [libsql-client-go](https://github.com/tursodatabase/libsql-client-go) - Turso Go client
