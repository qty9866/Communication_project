package process

import "fmt"

type UserMgr struct {
	OnlineUsers map[int]*UserProcess
}

// USERMgr 因为UserMgr实例在服务器端有且仅有一个
// 在很多地方都会使用到，将其定义为一个全局变量
var (
	USERMgr *UserMgr
)

// InitUserMgr 完成对userMgr的初始化工作
func InitUserMgr() {
	USERMgr = &UserMgr{
		OnlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// AddOnlineUser 完成对onlineUsers添加
func (usermgr *UserMgr) AddOnlineUser(userprocess *UserProcess) {
	usermgr.OnlineUsers[userprocess.UserId] = userprocess
}

// DelOnlineUser 假如有用户下线了，进行删除
func (usermgr *UserMgr) DelOnlineUser(userId int) {
	// 用内置函数删除map的[key]value
	delete(usermgr.OnlineUsers, userId)
}

// GetAllOnlineUsers 查询，返回当前所有在线的用户
func (usermgr *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return usermgr.OnlineUsers
}

//根据id 返回对应的userprocess
func (usermgr *UserMgr) GetOnlineUserById(userId int) (userprocess *UserProcess, err error) {
	// 从map中取出一个值，带检测方法
	userprocess, ok := usermgr.OnlineUsers[userId]
	if !ok {
		// 说明想要查找的这个用户当前不在线
		err = fmt.Errorf("用户%d不存在", userId) // 返回一个格式化的error
		return
	}
	return
}
