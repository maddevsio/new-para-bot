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
		// idex
		dces := [][]string{
			{"ZB"},
			{"Yobit", "https://yobit.io/", "https://yobit.io/en/trade/%v/%v"},
			{"Yunbi"},
			{"Waves"},
			{"Wex"},
			{"tux_exchange"},
			{"Upbit"},
			{"trade_ogre", "https://tradeogre.com/", "https://tradeogre.com/exchange/%v-%v"},
			{"trade_satoshi", "https://tradesatoshi.com/", "https://tradesatoshi.com/Exchange/?market=%v_%v"},
			{"therocktrading"},
			{"tidex"},
			{"SZZC"},
			{"stocks_exchange", "https://stocks.exchange/", "https://stocks.exchange/trade/%v/%v"},
			{"rightbtc"},
			{"south_xchange"},
			{"qryptos"},
			{"Quoine"},
			{"Paymium"},
			{"Poloniex", "https://poloniex.com/", "https://poloniex.com/exchange#%v_%v"},
			{"novaexchange"},
			{"paribu"},
			{"Neraex"},
			{"nlexch"},
			{"Lykke", "https://www.lykke.com/"},
			{"Mercatox", "https://mercatox.com/", "https://mercatox.com/exchange/%v/%v"},
			{"Luno"},
			{"Nanex"},
			{"litebiteu"},
			{"Livecoin", "https://www.livecoin.net/"},
			{"Liqui"},
			{"lbank"},
			{"Latoken", "https://latoken.com/"},
			{"Jubi"},
			{"Koinex"},
			{"k_kex"},
			{"kyber_network"},
			{"infinity_coin"},
			{"hitbtc", "https://hitbtc.com/", "https://hitbtc.com/%v-to-%v"},
			{"Huobi"},
			{"gopax"},
			{"Gemini", "https://gemini.com/"},
			{"Gate"},
			{"Gatecoin", "https://gatecoin.com/"}, // https://gatecoin.com/markets/btcusd lowercase
			{"Extstock"},
			{"Fisco"},
			{"Ethfinex"},
			{"Exmo"},
			{"Cryptopia"},
			{"crypto_bridge", "https://crypto-bridge.org/", "https://wallet.crypto-bridge.org/market/BRIDGE.%v_BRIDGE.%v"},
			{"crypto_hub"},
			{"crex24", "https://crex24.com/", "https://crex24.com/exchange/%v-%v"},
			{"crxzone"},
			{"coins_markets"},
			{"COSS", "https://exchange.coss.io/"}, // https://exchange.coss.io/exchange/coss-eth lowecase
			{"Coinone"},
			{"Coinrail", "https://coinrail.co.kr/"},
			{"Coinhouse"},
			{"Coinroom"},
			{"coin_exchange", "https://www.coinexchange.io/", "https://www.coinexchange.io/market/%v/%v"},
			{"Coinbene", "https://www.coinbene.com/"},
			{"Coinex", "https://www.coinex.com/"}, // https://www.coinex.com/trading?currency=bch&dest=xmr#limit lowecase
			{"Coinfalcon"},
			{"btc_alpha", "https://btc-alpha.com/", "https://btc-alpha.com/exchange/%v_%v/"},
			{"Cobinhood"},
			{"ccex"},
			{"cex"},
			{"bx_thailand"},
			{"Buyucoin"},
			{"btcc"},
			{"Abucoins"},
			{"ACX"},
			{"Bancor", "https://www.bancor.network/"},
			{"Bibox", "https://www.bibox.com/"},
			{"BigONE"},
			{"bit_z", "https://www.bit-z.com/"},
			{"Binance"},
			{"Bitfinex"},
			{"Bitflyer"},
			{"Bithumb", "https://www.bithumb.com/"},
			{"Bitmex"},
			{"bits_blockchain"},
			{"Bitso", "https://bitso.com/"},
			{"Bittrex"},
			{"Bleutrade", "https://bleutrade.com/"},
		}

		for _, dce := range dces {
			checkDCEUsingCryptoexchangeAndAlert(dce)
			time.Sleep(2 * time.Second)
		}

		log.Print("Sleeping...")
		time.Sleep(60 * time.Second)
	}
}

func checkDCEUsingCryptoexchangeAndAlert(dceInfo []string) {
	name := dceInfo[0]
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

			err = bot.SendMessageToTelegramChannel(config, bot.FormatMessage(dceInfo, diff))
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
