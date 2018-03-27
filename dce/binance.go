package dce

import (
	"bytes"
	"errors"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	"github.com/sergi/go-diff/diffmatchpatch"
	resty "gopkg.in/resty.v1"
)

// Binance is a data struct for GORM to store prevours pairs from Binance
type Binance struct {
	gorm.Model
	LastPairs string
	URL       string `gorm:"-"`
	DBPath    string `gorm:"-"`
}

// NewBinance is a Binance struct constructor
func NewBinance() *Binance {
	return &Binance{
		URL:    "https://api.binance.com/api/v1/exchangeInfo",
		DBPath: "/tmp/test.db", // TODO: need to handle this
	}
}

// GetListOfActualPairs makes a call to API and returns \n separated pairs from api.binance.com
func (b *Binance) GetListOfActualPairs() (string, error) {
	// iterate throught all active pairs
	resp, err := resty.R().Get(b.URL)
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

// GetListOfSavedPairs returns the list of previously saved pairs, stored in sqlite
func (b *Binance) GetListOfSavedPairs() (string, error) {
	db, err := b.getDB()
	if err != nil {
		return "", err
	}
	defer db.Close()
	db.Last(b)
	return b.LastPairs, db.Error
}

// Count returns the number of all Binance records. This should be 0 or 1
func (b *Binance) Count() (int, error) {
	db, err := b.getDB()
	defer db.Close()
	if err != nil {
		return 0, err
	}
	var count int
	db.Model(&Binance{}).Count(&count)
	return count, db.Error
}

// DeleteAll deletes all data from Binance table
func (b *Binance) DeleteAll() error {
	db, err := b.getDB()
	defer db.Close()
	if err != nil {
		return err
	}
	db.DropTable(&Binance{})
	return db.Error
}

// UpdatePairs deletes all data from Binance table and create new record with new pairs
func (b *Binance) UpdatePairs(pairs string) error {
	err := b.DeleteAll()
	if err != nil {
		return err
	}
	db, err := b.getDB()
	defer db.Close()
	if err != nil {
		return err
	}
	db.Create(&Binance{LastPairs: pairs})
	return db.Error
}

// Diff returns differencies between two set of pairs
// * if pairs are equal the result is "" string
// * if we have new pair, the ADDED PAIRNAME\n shall be added to the return set
// * if some pairs were deleted, that DELETED PAIRNAME\n shall be added to the return set
// * example (note that several diffs are splitted by newlines, e.g. \n):
// ADDED: KGZBTC
// ADDED: BTCKGZ
// DELETED: MAVROETH
func (b *Binance) Diff(savedPairs string, actualPairs string) (string, error) {
	if savedPairs[len(savedPairs)-1:] != "\n" {
		return "", errors.New("pairs should have a newline in the end")
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(savedPairs, actualPairs, true)
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
	return buff.String(), nil
}

func (b *Binance) getDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", b.DBPath)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(b)
	return db, nil
}
