# MKV Mender - Quick Start Guide

Get up and running with MKV Mender in 5 minutes.

## Step 1: Build the Project

```bash
make build
```

This creates two binaries:
- `bin/mkvmender` - CLI tool
- `bin/mkvmender-server` - API server

## Step 2: Set Up Turso Database

If you haven't already, sign up for [Turso](https://turso.tech) (free tier available).

```bash
# Install Turso CLI
curl -sSfL https://get.tur.so/install.sh | bash

# Create a new database
turso db create mkvmender

# Get your database URL and auth token
turso db show mkvmender
```

## Step 3: Configure Environment

Create a `.env` file:

```bash
cp .env.example .env
```

Edit `.env` with your Turso credentials:

```
TURSO_DATABASE_URL=libsql://your-database-name.turso.io
TURSO_AUTH_TOKEN=your-auth-token-here
PORT=8080
```

## Step 4: Start the Server

```bash
# Load environment variables
export $(cat .env | xargs)

# Start the server
./bin/mkvmender-server
```

The server will:
1. Connect to your Turso database
2. Run migrations automatically
3. Start listening on port 8080

## Step 5: Register a User Account

In a new terminal:

```bash
./bin/mkvmender register
```

Enter a username when prompted. Save the API key that's displayed!

## Step 6: Configure the CLI

```bash
./bin/mkvmender login
```

Enter your API key when prompted.

## Step 7: Try It Out

### Hash a file

```bash
./bin/mkvmender hash your-movie.mkv
```

### Upload naming for a file

```bash
./bin/mkvmender upload your-movie.mkv \
  --type movie \
  --name "The Matrix (1999) 1080p BluRay.mkv" \
  --title "The Matrix" \
  --year 1999 \
  --quality 1080p \
  --source Blu-ray
```

### Look up existing naming options

```bash
./bin/mkvmender lookup your-movie.mkv
```

### Rename a file interactively

```bash
./bin/mkvmender rename your-movie.mkv
```

### Vote on naming submissions

```bash
./bin/mkvmender vote your-movie.mkv
```

This opens an interactive menu to:
- View all submissions for the file
- Select which one to vote on
- Choose upvote or downvote
- See updated rankings

## Common Commands

```bash
# Get help
./bin/mkvmender --help

# Get help for a specific command
./bin/mkvmender upload --help

# Process all files in a directory
./bin/mkvmender batch /path/to/movies

# Vote on a file (interactive)
./bin/mkvmender vote your-movie.mkv
```

## Tips

1. **Install globally**: Copy the CLI to your PATH for easier access:
   ```bash
   sudo cp bin/mkvmender /usr/local/bin/
   ```

2. **API URL**: If your server is on a different machine, set the URL:
   ```bash
   ./bin/mkvmender login --url http://your-server:8080
   ```

3. **Batch processing**: Use dry-run mode to preview changes:
   ```bash
   ./bin/mkvmender batch /path/to/movies --dry-run
   ```

4. **File extensions**: Specify which extensions to process:
   ```bash
   ./bin/mkvmender batch /path/to/movies --ext .mkv,.mp4
   ```

## Troubleshooting

### "Failed to connect to database"
- Check your `TURSO_DATABASE_URL` and `TURSO_AUTH_TOKEN`
- Verify your Turso database is active: `turso db show mkvmender`

### "Authentication required"
- Run `./bin/mkvmender login` to set your API key
- Check `~/.mkvmender/config.yaml` exists and contains your API key

### "No naming submissions found"
- Be the first! Upload your own naming with `./bin/mkvmender upload`
- The database starts empty - it's community-driven

## Next Steps

- Invite others to use your server and contribute
- Set up a public server for the community
- Contribute to the project on GitHub

## Development

Build and run in development mode:

```bash
# Build
make build

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean
```

## Need Help?

- Check the [README.md](README.md) for detailed documentation
- Report issues on GitHub
- Join the community discussion

Happy renaming! 🎬
