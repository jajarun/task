package main

import (
	"net/http"
	"task/route"
)

func main() {
	//ws.Run()
	route.WebRouteInit()

	_ = http.ListenAndServe(":8181", nil)
}
