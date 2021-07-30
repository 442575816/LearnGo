package servlet

import (
	"LearnGo/src/buffer"
	"errors"
	"github.com/panjf2000/gnet"
	"log"
)

type TcpServer struct {
	*gnet.EventServer
	onNewConn func(conn gnet.Conn)
	InitConn func(conn gnet.Conn)
}

type ConnPipeline struct {
	conn gnet.Conn
	Head *ConnHandlerContext
	Tail *ConnHandlerContext
}

type ConnHandlerContext struct {
	Pipeline *ConnPipeline
	Name string
	Handler interface{}
	Next *ConnHandlerContext
	Prev *ConnHandlerContext
}

func NewTcpServer() *TcpServer {
	tcpServer := TcpServer{}
	tcpServer.onNewConn = func(conn gnet.Conn) {
		var pipeline = NewConnPipeline(conn)
		conn.SetContext(pipeline)
	}
	return &tcpServer
}

func isConnHandler(handler interface{}) bool {
	switch handler.(type) {
	case InboundHandler:
	case OutboundHandler:
	default:
		return false
	}
	return true
}

func NewConnPipeline(conn gnet.Conn) *ConnPipeline {
	var pipeline ConnPipeline
	pipeline.conn = conn
	pipeline.Head = &ConnHandlerContext{Name: "internal_head_handler", Handler: nil, Pipeline: &pipeline}
	pipeline.Tail = &ConnHandlerContext{Name: "internal_tail_handler", Handler: nil, Pipeline: &pipeline}

	pipeline.Head.Next = pipeline.Tail
	pipeline.Tail.Prev = pipeline.Head

	return &pipeline
}

func (pipeline *ConnPipeline) AddFirst(name string, handler interface{}) error  {
	if !isConnHandler(handler) {
		return errors.New("only support inboundhandler or outboundhandler")
	}

	context := ConnHandlerContext{Name: name, Handler: handler, Pipeline: pipeline}
	context.Next = pipeline.Head.Next
	context.Prev = pipeline.Head
	pipeline.Head.Next.Prev = &context
	pipeline.Head.Next = &context

	return nil
}

func (pipeline *ConnPipeline) AddLast(name string, handler interface{}) error {
	if !isConnHandler(handler) {
		return errors.New("only support inboundhandler or outboundhandler")
	}

	context := ConnHandlerContext{Name: name, Handler: handler, Pipeline: pipeline}
	context.Prev = pipeline.Tail.Prev
	context.Next = pipeline.Tail
	pipeline.Tail.Prev.Next = &context
	pipeline.Tail.Prev = &context

	return nil
}

func (pipeline *ConnPipeline) Remove(name string) {
	context := pipeline.Head
	for context != nil && context.Name != name {
		context = context.Next
	}

	if context != nil {
		var prev = context.Prev
		var next = context.Next
		if prev != nil {
			prev.Next = next
		}
		if next != nil {
			next.Prev = prev
		}
	}
}

func (pipeline *ConnPipeline) RemoveByHandler(handler interface{}) error {
	if !isConnHandler(handler) {
		return errors.New("only support inboundhandler or outboundhandler")
	}
	
	context := pipeline.Head
	for context != nil && context.Handler != handler {
		context = context.Next
	}

	if context != nil {
		var prev = context.Prev
		var next = context.Next
		if prev != nil {
			prev.Next = next
		}
		if next != nil {
			next.Prev = prev
		}
	}
	return nil
}

type InboundHandler interface {
	FireConnOpen(context ConnHandlerContext)
	FireMessageRead(context ConnHandlerContext, msg interface{})
	FireConnClose(context ConnHandlerContext, err error)
}

type OutboundHandler interface {
	FireWrite(context ConnHandlerContext, msg interface{})
}

type InboundHandlerAdapter struct {
	ByteBuf *buffer.ByteBuf

}

func (i *InboundHandlerAdapter) FireConnOpen(context ConnHandlerContext) {
	context.FireConnOpen()
}

func (i *InboundHandlerAdapter) FireMessageRead(context ConnHandlerContext, msg interface{}) {
	context.FireMessageRead(msg)
}

