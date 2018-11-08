package main

import (
	"github.com/OhBonsai/yogo/app"
	"github.com/OhBonsai/yogo/api"
)

func main() {
	configFileLocation := "./config/default.json"
	a, err := app.New(configFileLocation)

	if err != nil{
		print(err.Error())
		panic(err)
	}

	api.Init(a, a.Srv.Router)
	serverErr := a.StartServer()

	if serverErr != nil {
		print(serverErr.Error())
		panic(err)
	}
}

