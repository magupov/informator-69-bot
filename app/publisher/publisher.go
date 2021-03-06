package publisher

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"informator-69-bot/app/storage"
	"informator-69-bot/app/wiki"
	"log"
	"time"
)

//func Notifier(store *storage.ActiveChatsStore, bot *tgbotapi.BotAPI) {
//	location, _ := time.LoadLocation("Europe/Moscow")
//	log.Print(location)
//	now := time.Now().In(location)
//	todayNotif := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, location)
//	var nextNotif = todayNotif
//	if todayNotif.Before(now) {
//		nextNotif = nextNotif.AddDate(0, 0, 1)
//	}
//
//	for {
//		timer := time.NewTimer(time.Until(nextNotif))
//		<-timer.C
//		report := wiki.GetTodaysReport()
//		store.RLock()
//		for _, chat := range store.Cache {
//			SendMessage(chat.Id, "Доброе утро", bot)
//			SendMessage(chat.Id, report, bot)
//		}
//		store.RUnlock()
//		nextNotif = nextNotif.AddDate(0, 0, 1)
//	}
//}

func Notifier(store *storage.ActiveChatsStore, holidays *wiki.Holidays, bot *tgbotapi.BotAPI) {
	location, _ := time.LoadLocation("Europe/Moscow")
	log.Print(location)
	now := time.Now().In(location)
	todayNotif := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, location)
	var nextNotif = todayNotif
	if todayNotif.Before(now) {
		nextNotif = nextNotif.AddDate(0, 0, 1)
	}

	for {
		timer := time.NewTimer(time.Until(nextNotif))
		<-timer.C
		report := wiki.GetTodaysReport(holidays)
		store.RLock()
		for _, chat := range store.Cache {
			SendMessage(chat.Id, "Доброе утро", bot)
			SendMessage(chat.Id, report, bot)
		}
		store.RUnlock()
		nextNotif = nextNotif.AddDate(0, 0, 1)
	}
}

func SendMessage(chatId int64, report string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatId, report)
	msg.ParseMode = "markdown"
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}
