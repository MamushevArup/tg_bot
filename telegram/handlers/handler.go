package handlers

import (
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/MamushevArup/krisha-scraper/telegram/inline"
	"github.com/MamushevArup/krisha-scraper/utils/files"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (b *Bot) HandleUpdate(update *tgbotapi.Update, filter *models.Filter, sentSecondInlineKeyboard map[int64]bool) {
	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			b.handleMessageCommand(update)
		}
	} else if update.CallbackQuery != nil {
		b.handleCallbackQuery(update, filter, sentSecondInlineKeyboard)
	}
}

func (b *Bot) handleMessageCommand(update *tgbotapi.Update) {
	if update.Message != nil { // If we got a message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			handleCommand(&msg, update)
			b.sendMessage(&msg)
		}
	}
}

func handleCommand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		msg.Text = "Рад видеть тебя здесь.\nДля начала тебе нужно выбрать следущее"
		msg.ReplyMarkup = inline.BuyOrRent()
	case "help":
		val, err := files.ReadTXT("utils/texts/text.txt")
		if err != nil {
			log.Println("Error with something", err)
		}
		msg.Text = val
	case "city":
		val := update.Message.CommandArguments()
		msg.Text = "You are choose the city" + val
	case "run":

	}
}

func (b *Bot) handleCallbackQuery(update *tgbotapi.Update, filter *models.Filter, sentSecondInlineKeyboard map[int64]bool) {
	chatID := update.CallbackQuery.Message.Chat.ID
	if !sentSecondInlineKeyboard[chatID] {
		b.sendSecondInlineKeyboard(chatID)
		sentSecondInlineKeyboard[chatID] = true
		filter.BuyOrRent = update.CallbackQuery.Data
	}
	filter.TypeItem = update.CallbackQuery.Data
}
func (b *Bot) sendMessage(msg *tgbotapi.MessageConfig) {
	b.bot.Send(msg)
}
