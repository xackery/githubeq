package main

import (
	"log"
	"os"

	"github.com/xackery/githubeq/service"
)

func main() {
	err := service.Start()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
