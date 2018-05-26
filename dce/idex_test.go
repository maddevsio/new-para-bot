package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdex(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	idex := NewIdex(&dao)
	pairs, err := idex.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "ETH_XNN")
	assert.Contains(t, pairs, "ETH_XOXO")
}
