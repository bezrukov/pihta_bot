package main

import (
	"flag"
	"log"
)

func main() {

	var (
		pidFileFlag = flag.String("pid-file", "", "Pid file name")
	)

	flag.Parse()

	if len(*pidFileFlag) == 0 {
		log.Fatalln("Not pid")
	}

	// context := &daemon.Context{PidFileName: pidFileFlag}
	// child, err := context.Reborn()

	initRoutes()
}
