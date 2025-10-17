# Search Feature - Documentation

## Overview

The search feature allows users to discover movies and TV shows in the database by searching for titles, without needing to have the actual file. This is perfect for:
- Discovering what's already in the database
- Finding naming conventions before ripping your own media
- Browsing TV show episodes and seasons
- Seeing what the community has submitted

## Usage

```bash
mkvmender search "title"
```

## Interactive Workflow

### For Movies

1. **Search**: Enter a movie title
2. **View Results**: See all matching movies with year and media type
3. **Select**: Choose a movie from the list
4. **View Submissions**: See all naming submissions with:
   - Filename
   - Submitter
   - Vote counts
   - Quality and source information

### For TV Shows

1. **Search**: Enter a TV show title
2. **View Results**: See all matching shows with year
3. **Select Show**: Choose a show from the list
4. **View Seasons**: Browse available seasons
5. **Select Season**: Choose a season
6. **View Episodes**: Browse episodes in that season
7. **Select Episode**: Choose an episode
8. **View Submissions**: See all naming submissions for that episode

## Examples

### Movie Search

```bash
$ mkvmender search "The Matrix"

Searching for 'The Matrix'...

Found 1 result(s):

[1] 🎬 The Matrix (1999) - movie

Select a title (1-1) or 'q' to quit: 1

The Matrix (1999)
Hash: 30e14955ebf1352266dc2ff8067e68104607e750abb9d3b36582b8af909fcb58
Size: 1.0 MB

Found 2 naming submission(s):

[1] The Matrix (1999) 1080p BluRay.mkv
    Submitted by: testuser
    Votes: 2 (↑2 ↓0)

[2] Matrix.1999.1080p.BluRay.x264.mkv
    Submitted by: seconduser
    Votes: 1 (↑2 ↓1)
```

### TV Show Search

```bash
$ mkvmender search "Breaking Bad"

Searching for 'Breaking Bad'...

Found 1 result(s):

[1] 📺 Breaking Bad (2008) - tv

Select a title (1-1) or 'q' to quit: 1

Breaking Bad (2008) - TV Show

Found 1 season(s):

[1] Season 1 (1 episode(s))

Select a season (1-1) or 'q' to quit: 1

Breaking Bad - Season 1

Found 1 episode(s):

[1] Episode 1 (1 submission(s))

Select an episode (1-1) or 'q' to quit: 1

Breaking Bad - S01E1
Hash: 07854d2fef297a06ba81685e660c332de36d5d18d546927d30daad6d7fda1541
Size: 512.0 KB

Found 1 naming submission(s):

[1] Breaking Bad - S01E01.mkv
    Submitted by: testuser
    Votes: 0 (↑0 ↓0)
```

## Features

- **Smart Grouping**: Results are automatically grouped by title and year
- **Visual Icons**: 🎬 for movies, 📺 for TV shows
- **Hierarchical Navigation**: TV shows organized by seasons and episodes
- **Vote Information**: See community consensus on naming conventions
- **Metadata Display**: Quality, source, and other details shown when available
- **Easy Navigation**: Type 'q' at any point to quit

## Use Cases

### Before Ripping
Check if someone has already uploaded naming for the media you're about to rip:
```bash
mkvmender search "Inception"
```

### Discovering Conventions
See what naming conventions the community prefers:
```bash
mkvmender search "The Dark Knight"
```

### TV Show Planning
Browse all episodes of a show before starting to rip:
```bash
mkvmender search "Game of Thrones"
```

### Quality Comparison
Compare different rips and quality submissions:
```bash
mkvmender search "Blade Runner"
```

## Technical Details

- **Search Algorithm**: Case-insensitive partial match on title
- **Deduplication**: Results with same hash/season/episode are deduplicated
- **Sorting**: 
  - Movies: By title, then year (descending)
  - TV Shows: By season (ascending), then episode (ascending)
- **Performance**: Search uses indexed database queries for fast results

## Tips

1. **Partial Matches**: You can search with partial titles:
   - `mkvmender search "Matrix"` finds "The Matrix"
   
2. **Year Disambiguation**: If multiple versions exist, year helps differentiate:
   - Results show "(1999)" vs "(2021)" for remakes

3. **Quick Navigation**: Use number keys to quickly navigate through menus

4. **Quit Anytime**: Press 'q' at any prompt to exit search

## API Endpoint

The search feature uses the `/api/search?q=<query>` endpoint:

```bash
curl "http://localhost:8080/api/search?q=Matrix"
```

Response format:
```json
{
  "query": "Matrix",
  "results": [
    {
      "title": "The Matrix",
      "year": 1999,
      "media_type": "movie",
      "hash": "...",
      "file_size": 1048576,
      "submissions": [...]
    }
  ]
}
```

## Future Enhancements

- [ ] Search filters (year, quality, media type)
- [ ] Sort options (by votes, date, submitter)
- [ ] Fuzzy matching for misspellings
- [ ] Search by multiple titles at once
- [ ] Export search results to file
- [ ] Search history

The search feature makes discovering and browsing community submissions effortless!
