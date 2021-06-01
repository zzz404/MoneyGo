package db

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
