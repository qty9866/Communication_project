package process

import (
	"fmt"
	"pro_connection/client/model"

	"pro_connection/common/message"
)

var onlineusers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser // 我们在用户登录成功之后，完成对CurUser的初始化

// 在客户端显示当前的用户
func OutputOnlineUser() {
	// 遍历一把onlineuser
	fmt.Println("显示当前在线的用户列表")
	for id, _ := range onlineusers {
		fmt.Println("用户id:\t", id)

	}
}

// 编写一个方法，处理返回的NotifyUserStatusMes
func UpdateUserStatus(notifyuserstatusmes *message.NotifyUserStatusMes) {
	// 看看原先是否有
	user, ok := onlineusers[notifyuserstatusmes.UserId]
	//  如果没有才去创建
	if !ok {
		user = &message.User{
			UserId: notifyuserstatusmes.UserId,
		}
	}
	// 改变状态
	user.UserStatus = notifyuserstatusmes.Status
	onlineusers[notifyuserstatusmes.UserId] = user
	// OutputOnlineUser()
}
