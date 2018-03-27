package main

import (
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/maddevsio/new-para-bot/bot"
	"github.com/maddevsio/new-para-bot/dce"
)

func main() {
	binance := dce.NewBinance()

	// get actual pairs and check
	actualPairs, err := binance.GetListOfActualPairs()
	if err != nil {
		log.Panic(err)
	}

	savedPairs, err := binance.GetListOfSavedPairs()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Pairs length: %v, %v", len(actualPairs), len(savedPairs))

	if savedPairs == "" {
		err = binance.UpdatePairs(actualPairs)
		if err != nil {
			log.Panic(err)
		}
		log.Print("No saved data. Seems the first run")
	} else {
		diff, err := binance.Diff(savedPairs, actualPairs)
		if err != nil {
			log.Panic(err)
		}
		if diff != "" {
			log.Print("We have diffs")
			err = binance.UpdatePairs(actualPairs)
			if err != nil {
				log.Panic(err)
			}
			log.Print("Pairs updated")
			config, err := bot.GetTelegramConfig("../.env")
			if err != nil {
				log.Panic(err)
			}
			err = bot.SendMessageToTelegramChannel(config, diff)
			if err != nil {
				log.Panic(err)
			}
			log.Print("Bot message sent")
		} else {
			log.Print("No diffs")
		}
	}
	os.Exit(0)
}
