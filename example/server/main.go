/**
 * @Author: Jason
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2022/4/8 9:56
 */

package main

import "github.com/dyowoo/easysocket"

func main() {
	s := easysocket.NewServer("server", easysocket.WsServer, "0.0.0.0", 29000)

	go s.Serve()

	select {}
}
