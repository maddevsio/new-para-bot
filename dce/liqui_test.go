package dce

import (
	"testing"

	"github.com/antonholmquist/jason"
	"github.com/stretchr/testify/assert"
	resty "gopkg.in/resty.v1"
)

// https://api.liqui.io/api/3/info

func TestLiqui(t *testing.T) {
	resp, err := resty.R().Get("https://api.liqui.io/api/3/info")
	assert.NoError(t, err)
	var pairs string
	v, err := jason.NewObjectFromBytes(resp.Body())
	assert.NoError(t, err)
	pairsObject, err := v.GetObject("pairs")
	for key, _ := range pairsObject.Map() {
		pairs += key + "\n"
	}
	t.Log(pairs) // TODO: need to sort
	t.Fail()
}
