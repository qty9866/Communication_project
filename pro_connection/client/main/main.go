package main

import (
	"fmt"
	"os"
	"pro_connection/client/process"
)

// 定义两个全局变量，一个表示用户名，一个表示用户密码
var userid int
var upwd string
var uname string

func main() {
	// 接收客户的选择
	var key int
	// 判断是否继续显示菜单
	// var loop = true
	for {
		fmt.Println("---------------------欢迎登录多人聊天系统---------------------")
		fmt.Println("\t\t\t 1. 登录聊天室")
		fmt.Println("\t\t\t 2. 用户注册")
		fmt.Println("\t\t\t 3. 退出系统")
		fmt.Println("\t\t\t 请选择(1~3)")
		fmt.Println("输入:")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户ID")
			fmt.Scanf("%d\n", &userid)
			fmt.Println("请输入用户密码")
			fmt.Scanf("%s\n", &upwd)
			// 完成登录
			// 1.创建一个userproceess实例
			userprocess := &process.UserProcess{}
			userprocess.Login(userid, upwd)

		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户ID")
			fmt.Scanf("%d\n", &userid)
			fmt.Println("请输入用户密码")
			fmt.Scanf("%s\n", &upwd)
			fmt.Println("请输入您的名字")
			fmt.Scanf("%s\n", &uname)

			userprocrss := &process.UserProcess{}
			userprocrss.Register(userid, upwd, uname)

		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("输入有误，请重新输入")
		}
	}

	// 根据用户的输入，显示更新的提示信息

}
