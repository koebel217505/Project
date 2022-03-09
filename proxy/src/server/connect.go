// Package serverHandler 連接事件
package serverHandler

import (
	"fmt"

	"github.com/koebel217505/Project/projCommon/projTcp"
	"github.com/koebel217505/Project/projCommon/projVar"
)

type Connect struct{}

// OnUserConnect 客户端连接事件
func (ser *Connect) OnUserConnect(s *projTcp.Session) {
	if projVar.UI != nil {
		projVar.UI.Eval(fmt.Sprintf(`app.$data.items.push({ID:%d,IP:"%s",Status:"已連線"})`, s.GetID(), s.GetConn().RemoteAddr().String()))
	}
}

// OnUserDisconnect 客户端断开连接事件
func (ser *Connect) OnUserDisconnect(s *projTcp.Session) {
	if projVar.UI != nil {
		projVar.UI.Eval(fmt.Sprintf(`app.setItems({ID:%d,IP:"%s",Status:"已斷線"})`, s.GetID(), s.GetConn().RemoteAddr().String()))
	}
}

// OnServerInit bla-bla
func (ser *Connect) OnServerInit(s *projTcp.BaseServer) {

}

// OnServerDestroy bla-bla
func (ser *Connect) OnServerDestroy(s *projTcp.BaseServer) {

}
