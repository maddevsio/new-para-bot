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
	cryptopia := dce.NewCryptopia(&dao)
	idex := dce.NewIdex(&dao)
	for {
		log.Print("Checking...")
		checkDCEAndAlert(binance)
		checkDCEAndAlert(hibtc)
		checkDCEAndAlert(liqui)
		checkDCEAndAlert(ethfinex)
		//checkDCEAndAlert(kucoin)
		checkDCEAndAlert(livecoin)
		checkDCEAndAlert(tidex)
		checkDCEAndAlert(okex)
		checkDCEAndAlert(huobi)
		checkDCEAndAlert(kraken)
		checkDCEAndAlert(bitz)
		checkDCEAndAlert(wex)
		checkDCEAndAlert(cex)
		checkDCEAndAlert(cryptopia)
		checkDCEAndAlert(idex)
		log.Print("Sleeping...")
		time.Sleep(60 * time.Second)
	}
}

func checkDCEAndAlert(dce dce.DCEChecker) {
	// get actual pairs and check
	actualPairs, err := dce.GetListOfActualPairs()
	if err != nil {
		log.Panic(err)
	}

	savedPairs, err := dce.GetDAO().GetListOfSavedPairs(dce)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%v: Pairs length: %v, %v", dce.GetName(), len(actualPairs), len(savedPairs))

	if actualPairs == "" {
		log.Printf("%v: actual pairth length is 0. seems did not get the data from API, skipping...", dce.GetName())
		return
	}

	if savedPairs == "" {
		err = dce.GetDAO().UpdatePairsAndSave(dce, actualPairs)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("%v: No saved data. Seems the first run", dce.GetName())
	} else {
		utils.SaveNonEqualStringsToFiles(dce.GetName(), savedPairs, actualPairs)
		diff, err := utils.DiffSets(savedPairs, actualPairs)
		if err != nil {
			log.Panic(err)
		}
		if diff != "" {
			log.Printf("%v: We have diffs", dce.GetName())
			err = dce.GetDAO().UpdatePairsAndSave(dce, actualPairs)
			if err != nil {
				log.Panic(err)
			}
			log.Printf("%v: Pairs updated", dce.GetName())
			config, err := bot.GetTelegramConfig("")
			if err != nil {
				log.Panic(err)
			}
			err = bot.SendMessageToTelegramChannel(config, dce.GetName()+"\n"+dce.GetWebsite()+"\n"+diff)
			if err != nil {
				log.Panic(err)
			}
			log.Printf("%v: Bot message sent", dce.GetName())
		} else {
			log.Printf("%v No diffs", dce.GetName())
		}
	}
}
