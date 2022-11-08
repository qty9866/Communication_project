package process

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"pro_connection/client/utils"
	"pro_connection/common/message"
)

// 显示登录成功后的界面
func ShowMenu() {
	fmt.Printf("----------------------%v,您好,登录成功!----------------------\n", CurUser.Uname)
	fmt.Println("\t\t\t 1.显示在线用户")
	fmt.Println("\t\t\t 2.发送消息")
	fmt.Println("\t\t\t 3.消息列表")
	fmt.Println("\t\t\t 4.退出系统")
	fmt.Println("请输入1~4:")
	// 获取用户输入
	var key int
	var content string
	//  因为我们总会用到SmsProcess实例，因此我们将其定义在switch外部
	smsprocess := &SmsProcess{}
	fmt.Scanf("%v\n", &key)
	switch key {
	case 1:
		OutputOnlineUser()
	case 2:
		fmt.Println("请输入需要群发的信息")
		fmt.Scanf("%v\n", &content)
		err := smsprocess.BroadCast(content)
		if err != nil {
			fmt.Println("/server.go: Broadcast err", err.Error())
		}
	case 3:
		fmt.Println("消息列表")
	case 4:
		fmt.Println("您退出了系统")
		os.Exit(0)
	default:
		fmt.Println("输入错误,请重新输入")
	}
}

// 该协程保持和服务端的通讯，如果服务端有数据推送
// 则接收并显示在客户端的终端
func ProcessServerMes(conn net.Conn) {
	// 创建一个Transfer实例，不停的读数据
	transfer := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取")
		mes, err := transfer.ReadPkg()
		if err != nil {
			fmt.Printf("循环读取服务器端推送数据出错: %v\n", err)
			return
		}
		// 如果读到了消息，进行处理下一步逻辑
		switch mes.Type {
		case message.NotifyUserStatusMesType:
			// 1.取出NotifyUserStatusMes
			var notifyuserstatusmes message.NotifyUserStatusMes
			err = json.Unmarshal([]byte(mes.Data), &notifyuserstatusmes)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			// 2.将这个用户的信息，保存到客户端的map中去
			UpdateUserStatus(&notifyuserstatusmes)

		// 收到了有客户端发来的群发消息
		case message.SmsMesType:
			OutputGroupMes(&mes)
		default:
			fmt.Println("无法识别服务器发送的消息类型")
		}
	}
}
