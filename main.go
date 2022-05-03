package main

import (
	"github.com/hexcraft-biz/base-account-service/service"
)

func main() {
	app, err := service.New()
	MustNot(err)

	app.Run()
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}
