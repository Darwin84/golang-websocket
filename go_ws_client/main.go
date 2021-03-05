package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
	"ws_client/ws_client"
	"ws_client/shooter"
)

var addr = flag.String("addr", "0.0.0.0:3030", "http service address")


func writeTime(wbuff chan string, rbuff chan string, quit chan string){

    ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select{
		case t:=<-ticker.C:
			wbuff <- t.String()
		case output:=<-rbuff:
			log.Printf(output)
		case _,ok := <-quit:
			if !ok {
				break
			}
		}
	}
}

func main() {

	flag.Parse()
	log.SetFlags(0)

	quit := make(chan string)
	shooterObj := shooter.NewShooter()
	wsClient := ws_client.NewWsClient(*addr, "/echo", shooterObj.RecvBuff, shooterObj.SendBuff)
	go shooterObj.Run()
	wsClient.Run()

	log.Printf("wsClient Run")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
    // go writeTime(shooterObj.SendBuff, rbuff, quit)

	select{
	case <-interrupt:
		log.Printf("interrupt")
		wsClient.Close()
		close(quit)
	}
}

