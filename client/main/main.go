package main

import (
	"fmt"
	"go_code/chatroom/client/process"
	"os"
)

var userId int
var userPwd string
var userName string

func main() {
	var key int //接受用户的选择
	//var loop = true //判断是否继续显示菜单

	for true {
		fmt.Println("------欢迎登录多人聊天系统------")
		fmt.Println("1、登录聊天系统")
		fmt.Println("2、注册用户")
		fmt.Println("3、退出系统")
		fmt.Println("请选择(1-3):")
		//fmt.Scanf("%d\n", &key)
		//fmt.Println("key=", key)
		//switch key {
		//
		//case 1:
		//	fmt.Println("登录聊天室")
		//
		//	//fmt.Printf("你输入的 userid=%d pwd=%v\n", userId, userPwd)
		//	loop = false
		//case 2:
		//	fmt.Println("注册用户...")
		//	loop = false
		//case 3:
		//	fmt.Println("退出系统...")
		//	os.Exit(0)
		//default:
		//	fmt.Println("输入有误，请重新输入")
		//}
		//}
		//根据用户的输入，显示新的提示信息

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			//说明用户要登录
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户的id")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的Pwd")
			fmt.Scanf("%s", &userPwd)

			//完成登录
			//1. 创建一个UserProcess的实例
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名字(nickName):")
			fmt.Scanf("%s\n", &userName)
			////2.调用UserProcess完成注册请求
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)

		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("你的输入有误，请重新输入")
		}
	}
}
