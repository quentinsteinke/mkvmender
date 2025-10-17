# Railway Deployment Guide

This guide walks you through deploying the MKV Mender API server on Railway.

## Prerequisites

- A Railway account (sign up at [railway.app](https://railway.app))
- A Turso database (or you can create one during setup)

## Step-by-Step Setup

### 1. Connect Your GitHub Repository

1. Go to [railway.app](https://railway.app) and log in
2. Click **"New Project"**
3. Select **"Deploy from GitHub repo"**
4. Choose **`quentinsteinke/mkvmender`**

### 2. Configure Build Settings

Railway should auto-detect your Go project, but verify:

- **Build Command**: `make build-server` (or Railway will auto-detect)
- **Start Command**: `./bin/mkvmender-server`
- **Root Directory**: `/` (leave blank)

### 3. Set Environment Variables

Click on your deployed service, then go to **"Variables"** tab and add:

#### Required Variables:

```bash
TURSO_DATABASE_URL=libsql://your-database.turso.io
TURSO_AUTH_TOKEN=your-auth-token-here
PORT=8080
```

#### Get Your Turso Credentials:

```bash
# Create database
turso db create mkvmender

# Get database URL
turso db show mkvmender

# Create auth token
turso db tokens create mkvmender
```

Copy the URL and token to Railway's environment variables.

### 4. Deploy

1. Railway will automatically deploy when you push to `main`
2. Or click **"Deploy"** manually in the Railway dashboard

### 5. Get Your Public URL

1. Go to **"Settings"** tab in Railway
2. Under **"Networking"**, click **"Generate Domain"**
3. Your API will be available at: `https://your-app.up.railway.app`

### 6. Test Your Deployment

```bash
# Check health
curl https://your-app.up.railway.app/api/health

# Expected response:
# {"status":"ok"}
```

## Configure CLI to Use Railway API

Update your local CLI to point to Railway:

```bash
# Option 1: Using environment variable
export MKVMENDER_API_URL=https://your-app.up.railway.app

# Option 2: Update config file
# Edit ~/.mkvmender/config.yaml
base_url: https://your-app.up.railway.app
```

## Automatic Deployments

Railway automatically deploys when you push to GitHub:

```bash
git add .
git commit -m "Update server"
git push origin main
```

Railway will:
1. Pull latest code
2. Build the server
3. Run migrations
4. Start the service
5. Zero-downtime deployment

## Monitoring

### View Logs

1. Go to your Railway dashboard
2. Click on your service
3. Click **"Deployments"** tab
4. Click on latest deployment
5. View real-time logs

### Metrics

Railway provides:
- CPU usage
- Memory usage
- Network traffic
- Request metrics

## Troubleshooting

### Build Fails

Check that your Go version is compatible:
- Railway uses Go 1.21+ by default
- Your project requires Go 1.24+

### Server Won't Start

1. Check environment variables are set correctly
2. Verify Turso database is accessible
3. Check logs for errors: `Migration failed` usually means database issues

### Can't Connect to API

1. Verify domain is generated in Settings → Networking
2. Check firewall/security settings
3. Test health endpoint: `/api/health`

## Cost

Railway pricing:
- **Free Tier**: $5 credit/month (good for testing)
- **Developer Plan**: $5/month + usage
- **Team Plan**: $20/month + usage

Your API server should cost ~$1-5/month depending on usage.

## Environment Variables Reference

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `TURSO_DATABASE_URL` | Yes | Turso database URL | `libsql://db.turso.io` |
| `TURSO_AUTH_TOKEN` | Yes | Turso auth token | `eyJhbGci...` |
| `PORT` | No | Server port (Railway sets this) | `8080` |
| `MIGRATION_PATH` | No | Migration file path | `migrations/001_initial_schema.sql` |

## Next Steps

1. Update README.md with your Railway URL
2. Configure CLI clients to use Railway API
3. Set up custom domain (optional)
4. Enable auto-scaling (optional)

## Support

- Railway Docs: https://docs.railway.app
- Turso Docs: https://docs.turso.tech
- MKV Mender Issues: https://github.com/quentinsteinke/mkvmender/issues
