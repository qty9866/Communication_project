package process

import (
	"encoding/json"
	"fmt"
	"pro_connection/common/message"
)

func OutputGroupMes(mes *message.Message) {
	// 显示即可
	// 1. 反序列化
	var smsmes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsmes)
	if err != nil {
		fmt.Println("/smsMgr.go:mes Unmashal err", err.Error())
		return
	}
	// 显示
	fmt.Printf("用户id:\t%d 对大家说了一句：\t%s", smsmes.UserId, smsmes.Content)

}
