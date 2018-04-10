package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKucoin(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	kucoin := NewKucoin(&dao)
	pairs, err := kucoin.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "ACAT-ETH")
	assert.Contains(t, pairs, "KEY-BTC")
}
