package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/artyomkorchagin/storeyourimages/internal/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleMediaMessage(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, svcs *AllServices) tgbotapi.MessageConfig {
	var fileID string
	var fileSize int
	var mimeType string
	var mediaType string
	var format string
	var successMessage string

	switch {
	case update.Message.Photo != nil:
		photo := update.Message.Photo[len(update.Message.Photo)-1]
		fileID = photo.FileID
		fileSize = photo.FileSize
		mediaType = types.Photo
		format = "jpg"
		successMessage = "Photo saved"
	case update.Message.Video != nil:
		video := update.Message.Video
		fileID = video.FileID
		fileSize = video.FileSize
		mimeType = video.MimeType
		mediaType = types.Video
		format = "mp4"
		if mimeType != "" {
			parts := strings.Split(mimeType, "/")
			if len(parts) == 2 {
				format = parts[1]
			}
		}
		successMessage = "Video saved"
	default:
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Unsupported media type")
	}

	if fileSize > MAX_FILE_SIZE {
		return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("File too large (max %dMB)", MAX_FILE_SIZE/1024/1024))
	}

	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		log.Printf("Failed to get file info: %v", err)
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process file")
	}

	if file.FileSize > MAX_FILE_SIZE {
		return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("File too large (max %dMB)", MAX_FILE_SIZE/1024/1024))
	}

	mediaURL, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		log.Printf("Failed to get media URL: %v", err)
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process file")
	}

	path, err := downloadContent(mediaURL, fileID, update.Message.Chat.ID, format)
	if err != nil {
		log.Printf("Failed to download media: %v", err)
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to save file")
	}

	wdr := types.NewWriteDataRequest(update.Message.Chat.ID, path, mediaType)
	if err := svcs.Content.WriteContent(ctx, wdr); err != nil {
		log.Printf("Failed to save media to database: %v", err)
		os.Remove(path)
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to save file")
	}

	log.Printf("Media saved successfully: %s", mediaType)
	return tgbotapi.NewMessage(update.Message.Chat.ID, successMessage)
}
