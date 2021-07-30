package servlet

import (
	internalErrors "LearnGo/src/errors"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/panjf2000/gnet"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	// ActionCompress 是否启用压缩了
	ActionCompress string = "compress"
)

type ServerProtocol int

const (
	TCP ServerProtocol = iota
	HTTP
	UDP
	WEBSOCKET
)

// ServletConfig Servlet配置接口
type ServletConfig interface {
	Get(key string) interface{}
	GetSessionTickTime() int
	GetSessionTimeoutMillis() int
	GetSessionEmptyTimeoutMillis() int
	GetSessionInvalidateMillis() int
	GetSessionNextDayInvalidateMillis() int
}

// ServletContext Servlet上下文接口
type ServletContext interface {
	Get(key string, defaultValue interface{}) interface{}
	Set(key string, value interface{})
	Delete(key string)
}

// Request 请求Request接口
type Request interface {
	Command() string
	RequestId() int
	Content() []byte
	Ip() string
	GetParameterValues(key string) []string
	ParameterMap() map[string][]string
	CreateTime() time.Time
	GetHeader(key string) (string, error)
	SetSessionId(key string)
	Protocol() ServerProtocol
	GetSession(allowCreate bool) *Session
	GetNewSession() (*Session, error)
}

// Response 请求Response接口
type Response interface {
	Write(buff []byte)
	AddHeader(name string, value string) error
	AddCookie(name string, value string) error
	SetHttpStatus(status int) error
	Protocol() ServerProtocol
	MarkClose()
}

// Session Session定义
type Session interface {
	Id() string
	Get(key string) interface{}
	Set(key string, value interface{})
	Delete(key string)
	Access()
	IsValid() bool
	IsExpire() bool
	IsActive() bool
	IsInvalidate() bool
	ReActive()
	IsEmpty() bool
	CheckAlive() bool
	SetPush(push *Push)
	GetPush() *Push
	MarkDiscard()
}

// Push 推送通道的定义
type Push interface {
	Push(command string, bytes []byte)
	IsPushable() bool
	Discard()
	Heartbeat()
	Protocol() ServerProtocol
}

// Servlet 标准Servlet定义
type Servlet interface {
	// Init 初始化Servlet
	Init(config ServletConfig, context ServletContext)

	// Service 处理请求
	Service(request Request, response Response) error

	// AddHandler 增加handler
	AddHandler(command string, handler func(Request, Response)) error
}

type DispatchServlet struct {
	config    ServletConfig
	context   ServletContext
	handleMap map[string]func(Request, Response)
	compress  bool
}

func (servlet *DispatchServlet) Init(config ServletConfig, context ServletContext) {
	servlet.config = config
	servlet.context = context
	servlet.handleMap = make(map[string]func(Request, Response))

	servlet.initCompress()
}

func (servlet *DispatchServlet) AddHandler(command string, handler func(Request, Response)) (err error) {
	_, ok := servlet.handleMap[command]
	if ok {
		return internalErrors.HandleAlreadyExists
	}

	servlet.handleMap[command] = handler
	return nil
}

func (servlet *DispatchServlet) Service(request Request, response Response) (err error) {
	command := request.Command()
	handler, ok := servlet.handleMap[command]
	if ok {
		handler(request, response)
		return nil
	}

	return errors.New(fmt.Sprintf("%s does not hava handler", command))
}

func (servlet *DispatchServlet) initCompress() {
	value := servlet.config.Get(ActionCompress)
	if value != "" {
		servlet.compress = "true" == value
	}
}

// XmlServletConfig xml类型配置文件
type XmlServletConfig struct {
	config map[string]interface{}
}

// NewXmlServletConfig 创建新的xml配置文件
func NewXmlServletConfig(path string) *XmlServletConfig {
	config := XmlServletConfig{}
	config.config = make(map[string]interface{})
	config.parse(path)

	return &config
}

func (servlet *XmlServletConfig) parse(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	type Property struct {
		XMLName xml.Name `xml:"property"`
		Name    string   `xml:"name,attr"`
		Type    string   `xml:"type,attr"`
		Value   string   `xml:",chardata"`
	}
	type InitParam struct {
		Props []Property `xml:"props>property"`
	}

	type Params struct {
		XMLName   xml.Name  `xml:"servlet"`
		InitParam InitParam `xml:"init-param"`
	}

	params := Params{}
	err = xml.Unmarshal(content, &params)
	if err != nil {
		return err
	}

	for _, p := range params.InitParam.Props {
		switch p.Type {
		case "bool":
			v, err := strconv.ParseBool(p.Value)
			if err == nil {
				servlet.config[p.Name] = v
			}
		case "string":
			servlet.config[p.Name] = p.Value
		case "int":
			v, err := strconv.Atoi(p.Value)
			if err == nil {
				servlet.config[p.Name] = v
			}
		}
	}

	return nil
}

