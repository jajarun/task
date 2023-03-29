package base

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	ERROR_NONE    = 0
	ERROR_DEFAULT = 1
	ERROR_PARAM   = 101
)

type ControllerHandle struct {
	W    http.ResponseWriter
	R    *http.Request
	Path string
	//HandleFunc func()
}

type responseData struct {
	ErrorCode int
	Message   string
	Data      interface{}
}

func (c *ControllerHandle) Init() error {
	return nil
}

func (c *ControllerHandle) GetQuery(key string, defaultValue ...string) string {
	query := c.R.URL.Query()
	value := query.Get(key)
	if value == "" && defaultValue != nil {
		value = defaultValue[0]
	}
	return value
}

func (c *ControllerHandle) GetRawBody() string {
	s, _ := io.ReadAll(c.R.Body) //把	body 内容读入字符串 s
	return string(s)
}

func (c *ControllerHandle) GetMapBody() map[string]interface{} {
	body := c.GetRawBody()
	mapBody := make(map[string]interface{})
	_ = json.Unmarshal([]byte(body), &mapBody)
	return mapBody
}

func (c *ControllerHandle) ReturnError(errCode int, errMsg string) {
	r := responseData{
		ErrorCode: errCode,
		Message:   errMsg,
	}
	result, _ := json.Marshal(r)
	_, _ = io.WriteString(c.W, string(result))
}

func (c *ControllerHandle) ReturnData(data ...interface{}) {
	r := responseData{
		ErrorCode: 0,
	}
	if len(data) > 0 {
		r.Data = data[0]
	}
	result, _ := json.Marshal(r)
	_, _ = io.WriteString(c.W, string(result))
}
