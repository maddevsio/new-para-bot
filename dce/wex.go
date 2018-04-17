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
	Name      string
	LastPairs string
	Website   string
	dao       *DAO   `gorm:"-"`
	URL       string `gorm:"-"`
}

// NewWex is a Wex struct constructor
func NewWex(dao *DAO) *Wex {
	return &Wex{
		URL:     "https://wex.nz/api/3/info",
		Name:    "Wex",
		Website: "https://wex.nz/",
		dao:     dao,
	}
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

// GetListOfSavedPairs returns the list of previously saved pairs, stored in sqlite
func (w *Wex) GetListOfSavedPairs() (string, error) {
	err := w.dao.GetLast(w)
	return w.LastPairs, err
}

// UpdatePairs returns the list of previously saved pairs, stored in sqlite
func (w *Wex) UpdatePairs(pairs string) error {
	w.LastPairs = pairs
	err := w.dao.DeleteAllAndCreate(w)
	return err
}
