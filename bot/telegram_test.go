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

func TestFormateMessage(t *testing.T) {
	dceInfo := []string{"Test DCE"}
	diff := "+USDT-BTC\n+USDT-ETH\n"
	message := FormatMessage(dceInfo, diff)
	assert.Equal(t, "Test DCE \n+USDT-BTC\n+USDT-ETH\n", message)

	dceInfo = []string{"Test DCE", "http://testdce.stock"}
	diff = "+USDT-BTC\n+USDT-ETH\n"
	message = FormatMessage(dceInfo, diff)
	assert.Equal(t, "Test DCE http://testdce.stock\n+USDT-BTC\n+USDT-ETH\n", message)

	dceInfo = []string{"Test DCE", "http://testdce.stock", "http://testdce.stock/trade/#%v-%v"}
	diff = "+USDT-BTC\n"
	message = FormatMessage(dceInfo, diff)
	assert.Equal(t, "Test DCE http://testdce.stock\n+USDT-BTC http://testdce.stock/trade/#USDT-BTC\n", message)

	dceInfo = []string{"Test DCE", "http://testdce.stock", "http://testdce.stock/trade/%v-%v"}
	diff = "+USDT-BTC\n+USDT-ETH"
	message = FormatMessage(dceInfo, diff)
	assert.Equal(t, "Test DCE http://testdce.stock\n+USDT-BTC http://testdce.stock/trade/USDT-BTC\n+USDT-ETH http://testdce.stock/trade/USDT-ETH\n", message)

	config, err := GetTelegramConfig("../.env")
	assert.NoError(t, err)
	err = SendMessageToTelegramChannel(config, message)
	assert.NoError(t, err)
}

func TestBugBtcAlpha(t *testing.T) {
	dceInfo := []string{"btc_alpha", "https://btc-alpha.com/", "https://btc-alpha.com/exchange/%v_%v/"}
	diff := "+ SPD-ETH\n"
	message := FormatMessage(dceInfo, diff)
	assert.Equal(t, "btc_alpha https://btc-alpha.com/\n+ SPD-ETH https://btc-alpha.com/exchange/SPD_ETH/\n", message)
}
