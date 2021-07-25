package bank

import (
	"fmt"
	"strings"

	"github.com/zzz404/MoneyGo/internal/db"
	"github.com/zzz404/MoneyGo/internal/utils"
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
	defer func() {
		utils.Must(rows.Close())
	}()

	for rows.Next() {
		bank := Bank{}
		utils.Must(rows.Scan(&bank.Id, &bank.Name))
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

type BankAccount struct {
	BankId  int
	Account string
}

func (a1 *BankAccount) Compare(a2 *BankAccount) int {
	diff := a1.BankId - a2.BankId
	if diff > 0 {
		return 1
	} else if diff < 0 {
		return -1
	} else {
		return strings.Compare(a1.Account, a2.Account)
	}
}

func (a *BankAccount) BankName() string {
	return GetBank(a.BankId).Name
}

var BankAccounts []*BankAccount

func loadBankAccounts() {
	rows, err := db.DB.Query("SELECT bankId, account FROM BankAccount ORDER BY bankId, account")
	if err != nil {
		panic(err)
	}
	defer func() {
		utils.Must(rows.Close())
	}()

	for rows.Next() {
		a := &BankAccount{}
		utils.Must(rows.Scan(&a.BankId, &a.Account))
		BankAccounts = append(BankAccounts, a)
	}
}

func AddBankAccount(account *BankAccount) error {
	sql := "INSERT INTO BankAccount (account, bankId) VALUES (?, ?)"
	_, err := db.DB.Exec(sql, account.Account, account.BankId)
	if err != nil {
		return err
	}

	if len(BankAccounts) == 0 {
		BankAccounts = append(BankAccounts, account)
	} else {
		for i, a := range BankAccounts {
			if account.Compare(a) < 0 {
				BankAccounts = append(BankAccounts, nil)
				copy(BankAccounts[i+1:], BankAccounts[i:])
				BankAccounts[i] = account
				break
			}
		}
	}
	return nil
}

func DeleteBankAccount(account string) error {
	sql := "DELETE FROM BankAccount WHERE account=?"
	_, err := db.DB.Exec(sql, account)
	if err != nil {
		return err
	}

	for i, a := range BankAccounts {
		if a.Account == account {
			BankAccounts = append(BankAccounts[:i], BankAccounts[i+1:]...)
			break
		}
	}
	return nil
}
