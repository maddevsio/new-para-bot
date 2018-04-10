package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLivecoin(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	livecoin := NewLivecoin(&dao)
	pairs, err := livecoin.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "ZBC/BTC")
	assert.Contains(t, pairs, "USC/ETH")
}
