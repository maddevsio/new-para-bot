package dce

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jinzhu/gorm"
)

type TestDaoObject struct {
	gorm.Model
	LastPairs string
}

type TestAnotherDaoObject struct {
	gorm.Model
	LastPairs string
}

func TestDao(t *testing.T) {
	dao := NewDAO("/tmp/test.db")

	err := dao.DeleteAll(&TestDaoObject{})
	assert.NoError(t, err)

	c, err := dao.Count(&TestDaoObject{})
	assert.NoError(t, err)
	assert.Equal(t, 0, c)

	err = dao.DeleteAllAndCreate(&TestDaoObject{LastPairs: "lastpairs"})
	assert.NoError(t, err)

	err = dao.DeleteAllAndCreate(&TestAnotherDaoObject{LastPairs: "lastpairs"})
	assert.NoError(t, err)

	c, err = dao.Count(&TestDaoObject{})
	assert.NoError(t, err)
	assert.Equal(t, 1, c)
}
