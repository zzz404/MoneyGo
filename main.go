package main

import (
    "github.com/zzz404/MoneyGo/internal/web"
    "log"
    "net/http"
)

func main() {
	http.HandleFunc("/", web.PersonsView)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
