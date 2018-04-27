package ext

import (
	"os"
	"os/exec"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/maddevsio/new-para-bot/dce"
)

type Cryptoexchange struct {
	gorm.Model
	Name      string
	LastPairs string
	ExePath   string
	DAO       *dce.DAO `gorm:"-"`
}

func NewCryptoexchange(dao *dce.DAO, name string) *Cryptoexchange {
	godotenv.Load(".env")
	cryptoexchange := &Cryptoexchange{}
	cryptoexchange.Name = name
	cryptoexchange.DAO = dao
	cryptoexchange.ExePath = os.Getenv("CRYPTOEXCHANGE")
	return cryptoexchange
}

func (c *Cryptoexchange) GetListOfActualPairs() string {
	return exe(c.ExePath, []string{c.Name})
}

func (c *Cryptoexchange) UpdatePairsAndSave(pairs string) error {
	c.LastPairs = pairs
	db := c.DAO.GetDB(c)
	defer db.Close()
	db.Where(Cryptoexchange{Name: c.Name}).Assign(Cryptoexchange{LastPairs: pairs}).FirstOrCreate(c)
	return db.Error
}

func (c *Cryptoexchange) GetListOfSavedPairs() (string, error) {
	db := c.DAO.GetDB(c)
	defer db.Close()
	db.First(c, "name = ?", c.Name)
	return c.LastPairs, db.Error
}

func exe(cmdName string, cmdArgs []string) string {
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		return ""
	}
	return string(cmdOut)
}
