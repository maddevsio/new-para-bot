package dce

import (
	"testing"

	"github.com/maddevsio/new-para-bot/utils"
	"github.com/stretchr/testify/assert"
)

func TestHitbtc(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	hibtc := NewHibtc(&dao)
	actualPairs, err := hibtc.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, actualPairs, "BTCUSD") // check popular pairs
	assert.Contains(t, actualPairs, "XDNETH")

	// delete all data before main logic
	// this deletion only for testing, we need to empty data set
	err = dao.DeleteAll(hibtc)
	assert.NoError(t, err)
	count, err := dao.Count(hibtc)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// update binance struct with newly get pairs from API
	err = hibtc.UpdatePairs(actualPairs)
	assert.NoError(t, err)

	// simulate the case when we are getting the data from storage
	savedPairs, err := hibtc.GetListOfSavedPairs()
	assert.NoError(t, err)
	assert.Equal(t, actualPairs, savedPairs)

	// add new pair and check diff with stored pairs
	// in main bot logic this can be changed with Binance API call
	actualPairs += "KGZBTC\n"

	diff, err := utils.Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "ADDED: KGZBTC\n", diff)
}
