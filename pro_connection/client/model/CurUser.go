// CurrentUser 不是公用的 所以定义在这个包里面
package model

import (
	"net"
	"pro_connection/common/message"
)

// 因为在客户端很多地方很多地方会用到CurUser，我们将其作为一个全局变量
type CurUser struct {
	Conn net.Conn
	message.User
}
