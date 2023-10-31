package handlers

import (
	"github.com/MamushevArup/krisha-scraper/database/postgres"
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/MamushevArup/krisha-scraper/telegram/inline"
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
	b.bot.Debug = true
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	user := new(models.User)
	for update := range updates {
		b.HandleUpdate(&update, user, db)
	}

}

func (b *Bot) sendSecondInlineKeyboard(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Отлично теперь выбери одно из следущих")
	msg.ReplyMarkup = inline.TypeItem()
	b.sendMessage(&msg)
}
