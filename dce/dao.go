package dce

import (
	"bytes"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// DAO is data access object, which consists of storage and data related functions
type DAO struct {
	DBPath string
}

// NewDAO is a DAO constructor
func NewDAO(DBPath string) DAO {
	return DAO{DBPath: DBPath}
}

// Count returns the number of all Binance records. This should be 0 or 1
func (dao *DAO) Count(obj interface{}) (int, error) {
	db := dao.getDB(obj)
	defer db.Close()
	var count int
	db.Model(obj).Count(&count)
	return count, db.Error
}

// DeleteAll deletes all data from Binance table
func (dao *DAO) DeleteAll(obj interface{}) error {
	db := dao.getDB(obj)
	defer db.Close()
	db.DropTable(obj)
	return db.Error
}

// DeleteAllAndCreate deletes all data from Binance table and create new record with new pairs
func (dao *DAO) DeleteAllAndCreate(obj interface{}) error {
	err := dao.DeleteAll(obj)
	if err != nil {
		return err
	}
	db := dao.getDB(obj)
	defer db.Close()
	db.Create(obj)
	return db.Error
}

// GetLast returns the last object form table
func (dao *DAO) GetLast(obj interface{}) error {
	db := dao.getDB(obj)
	defer db.Close()
	db.Last(obj)
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
func (dao *DAO) Diff(savedPairs string, actualPairs string) (string, error) {
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

func (dao *DAO) getDB(obj interface{}) *gorm.DB {
	db, err := gorm.Open("sqlite3", dao.DBPath)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(obj)
	return db
}
