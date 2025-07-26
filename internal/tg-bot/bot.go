package bot

import (
	"fmt"
	"log"

	"github.com/artyomkorchagin/storeyourimages/internal/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InitBot() {
	bot, err := tgbotapi.NewBotAPI(helpers.GetEnv("BOT", ""))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	bot.Send(tgbotapi.NewSetMyCommands(botCommands...))
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var msg tgbotapi.MessageConfig
		if update.Message == nil {
			continue
		}
		fmt.Println(update.Message.Text)
		if !update.Message.IsCommand() {
			continue
		}
		msg = serveCommands(update)
		fmt.Println(msg)
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func serveCommands(upd tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
	switch upd.Message.Command() {
	case "start":
		msg.ReplyMarkup = newKeyboard
		msg.Text = "okay"
	default:
		msg.Text = "I don't know that command"
	}
	return msg
}
