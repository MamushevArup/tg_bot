package main

import (
	"github.com/MamushevArup/krisha-scraper/krisha/scrap"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	scrap.Scrap()
}
