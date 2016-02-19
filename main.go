package main

import (
	"github.com/xackery/githubeq/service"
	//"github.com/xackery/githubeq/github"
	"log"
)

func main() {
	err := service.Start()
	if err != nil {
		log.Println(err.Error())
		return
	}
	/*client, err := github.GetClient()
	if err != nil {
		log.Println(err.Error())
		return
	}
	issues, resp, err := client.Issues.List(true, nil)
	if err != nil {
		log.Println("Error with listing issues:", err.Error())
		return
	}*/

	//log.Println(resp)
	//log.Println(issues)

	//for _, issue := range issues {
	//		issue.Number
	//	}

}
