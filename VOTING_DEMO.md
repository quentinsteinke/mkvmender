# Interactive Voting - User Experience Improvement

## Problem (Before)
Users had to manually note submission IDs and run separate commands:

```bash
# Step 1: Look up file to see submissions
$ mkvmender lookup movie.mkv
Found 3 naming options:
[1] The.Matrix.1999.mkv (ID: 123)
[2] Matrix.1999.1080p.mkv (ID: 124)
[3] The Matrix (1999).mkv (ID: 125)

# Step 2: Remember/copy ID and vote separately
$ mkvmender vote 124 up
Vote recorded
```

**Issues:**
- ❌ Tedious multi-step process
- ❌ Must remember or copy submission IDs
- ❌ No immediate feedback on ranking changes
- ❌ Not user-friendly

## Solution (After)
Single interactive command with guided workflow:

```bash
$ mkvmender vote movie.mkv

Hashing file...
Looking up naming options...

Found 2 naming option(s):

[1] Matrix.1999.1080p.BluRay.x264.mkv
    Submitted by: seconduser
    Votes: 1 (↑2 ↓1)
    Title: The Matrix (1999)

[2] The Matrix (1999) 1080p BluRay.mkv
    Submitted by: testuser
    Votes: 1 (↑1 ↓0)
    Title: The Matrix (1999)

Select a submission to vote on (1-2) or 'q' to quit: 2

Upvote or downvote? (up/down): up

✓ Successfully upvoted: The Matrix (1999) 1080p BluRay.mkv

Fetching updated vote counts...

Updated rankings:
→ [1] The Matrix (1999) 1080p BluRay.mkv - Votes: 2 (↑2 ↓0)
  [2] Matrix.1999.1080p.BluRay.x264.mkv - Votes: 1 (↑2 ↓1)
```

**Benefits:**
- ✅ Single command workflow
- ✅ Interactive selection (no ID memorization)
- ✅ Shows all context (votes, metadata, submitter)
- ✅ Immediate feedback with updated rankings
- ✅ Visual indicator (→) shows voted item
- ✅ User-friendly prompts

## Key Features

1. **Automatic file hashing** - Just provide the file path
2. **Interactive selection** - Choose from a list, not IDs
3. **Clear context** - See votes, metadata, and submitters
4. **Immediate feedback** - Updated rankings shown after voting
5. **Visual cues** - Arrow indicator shows your voted submission
6. **Easy to quit** - Press 'q' to cancel anytime

## Real-World Usage

```bash
# Community member finds a file with multiple naming options
$ mkvmender vote "movie.mkv"

# Sees 5 different naming conventions
# Votes for the one that follows their preferred style
# Immediately sees it move up in rankings
# Other users benefit from the community consensus
```

This makes voting accessible and encourages community participation!
