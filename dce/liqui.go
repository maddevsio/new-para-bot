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
	Name      string
	LastPairs string
	dao       *DAO   `gorm:"-"`
	URL       string `gorm:"-"`
}

// NewLiqui is a Liqui struct constructor
func NewLiqui(dao *DAO) *Liqui {
	return &Liqui{
		URL:  "https://api.liqui.io/api/3/info",
		Name: "Liqui",
		dao:  dao,
	}
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs from api.liqui.io
func (l *Liqui) GetListOfActualPairs() (string, error) {
	resp, err := resty.R().Get("https://api.liqui.io/api/3/info")
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

// GetListOfSavedPairs returns the list of previously saved pairs, stored in sqlite
func (l *Liqui) GetListOfSavedPairs() (string, error) {
	err := l.dao.GetLast(l)
	return l.LastPairs, err
}

// UpdatePairs returns the list of previously saved pairs, stored in sqlite
func (l *Liqui) UpdatePairs(pairs string) error {
	l.LastPairs = pairs
	err := l.dao.DeleteAllAndCreate(l)
	return err
}
