package ext

import (
	"os"
	"os/exec"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/maddevsio/new-para-bot/dce"
	resty "gopkg.in/resty.v1"
)

type Cryptoexchange struct {
	gorm.Model
	Name      string
	LastPairs string
	ExePath   string
	URL       string
	DAO       *dce.DAO `gorm:"-"`
}

func NewCryptoexchange(dao *dce.DAO, name string) *Cryptoexchange {
	godotenv.Load(".env")
	cryptoexchange := &Cryptoexchange{}
	cryptoexchange.Name = name
	cryptoexchange.DAO = dao
	cryptoexchange.URL = os.Getenv("CRYPTOEXCHANGE")
	return cryptoexchange
}

func (c *Cryptoexchange) GetListOfActualPairs() (string, error) {
	// iterate throught all active pairs
	resp, err := resty.R().Get(c.URL + "?dce=" + c.Name)
	if err != nil {
		return "", err
	}
	var pairs string
	jsonparser.ArrayEach([]byte(resp.String()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pair, _ := jsonparser.GetString(value, "symbol")
		pairs += pair + "\n"
	}, "pairs")

	return pairs, nil
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
