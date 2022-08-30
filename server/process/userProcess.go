package process3

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"go_code/chatroom/server/model"
	"go_code/chatroom/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加一个字段，表示该Conn是哪个用户
	UserId int
}

//编写一个serverProcessLogin函数，专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {

	//1.先从mes中取出mes.Data并直接序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarsha1 fail err=", err)
		return
	}
	//1.声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//2.再声明一个LoginResMes并完成赋值
	var loginResMes message.LoginResMes

	//我们到redis数据库完成验证
	//1.使用mode.MyUserDao到redis完成验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {
		if err == model.ERROR_USER_NOTXISTS {
			loginResMes.Code = 500
			fmt.Println(user, "500")
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			fmt.Println(user, "403")
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}
		//loginResMes.Code = 500
		//loginResMes.Error = "用户不存在，请注册"
	} else {

		loginResMes.Code = 200

		//这里因为用户登录成功，我们需要把该登录成功的用户放到userMgr中
		//将登录成功的用户的UserID赋给this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)

		//通知其他在线用户我上线了
		this.NotifyOnthersOnlineUser(loginMes.UserId,loginResMes.Code)
		//将当前的在线用户id放到loginResMes.UserId
		//遍历userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登录成功")

	}

	//loginResMes.Code = 500
	//loginResMes.Error = "BUCUNZAI"
	//} else {
	//loginResMes.Code = 200
	//fmt.Println(user, "登录成功")
	////这里因为用户登录成功，我们需要把该登录成功的用户放到userMgr中
	////将登录成功的用户的UserID赋给this
	//this.UserId = loginMes.UserId
	//}

	////如果用户id=100，密码=123456即认为合法否则不合法
	//if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	//	//合法
	//	loginResMes.Code = 200
	//} else {
	//	//不合法
	//	loginResMes.Code = 500 //500状态码表示该用户不存在
	//	loginResMes.Error = "该用户不存在，请注册再使用。。。"
	//}

	//3.将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}
	//fmt.Println("loginResMes data=", data)
	//return

	//4.将data赋值给resMes
	resMes.Data = string(data)

	//5.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//6.发送data，我们将其封装到writePkg函数
	//因为时用了分层模式（mvc）,我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	return

}

//服务端处理注册
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//1.先从mes中取出mes.Data,并直接反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarsha fail err=", err)
		return
	}

	//1.先申明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//我们需要到redis数据库完成注册
	//1.使用mode.MyUserDao到redis验证
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误..."
		}

	} else {
		registerResMes.Code = 200
		fmt.Println("注册成功")
	}
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marsh fail", err)
		return
	}

	//4. 将data赋值给resMes
	resMes.Data = string(data)

	//5. 对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//6.发送data，我们将其封装到writePkg函数
	//因为使用分层模式（mvc）,我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

//这里我们编写通知所有在线用户的方法
//userId要通知其他的在线用户,上线
func (this *UserProcess) NotifyOnthersOnlineUser(userId int,code int) {
	//遍历onlineUsers,然后我们一个个发送NotifyUsersStatusMes
	for id, up := range userMgr.onlineUsers {
		//过滤掉自己
		fmt.Println("NotifyOnthersOnlineUser 在线用户id",id)

		if id == userId {
			continue
		}
		//开始通知
		if code == 200 {
		up.NotifyMeOnline(userId)
		fmt.Println("NotifyOnthersOnlineUser Code",code)
	} else if code == 400 {
		up.NotifyOnthersOfflineUser(userId)
			fmt.Println("NotifyOnthersOnlineUser Code",code)
		}
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {
	//组装我们的NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marsh err=", err)
		return
	}

	//将序列化后的notifyUserStatusMes赋值给mes.Data
	mes.Data = string(data)

	//对mes再次序列化准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marsh err=", err)
		return
	}

	//发送--我们通过TransFer实例发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
	}
	fmt.Printf("服务端发送的NotifyMeOnline消息消息长度=%d 内容=%s\n", len(data), string(data))
	fmt.Println("NotifyMeOnline err=", err)
	return

}

func (this *UserProcess) NotifyOnthersOfflineUser(userId int){
	//组装我们的NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.ExitLoginResMesType

	var notifyUserStatusMes message.ExitLoginResMes
	notifyUserStatusMes.UsersId = userId
	notifyUserStatusMes.UserStatus = message.UserOffline

	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marsh err=", err)
		return
	}

	//将序列化后的notifyUserStatusMes赋值给mes.Data
	mes.Data = string(data)

	//对mes再次序列化准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marsh err=", err)
		return
	}

	//发送--我们通过TransFer实例发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
	}
	fmt.Printf("服务端发送的NotifyOnthersOfflineUser消息消息长度=%d 内容=%s\n", len(data), string(data))
	fmt.Println("NotifyOnthersOfflineUser err=", err)
	return
}
//编写一个serverProcessExitLogin函数，专门处理退出登录请求
func (this *UserProcess) ServerProcessExitLogin(mes *message.Message) (err error) {

	//1.先从mes中取出mes.Data并直接序列化成LoginMes
	var exitLoginMes message.ExitLoginMes
	err = json.Unmarshal([]byte(mes.Data), &exitLoginMes)
	if err != nil {
		fmt.Println("json.Unmarsha1 fail err=", err)
		return
	}
	//1.声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//2.再声明一个LoginResMes并完成赋值
	var resLoginMes message.LoginResMes
   resLoginMes.Code = 400

       this.UserId = exitLoginMes.UserId
	   userMgr.DelOnlineUser(exitLoginMes.UserId)
	userMgr.onlineUsers = userMgr.GetAllOnlineUser()

		for id, _ := range userMgr.onlineUsers {
			resLoginMes.UsersId = append(resLoginMes.UsersId, id)

		}

	this.NotifyOnthersOnlineUser(exitLoginMes.UserId,resLoginMes.Code)
		fmt.Println(exitLoginMes.UserId, "下线成功")



	//3.将loginResMes序列化
	data, err := json.Marshal(resLoginMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//4.将data赋值给resMes
	resMes.Data = string(data)

	//5.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//6.发送data，我们将其封装到writePkg函数
	//因为时用了分层模式（mvc）,我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	fmt.Printf("服务端发送的ServerProcessExitLogin消息消息长度=%d 内容=%s\n", len(data), string(data))

	err = tf.WritePkg(data)
	return



}
