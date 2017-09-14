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

	//Инициализация "плагинов"
	telegramBot.AddPlugins(
		&telebotframe.SimplePlugin{},
	)

	telegramBot.Start()
}
```


Sample plugin:
```go
import (
	"github.com/coldfs/telebotframe"
)

type TestPlugin struct {
}

func (tp *TestPlugin) Register(bot *telebotframe.TelegramBot) error {

	bot.AddCommand("еще", func(sendChan telebotframe.SendChannel, Message *tgbotapi.Message) error {
		msg := tgbotapi.NewMessage(Message.Chat.ID, "")
		msg.Text = "Вот вам еще"
		sendChan <- msg
		return nil
	})

	return nil
}

func (tp *TestPlugin) Buttons() [][]string {

	return [][]string{
		{"еще", "ничего"},
		{"еще один ряд клавиатуры"}}
}
```

