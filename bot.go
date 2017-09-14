package telebotframe

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kovetskiy/lorg"
	"log"
)

type SendChannel chan tgbotapi.Chattable
type MsgFunc func(SendChannel, *tgbotapi.Message) error
type CallbackFunc func(SendChannel, *tgbotapi.CallbackQuery) error

type BotPlugin interface {
	Register(bot *TelegramBot) error
	Buttons() [][]string
}

//Зачем?
func NewTelegramBot() TelegramBot {
	return TelegramBot{}
}

type TelegramBot struct {
	API         *tgbotapi.BotAPI        // API телеграмма
	Updates     tgbotapi.UpdatesChannel // Канал обновлений
	SendChannel SendChannel             //Канал отправки сообщений
	commands    map[string]MsgFunc
	callbacks   map[string]CallbackFunc
	Keyboard    tgbotapi.ReplyKeyboardMarkup
}

func (telegramBot *TelegramBot) AddPlugins(plugins ...BotPlugin) {
	ReplyKeyboardRows := make([][]tgbotapi.KeyboardButton, 0)

	for _, plgn := range plugins {
		plgn.Register(telegramBot)

		btns := plgn.Buttons()
		for _, row := range btns {
			buttonRow := make([]tgbotapi.KeyboardButton, 0)
			for _, text := range row {
				buttonRow = append(buttonRow, tgbotapi.NewKeyboardButton(text))
			}
			ReplyKeyboardRows = append(ReplyKeyboardRows, buttonRow)
		}
	}

	telegramBot.Keyboard = tgbotapi.NewReplyKeyboard(ReplyKeyboardRows...)
}

func (telegramBot *TelegramBot) Init(token string, senders int, debug bool) {
	botAPI, err := tgbotapi.NewBotAPI(token) // Инициализация API
	if err != nil {
		log.Fatal(err)
	}
	telegramBot.API = botAPI

	lorg.Infof("Authorized on account %s", telegramBot.API.Self.UserName)

	telegramBot.API.Debug = debug

	botUpdate := tgbotapi.NewUpdate(0) // Инициализация канала обновлений
	botUpdate.Timeout = 25
	botUpdates, err := telegramBot.API.GetUpdatesChan(botUpdate)
	if err != nil {
		log.Fatal(err)
	}
	telegramBot.Updates = botUpdates

	telegramBot.SendChannel = make(chan tgbotapi.Chattable, senders)

	lorg.Infof("Initializing %d send workers", senders)
	for i := 0; i < senders; i++ {
		go func(id int, bot *TelegramBot) {
			for msg := range bot.SendChannel {
				bot.API.Send(msg)
			}
		}(i+1, telegramBot)
	}

	telegramBot.commands = make(map[string]MsgFunc)
	telegramBot.callbacks = make(map[string]CallbackFunc)

}

func (telegramBot *TelegramBot) Start() {
	for update := range telegramBot.Updates {
		// Если сообщение есть  -> начинаем обработку
		go telegramBot.analyzeUpdate(update)
	}
}

func (telegramBot *TelegramBot) analyzeUpdate(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		telegramBot.HandleCallback(update.CallbackQuery)
		return
	}

	// no message no cry
	if update.Message == nil {
		return
	}

	telegramBot.HandleMessage(update.Message)
}

func (telegramBot *TelegramBot) AddCommand(name string, f MsgFunc) {
	telegramBot.commands[name] = f
}

func (telegramBot *TelegramBot) AddCallback(name string, f CallbackFunc) {
	telegramBot.callbacks[name] = f
}

//Callbacks are called onetime only
func (telegramBot *TelegramBot) HandleCallback(callback *tgbotapi.CallbackQuery) {

	lorg.Infof("got calback %s", callback.Data)
	if com, ok := telegramBot.callbacks[callback.Data]; ok {
		//do something here
		com(telegramBot.SendChannel, callback)
		delete(telegramBot.callbacks, callback.Data)
		return
	}
	lorg.Error("Nothing to respond")
}

func (telegramBot *TelegramBot) HandleMessage(Message *tgbotapi.Message) {

	lorg.Infof("[u:%d:%s] %s",
		Message.From.ID,
		Message.From.UserName,
		Message.Text)

	if Message.IsCommand() {
		telegramBot.HandleCommand(Message)
		return
	}

	if com, ok := telegramBot.commands[Message.Text]; ok {
		com(telegramBot.SendChannel, Message)
		return
	}

	lorg.Error("Nothing to respond")
}

func (telegramBot *TelegramBot) HandleCommand(Message *tgbotapi.Message) error {
	//

	if !Message.IsCommand() {
		return errors.New("not a command")
	}
	msg := tgbotapi.NewMessage(Message.Chat.ID, "")

	switch Message.Command() {
	case "help":
		msg.Text = "type /start or /stop or /status."
	case "start":
		msg.ReplyMarkup = telegramBot.Keyboard
		msg.Text = "Добро пожаловать"
	case "stop":
		msg.Text = "Клавиатуру можно вернуть командой /start"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	case "status":
		msg.Text = "I'm ok."
	default:
		msg.Text = "Я не знаю этой команды =/"
	}

	telegramBot.SendChannel <- msg
	return nil
}
