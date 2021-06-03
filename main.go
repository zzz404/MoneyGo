package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/zzz404/MoneyGo/internal/web"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("aaa")
			os.Exit(0)
		}
	}()

	web.Start()
}
