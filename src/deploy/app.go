package deploy

import (
	"io"
	"net/http"
)

type AppServer struct {
}

func (index *AppServer) getAppList(w http.ResponseWriter, request *http.Request) {
	io.WriteString(w, "Hello World")
}

func (index *AppServer) Init() {
	AddHandler("/getAppList", index.getAppList)
}
