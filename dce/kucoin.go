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
	Base
}

// NewKucoin is a Kucoin struct constructor
func NewKucoin(dao *DAO) *Kucoin {
	kucoin := &Kucoin{}
	kucoin.URL = "https://api.kucoin.com/v1/open/tick"
	kucoin.Name = "Kucoin"
	kucoin.DAO = dao
	return kucoin
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
