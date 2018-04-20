package dce

import (
	"sort"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Kraken is a data struct for GORM to store prevours pairs from Kraken exchange service
// API manual https://www.kraken.com/en-us/help/api#public-market-data
type Kraken struct {
	gorm.Model
	Base
}

// NewKraken is a Kraken struct constructor
func NewKraken(dao *DAO) *Kraken {
	kraken := &Kraken{}
	kraken.URL = "https://api.kraken.com/0/public/AssetPairs"
	kraken.Name = "Kraken"
	kraken.Website = "https://www.kraken.com/"
	kraken.DAO = dao
	return kraken
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs from api.liqui.io
func (k *Kraken) GetListOfActualPairs() (string, error) {
	resp, err := resty.R().Get(k.URL)
	if err != nil {
		return "", err
	}

	var pairs string
	var pairsSlice []string

	v, err := jason.NewObjectFromBytes(resp.Body())
	if err != nil {
		return "", err
	}

	pairsObject, err := v.GetObject("result")
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
