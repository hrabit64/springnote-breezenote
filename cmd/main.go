package main

import (
	"github.com/hrabit64/springnote-breezenote/config"
	"github.com/hrabit64/springnote-breezenote/core"
	"github.com/hrabit64/springnote-breezenote/database"
	"github.com/hrabit64/springnote-breezenote/di"
	validationUtil "github.com/hrabit64/springnote-breezenote/pkg/utils/validation"
	"log"
)

func main() {
	err := config.SetupConfig()
	if err != nil {
		panic(err)
	}

	err = database.RunSetup()
	if err != nil {
		panic(err)
	}

	err = config.SetFirebaseAuth()
	if err != nil {
		panic(err)
	}

	di.InitApplicationContext()

	validationUtil.SetupValidator()

	r := core.SetupRouter()

	log.Printf("Server is running on port %s", "8080")
	r.Run(":8080")

}
