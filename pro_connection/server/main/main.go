package main

import (
	"fmt"
	"net"
	"pro_connection/server/model"
	"pro_connection/server/process"
	"time"
)

// 这里编写一个函数，对UserDao的初始化任务
func InitUserDao() {
	// 这里的pool 来自于redis.go声明的全局变量
	// 所以这个函数的调用需要在InitPool()的后面
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	process.InitUserMgr()
	// 服务器一开启，就初始化redis链接池
	InitPool("localhost:6379", 16, 0, 300*time.Second)
	// 初始化一个UserDao
	InitUserDao()
	// 提示信息
	fmt.Println("服务器在8889端口进行监听")
	listen, err := net.Listen("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Printf("服务器端监听端口失败: %v\n", err)
		return
	}
	defer listen.Close()
	// 一旦监听成功 ，就等待客户端进行连接
	for {
		fmt.Println("等待客户端连接服务器")
		conn, err2 := listen.Accept()
		if err2 != nil {
			fmt.Println("Accetp with wrong")
			fmt.Printf("err2: %v\n", err2)
		}
		// 一旦连接成功则起一个协程和客户端保持数据的通讯
		go Process(conn)
	}

}

//处理和客户端的通讯
func Process(conn net.Conn) {
	// 这里需要延时关闭conn
	defer conn.Close()
	// 这里调用总控，创建一个
	p := &Processor{
		Conn: conn,
	}
	p.Processing()

}
