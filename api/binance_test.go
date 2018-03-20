package api

import (
	"testing"

	"gopkg.in/resty.v1"
)

func TestBinanceApi(t *testing.T) {
	// just try to figure out how Binance API works
	resp, err := resty.R().Get("https://api.binance.com/api/v1/LOOKING")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("\nResponse Body: %v", resp)
}
