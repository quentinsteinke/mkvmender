# MKV Mender Frontend

Simple web interface for searching the MKV Mender database.

## Features

- 🔐 Login with API key
- 🔍 Search movies and TV shows
- 🎯 Fuzzy matching with typo tolerance
- 📊 Sort by relevance, votes, date, or title
- 📱 Responsive design

## Local Development

```bash
# Install dependencies
npm install

# Start development server
npm start
```

Visit http://localhost:3000

## Railway Deployment

This frontend should be deployed as a separate Railway service from the API.

### Step 1: Create New Service

1. Go to https://railway.app
2. Click **"New Project"**
3. Select **"Deploy from GitHub repo"**
4. Choose **`quentinsteinke/mkvmender`**
5. Select **"frontend"** as the root directory

### Step 2: Configure Build Settings

Railway will auto-detect the Node.js project.

**Root Directory**: `/frontend`

### Step 3: Configure Domain

1. Go to **"Settings"** → **"Networking"**
2. Add custom domain: **`mkvmender.org`**
3. Follow Railway's instructions to point your domain

The API should be on a different service with domain **`mkvmender.org`** or subdomain.

### Step 4: CORS Configuration

Make sure your API server (the Go backend) has CORS configured to allow requests from your frontend domain.

## Project Structure

```
frontend/
├── public/
│   ├── index.html     # Main HTML file
│   ├── css/
│   │   └── style.css  # Styles
│   └── js/
│       └── app.js     # JavaScript logic
├── server.js          # Express server
├── package.json       # Dependencies
└── railway.toml       # Railway config
```

## How It Works

1. **Login**: Users enter their API key (obtained via CLI `mkvmender register`)
2. **Search**: Query the API for movies/TV shows
3. **Results**: Display submissions with votes and metadata

## API Endpoints Used

- `GET /api/health` - Health check (with auth)
- `GET /api/search?q=query&sort=relevance&fuzzy=true` - Search database

## Registration

Users must register via CLI to get an API key:

```bash
mkvmender register
```

This is intentional - keeps the database quality high by requiring CLI usage.
