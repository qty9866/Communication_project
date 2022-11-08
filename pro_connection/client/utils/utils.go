// 工具的集合
package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"pro_connection/common/message"
)

// 将这些工具关联到相应的结构体中,这个结构体负责传输
// 将Transfer与readPkg/sendPkg绑定起来，只要新建一个这样的实例就可以调用方法
type Transfer struct {
	// 首先需要一个链接
	Conn net.Conn
	Buf  [4096]byte
}

func (transfer *Transfer) ReadPkg() (mes message.Message, err error) {
	// conn.Read 只有在conn没有被关闭的情况下才会阻塞
	// 如果说有任意一方(无论是client还是server)关闭了连接，马上就不会阻塞了
	// 这样也就会出现，读不到东西，出错

	_, err = transfer.Conn.Read(transfer.Buf[:4])
	fmt.Println("读取客户端发送的数据")
	if err != nil {
		return
	}
	// 根据读到的buf长度，转成一个uint32类型，因为需要知道需要读多少个字节
	pkgLen := binary.BigEndian.Uint32(transfer.Buf[0:4])
	fmt.Printf("pkgLen: %v\n", pkgLen)
	// 根据pkgLen读取消息内容,将conn中的pkgLen长度的内容读到buf[]里面去
	n, err := transfer.Conn.Read(transfer.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Printf("读取信息失败: %v\n", err)
		return
	}

	// 将pkglen反序列化-> message.Message
	// 这里如果不是&mes，mes会是空的
	err = json.Unmarshal(transfer.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Printf("反序列化出错%v:", err)
		return
	}
	return
}

// SendPkg 用于发送数据包
func (transfer *Transfer) SendPkg(data []byte) (err error) {
	// 1.先发送一个长度给对方
	pkgLen := uint32(len(data))
	binary.BigEndian.PutUint32(transfer.Buf[0:4], pkgLen) // 相当于把这个pkgLen转化为Byte
	// 现在发送长度
	n, err := transfer.Conn.Write(transfer.Buf[:4])
	if n != 4 || err != nil {
		fmt.Printf("长度发送失败: %v\n", err)
		return
	}

	// 2.发送消息本身
	n, err = transfer.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("服务器向客户端发送消息失败", err)
		return
	}
	return
}
