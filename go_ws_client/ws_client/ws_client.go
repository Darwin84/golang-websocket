package ws_client


import (
	"log"
	"net/url"
	"time"
	"github.com/gorilla/websocket"
)


type WsClient struct{
	addr string
	path string
	conn *websocket.Conn
	rBuff chan string
	wBuff chan string
	quit chan int
}


func NewWsClient(addr string, path string, readBuff chan string, writeBuff chan string ) *WsClient{

	return &WsClient{
		addr: addr,
		path: path,
		rBuff: readBuff,
		wBuff: writeBuff,
		conn:nil,
	}
}


func (self *WsClient) connect() (*websocket.Conn, error){

	u := url.URL{Scheme: "ws", Host: self.addr, Path: self.path}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    return conn, err
}


func (self *WsClient) deInit(conn *websocket.Conn, readConnect chan string, writeConnect chan string){
	close(readConnect)
	close(writeConnect)
	conn.Close()
}


func (self *WsClient) Run(){
	reconnect := make(chan string)
	for{
		readConnect := make(chan string)
		writeConnect := make(chan string)
		conn, err := self.connect()
		if err != nil{
			log.Printf("connect failed")
			continue
		}
		go self.readLoop(conn, reconnect, readConnect)
		go self.writeLoop(conn, reconnect, writeConnect)
		select{
		case <-reconnect:
			self.deInit(conn, readConnect, writeConnect)
		case _, ok:= <-self.quit:
			if !ok{
				self.deInit(conn, readConnect, writeConnect)
			}
			return
		}
	}
}


func (self * WsClient) Close(){

	close(self.quit)
}


func (self *WsClient) readLoop(conn *websocket.Conn, reconnect chan string, readConnect chan string){

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			reconnect <- "reconnect"
			return
		}
		select {
		case _, ok:=<- readConnect:
			if !ok {
				log.Println("readLoop quit")
				return
			}
		default:
		}
		self.rBuff <- string(message)
	}
}


func (self *WsClient) writeLoop(conn *websocket.Conn, reconnect chan string, writeConnect chan string){

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				reconnect <- "reconnect"
				return
			}
		case _, ok:= <- writeConnect:
			if !ok{
				log.Println("writeLoop quit")
				return
			}
		default:
		}
	}
}
