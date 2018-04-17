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
	Name        string
	LastPairs   string
	Website     string
	dao         *DAO   `gorm:"-"`
	URL         string `gorm:"-"`
	RandomParam string `gorm:"-"` // need to add this to the URL to avoid cached responces
}

// NewHibtc is a Hibtc struct constructor
func NewHibtc(dao *DAO) *Hibtc {
	return &Hibtc{
		URL:     "https://api.hitbtc.com/api/2/public/symbol",
		Name:    "Hibtc",
		Website: "https://hitbtc.com/",
		dao:     dao,
	}
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

// GetListOfSavedPairs returns the list of previously saved pairs, stored in sqlite
func (h *Hibtc) GetListOfSavedPairs() (string, error) {
	err := h.dao.GetLast(h)
	return h.LastPairs, err
}

// UpdatePairs returns the list of previously saved pairs, stored in sqlite
func (h *Hibtc) UpdatePairs(pairs string) error {
	h.LastPairs = pairs
	err := h.dao.DeleteAllAndCreate(h)
	return err
}
