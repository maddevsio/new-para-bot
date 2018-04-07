package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	var savedPairs, actualPairs, diff string
	dao := NewDAO("/tmp/test.db")

	// alwaus end with newline (\n)
	// no diff check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\nPAIR2\n"
	diff, err := dao.Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Empty(t, diff)

	// deleted line check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\n"
	diff, err = dao.Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "DELETED: PAIR2\n", diff)

	// added line check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\nPAIR2\nPAIR3\n"
	diff, err = dao.Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "ADDED: PAIR3\n", diff)

	// added lines check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\nPAIR2\nPAIR3\nPAIR4\nPAIR5\nPAIR6\nPAIR7\n"
	diff, err = dao.Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "ADDED: PAIR3\nPAIR4\nPAIR5\nPAIR6\nPAIR7\n", diff)

	// newline in the end existance check
	savedPairs, actualPairs = "PAIR1\nPAIR2", "PAIR1\nPAIR2\n"
	diff, err = dao.Diff(savedPairs, actualPairs)
	assert.Error(t, err, "pairs should have a newline in the end")
	assert.Empty(t, diff)
}
