# Deployment Guide

This guide explains how to deploy MKV Mender with the API backend and web frontend on Railway.

## Architecture

- **API Backend** (Go): Serves at `mkvmender.org/api/*`
- **Frontend** (Node.js): Serves at `mkvmender.org/`
- **Database**: Turso (distributed SQLite)

Both services run as separate Railway services.

---

## Prerequisites

1. Railway account ([railway.app](https://railway.app))
2. Turso database created
3. Domain `mkvmender.org` pointed to Railway

---

## Part 1: Deploy API Backend

### 1. Create API Service

1. Go to Railway dashboard
2. Click **"New Project"**
3. Select **"Deploy from GitHub repo"**
4. Choose **`quentinsteinke/mkvmender`**
5. Set **Root Directory** to `/` (root)

### 2. Configure Environment Variables

Add these variables in the **"Variables"** tab:

```
TURSO_DATABASE_URL=libsql://your-database.turso.io
TURSO_AUTH_TOKEN=your-auth-token
```

### 3. Configure Domain

1. Go to **"Settings"** → **"Networking"**
2. Add custom domain: **`mkvmender.org`**
3. Railway will provide DNS instructions

### 4. Verify Deployment

```bash
curl https://mkvmender.org/api/health
# Expected: {"status":"ok"}
```

---

## Part 2: Deploy Frontend

### 1. Create Frontend Service

**In the SAME Railway project** (to share the domain):

1. Click **"New Service"** (in the existing project)
2. Select **"Deploy from GitHub repo"**
3. Choose **`quentinsteinke/mkvmender`**
4. **IMPORTANT**: Set **Root Directory** to `/frontend`

### 2. No Environment Variables Needed

The frontend automatically uses `window.location.origin + '/api'` for API calls.

### 3. Configure Domain Routing

Railway should automatically route:
- `/` → Frontend service
- `/api/*` → API service

If not:
1. Go to project settings
2. Configure path-based routing:
   - Path `/` → Frontend service
   - Path `/api` → API service

### 4. Verify Deployment

Visit `https://mkvmender.org` in your browser. You should see the MKV Mender frontend.

---

## Testing the Full Stack

### 1. Register a User (CLI)

```bash
mkvmender register
# Save the API key provided
```

### 2. Login to Web Interface

1. Go to `https://mkvmender.org`
2. Enter your API key
3. Click "Login"

### 3. Search the Database

1. Enter a search query (e.g., "Matrix")
2. View results with submissions and votes

---

## Domain Configuration

### DNS Records

Point your domain to Railway:

1. Go to your domain registrar (e.g., Namecheap, Cloudflare)
2. Add these records:

```
Type: CNAME
Name: @
Value: [provided by Railway]

Type: CNAME
Name: www
Value: [provided by Railway]
```

Railway will provide the exact CNAME values in the domain settings.

### SSL Certificate

Railway automatically provisions SSL certificates via Let's Encrypt.

---

## Monitoring

### API Service Logs

1. Go to Railway project
2. Select API service
3. Click **"Deployments"** → View logs

### Frontend Service Logs

1. Go to Railway project
2. Select Frontend service
3. Click **"Deployments"** → View logs

### Health Checks

- API: `https://mkvmender.org/api/health`
- Frontend: `https://mkvmender.org` (should load the page)

---

## Troubleshooting

### Frontend Can't Connect to API

**Problem**: "Search failed" errors

**Solution**:
1. Check API is running: `curl https://mkvmender.org/api/health`
2. Check CORS configuration in API server
3. Verify domain routing is correct

### API Returns 401 Unauthorized

**Problem**: Login fails with valid API key

**Solution**:
1. Check API key is correct
2. Verify user exists in database
3. Check server logs for authentication errors

### Database Connection Fails

**Problem**: "Migration failed" or database errors

**Solution**:
1. Verify `TURSO_DATABASE_URL` is correct
2. Verify `TURSO_AUTH_TOKEN` is valid
3. Check Turso database is running:
   ```bash
   turso db show mkvmender
   ```

### Domain Not Working

**Problem**: Domain doesn't resolve to Railway

**Solution**:
1. Check DNS propagation: `dig mkvmender.org`
2. Wait up to 24 hours for DNS propagation
3. Verify CNAME records are correct

---

## Updating

### Update API

```bash
git add .
git commit -m "Update API"
git push origin main
```

Railway auto-deploys the API service.

### Update Frontend

```bash
git add frontend/
git commit -m "Update frontend"
git push origin main
```

Railway auto-deploys the frontend service.

---

## Cost Estimate

Railway pricing:
- **Hobby Plan**: $5/month (both services included)
- **Bandwidth**: ~$0.10 per GB
- **Build minutes**: Included

Turso pricing:
- **Starter**: Free (500 MB storage, 1B row reads/month)
- **Scaler**: $29/month (unlimited)

**Estimated monthly cost**: $5-10 for small to medium usage.

---

## Support

- Railway Docs: https://docs.railway.app
- Turso Docs: https://docs.turso.tech
- MKV Mender Issues: https://github.com/quentinsteinke/mkvmender/issues
