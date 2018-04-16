package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitz(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	bitz := NewBitz(&dao)
	pairs, err := bitz.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "pxc_eth")
	assert.Contains(t, pairs, "uct_btc")
}
