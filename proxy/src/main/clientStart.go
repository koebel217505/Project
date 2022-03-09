package main

import (
	"context"
	"github.com/koebel217505/Project/Proxy/src/client"
	"github.com/koebel217505/Project/projCommon/projChannel"
	"github.com/koebel217505/Project/projCommon/projConfig"
	"github.com/koebel217505/Project/projCommon/projConvert"
	"github.com/koebel217505/Project/projCommon/projTcp"
	"github.com/koebel217505/Project/projCommon/projType"
	"github.com/koebel217505/Project/projCommon/projVar"
	"time"
)

func clientStart(c context.Context, closeCh *projChannel.Channel) {
	time.Sleep(1 * time.Second)
	projVar.Client = projTcp.NewBaseClient()

	EventHandler := projTcp.NewEventHandler(clientHandler.Event{}, nil)
	ConnectHandler := &clientHandler.Connect{}
	projVar.Client.SetEventHandler(EventHandler)
	projVar.Client.SetUserHandler(ConnectHandler)
	projVar.Client.SetSessions(5)
	projVar.Client.SetServerAddrArray(append([]projType.Addr{},
		projType.Addr{
			Name: "ProxyClient",
			IP:   projConfig.Config.ProxyIP,
			Port: int(projConvert.ConvInt32(projConfig.Config.ProxyPort)),
		},
		projType.Addr{
			Name: "ProxyClient",
			IP:   projConfig.Config.ProxyIP,
			Port: int(projConvert.ConvInt32(projConfig.Config.ProxyPort)),
		}))

	projVar.Client.Connect(c)
}
