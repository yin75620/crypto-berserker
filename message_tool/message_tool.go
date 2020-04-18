package message_tool

import (
	"log"

	"github.com/alexcesaro/mail/gomail"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yin75620/crypto-berserker/setting"
)

/// message
func sendMail(content string) {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", "yin75620@gmail.com", "Golang")
	msg.SetHeader("To", "yin75620@gmail.com")
	msg.AddHeader("To", "yin75620@gmail.com")
	msg.SetHeader("Subject", "Hello!")
	msg.SetBody("text/plain", "Hello Has Profit")
	msg.AddAlternative("text/html", content)

	m := gomail.NewMailer("smtp.gmail.com", "yin75620", setting.GMAIL_PASSWORD, 25)
	if err := m.Send(msg); err != nil {
		log.Println(err)
	}
}

const (
	PAUSE_SEND_TELEGRAM = false
)

//get chatId, https://api.telegram.org/bot***REMOVED***/getUpdates

var bot *tgbotapi.BotAPI
var broadcasterBot *tgbotapi.BotAPI

const BroadcasterToken = "***REMOVED***"

func StartTelegram() {
	bot, _ = tgbotapi.NewBotAPI(setting.TELEGRAM_BOT_TOKEN)
	broadcasterBot, _ = tgbotapi.NewBotAPI(BroadcasterToken)
}

func SendTelegram(content string) {
	sendTo(bot, 945156610, content)
}

func SendBroadcastArcherGroup(content string) {
	sendTo(broadcasterBot, -358468104, content)
}

func SendWatcherGroup(content string) {
	sendTo(bot, -361312212, content)
}

func sendTo(bo *tgbotapi.BotAPI, groupId int64, content string) {
	if PAUSE_SEND_TELEGRAM {
		return
	}
	msg := tgbotapi.NewMessage(groupId, content)
	bo.Send(msg)
}
