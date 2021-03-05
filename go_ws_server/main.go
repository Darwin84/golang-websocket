package main

import (
	"fmt"
	"time"
	"flag"
	"log"
	"os"
	"encoding/json"
	"os/signal"
	"net/http"
	"ws_server/ws_server"
)


var addr = flag.String("addr", ":3030", "http service address")

type wsMsg struct{
	MsgType string `json:"msg_type"`
	Body interface{} `json:"body"`
}


func writeBuff(wsServerMap map[string]*ws_server.WsServer){

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	for {
		select{
		case t:=<-ticker.C:
			task :=  "task"
			wsmsg := wsMsg{
				MsgType: "test",
				Body: task,
			}
			sendMsg,_ := json.Marshal(wsmsg)
			for _, server:=range wsServerMap{
				server.Rbuff <- sendMsg
			}
		case <-interrupt:
			return
		}
	}
}


func wsServerTest(){
	wsServerMap := make(map[string]*ws_server.WsServer)
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		ws_server.WsHandler(w, r, wsServerMap)
	})
	go writeBuff(wsServerMap)
	log.Fatal(http.ListenAndServe(*addr, nil))
}


func main() {

	flag.Parse()
	fmt.Println("start")
	log.SetFlags(0)
	wsServerTest()
}
