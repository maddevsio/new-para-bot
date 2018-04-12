package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHuobi(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	huobi := NewHuobi(&dao)
	pairs, err := huobi.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "zec-btc")
	assert.Contains(t, pairs, "qsp-eth")
}
