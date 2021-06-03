package db

type Member struct {
	Id   int8
	Name string
}

func QueryMembers() ([]Member, error) {
	rows, err := db.Query("SELECT id, name FROM Member")
	if err != nil {
		return nil, err
	}

	var members []Member
	for rows.Next() {
		member := Member{}
		err = rows.Scan(&member.Id, &member.Name)
		assertSucc(err)
		members = append(members, member)
	}
	return members, nil
}
