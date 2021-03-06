package dce

import (
	"github.com/jinzhu/gorm"
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
	db := dao.GetDB(obj)
	defer db.Close()
	var count int
	db.Model(obj).Count(&count)
	return count, db.Error
}

// DeleteAll deletes all data from Binance table
func (dao *DAO) DeleteAll(obj interface{}) error {
	db := dao.GetDB(obj)
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
	db := dao.GetDB(obj)
	defer db.Close()
	db.Create(obj)
	return db.Error
}

// GetLast returns the last object form table
func (dao *DAO) GetLast(obj interface{}) error {
	db := dao.GetDB(obj)
	defer db.Close()
	db.Last(obj)
	return db.Error
}

// GetListOfSavedPairs returns the list of previously saved pairs, stored in sqlite
func (dao *DAO) GetListOfSavedPairs(obj DCEChecker) (string, error) {
	err := dao.GetLast(obj)
	return obj.GetPairs(), err
}

// UpdatePairsAndSave returns the list of previously saved pairs, stored in sqlite
func (dao *DAO) UpdatePairsAndSave(obj DCEChecker, pairs string) error {
	obj.UpdatePairs(pairs)
	err := dao.DeleteAllAndCreate(obj)
	return err
}

func (dao *DAO) GetDB(obj interface{}) *gorm.DB {
	db, err := gorm.Open("sqlite3", dao.DBPath)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(obj)
	return db
}
