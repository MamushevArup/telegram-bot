package main

import (
	"github.com/MamushevArup/telegram-bot-krisha/internal/texts"
	"github.com/MamushevArup/telegram-bot-krisha/scraping"
	"github.com/MamushevArup/telegram-bot-krisha/telegram/config"
	"github.com/MamushevArup/telegram-bot-krisha/telegram/database"
	"github.com/MamushevArup/telegram-bot-krisha/telegram/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	open, err := utils.OpenJsonFile(texts.PathToJSON)
	bot, err := tgbotapi.NewBotAPI(open)
	if err != nil {
		log.Fatal("Error", err.Error())
	}
	bot.Debug = true

	var startBot *config.Bot
	startBot = config.New(bot)
	val, err := startBot.Start()
	config.NewInfo(val.Username, val.Region, val.StartPrice, val.EndPrice, val.IsOwner, val.URL)
	if err != nil {
		log.Fatal("Err while starting the bot", err)
	}
	var s *database.Mongo
	//var f *database.Filter
	value := s.Connection(val.Username, val.Region, val.StartPrice, val.EndPrice, val.IsOwner, val.URL)
	if value != nil {
		log.Fatal("Error with dividing by part")
		return
	}
	//var f scraping.Houses
	_, err = database.FindByFilter(val.Username)
	if err != nil {
		log.Println("Error 909090909", err)
	}
	scraping.ScrapKrisha(bot, val)
}
