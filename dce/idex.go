package dce

import (
	"sort"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Idex is a data struct for GORM to store prevours pairs from Idex exchange service
// API manual https://github.com/AuroraDAO/idex-api-docs
type Idex struct {
	gorm.Model
	Base
}

// NewIdex is a Idex struct constructor
func NewIdex(dao *DAO) *Idex {
	idex := &Idex{}
	idex.URL = "https://api.idex.market/returnTicker"
	idex.Name = "Idex"
	idex.Website = "https://idex.market/"
	idex.DAO = dao
	return idex
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs
func (i *Idex) GetListOfActualPairs() (string, error) {
	// iterate throught all active pairs
	resp, err := resty.R().Post(i.URL)
	if err != nil {
		return "", err
	}

	var pairs string
	var pairsSlice []string

	v, err := jason.NewObjectFromBytes(resp.Body())
	if err != nil {
		return "", err
	}

	pairsObject, err := v.GetObject()
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
