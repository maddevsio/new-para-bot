package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptopia(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	cryptopia := NewCryptopia(&dao)
	pairs, err := cryptopia.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "PEPE/DOGE")
	assert.Contains(t, pairs, "SEND/BTC")
}
