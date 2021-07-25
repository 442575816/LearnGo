package deploy

import (
	"io"
	"net/http"
)

type IndexServer struct {
}

func (index *IndexServer) HandleRequest(w http.ResponseWriter, request *http.Request) {
	io.WriteString(w, "Hello World")
}

func (index *IndexServer) Init() {
	AddHandler("/index", index.HandleRequest)
}
