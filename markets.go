package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type BinanceResponse struct {
	Symbol string
	Price  string
}

func Cost(currency string) (cost float64, httpStatus int, err error) {
	//About Binance API: https://binance-docs.github.io/apidocs/
	marketEndpoint := "https://api3.binance.com/api/v3/ticker/price?"
	params := url.Values{}
	params.Set("symbol", currency)

	r, err := http.Get(marketEndpoint + params.Encode())
	defer r.Body.Close()
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("Failed market connect")
	}

	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("Failed market response")
	}

	respData := &BinanceResponse{}
	if err = json.Unmarshal(respBody, respData); err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("Failed market response")
	}

	price, err := strconv.ParseFloat(respData.Price, 64)
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("Failed market response")
	}

	return price, http.StatusOK, nil
}
