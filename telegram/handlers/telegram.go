package handlers

import (
	"github.com/MamushevArup/krisha-scraper/database/postgres"
	"github.com/MamushevArup/krisha-scraper/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(b *tgbotapi.BotAPI) *Bot {
	return &Bot{bot: b}
}

func (b *Bot) Start(db *postgres.Sql) {
	log.Println("Start the application")
	//b.bot.Debug = true
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	user := new(models.User)
	var lastTwo []string
	for update := range updates {
		b.HandleUpdate(&update, user, db, &lastTwo)

	}

}
