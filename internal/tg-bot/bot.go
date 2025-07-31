package bot

import (
	"context"
	"log"
	"time"

	"github.com/artyomkorchagin/storeyourimages/internal/helpers"
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
