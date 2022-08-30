package main

import (
	"fmt"
	"go_code/chatroom/server/model"
	"net"
	"time"
)

//将读取包的任务封装得一个函数中readPkg()
//func readPkg(conn net.Conn) (mes message.Message, err error) {
//buf := make([]byte, 8086)
//fmt.Println("读取客户端发送的数据...")
//n, err := conn.Read(buf[:4])
//if n != 4 || err != nil {
//	fmt.Println("conn.Read err=", err)
//	return
//}
//fmt.Println("读到的buf=", buf[:4])
//
////根据buf[:4]转成一个uint32类型
//var pkgLen uint32
//pkgLen = binary.BigEndian.Uint32(buf[:4])
////根据pkgLen读取消息内容
//n, err = conn.Read(buf[:pkgLen])
//if n != int(pkgLen) || err != nil {
//	return
//}
////把pkgLen反序列化--> message.Message
//err = json.Unmarshal(buf[:pkgLen], &mes)
//if err != nil {
//	fmt.Println("json.Unmarsha err=", err)
//	return
//}
//return
//}

//编写一个函数serverProcessLogin函数，专门处理登录请求
//func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
////1.先从mes中取出mes.Data并直接序列化成LoginMes
//var loginMes message.LoginMes
//err = json.Unmarshal([]byte(mes.Data), &loginMes)
//if err != nil {
//	fmt.Println("json.Unmarsha1 fail err=", err)
//	return
//}
////1.声明一个resMes
//var resMes message.Message
//resMes.Type = message.LoginResMesType
////2.再声明一个LoginResMes并完成赋值
//var loginResMes message.LoginResMes
//
////如果用户id=100，密码=123456即认为合法否则不合法
//if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
//	//合法
//	loginResMes.Code = 200
//} else {
//	//不合法
//	loginResMes.Code = 500 //500状态码表示该用户不存在
//	loginResMes.Error = "该用户不存在，请注册再使用。。。"
//}
//
////3.将loginResMes序列化
//data, err := json.Marshal(loginResMes)
//if err != nil {
//	fmt.Println("json.Marshal fail", err)
//	return
//}
//fmt.Println("loginResMes data=", data)
//return
//}

//编写一个serverProcessMes函数
//功能: 根据客户端发送消息种类不同决定调用哪个函数处理
//func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
//	switch mes.Type {
//	case message.LoginMesType:
//		//处理登录
//		err = serverProcessLogin(conn, mes)
//	case message.RegisterMesType:
//	//处理注册
//	default:
//		fmt.Println("消息类型不存在，无法处理...")
//
//	}
//	return
//}

//处理和客户端的通信
//func process(conn net.Conn) {
//	//这里需要延迟关闭conn
//	defer conn.Close()
//
//	//循环的等待客户端发送的消息
//	for {
//		//buf := make([]byte, 8086)
//		//fmt.Println("读取客户端发送的数据...")
//		//n, err := conn.Read(buf[:4])
//		//if n != 4 || err != nil {
//		//	fmt.Println("conn.Read err=", err)
//		//	return
//		//}
//		//fmt.Println("读到的buf=", buf[:4])
//
//		//这里我们将读取的数据包封装成一个喊出readPkg(),返回Message，Error
//		mes, err := readPkg(conn)
//		if err != nil {
//			if err == io.EOF {
//				fmt.Println("客户端退出，服务器也将退出。。。")
//				return
//			} else {
//				fmt.Println("readPkg err=", err)
//				return
//			}
//		}
//		fmt.Println("mes=", mes)
//		err = serverProcessMes(conn, &mes)
//		if err != nil {
//			return
//		}
//	}
//}

//func init() {
//	//当服务器启动时，我们就去初始化我们的redis连接池
//
//	initUserDao()
//}

//处理和客户端的通讯
func process(conn net.Conn) {
	//延时关闭
	defer conn.Close()

	//这里调用总控，创建一个
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯携程错误=err", err)
		return
	}

}

//这里我们编写一个函数，完成对UserDao的初始化任务
func initUserDao() {
	//这里的pool本身就是一个全局的变量
	//这里需要注意一个初始化顺序问题
	//initPool,在initUserDao
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	initPool("192.168.1.232:6379", 16, 0, 300*time.Second)
	initUserDao()
	//提示信息
	fmt.Println("服务器在8889端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	//一旦监听成功，就等待客户端来链接服务器

	for {
		fmt.Println("等待客户端来连接服务器....")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}

		//一旦连接成功，则启动一个协程和客户端保持通讯
		go process(conn)
	}
}
