package main

import (
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/githubeq/service"
	"log"
)

func main() {
	config, err := eqemuconfig.GetConfig()
	if err != nil {
		log.Println("Failed to get config:", err.Error())
		return
	}
	if config.Github.RepoUser == "" {
		log.Println("Github not set in eqemuconfig.xml!")
		return
	}
	//TODO: Sanity checks for all eqemuconfig options

	err = service.Start()
	if err != nil {
		log.Println(err.Error())
		return
	}
}
