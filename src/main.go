package main

import (
	"LearnGo/src/test"
	"fmt"

	"github.com/panjf2000/gnet"
)

type echoServer struct {
	*gnet.EventServer
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	return
}

func main() {
	// echo := new(echoServer)
	// log.Fatal(gnet.Serve(echo, "tcp://:9000", gnet.WithMulticore(true)))
	// utils := &deploy.DBUtils{}
	// utils.Open("wgp", "127.0.0.1", 3306, "root", "will")
	// results := utils.Query("select * from game")
	// fmt.Println(results)

	// deploy.InstallServer(&deploy.IndexServer{})
	// deploy.Startup(1234)
	arr01 := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("len:%d %T\n", len(arr01), arr01)
	arr02 := arr01[4:6] // cap是数组最大长度
	var arr03 = arr02[:cap(arr02)]
	fmt.Printf("value:%d", arr03[3])

	test.TestSlice()
	err := test.TestXml("test/servlet.xml")
	if err != nil {
		fmt.Println(err)
	}
}
