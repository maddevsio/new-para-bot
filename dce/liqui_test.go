package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://api.liqui.io/api/3/info

func TestLiqui(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	liqui := NewLiqui(&dao)
	pairs, err := liqui.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "1st_btc") // check popular pairs
	assert.Contains(t, pairs, "zrx_usdt")
}
