package telebotframe

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type SendChannel chan tgbotapi.Chattable
type ListenFunc func(tgbotapi.Update) error

type BotPlugin interface {
	Register(bot *TelegramBot) error
	Buttons() [][]string
	GetName() string
}

//Зачем?
func NewTelegramBot() TelegramBot {
	return TelegramBot{}
}

type TelegramBot struct {
	API         *tgbotapi.BotAPI        // API телеграмма
	Updates     tgbotapi.UpdatesChannel // Канал обновлений
	SendChannel SendChannel             //Канал отправки сообщений
	listeners   map[string]ListenFunc
	plugins     []BotPlugin
	Verbose     bool
}

func (telegramBot *TelegramBot) AddPlugins(plugins ...BotPlugin) {

	telegramBot.plugins = plugins
	for _, plgn := range plugins {
		plgn.Register(telegramBot)
		if telegramBot.Verbose {
			log.Printf("Plugin '%s' inited", plgn.GetName())
		}
	}
}

func (telegramBot *TelegramBot) GetKeyboard() tgbotapi.ReplyKeyboardMarkup {
	ReplyKeyboardRows := make([][]tgbotapi.KeyboardButton, 0)
	for _, plgn := range telegramBot.plugins {
		btns := plgn.Buttons()
		for _, row := range btns {
			buttonRow := make([]tgbotapi.KeyboardButton, 0)
			for _, text := range row {
				buttonRow = append(buttonRow, tgbotapi.NewKeyboardButton(text))
			}
			ReplyKeyboardRows = append(ReplyKeyboardRows, buttonRow)
		}
	}

	return tgbotapi.NewReplyKeyboard(ReplyKeyboardRows...)
}

func (telegramBot *TelegramBot) Init(token string, senders int, debug bool) {
	botAPI, err := tgbotapi.NewBotAPI(token) // Инициализация API
	if err != nil {
		log.Fatal(err)
	}
	telegramBot.API = botAPI

	if telegramBot.Verbose {
		log.Printf("Authorized on account %s", telegramBot.API.Self.UserName)
	}

	telegramBot.API.Debug = debug

	botUpdate := tgbotapi.NewUpdate(0) // Инициализация канала обновлений
	botUpdate.Timeout = 25
	botUpdates, err := telegramBot.API.GetUpdatesChan(botUpdate)
	if err != nil {
		log.Fatal(err)
	}
	telegramBot.Updates = botUpdates

	telegramBot.SendChannel = make(chan tgbotapi.Chattable, senders)

	if telegramBot.Verbose {
		log.Printf("Initializing %d send workers", senders)
	}

	for i := 0; i < senders; i++ {
		go func(id int, bot *TelegramBot) {
			for msg := range bot.SendChannel {
				bot.API.Send(msg)
			}
		}(i+1, telegramBot)
	}

	telegramBot.listeners = make(map[string]ListenFunc)

	//Default command listeners
	telegramBot.Listen("/command/start", func(update tgbotapi.Update) error {
		//
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyMarkup = telegramBot.GetKeyboard()
		msg.Text = "Добро пожаловать!"
		telegramBot.SendChannel <- msg
		return nil
	})

	telegramBot.Listen("/command/stop", func(update tgbotapi.Update) error {
		//
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.Text = "Клавиатуру можно вернуть командой /start"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		telegramBot.SendChannel <- msg
		return nil
	})
}

func (telegramBot *TelegramBot) Start() {
	for update := range telegramBot.Updates {
		// Если сообщение есть  -> начинаем обработку
		telegramBot.analyzeUpdate(update)
	}
}

func (telegramBot *TelegramBot) analyzeUpdate(update tgbotapi.Update) {
	var updateType string
	var key string

	if update.CallbackQuery != nil {
		updateType = "callback"
		key = "/" + updateType + "/" + update.CallbackQuery.Data
	}

	if update.Message != nil {
		if update.Message.IsCommand() {
			updateType = "command"
			key = "/" + updateType + "/" + update.Message.Command()
		} else {
			updateType = "message"
			key = "/" + updateType + "/" + update.Message.Text
		}
	}

	if updateType == "" {
		if telegramBot.Verbose {
			log.Print("Unkmown message type")
		}
		return
	}

	if com, ok := telegramBot.listeners[key]; ok {

		com(update)
		if updateType == "callback" {
			delete(telegramBot.listeners, key)
		}
		return
	}

	//Checking wildcard
	if telegramBot.Verbose {
		log.Printf("Looking for %s", "/"+updateType+"/*")
	}
	if com, ok := telegramBot.listeners["/"+updateType+"/*"]; ok {
		com(update)
		return
	}
	if telegramBot.Verbose {
		log.Printf("Can't find listener for '%s'", key)
	}
}

func (telegramBot *TelegramBot) Listen(name string, clbc ListenFunc) error {

	if telegramBot.Verbose {
		log.Printf("Adding listener to %s", name)
	}
	telegramBot.listeners[name] = clbc
	return nil
}
