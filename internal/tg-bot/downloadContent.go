package bot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func downloadContent(url, fileID string, userID int64, format string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("url is empty")
	}

	timestamp := time.Now().Unix()
	path := filepath.Join("uploads", fmt.Sprintf("%d", userID), fmt.Sprintf("%s_%d.%s", fileID[:35], timestamp, format))

	client := &http.Client{
		Timeout: DOWNLOAD_TIMEOUT,
	}

	response, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad response status: %s", response.Status)
	}

	if response.ContentLength > MAX_FILE_SIZE {
		return "", fmt.Errorf("file too large: %d bytes", response.ContentLength)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Failed to close file: %v", closeErr)
		}
	}()

	limitedReader := io.LimitReader(response.Body, MAX_FILE_SIZE+1)
	if _, err := io.Copy(file, limitedReader); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	if fileInfo, err := file.Stat(); err == nil && fileInfo.Size() > MAX_FILE_SIZE {
		os.Remove(path)
		return "", fmt.Errorf("downloaded file too large: %d bytes", fileInfo.Size())
	}

	log.Printf("Successfully downloaded file to: %s", path)
	return path, nil
}
