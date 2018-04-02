package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	var (
		pidFileFlag = flag.String("pid-file", "", "Pid file name")
		token = flag.String("token", "", "Telegram Bot Token")
		port = flag.String("port", "2000", "Port application")
	)

	flag.Parse()

	if len(*pidFileFlag) == 0 {
		log.Fatalln("Not pid")
	}

	if *token == "" {
		log.Print("-token is required")
		os.Exit(1)
	}

	// context := &daemon.Context{PidFileName: pidFileFlag}
	// child, err := context.Reborn()

	initRoutes(*port)
}
