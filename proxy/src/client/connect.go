package clientHandler

import (
	"github.com/koebel217505/Project/projCommon/projTcp"
)

type Connect struct{}

// OnUserConnect 客户端连接事件
func (ser *Connect) OnUserConnect(s projTcp.Session) {

}

// OnUserDisconnect 客户端断开连接事件
func (ser *Connect) OnUserDisconnect(s projTcp.Session) {

}

// OnServerInit bla-bla
func (ser *Connect) OnServerInit(s projTcp.BaseServer) {

}

// OnServerDestroy bla-bla
func (ser *Connect) OnServerDestroy(s projTcp.BaseServer) {

}
