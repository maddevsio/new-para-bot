package main

import (
	"log"
	"time"

	raven "github.com/getsentry/raven-go"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/maddevsio/new-para-bot/bot"
	"github.com/maddevsio/new-para-bot/dce"
	"github.com/maddevsio/new-para-bot/ext"
	"github.com/maddevsio/new-para-bot/utils"
)

func main() {
	raven.CapturePanicAndWait(func() {
		doExt()
	}, nil)
}

func doExt() {

	for {
		log.Print("Checking...")
		// bitconnect
		// aex
		// exx
		// forkdelta
		// gdax
		// Kraken
		// LocalBitcoins
		// nlexch
		// Switcheo
		// TrustDex
		// ether_delta
		// Zaif
		// lakebtc
		dces := []string{
			"ZB",
			"Yobit", "Yunbi",
			"Waves", "Wex",
			"tux_exchange", "Upbit",
			"trade_ogre", "trade_satoshi",
			"therocktrading", "tidex",
			"SZZC", "stocks_exchange",
			"rightbtc", "south_xchange",
			"qryptos", "Quoine",
			"Paymium", "Poloniex",
			"novaexchange", "paribu",
			"Neraex", "nlexch",
			"Lykke", "Mercatox",
			"Luno", "Nanex",
			"litebiteu", "Livecoin",
			"Liqui", "lbank",
			"Latoken",
			"Jubi", "Koinex",
			"k_kex", "kyber_network",
			"idex", "infinity_coin",
			"hitbtc", "Huobi",
			"gopax", "Gemini",
			"Gate", "Gatecoin",
			"Extstock", "Fisco",
			"Ethfinex", "Exmo",
			"Cryptopia",
			"crypto_bridge", "crypto_hub",
			"crex24", "crxzone",
			"coins_markets", "COSS",
			"Coinone", "Coinrail",
			"Coinhouse", "Coinroom",
			"coin_exchange", "Coinbene",
			"Coinex", "Coinfalcon",
			"btc_alpha", "Cobinhood",
			"ccex", "cex", "bx_thailand",
			"Buyucoin", "btcc",
			"Abucoins", "ACX",
			"Bancor",
			"Bibox", "BigONE",
			"bit_z", "Binance",
			"Bitfinex",
			"Bitflyer", "Bithumb",
			"Bitmex", "bits_blockchain",
			"Bitso", "Bittrex",
			"Bleutrade",
		}

		for _, dce := range dces {
			checkDCEUsingCryptoexchangeAndAlert(dce)
			time.Sleep(2 * time.Second)
		}

		log.Print("Sleeping...")
		time.Sleep(60 * time.Second)
	}
}

func checkDCEUsingCryptoexchangeAndAlert(name string) {
	dao := dce.NewDAO("/tmp/test.db")
	dce := ext.NewCryptoexchange(&dao, name)
	log.Printf("%v: starting...", dce.Name)
	actualPairs, err := dce.GetListOfActualPairs()
	if err != nil {
		log.Panic(err)
	}

	savedPairs, err := dce.GetListOfSavedPairs()
	if err != nil {
		log.Panic(err)
	}

	if actualPairs == "" {
		log.Printf("%v: actual pairth length is 0. seems did not get the data from API, skipping...", dce.Name)
		dce = nil
		return
	}

	if savedPairs == "" {
		err = dce.UpdatePairsAndSave(actualPairs)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("%v: No saved data. Seems the first run", dce.Name)
	} else {
		utils.SaveNonEqualStringsToFiles(dce.Name, savedPairs, actualPairs)
		diff, err := utils.DiffSets(savedPairs, actualPairs)
		if err != nil {
			log.Panic(err)
		}
		if diff != "" {
			log.Printf("%v: We have diffs", dce.Name)
			err = dce.UpdatePairsAndSave(actualPairs)
			if err != nil {
				log.Panic(err)
			}
			log.Printf("%v: Pairs updated", dce.Name)
			config, err := bot.GetTelegramConfig("")
			if err != nil {
				log.Panic(err)
			}

			err = bot.SendMessageToTelegramChannel(config, dce.Name+"\n"+diff)
			if err != nil {
				log.Panic(err)
			}

			log.Printf("%v: Bot message sent", dce.Name)
		} else {
			log.Printf("%v No diffs", dce.Name)
		}
	}

	log.Printf("%v: Pairs length: %v, %v", dce.Name, len(actualPairs), len(savedPairs))
	dce = nil
}
