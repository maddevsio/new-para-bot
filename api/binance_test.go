package api

import (
	"testing"

	"github.com/buger/jsonparser"
	"gopkg.in/resty.v1"
)

func TestBinanceApi(t *testing.T) {
	// iterate throught all active pairs
	resp, err := resty.R().Get("https://api.binance.com/api/v1/exchangeInfo")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pair, _ := jsonparser.GetString(value, "symbol")
		t.Log(pair)
	}, "symbols")
}
