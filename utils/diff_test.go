package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	var savedPairs, actualPairs, diff string

	// alwaus end with newline (\n)
	// no diff check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\nPAIR2\n"
	diff, err := Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Empty(t, diff)

	// deleted line check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\n"
	diff, err = Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "DELETED: PAIR2\n\n", diff)

	// added line check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\nPAIR2\nPAIR3\n"
	diff, err = Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "ADDED: PAIR3\n\n", diff)

	// added lines check
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\nPAIR2\nPAIR3\nPAIR4\nPAIR5\nPAIR6\nPAIR7\n"
	diff, err = Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "ADDED: PAIR3\nPAIR4\nPAIR5\nPAIR6\nPAIR7\n\n", diff)

	// newline in the end existance check
	savedPairs, actualPairs = "PAIR1\nPAIR2", "PAIR1\nPAIR2\n"
	diff, err = Diff(savedPairs, actualPairs)
	assert.Error(t, err, "pairs should have a newline in the end")
	assert.Empty(t, diff)

	// check additions and deletion in one time
	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR3\nPAIR2\nPAIR4\n"
	diff, err = Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "DELETED: 1\nADDED: 3\nADDED: PAIR4\n\n", diff)

	// check one empty
	savedPairs, actualPairs = "\n", "1\n2\n3\n"
	diff, err = Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "ADDED: 1\n2\n3\n", diff)

	// check one empty
	savedPairs, actualPairs = "1\n2\n3\n", "\n"
	diff, err = Diff(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "DELETED: 1\n2\n3\n", diff)
}

func TestDiffSet(t *testing.T) {
	savedPairs, actualPairs := "PAIR1\nPAIR2\n", "PAIR1\nPAIR2\n"
	diff, err := DiffSets(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Empty(t, diff)

	savedPairs, actualPairs = "PAIR1\nPAIR3\nPAIR2\n", "PAIR1\nPAIR2\n"
	diff, err = DiffSets(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "DELETED: PAIR3\n", diff)

	savedPairs, actualPairs = "PAIR1\nPAIR2\n", "PAIR1\nPAIR2\nPAIR3\n"
	diff, err = DiffSets(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Equal(t, "ADDED: PAIR3\n", diff)
}

func TestDiffSetEmptyRight(t *testing.T) {
	savedPairs, actualPairs := "PAIR1\nPAIR2\n", ""
	diff, err := DiffSets(savedPairs, actualPairs)
	assert.NoError(t, err)
	assert.Contains(t, diff, "DELETED: PAIR1\n")
	assert.Contains(t, diff, "DELETED: PAIR2\n")
}

func TestSaveDiffedDataToFiles(t *testing.T) {
	// if text data are not equal
	// save each text to separate file to /tmp
	var name = "Kucoin"
	var strings1 = "1\n2\n3\n"
	var strings2 = "1\n2\n3\n4\n"
	SaveNonEqualStringsToFiles(name, strings1, strings2)
}
