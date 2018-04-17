package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCex(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	cex := NewCex(&dao)
	pairs, err := cex.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "DASH-USD")
	assert.Contains(t, pairs, "ZEC-USD")
}
