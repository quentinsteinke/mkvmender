# 🎬 MKV Mender

**Stop manually renaming your media files. Let the community do it for you.**

MKV Mender is a community-driven tool that automatically names your movie and TV show files using crowd-sourced naming conventions. No more guessing—just hash your file and instantly see what everyone else is calling it.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/)

---

## 🤔 The Problem

You've just ripped your entire Blu-ray collection. Now what?

```
├── movie1.mkv
├── movie2.mkv
├── show_s01e01.mkv
├── show_s01e02.mkv
└── ...
```

Manually renaming hundreds of files is tedious. Different naming conventions make organization a nightmare. IMDB lookups by duration are unreliable. Filename parsing fails on edge cases.

## ✨ The Solution

MKV Mender uses **file hashing** to uniquely identify your files, then matches them against a **community database** of naming submissions.

```bash
$ mkvmender lookup movie1.mkv

Found 3 naming submissions:

[1] The Matrix (1999) 1080p BluRay.mkv
    👤 moviefan23  |  👍 156  👎 3

[2] Matrix.1999.1080p.BluRay.x264-SPARKS.mkv
    👤 encoder_pro  |  👍 89  👎 12

[3] The.Matrix.1999.REMASTERED.1080p.BluRay.mkv
    👤 quality_first  |  👍 45  👎 1
```

Pick your favorite, hit enter, and you're done.

---

## 🚀 Quick Start

### Installation

```bash
# Clone and build
git clone https://github.com/quentinsteinke/mkvmender.git
cd mkvmender
go build -o mkvmender ./cmd/cli

# Or use the Makefile
make build
```

### Basic Usage

```bash
# Register (one time)
mkvmender register

# Look up what the community calls your file
mkvmender lookup movie.mkv

# Rename it interactively
mkvmender rename movie.mkv

# Contribute your naming back to the community
mkvmender upload movie.mkv --name "Your Preferred Name.mkv"
```

---

## ✨ Features

### 🔍 **Smart Search**
Search for titles without having the file. Browse what's in the database before you rip.

```bash
$ mkvmender search "Breaking Bad"

[1] 📺 Breaking Bad (2008)
    → Season 1 (7 episodes)
    → Season 2 (13 episodes)
    → Season 3 (13 episodes)
```

Supports **fuzzy matching** too—typos won't stop you:
```bash
$ mkvmender search "Breking Bad"  # Still finds "Breaking Bad"
```

### 🗳️ **Community Voting**
The best names rise to the top through community votes.

```bash
$ mkvmender vote movie.mkv

[1] The Matrix (1999) 1080p BluRay.mkv  ← 👍 156  👎 3
[2] Matrix.1999.1080p.BluRay.x264.mkv   ← 👍 89   👎 12

Select [1-2] to vote:
```

### 📦 **Batch Processing**
Got hundreds of files? Process them all at once.

```bash
$ mkvmender batch /path/to/movies/

Processing 47 files...
✓ The Matrix (1999).mkv
✓ The Matrix Reloaded (2003).mkv
✓ The Matrix Revolutions (2003).mkv
...
```

### 🎯 **Hash-Based Matching**
Files are identified by SHA-256 hash—not filename, not duration, not file size. If the content matches, you'll get results.

### 📺 **TV Show Navigation**
Browse shows by season and episode:
```
Breaking Bad (2008)
  └─ Season 1
      ├─ Episode 1: Pilot
      ├─ Episode 2: Cat's in the Bag...
      └─ ...
```

### 🏷️ **Rich Metadata**
Store and display quality info, source (Blu-ray/DVD), release group, and more.

---

## 🎯 Use Cases

- **Media Hoarders**: Organize your massive movie/TV collection
- **Plex/Jellyfin Users**: Get Plex-friendly naming instantly
- **Release Groups**: Share standardized naming conventions
- **Archivists**: Preserve correct titles for rare releases
- **Seeders**: Help others identify what they're downloading

---

## 🏗️ How It Works

1. **Hash**: Generate SHA-256 hash of your file
2. **Lookup**: Query the community database
3. **Vote**: See which names the community prefers
4. **Apply**: Rename your file with one command
5. **Contribute**: Upload your naming to help others

All powered by:
- **Go** for performance
- **Turso** (distributed SQLite) for the database
- **RESTful API** for flexibility

---

## 📚 Documentation

- **[Full Documentation](DOCUMENTATION.md)** - Installation, API, configuration
- **[Quick Start Guide](QUICKSTART.md)** - Get up and running in 5 minutes
- **[Search Feature](SEARCH_FEATURE.md)** - Advanced search capabilities
- **[Project Structure](PROJECT_STRUCTURE.md)** - For contributors

---

## 🤝 Contributing

Found a bug? Want to add a feature? Contributions are welcome!

```bash
# Fork the repo, make your changes, then:
git add .
git commit -m "feat: Add awesome feature"
git push origin feature-branch
# Open a PR!
```

---

## 📜 License

MIT License - see [LICENSE](LICENSE) for details

---

## 🙏 Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - Powerful CLI framework
- [Turso](https://turso.tech) - Edge-native database
- [fuzzysearch](https://github.com/lithammer/fuzzysearch) - Fuzzy string matching

---

**Stop renaming files manually. Start mending with the community.** 🎬✨
