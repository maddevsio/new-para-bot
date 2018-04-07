package utils

import (
	"bytes"
	"errors"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Diff returns differencies between two set of pairs
// * if pairs are equal the result is "" string
// * if we have new pair, the ADDED PAIRNAME\n shall be added to the return set
// * if some pairs were deleted, that DELETED PAIRNAME\n shall be added to the return set
// * example (note that several diffs are splitted by newlines, e.g. \n):
// ADDED: KGZBTC
// ADDED: BTCKGZ
// DELETED: MAVROETH
func Diff(savedPairs string, actualPairs string) (string, error) {
	if savedPairs[len(savedPairs)-1:] != "\n" {
		return "", errors.New("pairs should have a newline in the end")
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(savedPairs, actualPairs, true)
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString("ADDED: " + text)
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString("DELETED: " + text)
		}
	}
	return buff.String(), nil
}
