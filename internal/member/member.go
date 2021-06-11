package member

import "github.com/zzz404/MoneyGo/internal/db"

type Member struct {
	Id   int
	Name string
}

var Members []*Member

func init() {
	rows, err := db.DB.Query("SELECT id, name FROM Member")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		member := Member{}
		err = rows.Scan(&member.Id, &member.Name)
		if err != nil {
			panic(err)
		}
		Members = append(Members, &member)
	}
}

func GetMember(id int) *Member {
	for _, member := range Members {
		if member.Id == id {
			return member
		}
	}
	return nil
}
