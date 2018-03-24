package dce

import (
	"bytes"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"
)

func TestBinanceApi(t *testing.T) {
	// iterate throught all active pairs
	resp, err := resty.R().Get("https://api.binance.com/api/v1/exchangeInfo")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	var pairs string
	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pair, _ := jsonparser.GetString(value, "symbol")
		pairs += pair + "\n"
	}, "symbols")

	// check popular pairs
	assert.Contains(t, pairs, "ETHBTC")
	assert.Contains(t, pairs, "XEMBNB")

	// TODO: if no pairs in storage than do not alert! consider this as a first run
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.DropTable(&Binance{})
	db.AutoMigrate(&Binance{})
	var count int
	db.Model(&Binance{}).Count(&count)
	assert.Equal(t, 0, count)
	db.Create(&Binance{LastPairs: pairs})
	var binance Binance
	db.Last(&binance)
	// this check is for GORM mainly
	assert.Equal(t, pairs, binance.LastPairs)

	// add new pair and check diff with stored pairs
	pairs += "KGZBTC\n"
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(binance.LastPairs, pairs, true)
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString("ADDED: " + text)
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString("DELETED: " + text)
		}
	}
	assert.Equal(t, "ADDED: KGZBTC\n", buff.String())
	db.Model(&binance).Update("LastPair", pairs)
}
