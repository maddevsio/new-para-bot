package dce

import (
	"sort"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Okex is a data struct for GORM to store prevours pairs from Okex exchange service
// API manual https://github.com/okcoin-okex/API-docs-OKEx.com
type Okex struct {
	gorm.Model
	Name      string
	LastPairs string
	dao       *DAO   `gorm:"-"`
	URL       string `gorm:"-"`
}

// NewOkex is a Hibtc struct constructor
func NewOkex(dao *DAO) *Okex {
	return &Okex{
		URL:  "https://www.okex.com/v2/spot/markets/index-tickers?limit=100000000",
		Name: "Okex",
		dao:  dao,
	}
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs from api.hitbtc.com
func (o *Okex) GetListOfActualPairs() (string, error) {
	// iterate throught all active pairs
	resp, err := resty.R().Get(o.URL)
	if err != nil {
		return "", err
	}
	var pairs string
	var pairsSlice []string

	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pair, _ := jsonparser.GetString(value, "symbol")
		pairsSlice = append(pairsSlice, pair+"\n")
	}, "data")

	sort.Strings(pairsSlice)
	pairs = strings.Join(pairsSlice, "")

	return pairs, nil
}

// GetListOfSavedPairs returns the list of previously saved pairs, stored in sqlite
func (o *Okex) GetListOfSavedPairs() (string, error) {
	err := o.dao.GetLast(o)
	return o.LastPairs, err
}

// UpdatePairs returns the list of previously saved pairs, stored in sqlite
func (o *Okex) UpdatePairs(pairs string) error {
	o.LastPairs = pairs
	err := o.dao.DeleteAllAndCreate(o)
	return err
}
