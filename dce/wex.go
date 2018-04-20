package dce

import (
	"sort"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Wex is a data struct for GORM to store prevours pairs from Wex exchange service
// API manual https://wex.nz/api/3/docs
type Wex struct {
	gorm.Model
	Base
}

// NewWex is a Wex struct constructor
func NewWex(dao *DAO) *Wex {
	wex := &Wex{}
	wex.URL = "https://wex.nz/api/3/info"
	wex.Name = "Wex"
	wex.Website = "https://wex.nz/"
	wex.DAO = dao
	return wex
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs
func (w *Wex) GetListOfActualPairs() (string, error) {
	resp, err := resty.R().Get(w.URL)
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
