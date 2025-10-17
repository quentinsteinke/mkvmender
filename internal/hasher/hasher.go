package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// HashResult contains the hash and file size of a hashed file
type HashResult struct {
	Hash     string
	FileSize int64
}

// HashFile computes the SHA-256 hash of a file
func HashFile(filePath string) (*HashResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	fileSize := fileInfo.Size()

	// Compute SHA-256 hash
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, fmt.Errorf("failed to hash file: %w", err)
	}

	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return &HashResult{
		Hash:     hashString,
		FileSize: fileSize,
	}, nil
}

// HashFilePartial computes a hash based on the first and last chunks of a file
// This is faster for very large files, but less secure
// chunkSize is the number of bytes to read from start and end (e.g., 65536 for 64KB)
func HashFilePartial(filePath string, chunkSize int64) (*HashResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	fileSize := fileInfo.Size()

	hasher := sha256.New()

	// If file is smaller than 2*chunkSize, just hash the whole file
	if fileSize <= chunkSize*2 {
		if _, err := io.Copy(hasher, file); err != nil {
			return nil, fmt.Errorf("failed to hash file: %w", err)
		}
	} else {
		// Hash first chunk
		firstChunk := make([]byte, chunkSize)
		n, err := file.Read(firstChunk)
		if err != nil {
			return nil, fmt.Errorf("failed to read first chunk: %w", err)
		}
		hasher.Write(firstChunk[:n])

		// Hash last chunk
		if _, err := file.Seek(-chunkSize, io.SeekEnd); err != nil {
			return nil, fmt.Errorf("failed to seek to end: %w", err)
		}
		lastChunk := make([]byte, chunkSize)
		n, err = file.Read(lastChunk)
		if err != nil {
			return nil, fmt.Errorf("failed to read last chunk: %w", err)
		}
		hasher.Write(lastChunk[:n])

		// Also include file size in the hash to prevent collisions
		fmt.Fprintf(hasher, "%d", fileSize)
	}

	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return &HashResult{
		Hash:     hashString,
		FileSize: fileSize,
	}, nil
}

// FormatFileSize formats a file size in bytes to a human-readable string
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
