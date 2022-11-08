package model

import (
	"encoding/json"
	"fmt"
	"pro_connection/common/message"

	"github.com/garyburd/redigo/redis"
)

// 定义一个UserDao结构体
// 完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

//在服务器启动时 初始化一个userdao实例
//把他做成全局变量，在需要和redis交互时，直接使用即可
var MyUserDao *UserDao

// 使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userdao *UserDao) {
	userdao = &UserDao{
		pool: pool,
	}
	return
}

// 方法一: 根据用户id，返回一个User实例+Err
func (userdao *UserDao) GetUserById(conn redis.Conn, id int) (user *message.User, err error) {
	mes, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	// 这里我们需要把res反序列化成一个Users实例
	user = &message.User{}
	err = json.Unmarshal([]byte(mes), &user)
	if err != nil {
		fmt.Printf("userDao.go:数据反序列化出错:%v\n", err)

	}
	return
}

// 完成一个登录校验 即完成对用户的校验
// Login 如果用户的id和pwd都正确，则返回一个user实例
// 如果有错误：返回对应的错误信息
func (userdao *UserDao) CheckLogin(userid int, upwd string) (user *message.User, err error) {
	// 先从UserDao的链接池中取出一个链接
	conn := userdao.pool.Get()
	defer conn.Close()
	user, err = userdao.GetUserById(conn, userid)
	if err != nil {
		return
	}
	// 此时证明获取到了这个用户，但是密码不一定正确
	// 这个时候再做判断
	if user.Upwd != upwd {
		err = ERROR_USER_PWDW
		return
	}
	return
}

// 使用一个UserDao实例完成一次对用户注册的操作
// 1.先对用户id进行查询，如果出现错误(nil)代表该用户不存在，也就是可以进行注册,否则表示用户已存在，不能注册
// 2.将传进来的mes进行序列化之后扔进redis进行存储
func (userdao *UserDao) Register(user *message.User) (err error) {
	// 先从UserDao的链接池中取出一个链接
	conn := userdao.pool.Get()
	defer conn.Close()
	_, err = userdao.GetUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_ALREADYEXISTS
		return
	}
	// 此时说明id在redis中还不存在，可以进行注册
	// 将user进行序列化，扔进redis中存储
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("/userDao.go:Register() user Marshal error:", err)
		return
	}
	// 扔进redis
	_, err = conn.Do("hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("/userDao.go:Register() hset into redis error:", err)
		return
	}
	return
}