func (servlet *XmlServletConfig) Get(key string) interface{} {
	v, ok := servlet.config[key]
	if ok {
		return v
	}
	return nil
}
func (servlet *XmlServletConfig) GetSessionTickTime() int {
	v := servlet.Get("sessionTickTime")
	if v != nil {
		return v.(int)
	}
	return 20000
}
func (servlet *XmlServletConfig) GetSessionTimeoutMillis() int {
	v := servlet.Get("sessionTimeoutTime")
	if v != nil {
		return v.(int)
	}
	return 180000
}
func (servlet *XmlServletConfig) GetSessionEmptyTimeoutMillis() int {
	v := servlet.Get("sessionEmptyTimeoutTime")
	if v != nil {
		return v.(int)
	}
	return 40000
}
func (servlet *XmlServletConfig) GetSessionInvalidateMillis() int {
	v := servlet.Get("sessionInvalidateMillis")
	if v != nil {
		return v.(int)
	}
	return 86400000
}
func (servlet *XmlServletConfig) GetSessionNextDayInvalidateMillis() int {
	v := servlet.Get("sessionNextDayInvalidateMillis")
	if v != nil {
		return v.(int)
	}
	return 1800000
}

type DefaultServletContext struct {
	contextMap map[string]interface{}
}

func (d *DefaultServletContext) Get(key string, defaultValue interface{}) interface{} {
	value, ok := d.contextMap[key]
	if ok {
		return value
	}
	return defaultValue
}

func (d *DefaultServletContext) Set(key string, value interface{}) {
	d.contextMap[key] = value
}

func (d *DefaultServletContext) Delete(key string) {
	delete(d.contextMap, key)
}

// RequestMessage 请求消息
type RequestMessage struct {
	RequestId int
	Command string
	Content []byte
	SessionId string
}

// TcpRequest tcp请求
type TcpRequest struct {
	command string
	requestId int
	content []byte
	createTime time.Time
	sessionId string
	parseFlag bool
	paramMap map[string][]string
	conn gnet.Conn
}

func (t *TcpRequest) Command() string {
	return t.command
}

func (t *TcpRequest) RequestId() int {
	return t.requestId
}

func (t *TcpRequest) Content() []byte {
	return t.content
}

func (t *TcpRequest) Ip() string {
	return t.conn.RemoteAddr().String()
}

func (t *TcpRequest) GetParameterValues(key string) []string {
	if !t.parseFlag {
		t.parseParam()
	}
	return t.paramMap[key]
}

func (t *TcpRequest) ParameterMap() map[string][]string {
	if !t.parseFlag {
		t.parseParam()
	}
	return t.paramMap
}

func (t *TcpRequest) CreateTime() time.Time {
	return t.createTime
}

func (t *TcpRequest) GetHeader(key string) (string, error) {
	return "", internalErrors.HandleAlreadyExists
}

func (t *TcpRequest) SetSessionId(key string) {
	t.sessionId = key
	t.conn.SetContext(key)
}

func (t *TcpRequest) Protocol() ServerProtocol {
	return TCP
}

func (t *TcpRequest) GetSession(allowCreate bool) *Session {
	return nil
}

func (t *TcpRequest) GetNewSession() (*Session, error) {
	return nil, internalErrors.NotSupport
}

func (t *TcpRequest) parseParam() {
	if t.parseFlag {
		return
	}

	if t.content == nil {
		return
	}

	str := string(t.content)
	defer func() {
		if err := recover(); err != nil {
			// 打印异常，关闭资源，退出此函数
			fmt.Println(err)
			// Ignore 忽略
		}
	}()
	t.parseParamInternal(str)
}

func (t *TcpRequest) parseParamInternal(content string)  {
	paramMap := make(map[string][]string)
	strs := strings.Split(content, "&")
	for _, value := range strs {
		index := strings.Index(value, "=")
		if index == -1 {
			paramMap[value] = nil
		} else {
			k := value[:index]
			v := value[index + 1:]
			tmpValue, ok := paramMap[k]
			if ok {
				if tmpValue == nil || len(tmpValue) == 0 {
					paramMap[k] = []string {v}
				} else {
					tmpValue = append(tmpValue, v)
					paramMap[k] = tmpValue
				}
			} else {
				paramMap[k] = []string {v}
			}
		}
	}
	t.paramMap = paramMap
}

func NewTcpquest(conn gnet.Conn, context ServletContext, message RequestMessage) Request  {
	var request TcpRequest
	request.requestId = message.RequestId
	request.command = message.Command
	request.content = message.Content
	request.sessionId = message.SessionId
	request.createTime = time.Now()
	request.conn = conn

	return &request
}

type TcpResponse struct {
	conn gnet.Conn
	closeFlag bool
}

func (t *TcpResponse) Write(buff []byte) {
	t.conn.AsyncWrite(buff)
}

func (t *TcpResponse) AddHeader(name string, value string) error {
	return internalErrors.NotSupport
}

func (t *TcpResponse) AddCookie(name string, value string) error {
	return internalErrors.NotSupport
}

func (t *TcpResponse) SetHttpStatus(status int) error {
	return internalErrors.NotSupport
}

func (t *TcpResponse) Protocol() ServerProtocol {
	return TCP
}

func (t *TcpResponse) MarkClose() {
	t.closeFlag = true
}

func NewTcpResponse(conn gnet.Conn) Response {
	var response TcpResponse

	response.conn = conn
	response.closeFlag = false

	return &response
}
