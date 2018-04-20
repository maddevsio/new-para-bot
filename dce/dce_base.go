package dce

// Base is a base struct for all DCEs
type Base struct {
	Name      string
	LastPairs string
	Website   string
	DAO       *DAO   `gorm:"-"`
	URL       string `gorm:"-"`
}

// GetPairs is a getter method for LastPairs field
func (base *Base) GetPairs() string {
	return base.LastPairs
}

// UpdatePairs updates LastPairs field
func (base *Base) UpdatePairs(pairs string) {
	base.LastPairs = pairs
}

func (base *Base) GetDAO() *DAO {
	return base.DAO
}

func (base *Base) GetName() string {
	return base.Name
}

func (base *Base) GetWebsite() string {
	return base.Website
}
