package main

import (
	"net/http"
	"task/route"
)

func main() {
	//ws.Run()
	route.WebRouteInit()
	route.AccRouteInit()

	_ = http.ListenAndServe(":8181", nil)
}
