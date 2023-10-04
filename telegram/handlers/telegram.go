package handlers

import (
	"github.com/MamushevArup/krisha-scraper/krisha/scrap"
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/MamushevArup/krisha-scraper/telegram/inline"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	c   *scrap.Krisha
}

func NewBot(b *tgbotapi.BotAPI, initKrisha *scrap.Krisha) *Bot {
	return &Bot{bot: b, c: initKrisha}
}

func (b *Bot) Start() {
	log.Println("Start the application")
	b.bot.Debug = true
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	sentSecondInlineKeyboard := make(map[int64]bool)
	filter := new(models.Filter)
	for update := range updates {
		b.HandleUpdate(&update, filter, sentSecondInlineKeyboard)
	}
}

func (b *Bot) sendSecondInlineKeyboard(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Теперь выбери один из этих вариантов")
	msg.ReplyMarkup = inline.TypeItem()
	b.sendMessage(&msg)
}
