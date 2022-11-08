package main

import (
	"fmt"
	"io"
	"net"
	"pro_connection/common/message"
	"pro_connection/server/process"
	"pro_connection/server/utils"
)

// 先创建一个Processor结构体
type Processor struct {
	Conn net.Conn
}

//Processing() 处理客户端通讯
func (P *Processor) Processing() {
	// 循环的读客户端的信息
	for {
		// 只读取有用的部分
		// 创建一个Transfer实例，完成读包的任务
		tf := &utils.Transfer{
			Conn: P.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端关闭了连接")
				return
			}
			fmt.Println("读取pkg错误", err)
			return
		}
		fmt.Printf("mes: %v\n", mes)

		err = P.ServerProcessMes(&mes)
		if err != nil {
			fmt.Println("processor.go/")
			return
		}
	}
}

// 编写一个SeverProcessMes 函数
// 功能：根据客户端发送的消息种类不同，决定调用哪个函数
// 相当于消息的总控
func (P *Processor) ServerProcessMes(mes *message.Message) (err error) {
	// 看看是否能接收到客户端群发的消息
	switch mes.Type {
	case message.LoginMesType:
		// 处理登录流程
		// 创建一个UserProcessor实例
		up := &process.UserProcess{
			Conn: P.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		// 处理注册流程
		// 创建一个UserProcessor实例
		up := &process.UserProcess{
			Conn: P.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		// 创建一个SmsProcess实例，完成转发群聊消息
		smsprocess := &process.SmsProcess{}
		smsprocess.SendGroupMes(mes)
	default:
		fmt.Println("消息类型不存在,无法处理")
	}
	return
}
