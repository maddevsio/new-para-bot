package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWex(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	wex := NewWex(&dao)
	pairs, err := wex.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "ltc_rur")
	assert.Contains(t, pairs, "usdt_usd")
}
