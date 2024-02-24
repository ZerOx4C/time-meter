package main

import (
	"os"
)

func main() {
	app := new(App)

	if err := app.initialize(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
