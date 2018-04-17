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
	Name      string
	LastPairs string
	Website   string
	dao       *DAO   `gorm:"-"`
	URL       string `gorm:"-"`
}

// NewHuobi is a Huobi struct constructor
func NewHuobi(dao *DAO) *Huobi {
	return &Huobi{
		URL:     "https://api.huobi.pro/v1/common/symbols",
		Name:    "Huobi",
		Website: "https://www.huobi.pro/",
		dao:     dao,
	}
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs from api.hitbtc.com
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

// GetListOfSavedPairs returns the list of previously saved pairs, stored in sqlite
func (h *Huobi) GetListOfSavedPairs() (string, error) {
	err := h.dao.GetLast(h)
	return h.LastPairs, err
}

// UpdatePairs returns the list of previously saved pairs, stored in sqlite
func (h *Huobi) UpdatePairs(pairs string) error {
	h.LastPairs = pairs
	err := h.dao.DeleteAllAndCreate(h)
	return err
}
