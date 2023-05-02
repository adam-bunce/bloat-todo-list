package main

import (
	"github.com/adam-bunce/grpc-todo/api"
	"github.com/adam-bunce/grpc-todo/bootstrap"
	"github.com/adam-bunce/grpc-todo/variables"
)

func main() {
	// read in config.hcl to set db and app info
	variables.InitConfig()

	// bootstrap the db (create connection and set the global db)
	bootstrap.InitDB()

	// start todo_service server
	api.CreateServer()
}
