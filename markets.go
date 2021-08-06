package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Money float64

func GetBinancePrice(currency string) (cost Money, err error) {
	//About Binance API: https://binance-docs.github.io/apidocs/
	params := url.Values{}
	params.Set("symbol", currency)

	r, err := http.Get(MarketEndpoint + params.Encode())
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

	return Money(priceFloat), nil
}

func Cost(currency string) (cost Money, err error) {
	return GetBinancePrice(currency)
}
