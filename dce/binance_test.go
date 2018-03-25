package dce

import (
	"bytes"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

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
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(savedPairs, actualPairs, true)
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString("ADDED: " + text)
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString("DELETED: " + text)
		}
	}
	assert.Equal(t, "ADDED: KGZBTC\n", buff.String())
}
