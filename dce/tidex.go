package dce

import (
	"sort"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Tidex is a data struct for GORM to store prevours pairs from Liqui exchange service
// API manual https://tidex.com/exchange/public-api
type Tidex struct {
	gorm.Model
	Base
}

// NewTidex is a Tidex struct constructor
func NewTidex(dao *DAO) *Tidex {
	tidex := &Tidex{}
	tidex.URL = "https://api.tidex.com/api/3/info"
	tidex.Name = "Tidex"
	tidex.Website = "https://tidex.com/"
	tidex.DAO = dao
	return tidex
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs from api.liqui.io
func (t *Tidex) GetListOfActualPairs() (string, error) {
	resp, err := resty.R().Get(t.URL)
	if err != nil {
		return "", err
	}

	var pairs string
	var pairsSlice []string

	v, err := jason.NewObjectFromBytes(resp.Body())
	if err != nil {
		return "", err
	}

	pairsObject, err := v.GetObject("pairs")
	if err != nil {
		return "", err
	}

	for key := range pairsObject.Map() {
		pairsSlice = append(pairsSlice, key+"\n")
	}
	sort.Strings(pairsSlice)
	pairs = strings.Join(pairsSlice, "")

	return pairs, nil
}
