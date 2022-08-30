package main

import (
	"fmt"
	"go_code/chatroom/common/message"
	"go_code/chatroom/server/process"
	"go_code/chatroom/server/utils"
	"io"
	"net"
)

//先创建一个Process的结构体
type Processor struct {
	Conn net.Conn
}

//编写一个ServerProcessMes函数
//功能: 根据客户端发送的消息种类不同，决定调用哪个函数处理
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登录
		//err = serverProcessLogin(this.Conn, mes)
		//创建一个UserProcess实例
		up := &process3.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)

	case message.RegisterMesType:
		//处理注册
		up := &process3.UserProcess{Conn: this.Conn}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		//处理群发消息
		smsProcess := &process3.SmsProcess{}
		smsProcess.SendGroupMes(mes)

	case message.SmsUserMesType:
		smsProcess := &process3.SmsProcess{}
		smsProcess.SendUserMes(mes)


	case message.ExitLoginMesType:
		up := &process3.UserProcess{
		Conn: this.Conn,
		}
		err = up.ServerProcessExitLogin(mes)

	default:
		fmt.Println("消息类型不存在，无法处理...")

	}
	return
}

func (this *Processor) process2() (err error) {
	//循环读取客户端发送的信息
	for {
		//这里我们将读取数据表，直接封装成一个函数readPkg(),返回Message,Err
		//创建一个Transfer实例完成读包任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端也退出。。。")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}
		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}

}
