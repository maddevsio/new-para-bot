package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBinanceApi represents a full flow and can be used in the
// bot as is. Just remowe asserts and you can get actual flow
func TestBinanceApi(t *testing.T) {
	// initiate Binance struct
	// TODO: need to pass custom params to the constructor
	binance := NewBinance()

	// get actual pairs and check
	actualPairs, err := binance.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, actualPairs, "ETHBTC") // check popular pairs
	assert.Contains(t, actualPairs, "XEMBNB")

	// delete all data before main logic
	// this deletion only for testing, we need to empty data set
	err = binance.DeleteAll()
	assert.NoError(t, err)
	count, err := binance.Count()
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

	diff, err := binance.Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "ADDED: KGZBTC\n", diff)
}

// TODO: need to add test cases for errors
// to get the full 100% coverage

func TestDiff(t *testing.T) {
	var savedPairs, actualPairs, diff string
	binance := NewBinance()

	// alwaus end with newline (\n)
	// no diff check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\nPAIR2\n"
	diff, err := binance.Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Empty(t, diff)

	// deleted line check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\n"
	diff, err = binance.Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "DELETED: PAIR2\n", diff)

	// added line check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\nPAIR2\nPAIR3\n"
	diff, err = binance.Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "ADDED: PAIR3\n", diff)

	// newline in the end existance check
	savedPairs, actualPairs = "PAIR1\nPAIR2", "PAIR1\nPAIR2\n"
	diff, err = binance.Diff(savedPairs, actualPairs)
	assert.Error(t, err, "pairs should have a newline in the end")
	assert.Empty(t, diff)
}
