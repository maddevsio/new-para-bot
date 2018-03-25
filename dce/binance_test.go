package dce

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

// TestBinanceApi represents a full flow and can be used in the
// bot as is. Just remowe asserts and you can get actual flow
func TestBinanceApi(t *testing.T) {
	b := NewBinance()
	actualPairs, err := b.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, actualPairs, "ETHBTC") // check popular pairs
	assert.Contains(t, actualPairs, "XEMBNB")
	err = b.DeleteAll()
	assert.NoError(t, err)
	count, err := b.Count()
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
	err = b.UpdatePairs(actualPairs)
	assert.NoError(t, err)
	// this check is for GORM mainly
	savedPairs, err := b.GetListOfSavedPairs()
	assert.NoError(t, err)
	assert.Equal(t, actualPairs, savedPairs)

	// add new pair and check diff with stored pairs
	actualPairs += "KGZBTC\n"

	diff := b.Diff(savedPairs, actualPairs)
	assert.Equal(t, "ADDED: KGZBTC\n", diff)
}
