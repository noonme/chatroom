package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
	"net"
	"os"
)

//显示登录成功后的界面
func ShowMenu(userId int) {
	fmt.Println("-------恭喜xxx登录成功---------")
	fmt.Println("-------1. 显示在线用户列表---------")
	fmt.Println("-------2. 发送消息---------")
	fmt.Println("-------3. 信息列表---------")
	fmt.Println("-------4. 退出系统---------")
	fmt.Println("请选择(1-4):")
	var key int
	var content string
	var num string
	var ohterUserId int

	//因为我们总会用到SmsProcess实例，因此我们将定义在swtich外部
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	fmt.Printf("您当前输入的是:%v\n", key)

	switch key {

	case 1:
		fmt.Println("显示在线用户列表")
		outputOnlineUser()
		//fmt.Printf("你输入的 userid=%d pwd=%v\n", userId, userPwd)
		//loop = false
	case 2:
		fmt.Println("发送消息...")
		fmt.Println("请选择需要发送消息的类型")
		fmt.Println("a、发送群聊消息")
		fmt.Println("b、发送私信")
		fmt.Scanf("%s\n", &num)
		if num == "a" {
			fmt.Println("a、发送群聊消息")
			fmt.Println("请输入你想说的话:")
			fmt.Scanf(" %s\n", &content)
			smsProcess.SendGroupMes(content)
		} else if num == "b" {
			fmt.Println("b、发送私信")
			fmt.Println("请输入你想说的话: ")
			fmt.Scanf("%s\n", &content)
			fmt.Println("请输入聊天的用户id: ")
			fmt.Scanf("%d\n", &ohterUserId)
			smsProcess.SendUserMes(content, ohterUserId)
		} else {
			fmt.Println("输入错误，请按照提示输入消息类型")
			return
		}

	case 3:
		fmt.Println("信息列表...")
	case 4:
		fmt.Println("退出系统...")
		smsProcess.putUserStatus(userId)
		fmt.Printf("用户%v\n下线成功", userId)
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不对请重新输入。。。")
		//os.Exit(0)

	}

}

//和服务端保持通讯
func serverProcessMes(conn net.Conn) {
	//创建一个transfer实例，不停的读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("客户端获取server发送消息 tf.ReadPkg err=", err)
			return
		}
		//如果读取的消息，又开始处理下一步逻辑
		//fmt.Printf("mes=%v\n", mes)
		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上线

			//1. 取出notifyUserstatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2. 把这个用户信息状态保存到客户端的map
			updateUserStatus(&notifyUserStatusMes)
		//处理

		case message.SmsMesType: //有人发送群聊
			outputGroupMes(&mes)
		case message.SmsUserMesType: //有人发私信
			outputUserMes(&mes)

		case message.ExitLoginResMesType: //有人下线了

			//1. 取出notifyUserstatusMes
			var notifyUserStatusMes message.ExitLoginResMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2. 把这个用户信息状态保存到客户端的map
			updateOffUserStatus(&notifyUserStatusMes)
		//处理

		default:
			fmt.Println("服务端返回无法处理的消息")

		}
	}

}
