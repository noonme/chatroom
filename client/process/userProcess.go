package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
	"net"
	"os"
)

type UserProcess struct {
	//。。。。
	//确定字段信息

}

func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	//1.连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//延迟关闭
	defer conn.Close()

	//2.准备通过conn发送消息服务
	var mes message.Message
	mes.Type = message.RegisterMesType

	//3.创建一个LoginMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	//4.registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//5. 把data赋值给mes.Data字段
	mes.Data = string(data)

	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	//发送data给服务器
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送消息失败 err=", err)
	}

	mes, err = tf.ReadPkg() //mes就是RegisterResMes
	if err != nil {
		fmt.Println("readPkg(conn)err=", err)
		return
	}

	//将mes的Data部分反序列化为RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，登录一把")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}

//给关联一个用户登录的方法
//写一个函数，完成登录
func (this *UserProcess) Login(userId int, userPwd string) (err error) {
	//下面要开始定协议
	//说明该文件就是原来的login.go做一个改进，即封装到UserProcess结构体
	//fmt.Printf("userId = %d userPwd=%s\n",userId,userPwd)
	//return nil

	//1.连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.准备通过conn发送消息给服务端
	var mes message.Message
	mes.Type = message.LoginMesType

	//3.创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4. 将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//4-5步的必要性，因为message包的Data是string类型，string()无法直接转换loginMes;
	//这里使用json.Marsh过度一下

	//5. 把data赋给mes.Date字段
	mes.Data = string(data)

	//6. 将mes （5这里导致mes是string需要再次转换）进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marsh err=", err)
		return
	}
	//7.到这个时候data就是我们需要发送的消息

	//7.1 先把data的长度发送给服务端
	//将获取到的data长度-->转化为一个表示长度的byte切片？？
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//fmt.Println("pkgLen=", pkgLen)
	var buf [4]byte
	//fmt.Printf("buf=%v,buf[:4]=%v\n", buf, buf[:4])
	binary.BigEndian.PutUint32(buf[0:4], pkgLen) //ByteOrder规定了如何将字节序列和 16、32或64比特的无符号整数互相转化。即格式转换
	//fmt.Printf("buf=%v,buf[:4]=%v\n", buf, buf[:4])

	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) failed", err)
		return
	}

	fmt.Printf("客户端发送的消息长度=%d 内容=%s\n", len(data), string(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}

	//休眠20s
	//time.Sleep(20 * time.Second)
	//fmt.Println("休眠20s...")

	//这里还需要处理服务器返回的消息，==
	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}

	//将mes的Data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {

		//初始化CurUser==SmsProcess相关
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		//可以显示当前在线用户列表，遍历loginResMes.UserId
		fmt.Println("当前在线用户列表如下:")
		for _, v := range loginResMes.UsersId {
			//如果我们要求不显示自己在线，下面我们增加一个代码
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)
			//完成客户单的onlineUsers的初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")

		//这里我们还需要在客户端启动一个协程
		//该协程保持和服务器端的通讯，如果服务器端有消息推送给客户端
		//则接收并显示在客户端终端
		go serverProcessMes(conn)

		//1.显示登录成功的菜单
		for {
			ShowMenu()
		}
		//return
		////初始化CurUser
		//CurUser.Conn = conn
		//CurUser.UserId = userId
		//CurUser.UserStatus = message.UserOnline
		////fmt.Println("登录成功")
		////可以显示当前在线用户列表，便利loginResMes.UserId
		//for _, v := range loginResMes.UsersId {
		//	//如果在线不显示自己在线，我们增加如下代码
		//	if v == userId {
		//		continue
		//	}
		//	fmt.Println("用户id:\t", v)
		//	//完成客户端的的onlineUsers完成初始化
		//	user := &message.User{
		//		UserId:     v,
		//		UserStatus: message.UserOnline,
		//	}
		//	onlineUsers[v] = user
		//}
		//fmt.Println("\n\n")
		//
		////} else if loginResMes.Code == 500 {
		////	fmt.Println(loginResMes.Error)
		////}
		//go serverProcessMes(conn)
		//
		////1. 显示我们的登录成功菜单循环。。。
		//for {
		//	showMenu()
		//}
		//fmt.Println("登录成功")

		//} else if loginResMes.Code == 500 {
		//	fmt.Println(loginResMes.Error)
		//	os.Exit(0)
		//} else if loginResMes.Code == 403 {
		//	fmt.Println(loginResMes.Error)
		//	os.Exit(0)
		//}
	} else if loginResMes.Code == 400 {

		//初始化CurUser==SmsProcess相关
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOffline

		//可以显示当前在线用户列表，遍历loginResMes.UserId
		fmt.Println("当前在线用户列表如下:")
		for _, v := range loginResMes.UsersId {


			fmt.Println("用户id:\t", v)
			//完成客户单的onlineUsers的初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOffline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")

		//这里我们还需要在客户端启动一个协程
		//该协程保持和服务器端的通讯，如果服务器端有消息推送给客户端
		//则接收并显示在客户端终端
		go serverProcessMes(conn)
}else {
		fmt.Println(loginResMes.Error)
		//os.Exit(0)
	}

	return
}



