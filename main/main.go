package main

import (
	"log"
	"time"

	raven "github.com/getsentry/raven-go"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/maddevsio/new-para-bot/bot"
	"github.com/maddevsio/new-para-bot/dce"
	"github.com/maddevsio/new-para-bot/utils"
)

// DCEChecker used in checkBinanceAndAlert
type DCEChecker interface {
	GetListOfActualPairs() (string, error)
	GetListOfSavedPairs() (string, error)
	UpdatePairs(pairs string) error
}

func main() {
	raven.CapturePanicAndWait(func() {
		do()
	}, nil)
}

func do() {
	// we can use db inside the container
	// because this is working table, no need
	// to have historical data
	dao := dce.NewDAO("/tmp/test.db")
	binance := dce.NewBinance(&dao)
	hibtc := dce.NewHibtc(&dao)
	liqui := dce.NewLiqui(&dao)
	ethfinex := dce.NewEthfinex(&dao)
	//kucoin := dce.NewKucoin(&dao)
	livecoin := dce.NewLivecoin(&dao)
	tidex := dce.NewTidex(&dao)
	okex := dce.NewOkex(&dao)
	huobi := dce.NewHuobi(&dao)
	kraken := dce.NewKraken(&dao)
	bitz := dce.NewBitz(&dao)
	wex := dce.NewWex(&dao)
	cex := dce.NewCex(&dao)
	for {
		log.Print("Checking...")
		checkDCEAndAlert(binance, binance.Name, binance.Website)
		checkDCEAndAlert(hibtc, hibtc.Name, hibtc.Website)
		checkDCEAndAlert(liqui, liqui.Name, liqui.Website)
		checkDCEAndAlert(ethfinex, ethfinex.Name, ethfinex.Website)
		//checkDCEAndAlert(kucoin, kucoin.Name)
		checkDCEAndAlert(livecoin, livecoin.Name, livecoin.Website)
		checkDCEAndAlert(tidex, tidex.Name, tidex.Website)
		checkDCEAndAlert(okex, okex.Name, okex.Website)
		checkDCEAndAlert(huobi, huobi.Name, huobi.Website)
		checkDCEAndAlert(kraken, kraken.Name, kraken.Website)
		checkDCEAndAlert(bitz, bitz.Name, bitz.Website)
		checkDCEAndAlert(wex, wex.Name, wex.Website)
		checkDCEAndAlert(cex, cex.Name, cex.Website)
		log.Print("Sleeping...")
		time.Sleep(60 * time.Second)
	}
}

func checkDCEAndAlert(dce DCEChecker, name string, url string) {
	// get actual pairs and check
	actualPairs, err := dce.GetListOfActualPairs()
	if err != nil {
		log.Panic(err)
	}

	savedPairs, err := dce.GetListOfSavedPairs()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%v: Pairs length: %v, %v", name, len(actualPairs), len(savedPairs))

	if actualPairs == "" {
		log.Printf("%v: actual pairth length is 0. seems did not get the data from API, skipping...", name)
		return
	}

	if savedPairs == "" {
		err = dce.UpdatePairs(actualPairs)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("%v: No saved data. Seems the first run", name)
	} else {
		utils.SaveNonEqualStringsToFiles(name, savedPairs, actualPairs)
		diff, err := utils.DiffSets(savedPairs, actualPairs)
		if err != nil {
			log.Panic(err)
		}
		if diff != "" {
			log.Printf("%v: We have diffs", name)
			err = dce.UpdatePairs(actualPairs)
			if err != nil {
				log.Panic(err)
			}
			log.Printf("%v: Pairs updated", name)
			config, err := bot.GetTelegramConfig("")
			if err != nil {
				log.Panic(err)
			}
			err = bot.SendMessageToTelegramChannel(config, name+"\n"+url+"\n"+diff)
			if err != nil {
				log.Panic(err)
			}
			log.Printf("%v: Bot message sent", name)
		} else {
			log.Printf("%v No diffs", name)
		}
	}
}
