package shooter


import (
	"fmt"
	"encoding/json"
	"github.com/goinggo/mapstructure"
)


type RedisDeployTask struct{

	Id string
	TaskType int
	TaskName string
	TaskStatus int
	TaskProgress string
	TaskConfig string
}


type TaskInfo struct {
	TaskId string
	Task interface{}
}


// 1) recve msg from fbm service
// 2) maintain the task info
type Shooter struct{
	RecvBuff chan string
	SendBuff chan string
	TaskPool map[string]*TaskInfo
}


type wsMsg struct{
	MsgType string `json:"msg_type"`
	Body interface{} `json:"body"`
}


var taskMap = make(map[string]func(msgBody interface{}))


func NewShooter() *Shooter{

	recvBuff := make(chan string)
	sendBuff := make(chan string)

	return &Shooter{
		RecvBuff:recvBuff,
		SendBuff:sendBuff,
	}
}


func (self *Shooter) Run(){

	for{
		select{
		case msg := <- self.RecvBuff:
			self.TaskDelivery(msg)
		}
	}

}


func (self *Shooter) HardwareTest(){


}


func (self *Shooter) HardwareCheck(){


}


func (self *Shooter) Deploy(msgBody interface{}){

	var taskInfo RedisDeployTask
	mapstructure.Decode(msgBody, &taskInfo)
	fmt.Printf("taskConfig %#v\n", taskInfo)
}


func (self *Shooter) ShiftWorker(){


}


func (self *Shooter) TaskDelivery(msgBody string){

	wsmsg := &wsMsg{}
	err := json.Unmarshal([]byte(msgBody), wsmsg)
	if err != nil {
		fmt.Println(err)
	}
    if wsmsg.MsgType == "test"{
		self.Deploy(wsmsg.Body)
	}

}

