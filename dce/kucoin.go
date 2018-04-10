package dce

import (
	"sort"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Kucoin is a data struct for GORM to store prevours pairs from HiBTC exchange service
// API manual https://kucoinapidocs.docs.apiary.io
type Kucoin struct {
	gorm.Model
	Name      string
	LastPairs string
	dao       *DAO   `gorm:"-"`
	URL       string `gorm:"-"`
}

// NewKucoin is a Kucoin struct constructor
func NewKucoin(dao *DAO) *Kucoin {
	return &Kucoin{
		URL:  "https://api.kucoin.com/v1/open/tick",
		Name: "Kucoin",
		dao:  dao,
	}
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs
func (k *Kucoin) GetListOfActualPairs() (string, error) {
	// iterate throught all active pairs
	resp, err := resty.R().Get(k.URL)
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
func (k *Kucoin) GetListOfSavedPairs() (string, error) {
	err := k.dao.GetLast(k)
	return k.LastPairs, err
}

// UpdatePairs returns the list of previously saved pairs, stored in sqlite
func (k *Kucoin) UpdatePairs(pairs string) error {
	k.LastPairs = pairs
	err := k.dao.DeleteAllAndCreate(k)
	return err
}
