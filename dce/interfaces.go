package dce

// DCEChecker used in checkBinanceAndAlert and in DAO
type DCEChecker interface {
	GetListOfActualPairs() (string, error)
	GetPairs() string
	UpdatePairs(pairs string)
	GetDAO() *DAO
	GetName() string
	GetWebsite() string
}
