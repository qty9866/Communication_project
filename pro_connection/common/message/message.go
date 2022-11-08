package message

// 定义消息类型常量
const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

// 这里我们定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusy
)

type Message struct {
	Type string `json:"type"` // 消息类型
	Data string `json:"data"` // 消息的类型
}

// LoginMes 先定义两个消息
type LoginMes struct {
	UserId int    `json:"userid"` // 用户ID
	Upwd   string `json:"upwd"`   // 用户密码
	Uname  string `json:"uname"`  // 用户名
}

// LoginResMes 定义返回消息
type LoginResMes struct {
	Code    int    `json:"code"` // 类似于状态码 500 表示该用户还没注册 200 表示登录成功
	UsersID []int  // 保存用户id的切片
	Error   string `json:"error"` // 返回错误信息
}

type RegisterMes struct {
	User User `json:"user"`
}

type RegisterResMes struct {
	Code  int    `json:"code"` // 400 表示该用户已占用 200 表示注册成功
	Error string `json:"error"`
}

// NotifyUserStatusMes 为了配合服务器端推送用户状态变化的消息，新定义一个类型
type NotifyUserStatusMes struct {
	UserId int `json:"userid"`
	Status int `json:"status"`
}

//  增加一个SmsMes
type SmsMes struct {
	Content string `json:"content"` // 内容
	User
}
