package dce

import (
	"sort"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Livecoin is a data struct for GORM to store prevours pairs from Liqui exchange service
// API manual https://www.livecoin.net/api?lang=ru
type Livecoin struct {
	gorm.Model
	Name      string
	LastPairs string
	dao       *DAO   `gorm:"-"`
	URL       string `gorm:"-"`
}

// NewLivecoin is a Liqui struct constructor
func NewLivecoin(dao *DAO) *Livecoin {
	return &Livecoin{
		URL:  "https://api.livecoin.net/exchange/ticker",
		Name: "Livecoin",
		dao:  dao,
	}
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs
func (l *Livecoin) GetListOfActualPairs() (string, error) {
	resp, err := resty.R().Get(l.URL)
	if err != nil {
		return "", err
	}

	var pairs string
	var pairsSlice []string

	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pair, _ := jsonparser.GetString(value, "symbol")
		pairsSlice = append(pairsSlice, pair+"\n")
	})

	sort.Strings(pairsSlice)
	pairs = strings.Join(pairsSlice, "")

	return pairs, nil
}

// GetListOfSavedPairs returns the list of previously saved pairs, stored in sqlite
func (l *Livecoin) GetListOfSavedPairs() (string, error) {
	err := l.dao.GetLast(l)
	return l.LastPairs, err
}

// UpdatePairs returns the list of previously saved pairs, stored in sqlite
func (l *Livecoin) UpdatePairs(pairs string) error {
	l.LastPairs = pairs
	err := l.dao.DeleteAllAndCreate(l)
	return err
}
