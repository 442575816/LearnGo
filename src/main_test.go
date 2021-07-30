package main_test

import (
	"LearnGo/src/buffer"
	"LearnGo/src/servlet"
	"testing"
	"time"

	//"fmt"
	"strings"

	//"LearnGo/src/test"
	//"fmt"

	"github.com/panjf2000/gnet"
)

type echoServer struct {
	*gnet.EventServer
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	return
}

type MyServerHandler struct {
	Servlet servlet.Servlet
	ServletConfig servlet.ServletConfig
	ServletContext servlet.ServletContext
	Handler servlet.InboundHandler
}

type MessageDecoder struct {
}

type MessageHandler struct {
	Servlet servlet.Servlet
	ServletConfig servlet.ServletConfig
	ServletContext servlet.ServletContext
	servlet.InboundHandlerAdapter
}

func (m *MessageDecoder) CallDecode(ctx servlet.ConnHandlerContext, in *buffer.ByteBuf, output *[]interface{}) {
	if in.ReadableBytes() < 4 {
		return
	}

	dataLen := in.GetUInt32()
	if in.ReadableBytes() < dataLen + 4 {
		return
	}

	in.SkipBytes(4)

	var message servlet.RequestMessage
	bytes := in.ReadBytes(32)
	message.Command = strings.Trim(string(bytes), "\x00")
	message.RequestId = int(in.ReadInt32())
	message.Content = in.ReadBytes(dataLen - 36)

	*output = append(*output, message)
}

func (i *MessageHandler) FireMessageRead(context servlet.ConnHandlerContext, msg interface{}) {
	context.FireWrite([]byte("helloworld"))
}

func (m *MyServerHandler) Init(servlet servlet.Servlet, config servlet.ServletConfig, ctx servlet.ServletContext) {
	m.Servlet = servlet
	m.ServletConfig = config
	m.ServletContext = ctx
	m.Handler = &MessageHandler{Servlet: servlet, ServletConfig: config, ServletContext: ctx}
}

func (m *MyServerHandler) InitConn(conn gnet.Conn) {
	v, ok := conn.Context().(*servlet.ConnPipeline)
	if ok {
		var coder servlet.InboundHandler = &servlet.ByteToMessageDecoder {Decoder: &MessageDecoder{}}
		v.AddLast("decoder", coder)
		v.AddLast("messageHandler", m.Handler)
	}
}


func TestTcp(t *testing.T) {
	// echo := new(echoServer)
	// log.Fatal(gnet.Serve(echo, "tcp://:9000", gnet.WithMulticore(true)))
	// utils := &deploy.DBUtils{}
	// utils.Open("wgp", "127.0.0.1", 3306, "root", "will")
	// results := utils.Query("select * from game")
	// fmt.Println(results)

	// deploy.InstallServer(&deploy.IndexServer{})
	// deploy.Startup(1234)

	t.Run("tcp", func(t *testing.T) {
		var handler MyServerHandler


		go servlet.StartTcpServer(&handler)
	})
	t.Run("finish", func(t *testing.T) {
		time.Sleep(time.Minute * 2)
		t.FailNow()
	})
}
