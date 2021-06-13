package coin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/zzz404/MoneyGo/internal/db"
)

type CoinType struct {
	Code   string
	Name   string
	ExRate float64
}

const TWD string = "TWD"

func init() {
	err := loadCoinTypes()
	if err != nil {
		panic(err)
	}

	m, err := queryExRateFromWeb()
	if err != nil {
		fmt.Printf("網路匯率失敗! 跳過匯率更新。錯誤訊息:  %s", err.Error())
	} else {
		err = updateExRateToDb(m)
		if err != nil {
			panic(err)
		}
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

func updateExRateToDb(m map[string]float64) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sql := "UPDATE CoinType SET exchangeRate=? WHERE code=?"
	pstmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer pstmt.Close()

	for _, coinType := range CoinTypes {
		rate, ok := m[coinType.Code]
		if ok {
			_, err = pstmt.Exec(rate, coinType.Code)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

var CoinTypes []*CoinType

func loadCoinTypes() error {
	rows, err := db.DB.Query("SELECT code, name, exchangeRate FROM CoinType")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		coinType := CoinType{}
		err = rows.Scan(&coinType.Code, &coinType.Name, &coinType.ExRate)
		if err != nil {
			return err
		}
		CoinTypes = append(CoinTypes, &coinType)
	}
	return nil
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

func queryExchangeRate(coinTypeCode string) (float64, error) {
	url := "https://api.coinbase.com/v2/exchange-rates?currency=" + coinTypeCode
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result = new(map[string]interface{})
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("state - %d", resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return 0, err
	}
	data := (*result)["data"].(map[string]interface{})
	rates := data["rates"].(map[string]interface{})
	exRateStr := rates["TWD"].(string)

	exRate, err := strconv.ParseFloat(exRateStr, 64)
	if err != nil {
		return 0, err
	}

	return exRate, nil
}

func GetCoinTypeByCode(code string) *CoinType {
	for _, coinType := range CoinTypes {
		if code == coinType.Code {
			return coinType
		}
	}
	panic(fmt.Errorf("CoinType %s 不存在", code))
}
