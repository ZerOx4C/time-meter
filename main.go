package main

import (
	"os"
)

func main() {
	app := new(App)

	if err := app.Run(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
