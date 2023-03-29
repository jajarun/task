package route

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"task/controller/base"
	"task/controller/index"
	"task/controller/user"
)

var routeHandles = make(map[string]func(ch *base.ControllerHandle))

func setHandle(route string, handleFunc func(ch *base.ControllerHandle)) {
	routeHandles[route] = handleFunc
}

func initHandles() {
	setHandle("index/index", index.Index)
	setHandle("user/info", user.Info)
	setHandle("user/login", user.Login)
	setHandle("user/register", user.Register)
}

func WebRouteInit() {

	initHandles()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//value := r.URL.Query()
		path := strings.TrimLeft(r.URL.Path, "/")
		path = strings.ToLower(path)
		if path == "" {
			path = "index/index"
		}

		handle, ok := routeHandles[path]
		if !ok {
			_, _ = io.WriteString(w, "404")
		}
		controllerHandle := base.ControllerHandle{
			W:    w,
			R:    r,
			Path: path,
			//HandleFunc: handle
		}
		err := controllerHandle.Init()
		if err != nil {
			fmt.Println(err)
		} else {
			handle(&controllerHandle)
		}
	})
}
