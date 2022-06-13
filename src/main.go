package main

import (
	"comedians/src/cfg"
	"comedians/src/db"
	"comedians/src/router"
	"fmt"
	"log"
	"time"
)

func init() {
	db.Connect()
}

func main() {
	time.Sleep(2 * time.Second)

	server := router.InitRoutes()
	server.MaxMultipartMemory = 2048 * 10 << 20

	port := cfg.Config().ServerPort

	server.Run(fmt.Sprintf(":%s", port))
	log.Println("SERVER RUNNED ON PORT", port)
}
