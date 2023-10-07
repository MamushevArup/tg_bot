package handlers

import (
	"fmt"
	"github.com/MamushevArup/krisha-scraper/krisha/scrap"
	"github.com/MamushevArup/krisha-scraper/models"
	"github.com/MamushevArup/krisha-scraper/telegram/inline"
	"github.com/MamushevArup/krisha-scraper/utils/files"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
	"log"
)

func (b *Bot) HandleUpdate(update *tgbotapi.Update, filter *models.Filter, sentSecondInlineKeyboard map[int64]bool) {
	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			b.handleMessageCommand(update, filter)
		}
	} else if update.CallbackQuery != nil {
		b.handleCallbackQuery(update, filter, sentSecondInlineKeyboard)
	}
}

func (b *Bot) handleMessageCommand(update *tgbotapi.Update, filter *models.Filter) {
	if update.Message != nil { // If we got a message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			handleCommand(&msg, update, filter)
			b.sendMessage(&msg)
		}
	}
}

func handleCommand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, filter *models.Filter) {
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
		// Here I need to implement the logic of taking data from the database for the user and
		// check is he/she fill the data that is required by a start command. Do it later
		if filter.TypeItem == "" && filter.BuyOrRent == "" {
			msg.Text = "Пожалуйста заполните необходимую информацию по команде /start"
			return
		}
		c := colly.NewCollector()
		init := scrap.New(c, filter)
		houses := init.NewScrap()
		fmt.Println(houses)
		msg.Text = houses
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
