package utils

import (
	"fmt"
	"testing"

	"github.com/PieterD/diff"
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
}

func TestAnotherDiffer(t *testing.T) {
	// Data
	l := []string{
		"linea",
		"lineb",
		"lineM",
		"linec",
	}
	r := []string{
		"linea",
		"lineQ",
		"linec",
		"lineM",
	}
	// Diff l and r using Strings
	d := diff.New(diff.Strings{
		Left:  l,
		Right: r,
	})
	// Print the diff
	for i := range d {
		switch d[i].Delta {
		case diff.Both:
			//fmt.Printf("  %s\n", l[d[i].Index])
		case diff.Left:
			fmt.Printf("- %s\n", l[d[i].Index])
		case diff.Right:
			fmt.Printf("+ %s\n", r[d[i].Index])
		}
	}
}
