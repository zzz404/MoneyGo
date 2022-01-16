package coin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/zzz404/MoneyGo/internal/db"
	"github.com/zzz404/MoneyGo/internal/utils"
)

type CoinType struct {
	Code   string
	Name   string
	ExRate float64
}

const TWD string = "TWD"

func init() {
	loadCoinTypes()

	m, err := queryExRateFromWeb()
	if err != nil {
		fmt.Printf("網路讀取匯率失敗! 跳過匯率更新。錯誤訊息:  %s", err.Error())
	} else {
		updateExRateToDb(m)
		updateExRateToCache(m)
	}
}

func updateExRateToCache(m map[string]float64) {
	for _, coinType := range CoinTypes {
		rate, ok := m[coinType.Code]
		if ok {
			coinType.ExRate = rate
		}
	}
}

func updateExRateToDb(m map[string]float64) {
	tx, err := db.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer func() {
		utils.Must(db.CommitOrRollback(tx, err))
	}()

	sql := "UPDATE CoinType SET exchangeRate=? WHERE code=?"
	pstmt, err := tx.Prepare(sql)
	if err != nil {
		return
	}
	defer func() {
		err = utils.CombineError(err, pstmt.Close())
	}()

	for _, coinType := range CoinTypes {
		rate, ok := m[coinType.Code]
		if ok {
			_, err1 := pstmt.Exec(rate, coinType.Code)
			if err1 != nil {
				err = utils.CombineError(err, err1)
				return
			}
		}
	}
}

var CoinTypes []*CoinType

func loadCoinTypes() {
	rows, err := db.DB.Query("SELECT code, name, exchangeRate FROM CoinType")
	if err != nil {
		panic(err)
	}
	defer func() {
		utils.Must(err, rows.Close())
	}()

	for rows.Next() {
		coinType := CoinType{}
		if err = rows.Scan(&coinType.Code, &coinType.Name, &coinType.ExRate); err != nil {
			return
		}
		CoinTypes = append(CoinTypes, &coinType)
	}
}

func queryExRateFromWeb() (map[string]float64, error) {
	result := map[string]float64{}
	for _, coinType := range CoinTypes {
		if coinType.Code == TWD {
			continue
		}
		rate, err := queryExchangeRate(coinType.Code)
		if err != nil {
			return nil, err
		} else {
			result[coinType.Code] = rate
		}
	}
	return result, nil
}

func queryExchangeRate(coinTypeCode string) (exRate float64, err error) {
	url := "https://api.coinbase.com/v2/exchange-rates?currency=" + coinTypeCode
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer func() {
		err = utils.CombineError(err, resp.Body.Close())
	}()

	var result = new(map[string]interface{})
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("state - %d", resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return
	}
	data := (*result)["data"].(map[string]interface{})
	rates := data["rates"].(map[string]interface{})
	exRateStr := rates["TWD"].(string)

	exRate, err = strconv.ParseFloat(exRateStr, 64)
	return
}

func GetCoinTypeByCode(code string) *CoinType {
	for _, coinType := range CoinTypes {
		if code == coinType.Code {
			return coinType
		}
	}
	panic(fmt.Errorf("CoinType %s 不存在", code))
}
