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
	Base
}

// NewEthfinex is a Ethfinex struct constructor
func NewEthfinex(dao *DAO) *Ethfinex {
	ethfinex := &Ethfinex{}
	ethfinex.URL = "https://api.ethfinex.com/v1/symbols"
	ethfinex.Name = "Ethfinex"
	ethfinex.Website = "https://www.ethfinex.com/"
	ethfinex.DAO = dao
	return ethfinex
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
