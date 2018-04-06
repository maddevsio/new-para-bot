package main

import (
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/maddevsio/new-para-bot/bot"
	"github.com/maddevsio/new-para-bot/dce"
)

func main() {
	for {
		log.Print("Checking...")
		checkBinanceAndAlert()
		log.Print("Sleeping...")
		time.Sleep(60 * time.Second)
	}
}

func checkBinanceAndAlert() {
	// we can use db inside the container
	// because this is working table, no need
	// to have historical data
	binance := dce.NewBinance()
	dao := dce.NewDAO("/tmp/test.db")

	// get actual pairs and check
	actualPairs, err := binance.GetListOfActualPairs()
	if err != nil {
		log.Panic(err)
	}

	savedPairs, err := binance.GetListOfSavedPairs(dao)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Pairs length: %v, %v", len(actualPairs), len(savedPairs))

	if savedPairs == "" {
		err = binance.UpdatePairs(dao, actualPairs)
		if err != nil {
			log.Panic(err)
		}
		log.Print("No saved data. Seems the first run")
	} else {
		// TODO: need to move DIFF to binance, but call dao diff func
		diff, err := dao.Diff(savedPairs, actualPairs)
		if err != nil {
			log.Panic(err)
		}
		if diff != "" {
			log.Print("We have diffs")
			err = binance.UpdatePairs(dao, actualPairs)
			if err != nil {
				log.Panic(err)
			}
			log.Print("Pairs updated")
			config, err := bot.GetTelegramConfig("")
			if err != nil {
				log.Panic(err)
			}
			err = bot.SendMessageToTelegramChannel(config, "Binance: "+diff)
			if err != nil {
				log.Panic(err)
			}
			log.Print("Bot message sent")
		} else {
			log.Print("No diffs")
		}
	}
}
