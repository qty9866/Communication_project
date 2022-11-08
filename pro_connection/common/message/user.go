package message

//定义一个用户的结构体
type User struct {
	UserId     int    `json:"userid"`
	Upwd       string `json:"upwd"`
	Uname      string `json:"uname"`
	UserStatus int    `json:"userstatus"` // 用户的状态
}
