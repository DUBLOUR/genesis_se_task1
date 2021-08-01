package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)


type CostType float64

func GetBinancePrice(currency string) (cost CostType, err error) {
	//About Binance API: https://binance-docs.github.io/apidocs/
	marketEndpoint := "https://api3.binance.com/api/v3/ticker/price?"
	params := url.Values{}
	params.Set("symbol", currency)

	r, err := http.Get(marketEndpoint + params.Encode())
	if err != nil {
		return 0, fmt.Errorf("failed market connect")
	}
	defer r.Body.Close()

	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, fmt.Errorf("incorrect market response")
	}

	type BinanceResponse struct {
		Symbol string
		Price  string
	}

	respData := &BinanceResponse{}
	if err = json.Unmarshal(respBody, respData); err != nil {
		return 0, fmt.Errorf("incorrect market response")
	}

	priceFloat, err := strconv.ParseFloat(respData.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("incorrect market response")
	}

	return CostType(priceFloat), nil
}

func Cost(currency string) (cost CostType, err error) {
	return GetBinancePrice(currency)
}
