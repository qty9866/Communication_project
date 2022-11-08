package process

import (
	"encoding/json"
	"fmt"
	"pro_connection/client/utils"
	"pro_connection/common/message"
)

type SmsProcess struct{}

// 发送群聊消息
func (smsprocess *SmsProcess) BroadCast(content string) (err error) {
	// 1.创建一个message
	var mes message.Message
	mes.Type = message.SmsMesType

	// 2.创建一个SmsMes实例
	var smsmes message.SmsMes
	smsmes.Content = content
	smsmes.UserId = CurUser.UserId
	smsmes.UserStatus = CurUser.UserStatus

	// 3.序列化smsems
	data, err := json.Marshal(smsmes)
	if err != nil {
		fmt.Println("/smsProcess.go: smsmes marshal error", err.Error())
		return
	}
	mes.Data = string(data)

	// 4.对mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("/smsProcess.go: mes marshal error", err.Error())
		return
	}

	// 5.将序列化后的mes发送给服务器
	transfer := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = transfer.SendPkg(data)
	if err != nil {
		fmt.Println("/smsProcess.go: mes send error", err.Error())
		return
	}
	return
}
