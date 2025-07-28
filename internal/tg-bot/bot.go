package bot

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/artyomkorchagin/storeyourimages/internal/helpers"
	"github.com/artyomkorchagin/storeyourimages/internal/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DEFAULT_TIMEOUT  = 60 * time.Second
	DOWNLOAD_TIMEOUT = 30 * time.Second
	MAX_FILE_SIZE    = 50 << 20 // 50 MB
)

func InitBot(svcs *AllServices) {
	token := helpers.GetEnv("BOT", "")
	if token == "" {
		log.Fatal("BOT environment variable is required")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	if _, err := bot.Request(tgbotapi.NewSetMyCommands(botCommands...)); err != nil {
		log.Printf("Failed to set bot commands: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = int(DEFAULT_TIMEOUT.Seconds())

	updates := bot.GetUpdatesChan(u)
	defer bot.StopReceivingUpdates()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go handleMessage(bot, update, svcs)
	}
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, svcs *AllServices) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var msg tgbotapi.MessageConfig

	switch {
	case update.Message.Photo != nil, update.Message.Video != nil:
		msg = handleMediaMessage(ctx, bot, update, svcs)
	case update.Message.IsCommand():
		msg = serveCommands(update, svcs)
	default:
		msg = serveButtons(update)
	}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

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

func serveCommands(upd tgbotapi.Update, svcs *AllServices) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")

	switch upd.Message.Command() {
	case "start":
		msg.ReplyMarkup = newKeyboard
		msg.Text = "Welcome! I can help you store your images and videos."

		if err := svcs.Users.CreateUser(context.Background(), upd.Message.Chat.ID); err != nil {
			log.Printf("Failed to create user: %v", err)
		}
	case "help":
		msg.Text = "Send me photos or videos and I'll save them for you!\n\n" +
			"Commands:\n" +
			"/start - Start the bot\n" +
			"/help - Show this help"
	default:
		msg.Text = "I don't know that command. Type /help for available commands."
	}

	return msg
}

func serveButtons(upd tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")

	switch upd.Message.Text {
	case "Просмотреть изображения":
		msg.Text = "Here you would see your images (feature coming soon)"
	default:
		msg.Text = "I don't understand. Try sending me a photo or video!"
	}

	return msg
}

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
