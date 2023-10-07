package main

import (
	"fmt"
	"github.com/MamushevArup/krisha-scraper/telegram/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	tgToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(tgToken)
	fmt.Println(bot)
	if err != nil {
		log.Panic(err)
	}
	init := handlers.NewBot(bot)
	init.Start()
}
