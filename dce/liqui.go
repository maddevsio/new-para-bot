package dce

import (
	"sort"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Liqui is a data struct for GORM to store prevours pairs from Liqui exchange service
// API manual https://api.liqui.io/
type Liqui struct {
	gorm.Model
	Base
}

// NewLiqui is a Liqui struct constructor
func NewLiqui(dao *DAO) *Liqui {
	liqui := &Liqui{}
	liqui.URL = "https://api.liqui.io/api/3/info"
	liqui.Name = "Liqui"
	liqui.Website = "https://liqui.io/"
	liqui.DAO = dao
	return liqui
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs from api.liqui.io
func (l *Liqui) GetListOfActualPairs() (string, error) {
	resp, err := resty.R().Get(l.URL)
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
