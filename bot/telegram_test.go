package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessageToTelegramChannel(t *testing.T) {
	config, err := GetTelegramConfig("../.env")
	assert.NoError(t, err)
	err = SendMessageToTelegramChannel(config, "Test Message")
	assert.NoError(t, err)
}

func TestFormatNoDCELink(t *testing.T) {
	dceInfo := []string{"Test DCE"}
	diff := "+USDT-BTC\n+USDT-ETH\n"
	message := FormatMessage(dceInfo, diff)
	assert.Equal(t, "[Test DCE]()* has pairs alerts:*\n\n+USDT-BTC\n+USDT-ETH\n", message)
}

func TestFormatDCELink(t *testing.T) {
	dceInfo := []string{"Test DCE", "http://testdce.stock"}
	diff := "+USDT-BTC\n+USDT-ETH\n"
	message := FormatMessage(dceInfo, diff)
	assert.Equal(t, "[Test DCE](http://testdce.stock)* has pairs alerts:*\n\n+USDT-BTC\n+USDT-ETH\n", message)
}

func TestFormatOnePairLink(t *testing.T) {
	dceInfo := []string{"Test DCE", "http://testdce.stock", "http://testdce.stock/trade/#%v-%v"}
	diff := "+USDT-BTC\n"
	message := FormatMessage(dceInfo, diff)
	assert.Equal(t, "[Test DCE](http://testdce.stock)* has pairs alerts:*\n\n[+USDT-BTC](http://testdce.stock/trade/#USDT-BTC)\n", message)
}

func TestFormatTwoPairsLink(t *testing.T) {
	dceInfo := []string{"Test DCE", "http://testdce.stock", "http://testdce.stock/trade/%v-%v"}
	diff := "+USDT-BTC\n+USDT-ETH"
	message := FormatMessage(dceInfo, diff)
	assert.Equal(t, "[Test DCE](http://testdce.stock)* has pairs alerts:*\n\n[+USDT-BTC](http://testdce.stock/trade/USDT-BTC)\n[+USDT-ETH](http://testdce.stock/trade/USDT-ETH)\n", message)
}

func TestBugBtcAlpha(t *testing.T) {
	dceInfo := []string{"btc_alpha", "https://btc-alpha.com/", "https://btc-alpha.com/exchange/%v_%v/"}
	diff := "+ SPD-ETH\n"
	message := FormatMessage(dceInfo, diff)
	assert.Equal(t, "[btc_alpha](https://btc-alpha.com/)* has pairs alerts:*\n\n[+ SPD-ETH](https://btc-alpha.com/exchange/SPD_ETH/)\n", message)
}
