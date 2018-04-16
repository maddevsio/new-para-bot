package dce

import (
	"sort"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Bitz is a data struct for GORM to store prevours pairs from Bitz exchange service
// API manual https://api.liqui.io/
type Bitz struct {
	gorm.Model
	Name      string
	LastPairs string
	dao       *DAO   `gorm:"-"`
	URL       string `gorm:"-"`
}

// NewBitz is a Bitz struct constructor
func NewBitz(dao *DAO) *Bitz {
	return &Bitz{
		URL:  "https://www.bit-z.com/api_v1/tickerall",
		Name: "Bitz",
		dao:  dao,
	}
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs from api.liqui.io
func (b *Bitz) GetListOfActualPairs() (string, error) {
	resp, err := resty.R().Get(b.URL)
	if err != nil {
		return "", err
	}

	var pairs string
	var pairsSlice []string

	v, err := jason.NewObjectFromBytes(resp.Body())
	if err != nil {
		return "", err
	}

	pairsObject, err := v.GetObject("data")
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
func (b *Bitz) GetListOfSavedPairs() (string, error) {
	err := b.dao.GetLast(b)
	return b.LastPairs, err
}

// UpdatePairs returns the list of previously saved pairs, stored in sqlite
func (b *Bitz) UpdatePairs(pairs string) error {
	b.LastPairs = pairs
	err := b.dao.DeleteAllAndCreate(b)
	return err
}
