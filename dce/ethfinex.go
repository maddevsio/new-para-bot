package dce

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Ethfinex is a data struct for GORM to store prevours pairs from HiBTC exchange service
// API manual https://www.ethfinex.com/api_docs
type Ethfinex struct {
	gorm.Model
	Name        string
	LastPairs   string
	dao         *DAO   `gorm:"-"`
	URL         string `gorm:"-"`
	RandomParam string `gorm:"-"` // need to add this to the URL to avoid cached responces
}

// NewEthfinex is a Ethfinex struct constructor
func NewEthfinex(dao *DAO) *Ethfinex {
	return &Ethfinex{
		URL:  "https://api.ethfinex.com/v1/symbols",
		Name: "Ethfinex",
		dao:  dao,
	}
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs
func (e *Ethfinex) GetListOfActualPairs() (string, error) {
	resp, err := resty.R().Get(e.URL)
	if err != nil {
		return "", err
	}

	var pairs string
	var pairsSlice []string

	err = json.Unmarshal([]byte(resp.Body()), &pairsSlice)
	if err != nil {
		return "", err
	}

	sort.Strings(pairsSlice)
	pairs = strings.Join(pairsSlice, "\n") + "\n"

	return pairs, nil
}

// GetListOfSavedPairs returns the list of previously saved pairs, stored in sqlite
func (e *Ethfinex) GetListOfSavedPairs() (string, error) {
	err := e.dao.GetLast(e)
	return e.LastPairs, err
}

// UpdatePairs returns the list of previously saved pairs, stored in sqlite
func (e *Ethfinex) UpdatePairs(pairs string) error {
	e.LastPairs = pairs
	err := e.dao.DeleteAllAndCreate(e)
	return err
}
