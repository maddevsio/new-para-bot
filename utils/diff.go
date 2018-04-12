package utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	set "gopkg.in/fatih/set.v0"
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
			_, _ = buff.WriteString("ADDED: " + text + "\n")
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString("DELETED: " + text + "\n")
		}
	}
	return buff.String(), nil
}

// DiffSets works as Diff, but with another algo
func DiffSets(savedPairs string, actualPairs string) (string, error) {
	savedPairs = strings.Trim(savedPairs, " \n")
	actualPairs = strings.Trim(actualPairs, " \n")

	savedPairsSlice := strings.Split(savedPairs, "\n")
	actualPairsSlice := strings.Split(actualPairs, "\n")

	savedPairsSet := set.New()
	if !(len(savedPairsSlice) == 1 && savedPairsSlice[0] == "") {
		for _, element := range strings.Split(savedPairs, "\n") {
			savedPairsSet.Add(element)
		}
	}

	actualPairsSet := set.New()
	if !(len(actualPairsSlice) == 1 && actualPairsSlice[0] == "") {
		for _, element := range strings.Split(actualPairs, "\n") {
			actualPairsSet.Add(element)
		}
	}

	added := set.Difference(actualPairsSet, savedPairsSet)   // ADDED
	deleted := set.Difference(savedPairsSet, actualPairsSet) // DELETED

	var result string
	if added.Size() > 0 {
		for _, item := range added.List() {
			result += "ADDED: " + item.(string) + "\n"
		}
	}
	for _, item := range deleted.List() {
		result += "DELETED: " + item.(string) + "\n"
	}
	return result, nil
}

// SaveNonEqualStringsToFiles should save two different stings to files
// this needs for debug
func SaveNonEqualStringsToFiles(name string, string1 string, string2 string) {
	if string1 != string2 {
		var rand = GetRandString()
		var filename1 = rand + "_1_" + name
		var filename2 = rand + "_2_" + name
		ioutil.WriteFile("/tmp/"+filename1, []byte(string1), 0777)
		ioutil.WriteFile("/tmp/"+filename2, []byte(string2), 0777)
	}
}
