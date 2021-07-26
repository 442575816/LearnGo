package servlet

import (
	internalErrors "LearnGo/src/errors"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

const (
	// 是否启用压缩了
	ActionCompress string = "compress"
)

type ServerProtocol int

const (
	TCP ServerProtocol = iota
	HTTP
	UDP
	WEBSOCKET
)

// Servlet配置接口
type ServletConfig interface {
	Get(key string) interface{}
	GetSessionTickTime() int
	GetSessionTimeoutMillis() int
	GetSessionEmptyTimeoutMillis() int
	GetSessionInvalidateMillis() int
	GetSessionNextDayInvalidateMillis() int
}

// Servlet上下文接口
type ServletContext interface {
	Get(key string, defaultValue interface{}) interface{}
	Set(key string, value interface{})
	Delete(key string)
}

// 请求Request接口
type Request interface {
	Command() string
	RequestId() int
	Content() []byte
	Ip() string
	GetParameterValues(key string) []string
	ParameterMap() map[string][]string
	CreateTime() time.Time
	GetHeader(key string) string
	SetSessionId(key string)
	Protocol() ServerProtocol
	GetSession(allowCreate bool) *Session
	GetNewSession() *Session
}

// 请求Response接口
type Response interface {
	Write(buff []byte)
	AddHeader(name string, value string)
	AddCookie(name string, value string)
	SetHttpStatus(status int)
	Protocol() ServerProtocol
	MarkClose()
}

// Session定义
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

// 推送通道的定义
type Push interface {
	Push(command string, bytes []byte)
	IsPushable() bool
	Discard()
	Heartbeat()
	Protocol() ServerProtocol
}

// 标准Servlet定义
type Servlet interface {
	// 初始化Servlet
	Init(config ServletConfig, context ServletContext)

	// 处理请求
	Service(request Request, response Response) error

	// 增加handler
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

type XmlServletConfig struct {
	config map[string]interface{}
}

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

func (config *XmlServletConfig) Get(key string) interface{} {
	v, ok := config.config[key]
	if ok {
		return v
	}
	return nil
}
func (config *XmlServletConfig) GetSessionTickTime() int {
	v := config.Get("sessionTickTime")
	if v != nil {
		return v.(int)
	}
	return 20000
}
func (config *XmlServletConfig) GetSessionTimeoutMillis() int {
	v := config.Get("sessionTimeoutTime")
	if v != nil {
		return v.(int)
	}
	return 180000
}
func (config *XmlServletConfig) GetSessionEmptyTimeoutMillis() int {
	v := config.Get("sessionEmptyTimeoutTime")
	if v != nil {
		return v.(int)
	}
	return 40000
}
func (config *XmlServletConfig) GetSessionInvalidateMillis() int {
	v := config.Get("sessionInvalidateMillis")
	if v != nil {
		return v.(int)
	}
	return 86400000
}
func (config *XmlServletConfig) GetSessionNextDayInvalidateMillis() int {
	v := config.Get("sessionNextDayInvalidateMillis")
	if v != nil {
		return v.(int)
	}
	return 1800000
}
