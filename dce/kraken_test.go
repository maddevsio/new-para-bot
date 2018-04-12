package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKraken(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	kraken := NewKraken(&dao)
	pairs, err := kraken.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "USDTZUSD")
	assert.Contains(t, pairs, "EOSEUR")
}
