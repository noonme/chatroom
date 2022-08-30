package process3

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
	"net"
)

type SmsProcess struct {
}

func (this *SmsProcess) SendGroupMes(mes *message.Message) {

	//遍历服务端的onlineUser map[int]*UserProcess
	//将消息转发取出
	//取出消息的内容SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal( err=", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		//这里过滤自己不转发
		if id == smsMes.UserId {
			continue
		}
		this.SendMesToEachOnlineUser(data, up.Conn)
	}
}

func (this *SmsProcess) SendUserMes(mes *message.Message) {

	//遍历服务端的onlineUser map[int]*UserProcess
	//将消息转发取出
	//取出消息的内容smsUserMes
	var smsUserMes message.SmsUserMes
	err := json.Unmarshal([]byte(mes.Data), &smsUserMes)
	if err != nil {
		fmt.Println("json.Unmarshal SmsProcess消息 err=", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal SmsProcess消息  err=", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		//这是用户转发消息
		if id == smsUserMes.OtherUserId {
			this.SendMesToEachOnlineUser(data, up.Conn)
		}

	}
}

func (this *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	//创建一个TransFer实例发送data
	tf := &utils.Transfer{
		Conn: conn,
	}

	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发雄安是失败 err=", err)
	}
}
