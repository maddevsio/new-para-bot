package dce

import (
	"sort"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Hibtc is a data struct for GORM to store prevours pairs from HiBTC exchange service
// API manual https://api.hitbtc.com/
type Hibtc struct {
	gorm.Model
	Base
}

// NewHibtc is a Hibtc struct constructor
func NewHibtc(dao *DAO) *Hibtc {
	hibtc := &Hibtc{}
	hibtc.URL = "https://api.hitbtc.com/api/2/public/symbol"
	hibtc.Name = "Hibtc"
	hibtc.Website = "https://hitbtc.com/"
	hibtc.DAO = dao
	return hibtc
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs from api.hitbtc.com
func (h *Hibtc) GetListOfActualPairs() (string, error) {
	// iterate throught all active pairs
	resp, err := resty.R().Get(h.URL)
	if err != nil {
		return "", err
	}
	var pairs string
	var pairsSlice []string

	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pair, _ := jsonparser.GetString(value, "id")
		pairsSlice = append(pairsSlice, pair+"\n")
	})

	sort.Strings(pairsSlice)
	pairs = strings.Join(pairsSlice, "")

	return pairs, nil
}
