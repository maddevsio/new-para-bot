package ext

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/maddevsio/new-para-bot/dce"
	"github.com/maddevsio/new-para-bot/utils"
	"github.com/stretchr/testify/assert"
)

func TestCryptoexchangeBinance(t *testing.T) {
	dao := dce.NewDAO("/tmp/test.db")

	c := NewCryptoexchange(&dao, "binance")
	c.URL = "http://localhost:4567/"
	actualPairs, err := c.GetListOfActualPairs()
	assert.NoError(t, err)
	err = c.UpdatePairsAndSave(actualPairs)
	assert.NoError(t, err)
	savedPairs, err := c.GetListOfSavedPairs()
	assert.NoError(t, err)
	actualPairs += actualPairs + "BLA-KZT"
	diff, err := utils.DiffSets(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Contains(t, diff, "+ BLA-KZT")

	c = NewCryptoexchange(&dao, "bitfinex")
	c.URL = "http://localhost:4567/"
	actualPairs, err = c.GetListOfActualPairs()
	assert.NoError(t, err)
	err = c.UpdatePairsAndSave(actualPairs)
	assert.NoError(t, err)
	savedPairs, err = c.GetListOfSavedPairs()
	assert.NoError(t, err)
	actualPairs += actualPairs + "BLA-RUB"
	diff, err = utils.DiffSets(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Contains(t, diff, "+ BLA-RUB")
}
