package deploy

import (
	"io"
	"net/http"
)

func test(w http.ResponseWriter, request *http.Request) {
	io.WriteString(w, "Test World")
}

func Init() {
	AddHandler("/test", test)
}
