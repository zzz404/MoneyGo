package main

import (
	"log"
	"net/http"

	"github.com/zzz404/MoneyGo/internal/web"
)

func main1() {
	http.HandleFunc("/", web.PersonsView)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
