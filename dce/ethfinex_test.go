package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEthfinex(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	ethfinex := NewEthfinex(&dao)
	pairs, err := ethfinex.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "trxbtc")
	assert.Contains(t, pairs, "avteth")
}
