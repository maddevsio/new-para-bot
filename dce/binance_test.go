package dce

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/maddevsio/new-para-bot/utils"
	"github.com/stretchr/testify/assert"
)

// TestBinanceApi represents a full flow and can be used in the
// bot as is. Just remowe asserts and you can get actual flow
func TestBinanceApi(t *testing.T) {
	// initiate Binance struct
	// TODO: need to pass custom params to the constructor
	dao := NewDAO("/tmp/test.db")
	binance := NewBinance(&dao)

	// get actual pairs and check
	actualPairs, err := binance.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, actualPairs, "ETHBTC") // check popular pairs
	assert.Contains(t, actualPairs, "XEMBNB")

	// delete all data before main logic
	// this deletion only for testing, we need to empty data set
	err = dao.DeleteAll(binance)
	assert.NoError(t, err)
	count, err := dao.Count(binance)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// update binance struct with newly get pairs from API
	err = binance.UpdatePairs(actualPairs)
	assert.NoError(t, err)

	// simulate the case when we are getting the data from storage
	savedPairs, err := binance.GetListOfSavedPairs()
	assert.NoError(t, err)
	assert.Equal(t, actualPairs, savedPairs)

	// add new pair and check diff with stored pairs
	// in main bot logic this can be changed with Binance API call
	actualPairs += "KGZBTC\n"

	diff, err := utils.Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "ADDED: KGZBTC\n\n", diff)
}
