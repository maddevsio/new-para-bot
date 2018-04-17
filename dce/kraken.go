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
	Name      string
	LastPairs string
	Website   string
	dao       *DAO   `gorm:"-"`
	URL       string `gorm:"-"`
}

// NewKraken is a Kraken struct constructor
func NewKraken(dao *DAO) *Kraken {
	return &Kraken{
		URL:     "https://api.kraken.com/0/public/AssetPairs",
		Name:    "Kraken",
		Website: "https://www.kraken.com/",
		dao:     dao,
	}
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

// GetListOfSavedPairs returns the list of previously saved pairs, stored in sqlite
func (k *Kraken) GetListOfSavedPairs() (string, error) {
	err := k.dao.GetLast(k)
	return k.LastPairs, err
}

// UpdatePairs returns the list of previously saved pairs, stored in sqlite
func (k *Kraken) UpdatePairs(pairs string) error {
	k.LastPairs = pairs
	err := k.dao.DeleteAllAndCreate(k)
	return err
}
