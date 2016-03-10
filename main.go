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
		log.Println("Github user not set in eqemuconfig.xml")
		return
	}
	if config.Github.RefreshRate < 1 {
		log.Println("Invalid or not placed Github RefreshRate entry in eqemuconfig.xml")
		return
	}
	//TODO: Sanity checks for all eqemuconfig options

	err = service.Start()
	if err != nil {
		log.Println(err.Error())
		return
	}
}
