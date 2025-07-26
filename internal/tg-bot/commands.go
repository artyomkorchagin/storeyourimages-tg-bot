package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var botCommands = []tgbotapi.BotCommand{
	{
		Command:     "start",
		Description: "Начать",
	},
}
