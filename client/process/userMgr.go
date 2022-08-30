package process

import (
	"fmt"
	"go_code/chatroom/client/model"
	"go_code/chatroom/common/message"
)

//客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser //我们在用户登录成功后完成对CurUser初始化
//在客户端显示在线的用户
func outputOnlineUser() {
	//遍历一把onlineUsers
	fmt.Println("当前在线用户列表:")
	for id, _ := range onlineUsers {
		fmt.Println("用户id:\t", id)
	}
}

//编写一个方法，处理返回来的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	//适当优化

	user, ok := onlineUsers[notifyUserStatusMes.UserId] //因为时updata可能之前有列表，需要先遍历一次判断是否有
	if !ok {
		//原来没有
		//user.OtherUserId = notifyUserStatusMes.Status
		user = &message.User{UserId: notifyUserStatusMes.UserId}
	}
	//user.UserId = notifyUserStatusMes.UserId ???
	//如果有更新状态即可
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user
	outputOnlineUser()
}



func updateOffUserStatus(notifyUserStatusMes *message.ExitLoginResMes) {
	//适当优化

	user, ok := onlineUsers[notifyUserStatusMes.UsersId] //因为时updata可能之前有列表，需要先遍历一次判断是否有
	if ok {
		//如果有更新状态即可
		user.UserStatus = notifyUserStatusMes.UserStatus
		onlineUsers[notifyUserStatusMes.UsersId] = user
		fmt.Printf("user.UserStatus=%v user.UsersId=%v\n  ",user.UserStatus,user.UserId)


		delete(onlineUsers,notifyUserStatusMes.UsersId)
		outputOnlineUser()
	}


}
