package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CoinType struct {
	Code   string
	Name   string
	ExRate float32
}

var CNY CoinType = CoinType{Code: "CNY", Name: "人民幣", ExRate: 1}
var TWD CoinType = CoinType{Code: "TWD", Name: "台幣", ExRate: 1}
var USD CoinType = CoinType{Code: "USD", Name: "美金", ExRate: 1}

func init() {
	fmt.Println("... 讀取匯率 ...")
	url := "https://api.coinbase.com/v2/exchange-rates?currency=" + CNY.Code
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var m = new(map[string]json.RawMessage)
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bodyBytes))
		err = json.Unmarshal(bodyBytes, &m)
		if err != nil {
			panic(err)
		}
	}

	// err = json.NewDecoder(r.Body).Decode(m)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Printf("m size: %d", len(*m))
}

var CoinTypes = []*CoinType{&TWD, &CNY, &USD}

func GetCoinTypeByCode(code string) (*CoinType, error) {
	switch code {
	case CNY.Code:
		return &CNY, nil
	case TWD.Code:
		return &TWD, nil
	case USD.Code:
		return &USD, nil
	}
	return nil, fmt.Errorf("不認識的 CoinType %s", code)
}

func Ccc() {
}
