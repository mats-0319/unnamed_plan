package main

import (
	"log"

	"github.com/mats0319/unnamed_plan/server/test/api"
)

// start server with flag '-t' to use test db connection
func main() {
	createTable()

	api.GetAccessToken()

	testApi("Register", api.Register)
	testApi("Login", api.Login)
	testApi("List User", api.ListUser)
	testApi("Modify User", api.ModifyUser)

	testApi("Create Note", api.CreateNote)
	testApi("List Note", api.ListNote)
	testApi("Modify Note", api.ModifyNote)
	testApi("Delete Note", api.DeleteNote)

	testApi("List Game Score", api.ListGameScore)
	testApi("Upload Game Score", api.UploadGameScore)

	//dropTable() // do not del testdata during dev

	log.Println("> All Test Passed! ^_^")
}

func testApi(name string, f func()) {
	log.Printf("> %s.\n", name)
	defer log.Println()

	f()
}
