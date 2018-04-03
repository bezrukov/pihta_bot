package main

import (
	"flag"
	"log"
	"os"
	// "pihta_bot/modules/dealMock"
)

type User struct {
	Id int
	Name string
}
func main() {

	var (
		pidFileFlag = flag.String("pid-file", "", "Pid file name")
		token = flag.String("token", "", "Telegram Bot Token")
		port = flag.String("port", "2000", "Port application")
	)

	userIds := map[int]User{}

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

	botCtrl.init(*token, userIds)

	// TODO пример того как работать со сделками
	//deal := dealMock.NewDeal()
	//
	//deal.Start(5)
	//result, isFinish := deal.Process()
}
