package bank

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/zzz404/MoneyGo/internal/db"
)

type Bank struct {
	Id   int
	Name string
}

var Banks []*Bank

func init() {
	loadBanks()
	loadBankAccounts()
}

func loadBanks() {
	rows, err := db.DB.Query("SELECT id, name FROM Bank")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		bank := Bank{}
		err = rows.Scan(&bank.Id, &bank.Name)
		if err != nil {
			panic(err)
		}
		Banks = append(Banks, &bank)
	}
}

func GetBank(id int) *Bank {
	for _, bank := range Banks {
		if bank.Id == id {
			return bank
		}
	}
	panic(fmt.Errorf("BankId %d 不存在", id))
}

var bankIdAccountsMap = map[string]([]string){}
var BankIdAccountsMapJson string

func addBankAccountToCache(bankId int, account string) {
	bankIdString := strconv.Itoa(bankId)
	accounts := bankIdAccountsMap[bankIdString]
	bankIdAccountsMap[bankIdString] = append(accounts, account)
}

func loadBankAccounts() {
	rows, err := db.DB.Query("SELECT bankId, account FROM BankAccount")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var bankId int
		var account string
		err = rows.Scan(&bankId, &account)
		if err != nil {
			panic(err)
		}
		addBankAccountToCache(bankId, account)
	}

	jsonByte, err := json.Marshal(bankIdAccountsMap)
	if err != nil {
		panic(err)
	}
	BankIdAccountsMapJson = string(jsonByte)
}

func AddBankAccount(account string, bankId int) error {
	if bankId <= 0 || account == "" {
		return errors.New("BankAccount 不能有空值")
	}

	sql := "INSERT INTO BankAccount (account, bankId) VALUES (?, ?)"
	_, err := db.ExecuteSql(sql, account, bankId)
	if err != nil {
		return err
	}
	addBankAccountToCache(bankId, account)

	return nil
}

func DeleteBankAccount(account string, bankId int) error {
	sql := "DELETE FROM BankAccount WHERE account=?"
	_, err := db.ExecuteSql(sql, account)
	if err != nil {
		return err
	}

	bankIdString := strconv.Itoa(bankId)
	accounts := bankIdAccountsMap[bankIdString]

	for i, a := range accounts {
		if a == account {
			bankIdAccountsMap[bankIdString] = append(accounts[:i], accounts[i+1:]...)
			break
		}
	}
	return nil
}
