package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOkex(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	okex := NewOkex(&dao)
	pairs, err := okex.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "viu_btc")
	assert.Contains(t, pairs, "ref_btc")
}
