package server

import (
	"github.com/emicklei/go-restful/v3"
)

type Server interface {
	Container() *restful.Container
	Start() //启动服务
	Logger
}
