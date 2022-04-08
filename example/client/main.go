/**
 * @Author: Jason
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2022/4/8 12:57
 */

package main

import (
	"fmt"
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/server/ProtoMsg"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:29000")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go func() {
		for {
			dp := easysocket.NewDataPack()

			ping := &ProtoMsg.C2S_Ping{
				Ping: "ping",
			}

			buffer, _ := proto.Marshal(ping)

			msg := dp.Pack(easysocket.NewMessage(int32(ProtoMsg.CMD_PING), buffer))

			_, err := conn.Write(msg)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			headData := make([]byte, dp.GetHeadLen())
			_, err = io.ReadFull(conn, headData)

			if err != nil {
				fmt.Println(err.Error())
				break
			}

			msgHead := dp.UnPack(headData)

			if msgHead.GetDataLen() > 0 {
				msg := msgHead.(*easysocket.Message)

				data := make([]byte, msg.GetDataLen())
				_, err := io.ReadFull(conn, data)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				msg.SetData(data)

				fmt.Println("===> ping recv msgId:", msg.GetMsgId(), ", len:", msg.GetDataLen(), ", data:", string(data))
			}

			time.Sleep(1 * time.Second)
		}
	}()

	select {}
}
