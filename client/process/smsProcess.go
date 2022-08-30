package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
)

type SmsProcess struct {
}

//发送群聊消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {

	//1.创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType

	//2.创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	//3.序列化smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail =", err.Error())
		return
	}

	mes.Data = string(data)

	//4.对mes再次进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail =", err.Error())
		return
	}
	//5.将mes发送给服务器
	tf := &utils.Transfer{Conn: CurUser.Conn}

	//6.发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGrup err=", err.Error())
		return
	}
	fmt.Printf("客户端发送的群聊消息长度=%d 内容=%s\n", len(data), string(data))

	return
}

//发送私聊消息
func (this *SmsProcess) SendUserMes(content string, ohterUserId int) (err error) {

	//1.创建一个Mes
	var mes message.Message
	mes.Type = message.SmsUserMesType

	//2.创建一个SmsMes实例
	var smsUserMes message.SmsUserMes
	smsUserMes.Content = content
	smsUserMes.UserId = CurUser.UserId
	smsUserMes.UserStatus = CurUser.UserStatus
	smsUserMes.OtherUserId = ohterUserId

	//3.序列化smsMes
	data, err := json.Marshal(smsUserMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail =", err.Error())
		return
	}

	mes.Data = string(data)

	//4.对mes再次进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail =", err.Error())
		return
	}
	//5.将mes发送给服务器
	tf := &utils.Transfer{Conn: CurUser.Conn}

	//6.发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGrup err=", err.Error())
		return
	}
	fmt.Printf("客户端发送的群聊消息长度=%d 内容=%s\n", len(data), string(data))

	return
}


//发送私聊消息
func (this *SmsProcess) putUserStatus(UserId int) (err error) {

	//1.创建一个Mes
	var exitLoginMes message.Message
	exitLoginMes.Type = message.ExitLoginMesType

	//2.创建一个SmsMes实例
	var exitLogin message.ExitLoginMes
	exitLogin.UserId = UserId

	//3.序列化smsMes
	data, err := json.Marshal(exitLogin)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail =", err.Error())
		return
	}

	exitLoginMes.Data = string(data)

	//4.对mes再次进行序列化
	data, err = json.Marshal(exitLoginMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail =", err.Error())
		return
	}
	//5.将mes发送给服务器
	tf := &utils.Transfer{Conn: CurUser.Conn}

	//6.发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGrup err=", err.Error())
		return
	}
	fmt.Printf("客户端发送的群聊消息长度=%d 内容=%s\n", len(data), string(data))

	return
}
