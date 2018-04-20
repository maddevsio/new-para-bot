package dce

import (
	"sort"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	resty "gopkg.in/resty.v1"
)

// Cryptopia is a data struct for GORM to store prevours pairs from Cex exchange service
// API manual https://www.cryptopia.co.nz/Forum/Thread/255
type Cryptopia struct {
	gorm.Model
	Base
}

// NewCryptopia is a Cryptopia struct constructor
func NewCryptopia(dao *DAO) *Cryptopia {
	сryptopia := &Cryptopia{}
	сryptopia.URL = "https://www.cryptopia.co.nz/api/GetTradePairs"
	сryptopia.Name = "Сryptopia"
	сryptopia.Website = "https://www.cryptopia.co.nz/"
	сryptopia.DAO = dao
	return сryptopia
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs
func (c *Cryptopia) GetListOfActualPairs() (string, error) {
	// iterate throught all active pairs
	resp, err := resty.R().Get(c.URL)
	if err != nil {
		return "", err
	}
	var pairs string
	var pairsSlice []string

	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pair, _ := jsonparser.GetString(value, "Label")
		pairsSlice = append(pairsSlice, pair+"\n")
	}, "Data")

	sort.Strings(pairsSlice)
	pairs = strings.Join(pairsSlice, "")

	return pairs, nil
}
