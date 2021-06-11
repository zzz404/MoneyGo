package bank

import "github.com/zzz404/MoneyGo/internal/db"

type Bank struct {
	Id   int
	Name string
}

var Banks []*Bank

func init() {
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
	return nil
}
