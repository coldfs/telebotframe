# telebotframe
#### Telegram Bot Frame
 
 Install:
 ```bash
 go get github.com/coldfs/telebotframe
 ```
 
 
 Usage:
 ```go
package main

import (
	"github.com/coldfs/telebotframe"
)

func main() {

	telegramBot := telebotframe.NewTelegramBot()
	telegramBot.Init(BOT_API_KEY, SENDERS_COUNT, DEBUG)

	//"Plugins"
	telegramBot.AddPlugins(
		&telebotframe.SimplePlugin{},
	)

	telegramBot.Start()
}
```


Sample plugin:
```go
package plugins

import (
	"github.com/coldfs/telebotframe"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type TestPlugin struct {
}

func (tp *TestPlugin) GetName() string {
	return "Test Plugin"
}

func (tp *TestPlugin) Register(bot *telebotframe.TelegramBot) error {

    bot.Listen("/message/test", func(update tgbotapi.Update) error {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
        msg.Text = "Test ok"
        bot.SendChannel <- msg
        return nil
    })

    return nil
}

func (tp *TestPlugin) Buttons() [][]string {

	return [][]string{
		{"test", "nothing"},
		{"another row"}}
}
```

Without plugins:
 ```go
package main

import (
	"github.com/coldfs/telebotframe"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	telegramBot := telebotframe.NewTelegramBot()
	telegramBot.Init(BOT_API_KEY, SENDERS_COUNT, DEBUG)

    telegramBot.Listen("/message/test", func(update tgbotapi.Update) error {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
        msg.Text = "Test ok"
        telegramBot.SendChannel <- msg
        return nil
    })
    
    telegramBot.Listen("/command/start", func(update tgbotapi.Update) error {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
        msg.Text = "You started bot with /start command"
        telegramBot.SendChannel <- msg
        return nil
    })

	telegramBot.Start()
}
```

Event types:

```telegramBot.Listen("/command/start" ...```  - commands like /start /stop
 
```telegramBot.Listen("/message/text" ...```  - regular messages

```telegramBot.Listen("/callback/data" ...```  - callbacks for inmessage keyboard actions (executed once))