func (i *InboundHandlerAdapter) FireConnClose(context ConnHandlerContext, err error) {
	context.FireConnClose(err)
}

type ByteToMessageDecoder struct {
	InboundHandlerAdapter
	Decoder
	outputList []interface{}
}

func (b *ByteToMessageDecoder) FireMessageRead(context ConnHandlerContext, msg interface{}) {
	v, ok := msg.(*buffer.ByteBuffer)
	if ok {
		if b.ByteBuf == nil {
			b.ByteBuf = buffer.New(64, buffer.BigEndian)
		}
		b.ByteBuf.WriteBytes(v.Data)

		b.CallDecode(context, b.ByteBuf, &b.outputList)

		if len(b.outputList) > 0 {
			for i:=0; i < len(b.outputList); i++ {
				context.FireMessageRead(b.outputList[i])
			}
			b.outputList = b.outputList[:0]
		}
		if b.ByteBuf.ReadableBytes() == 0 {
			b.ByteBuf.Reset()
		}
	}
}

type Decoder interface {
	CallDecode(ctx ConnHandlerContext, in *buffer.ByteBuf, output *[]interface{})
}

func (es *TcpServer) OnInitComplete(svr gnet.Server) (action gnet.Action) {
	return
}

func (es *TcpServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	es.onNewConn(c)
	es.InitConn(c)
	pipeline, ok := c.Context().(*ConnPipeline)
	if ok {
		pipeline.Head.FireConnOpen()
	}
	return
}

func (es *TcpServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	pipeline, ok := c.Context().(*ConnPipeline)
	if ok {
		pipeline.Head.FireConnClose(err)
	}
	return
}

func (es *TcpServer) PreWrite() {
}

func (es *TcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	pipeline, ok := c.Context().(*ConnPipeline)
	if ok {
		pipeline.Head.FireMessageRead(&buffer.ByteBuffer { Data: frame })
	}
	return
}

func (c *ConnHandlerContext) FireConnOpen() {
	h := c.getNextInboundHandlerContext()
	if h != nil {
		h.Handler.(InboundHandler).FireConnOpen(*h)
	}
}

func (c *ConnHandlerContext) FireMessageRead(msg interface{}) {
	h := c.getNextInboundHandlerContext()
	if h != nil {
		h.Handler.(InboundHandler).FireMessageRead(*h, msg)
	}
}

func (c *ConnHandlerContext) FireConnClose(err error) {
	h := c.getNextInboundHandlerContext()
	if h != nil {
		h.Handler.(InboundHandler).FireConnClose(*h, err)
	}
}

func (c *ConnHandlerContext) FireWrite(msg interface{}) {
	h := c.getPrevOutboundHandlerContext()
	if h != nil {
		h.Handler.(OutboundHandler).FireWrite(*h, msg)
		return
	}

	c.Pipeline.conn.AsyncWrite(msg.([]byte))
}

func (c *ConnHandlerContext) getNextInboundHandlerContext() *ConnHandlerContext {
	context := c.Next
	for context != nil {
		_, ok := context.Handler.(InboundHandler)
		if ok {
			return context
		}
		context = context.Next
	}
	return context
}

func (c *ConnHandlerContext) getPrevOutboundHandlerContext() *ConnHandlerContext {
	context := c.Prev
	for context != nil {
		_, ok := context.Handler.(OutboundHandler)
		if ok {
			return context
		}
		context = context.Prev
	}
	return context
}

type TcpServerHandler interface {
	Init(servlet Servlet, config ServletConfig, ctx ServletContext)
	InitConn(conn gnet.Conn)
}

func StartTcpServer(handler TcpServerHandler) {
	var servlet Servlet = &DispatchServlet{}
	var servletConfig = NewXmlServletConfig("server.xml")
	var context = &DefaultServletContext{}

	servlet.Init(servletConfig, context)
	handler.Init(servlet, servletConfig, context)

	var tcpServer = NewTcpServer()
	tcpServer.InitConn = handler.InitConn
	log.Fatal(gnet.Serve(tcpServer, "tcp://:9000", gnet.WithMulticore(true), gnet.WithReusePort(true)))
	log.Println("start suss")
}