// 处理用户相关操作
package process

import (
	"encoding/json"
	"fmt"
	"net"
	"pro_connection/common/message"
	"pro_connection/server/model"
	"pro_connection/server/utils"
)

type UserProcess struct {
	Conn net.Conn
	// 增加一个字段，表示该conn是哪个用户的
	UserId int
}

// NotifyOthersOnline 这里编写一个通知所有在线的用户的方法
// 这个id需要通知其他在线用户，我上线了
func (userprocess *UserProcess) NotifyOthersOnline(userId int) {
	// 遍历onlineUsers，然后一个一个发送(NotifyUserStatusMes)
	for id, userprocess := range USERMgr.OnlineUsers {
		// 过滤掉自己
		if id == userId {
			continue
		}
		// 开始通知:单独的写一个方法
		userprocess.NotifyToOthers(userId)

	}

}

func (userprocess *UserProcess) NotifyToOthers(userId int) {
	// 组装我们的NotifyUserStatusMes消息
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyuserstatusmes message.NotifyUserStatusMes
	notifyuserstatusmes.UserId = userId
	notifyuserstatusmes.Status = message.UserOnline

	// 将notifyuserstatusmes序列化
	data, err := json.Marshal(notifyuserstatusmes)
	if err != nil {
		fmt.Println("userProcess.go/NotifyToOthers() notifyuserstatusmes Marshal error", err)
		return
	}
	// 将序列化后的notifyuserstatusmes赋值给mes.Data
	mes.Data = string(data)

	// 再对mes进行序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("userProcess.go/NotifyToOthers() mes Marshal error", err)
		return
	}
	// 发送，创建一个transfer实例进行发送
	transfer := &utils.Transfer{
		Conn: userprocess.Conn,
	}
	err = transfer.SendPkg(data)
	if err != nil {
		fmt.Println("userProcess.go/NotifyToOthers() sendPkg error", err)
		return
	}
}

// ServerProcessLogin 编写一个方法 ServerProcessLogin 专门处理登录请求
func (userprocess *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 先从message中取出 mes.Data,并直接反序列化为LoginMes
	var Loginmes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &Loginmes)
	if err != nil {
		fmt.Println("LoginMes反序列化失败", err)
		return
	}

	// 1. 先声明一个要返回的resMes
	var ResMes message.Message
	ResMes.Type = message.LoginResMesType

	// 2. 再声明一个LoginResMes 并完成赋值
	var LoginResMes message.LoginResMes

	// 现在需要到Redis数据库去完成一个验证
	// 1.使用model.MyUserDao到redis中去进行验证
	user, err := model.MyUserDao.CheckLogin(Loginmes.UserId, Loginmes.Upwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			LoginResMes.Code = 500
			LoginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWDW {
			LoginResMes.Code = 403
			LoginResMes.Error = err.Error()
		} else {
			LoginResMes.Code = 404
			LoginResMes.Error = "服务器内部错误"
		}

	} else {
		LoginResMes.Code = 200
		// 这里因为用户已经登录成功了，于是就把该登录成功的用户放入到userMgr中
		// 将登录成功的用户ID赋给当前的userprocess
		userprocess.UserId = Loginmes.UserId
		USERMgr.AddOnlineUser(userprocess)
		userprocess.NotifyOthersOnline(Loginmes.UserId)
		fmt.Printf("user登录成功了: %v\n", user.Uname)
		// 将当前用户的id放入到LoginResMes的UseersID切片中
		// 遍历userMgr.onlineusers
		for id, _ := range USERMgr.OnlineUsers {
			LoginResMes.UsersID = append(LoginResMes.UsersID, id)
		}
	}

	// 3. 将loginResMes序列化
	data, err := json.Marshal(LoginResMes)
	if err != nil {
		fmt.Println("LoginResMes序列化失败")
		return
	}
	// 4.将序列化过后的data赋值给ResMes
	ResMes.Data = string(data)
	// 5.对ResMes进行序列化，准备发送
	data, err = json.Marshal(ResMes)
	if err != nil {
		fmt.Println("ResMes序列化失败")
		return
	}
	// 6.将发送操作封装到一个SendPkg函数里面去
	// 因为使用了分层模式(MVC)，先创建一个Transfer实例，然后读取
	transfer := &utils.Transfer{
		Conn: userprocess.Conn,
	}
	err = transfer.SendPkg(data)
	if err != nil {
		fmt.Println("userprocess发送数据失败")
		return
	}
	return
}

// 编写一个方法 ServerProcessRegister 专门处理用户的注册
func (userprocess *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 先从message中取出 mes.Data,并直接反序列化为RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("/userProcess.go:ServerProcessRegister() mes Unmarshal error ", err)
		return
	}

	// 1. 先声明一个要返回的resMes
	var ResMes message.Message
	ResMes.Type = message.RegisterResMesType

	// 2. 再声明一个RegisterResMes 并完成赋值
	var registerResMes message.RegisterResMes

	// 3. 现在在redis中完成对用户的注册
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_ALREADYEXISTS {
			registerResMes.Code = 400
			registerResMes.Error = err.Error()
		} else {
			registerResMes.Code = 404
			registerResMes.Error = "Register process error"
			fmt.Println("/userProcess.go:ServerProcessRegister() Register process error ", err)
			return
		}
	} else {
		// 没有错误
		registerResMes.Code = 200
	}

	// 4. 将loginResMes序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("/userProcess.go:ServerProcessRegister() registerResMes Marshal error", err)
		return
	}
	// 5.将序列化过后的data赋值给ResMes
	ResMes.Data = string(data)
	// 6.对ResMes进行序列化，准备发送
	data, err = json.Marshal(ResMes)
	if err != nil {
		fmt.Println("/userProcess.go:ServerProcessRegister() ResMes Marshal error", err)
		return
	}
	// 7.将发送操作封装到一个SendPkg函数里面去
	// 因为使用了分层模式(MVC)，先创建一个Transfer实例，然后读取
	transfer := &utils.Transfer{
		Conn: userprocess.Conn,
	}
	err = transfer.SendPkg(data)
	if err != nil {
		fmt.Println("userprocess发送数据失败")
		return
	}

	return
}
