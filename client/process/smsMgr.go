package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
)

func outputGroupMes(mes *message.Message) {

	//1.反序列化mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarsh err=", err.Error())
		return
	}

	//显示消息
	info := fmt.Sprintf("用户id:\t%d对大家说:\t%s", smsMes.UserId, smsMes.Content)

	fmt.Println(info)
	fmt.Println()
}

func outputUserMes(mes *message.Message) {

	//1.反序列化mes.Data
	var smsUserMes message.SmsUserMes
	err := json.Unmarshal([]byte(mes.Data), &smsUserMes)
	if err != nil {
		fmt.Println("json.Unmarsh err=", err.Error())
		return
	}

	//显示消息
	info := fmt.Sprintf("用户id:\t%d对用户id:\t%d说:\t%s", smsUserMes.UserId, smsUserMes.OtherUserId, smsUserMes.Content)

	fmt.Println(info)
	fmt.Println()
}
