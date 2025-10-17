# Simplified Deployment Guide

The API server now serves both the API endpoints and the frontend static files from a **single Railway service**.

## Quick Deployment

### 1. Deploy to Railway

1. Go to [railway.app](https://railway.app)
2. Click **"New Project"**
3. Select **"Deploy from GitHub repo"**
4. Choose **`quentinsteinke/mkvmender`**
5. Root directory: `/` (leave blank or set to root)

### 2. Add Environment Variables

In the **"Variables"** tab:

```
TURSO_DATABASE_URL=libsql://your-database.turso.io
TURSO_AUTH_TOKEN=your-auth-token
```

### 3. Configure Domain

1. Go to **"Settings"** → **"Networking"**
2. Add custom domain: **`mkvmender.org`**
3. Update your DNS as instructed by Railway

### 4. Test

```bash
# Test API
curl https://mkvmender.org/api/health

# Test Frontend
# Visit https://mkvmender.org in browser
```

## How It Works

The Go server now serves:
- **`/api/*`** - API endpoints
- **`/`** - Frontend static files (HTML, CSS, JS)

Everything runs in one service, one domain, no routing complexity!

## Advantages

✅ **Simpler**: One service instead of two
✅ **No CORS issues**: Same origin
✅ **Cheaper**: One Railway service
✅ **Faster**: No proxy/routing overhead
✅ **Easier maintenance**: One deployment

## Local Development

```bash
# Build and run
go build -o bin/mkvmender-server ./cmd/server
TURSO_DATABASE_URL="..." TURSO_AUTH_TOKEN="..." ./bin/mkvmender-server

# Frontend: http://localhost:8080/
# API: http://localhost:8080/api/health
```

## Removing the Separate Frontend Service

If you already deployed the frontend as a separate Railway service:

1. Go to your Railway project
2. Select the **frontend service**
3. Click **"Settings"** → **"Danger Zone"** → **"Remove Service"**

The API service now handles everything!
