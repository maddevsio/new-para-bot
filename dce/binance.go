package dce

import (
	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Binance is a data struct for GORM to store prevours pairs from Binance
// API manual https://github.com/binance-exchange/binance-official-api-docs
type Binance struct {
	gorm.Model
	Base
	RandomParam string `gorm:"-"` // need to add this to the URL to avoid cached responces
}

// NewBinance is a Binance struct constructor
func NewBinance(dao *DAO) *Binance {
	binance := &Binance{}
	binance.URL = "https://api.binance.com/api/v1/exchangeInfo"
	binance.Name = "Binance"
	binance.Website = "https://www.binance.com/"
	binance.DAO = dao
	return binance
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs from api.binance.com
func (b *Binance) GetListOfActualPairs() (string, error) {
	// need to change random param on each request
	// binane API does not accept any other symbol
	// instead of "?"
	if b.RandomParam == "" {
		b.RandomParam = "?"
	} else {
		b.RandomParam = ""
	}
	// iterate throught all active pairs
	resp, err := resty.R().Get(b.URL + b.RandomParam)
	if err != nil {
		return "", err
	}
	var pairs string
	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pair, _ := jsonparser.GetString(value, "symbol")
		pairs += pair + "\n"
	}, "symbols")

	return pairs, nil
}
