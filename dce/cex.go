package dce

import (
	"sort"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Cex is a data struct for GORM to store prevours pairs from Cex exchange service
// API manual https://cex.io/rest-api
type Cex struct {
	gorm.Model
	Base
}

// NewCex is a Cex struct constructor
func NewCex(dao *DAO) *Cex {
	cex := &Cex{}
	cex.URL = "https://cex.io/api/currency_limits"
	cex.Name = "Cex"
	cex.Website = "https://cex.io/"
	cex.DAO = dao
	return cex
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs
func (k *Cex) GetListOfActualPairs() (string, error) {
	// iterate throught all active pairs
	resp, err := resty.R().Get(k.URL)
	if err != nil {
		return "", err
	}
	var pairs string
	var pairsSlice []string

	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		symbol1, _ := jsonparser.GetString(value, "symbol1")
		symbol2, _ := jsonparser.GetString(value, "symbol2")
		pair := symbol1 + "-" + symbol2
		pairsSlice = append(pairsSlice, pair+"\n")
	}, "data", "pairs")

	sort.Strings(pairsSlice)
	pairs = strings.Join(pairsSlice, "")

	return pairs, nil
}
