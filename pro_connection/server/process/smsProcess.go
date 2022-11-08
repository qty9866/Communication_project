// 处理短消息相关
package process

import (
	"encoding/json"
	"fmt"
	"net"
	"pro_connection/client/utils"
	"pro_connection/common/message"
)

type SmsProcess struct{}

// 写一个方法 进行行消息的转发
func (smsprocess *SmsProcess) SendGroupMes(mes *message.Message) {
	//  遍历服务器端的onlineuser map[int]*UserProcess
	//  将消息转发取出

	// 取出mes的内容
	var smsmes message.SmsMes
	// 进行反序列化
	err := json.Unmarshal([]byte(mes.Data), &smsmes)
	if err != nil {
		fmt.Println("/smsProcess.go:mes Unmashal error", err)
		return
	}
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("/smsProcess.go:mes Mashal error", err.Error())
	}
	for id, userprocess := range USERMgr.OnlineUsers {
		// 这里还需要做个判断，过滤掉自己
		if id == smsmes.UserId {
			continue
		}
		smsprocess.SendMesToOnliner(data, userprocess.Conn)
	}
}

func (smsprocess *SmsProcess) SendMesToOnliner(content []byte, conn net.Conn) {
	// 创建一个Transfer实例，进行转发
	transfer := &utils.Transfer{
		Conn: conn,
	}
	err := transfer.SendPkg(content)
	if err != nil {
		fmt.Println("/smsProcess.go:Sendpkg err", err.Error())
	}
}
