package main

import (
	"github.com/OhBonsai/yogo/app"
	"github.com/OhBonsai/yogo/api"
)

func main() {
	a, err := app.New()

	if err != nil{
		panic(err)
	}

	api.Init(a, a.Srv.Router)
	serverErr := a.StartServer()

	if serverErr != nil {
		panic(err)
	}
}

