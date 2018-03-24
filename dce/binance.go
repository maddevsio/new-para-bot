package dce

import "github.com/jinzhu/gorm"

// Binance is a data struct for GORM to store prevours pairs from Binance
type Binance struct {
	gorm.Model
	LastPairs string
}
