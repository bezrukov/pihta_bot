package main

import (
	"flag"
	"log"
	"os"
)

type User struct {
	Id int64
	Name string
}
func main() {

	var (
		pidFileFlag = flag.String("pid-file", "", "Pid file name")
		token = flag.String("token", "", "Telegram Bot Token")
		port = flag.String("port", "2000", "Port application")
	)

	var userIds map[int64]User

	flag.Parse()

	if len(*pidFileFlag) == 0 {
		log.Fatalln("Not pid")
	}

	if *token == "" {
		log.Print("-token is required")
		os.Exit(1)
	}

	go initRoutes(*port)

	botCtrl := newBotCtrl()

	botCtrl.init(*token, &userIds)
}
