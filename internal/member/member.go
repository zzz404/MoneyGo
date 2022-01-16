package member

import (
	"fmt"

	"github.com/zzz404/MoneyGo/internal/db"
	"github.com/zzz404/MoneyGo/internal/utils"
)

type Member struct {
	Id   int
	Name string

	AllTotalTWD  float64
	TimeTotalTWD float64
}

func (m *Member) AllTotalTWDString() string {
	return fmt.Sprintf("%.2f", m.AllTotalTWD)
}

func (m *Member) TimeTotalTWDString() string {
	return fmt.Sprintf("%.2f", m.TimeTotalTWD)
}

var Members []*Member

func init() {
	rows, err := db.DB.Query("SELECT id, name FROM Member")
	if err != nil {
		panic(err)
	}
	defer func() {
		utils.Must(rows.Close())
	}()

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
	panic(fmt.Errorf("MemberId %d 不存在", id))
}
