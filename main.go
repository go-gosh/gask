package main

import "github.com/go-gosh/gask/app"

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	err = a.Run()
	if err != nil {
		panic(err)
	}
}
