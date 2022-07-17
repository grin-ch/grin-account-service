package main

import (
	"github.com/grin-ch/grin-account-service/pkg/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := runServer()
	if err != nil {
		log.Fatalln("run server ending, err:%s", err.Error())
	}
}

// run
func runServer() error {
	return cmd.RunServer()
}
