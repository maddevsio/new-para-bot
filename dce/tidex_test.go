package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTidex(t *testing.T) {
	dao := NewDAO("/tmp/test.db")
	tidex := NewTidex(&dao)
	pairs, err := tidex.GetListOfActualPairs()
	assert.NoError(t, err)
	assert.Contains(t, pairs, "ant_waves")
	assert.Contains(t, pairs, "rlc_weur")
}
