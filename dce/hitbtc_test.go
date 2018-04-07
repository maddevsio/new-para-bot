package dce

import (
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	resty "gopkg.in/resty.v1"
)

func TestHitbtc(t *testing.T) {
	resp, err := resty.R().Get("https://api.hitbtc.com/api/2/public/symbol")
	assert.NoError(t, err)
	var pairs string
	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pair, _ := jsonparser.GetString(value, "id")
		pairs += pair + "\n"
	})

	assert.Contains(t, pairs, "BTCUSD") // check popular pairs
	assert.Contains(t, pairs, "XDNETH")
}
