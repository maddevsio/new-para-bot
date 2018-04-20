package dce

import (
	"sort"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Huobi is a data struct for GORM to store prevours pairs from Huobi exchange service
// API manual https://api.huobi.pro/
type Huobi struct {
	gorm.Model
	Base
}

// NewHuobi is a Huobi struct constructor
func NewHuobi(dao *DAO) *Huobi {
	huobi := &Huobi{}
	huobi.URL = "https://api.huobi.pro/v1/common/symbols"
	huobi.Name = "Huobi"
	huobi.Website = "https://www.huobi.pro/"
	huobi.DAO = dao
	return huobi
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs
func (h *Huobi) GetListOfActualPairs() (string, error) {
	resp, err := resty.R().Get(h.URL)
	if err != nil {
		return "", err
	}

	var pairs string
	var pairsSlice []string

	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pair1, _ := jsonparser.GetString(value, "base-currency")
		pair2, _ := jsonparser.GetString(value, "quote-currency")
		pairsSlice = append(pairsSlice, pair1+"-"+pair2+"\n")
	}, "data")

	sort.Strings(pairsSlice)
	pairs = strings.Join(pairsSlice, "")

	return pairs, nil
}
