package telebotframe

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kovetskiy/lorg"
	"telesales/tbot"
)

type SimplePlugin struct {
}

func (sp *SimplePlugin) Register(bot *tbot.TelegramBot) error {

	lorg.Info("Simple plugin added")

	bot.AddCommand("команда1", func(sendChan tbot.SendChannel, Message *tgbotapi.Message) error {
		msg := tgbotapi.NewMessage(Message.Chat.ID, "")
		msg.Text = "Пытаемся что-то выполнить"
		sendChan <- msg
		return nil
	})

	bot.AddCommand("команда2", func(sendChan tbot.SendChannel, Message *tgbotapi.Message) error {
		msg := tgbotapi.NewMessage(Message.Chat.ID, "")
		msg.Text = "Странно, но ничего не происходит"
		sendChan <- msg

		msg2 := tgbotapi.NewMessage(Message.Chat.ID, "")
		msg2.Text = "Попробуйте другую комманду"
		sendChan <- msg2
		return nil
	})

	return nil
}

func (sp *SimplePlugin) Buttons() [][]string {
	return [][]string{{"команда1", "команда2"}}
}
