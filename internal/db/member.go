package db

type Member struct {
	Id   int
	Name string
}

type Bank struct {
	Id   int
	Name string
}

func QueryMembers() ([]Member, error) {
	rows, err := DB.Query("SELECT id, name FROM Member")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		member := Member{}
		err = rows.Scan(&member.Id, &member.Name)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}
