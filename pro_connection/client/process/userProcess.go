package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"pro_connection/client/utils"
	"pro_connection/common/message"
)

type UserProcess struct {
	// 暂时不需要字段
}

// 写一个方法，完成一个登录校验
func (userprocess *UserProcess) Login(userid int, upwd string) (err error) {
	//  下一个就要开始定协议

	// 1. 连接到服务器
	conn, err2 := net.Dial("tcp", "localhost:8889")
	if err2 != nil {
		fmt.Println("/UserProcess.go:客户端Dail失败")
		return err2
	}
	// 延时关闭
	defer conn.Close()
	// 2.准备通过conn发送消息给服务
	var mes message.Message //声明一个空的Message结构体
	mes.Type = message.LoginMesType
	// 3. 创建一个LoginMes 结构体
	var loginMes message.LoginMes
	loginMes.UserId = userid
	loginMes.Upwd = upwd
	// 这里注意不能直接把loginMes给mes.Data
	// 一个是结构体 一个是string字符串，肯定是不行的

	//4. 将loginMes进行序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Printf("/UserProcess.go:loginmes Json Marshal error: %v\n", err)
	}
	// 5. 把data赋值给mes.Data
	mes.Data = string(data)

	// 6.将Message结构体进行序列化
	data2, err4 := json.Marshal(mes)
	if err4 != nil {
		fmt.Printf("/UserProcess.go:mes Json Marshal error: %v\n", err4)
	}

	// 7. data2就是我们要发送的数据
	// 7.1 先把data2的长度发送给服务器
	// 先获取到data2的长度--> 转换成一个表示长度的byte切片
	pkgLen := uint32(len(data2))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen) // 相当于把这个pkgLen转化为Byte
	// 现在发送长度
	n, err := conn.Write(buf[:])
	if n != 4 || err != nil {
		fmt.Printf("/UserProcess.go:长度发送失败: %v\n", err)
		return
	}
	fmt.Printf("data2: %s\n", string(data2))

	// 发送消息本身
	_, err = conn.Write(data2)
	if err != nil {
		fmt.Printf("消息发送失败: %v\n", err)
		return
	}

	// 这里还需要处理服务器端返回的消息
	// 读包和写包都封装给Transfer
	// 创建一个Transfer连接
	transfer := &utils.Transfer{
		Conn: conn,
	}
	m, err := transfer.ReadPkg()
	if err != nil {
		fmt.Println("/UserProcess.go:读取服务端返回数据错误")
		return
	}
	var LoginResMes message.LoginResMes
	// 将message的data部分反序列化为LoginResMes
	err = json.Unmarshal([]byte(m.Data), &LoginResMes)
	if LoginResMes.Code == 200 {
		// 初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userid
		CurUser.UserStatus = message.UserOnline
		// 现在可以显示当前在线用户的列表，遍历LoginResMes里的UsersID切片
		fmt.Println("当前在线用户列表")
		for _, v := range LoginResMes.UsersID {
			if v == userid {
				continue
			}
			fmt.Println("用户id:\t", v)
			// 完成客户端的onlineusers的初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineusers[v] = user
		}

		fmt.Println("")
		fmt.Println("")
		// 这里我们还需再再客户端启动一个协程
		// 该协程保持和服务端的通讯，如果服务端有数据推送
		// 则接收并显示在客户端的终端
		go ProcessServerMes(conn)
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(LoginResMes.Error)
	}

	return
}

// 写一个方法，完成用户注册
func (userprocess *UserProcess) Register(userid int, upwd string, uname string) (err error) {
	//  下一个就要开始定协议

	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("/UserProcess.go:客户端Dail失败")
		return
	}
	// 延时关闭
	defer conn.Close()

	// 2.准备通过conn发送消息给服务
	var mes message.Message //声明一个空的Message结构体
	mes.Type = message.RegisterMesType
	// 3. 创建一个LoginMes 结构体
	var RegisterMessage message.RegisterMes
	RegisterMessage.User.UserId = userid
	RegisterMessage.User.Upwd = upwd
	RegisterMessage.User.Uname = uname

	//4. 将loginMes进行序列化
	data, err := json.Marshal(RegisterMessage)
	if err != nil {
		fmt.Printf("/UserProcess.go:RegisterMes.Data Json Marshal error: %v\n", err)
	}
	// 5. 把data赋值给mes.Data
	mes.Data = string(data)

	// 6.将Message结构体进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Printf("/UserProcess.go:RegisterMes Json Marshal error: %v\n", err)
		return
	}

	// 创建一个Transfer实例
	transfer := &utils.Transfer{
		Conn: conn,
	}

	// 发送data给服务器端
	err = transfer.SendPkg(data)
	if err != nil {
		fmt.Printf("/UserProcess.go:Register() Sendpkg error:%v", err)
		return
	}

	// 读取
	mes, err = transfer.ReadPkg() // 这个mes应该就是RegisterResMes
	if err != nil {
		fmt.Printf("/UserProcess.go:Register() ReadPkg() error:%v", err)
		return
	}

	// 将mes的Data部分反序列化为RegisterResMes
	var registerResMes message.RegisterResMes
	json.Unmarshal([]byte(mes.Data), &registerResMes)
	// 成功或者失败都退
	if registerResMes.Code == 200 {
		fmt.Println("注册成功了")
		os.Exit(0)
	} else {
		fmt.Println("注册失败，", registerResMes.Error)
		os.Exit(0)
	}
	return
}
