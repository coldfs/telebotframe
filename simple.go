package telebotframe

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type SimplePlugin struct {
}

func (sp *SimplePlugin) GetName() string {
	return "Simple Plugin"
}

func (sp *SimplePlugin) Buttons() [][]string {
	return [][]string{{"команда1", "команда2"}}
}

func (sp *SimplePlugin) Register(bot *TelegramBot) error {

	bot.Listen("/message/команда1", func(update tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.Text = "Пытаемся что-то выполнить"
		bot.SendChannel <- msg
		return nil
	})

	bot.Listen("/message/команда2", func(update tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.Text = "Странно, но ничего не происходит"
		bot.SendChannel <- msg

		msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg2.Text = "Попробуйте другую комманду"
		bot.SendChannel <- msg2
		return nil
	})

	return nil
}
