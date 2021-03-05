package ws_server


import (
	// "fmt"
	"log"
	"net/http"
	"strings"
	"github.com/gorilla/websocket"
)


type WsServer struct{

	ConnKey string
	Conn *websocket.Conn
	Wbuff chan string
	Rbuff chan string
	rquit chan int
	wquit chan int
	isError chan int
}


var upgrader = websocket.Upgrader{} // use default options


func newWsServer(w http.ResponseWriter, r *http.Request) *WsServer{

	conn, _ := upgrader.Upgrade(w, r, nil)
	rBuff := make(chan string, 1)
	wBuff := make(chan string, 1)
	rquit := make(chan int, 1)
	wquit := make(chan int, 1)
	isError := make(chan int, 1)
	connKey := strings.Split(conn.RemoteAddr().String(),":")[0]
	return &WsServer{
		Conn:conn,
		Wbuff:wBuff,
		Rbuff:rBuff,
		rquit:rquit,
		wquit:wquit,
		isError:isError,
		ConnKey:connKey,
	}
}


func WsHandler(w http.ResponseWriter, r *http.Request, wsServers map[string]*WsServer){

	wsServer := newWsServer(w,r)
	go wsServer.ReadLoop()
	go wsServer.WriteLoop()
	wsServers[wsServer.ConnKey] = wsServer
    <- wsServer.isError
	close(wsServer.rquit)
	close(wsServer.wquit)
	wsServer.Conn.Close()
	log.Println("conn quit", wsServer.ConnKey)
	delete(wsServers, wsServer.ConnKey)
}


func (self *WsServer) ReadLoop(){

	for{
		select{
		case _, ok := <-self.rquit:
			if !ok {
				log.Println("readquit:", self.ConnKey)
				return
			}
		default:
			_, message, err := self.Conn.ReadMessage()
		    if err != nil {
			    self.isError <- 1
			    log.Println("read:", err)
			    break
			}
			// self.Rbuff <- string(message)
			log.Printf("recv Message: %s", message)
		}
	}

}


func (self *WsServer) WriteLoop(){

    for{
		select{
		case message,_ := <- self.Rbuff:
			log.Printf("writeLoop send message")
			err := self.Conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				self.isError <- 1
				log.Println("write:", err)
				break
			}
		case _,ok := <- self.wquit:
			if !ok {
				log.Println("writequit:", self.ConnKey)
				return
			}
	    default:
		}
	}
}
