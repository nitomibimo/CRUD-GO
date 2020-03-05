package api

import (
	"fmt"
	"log"
	"os"

	controllers "github.com/nitomibimo/CRUD-GO/api/controller"

	"github.com/joho/godotenv"
	"github.com/nitomibimo/CRUD-GO/api/controller/sedd"
)

var server = controllers.Server{}

// Run - Script for running app
func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	sedd.Load(server.DB)

	server.Run(":8080")

}
