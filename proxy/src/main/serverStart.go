package main

import (
	"context"
	"github.com/koebel217505/Project/Proxy/src/server"
	"github.com/koebel217505/Project/projCommon/projChannel"
	"github.com/koebel217505/Project/projCommon/projConfig"
	"github.com/koebel217505/Project/projCommon/projConvert"
	"github.com/koebel217505/Project/projCommon/projTcp"
	"github.com/koebel217505/Project/projCommon/projType"
	"github.com/koebel217505/Project/projCommon/projVar"
	"github.com/koebel217505/Project/proxy/src/clientIP"
	"log"
)

func serverStart(c context.Context, closeCh *projChannel.Channel) {
	projVar.Server = projTcp.NewBaseServer()

	EventHandler := projTcp.NewEventHandler(serverHandler.Event{}, nil)
	ConnectHandler := &serverHandler.Connect{}
	projVar.Server.SetEventHandler(EventHandler)
	projVar.Server.SetUserHandler(ConnectHandler)
	projVar.Server.SetServerHandler(ConnectHandler)
	projVar.Server.SetSessions(5)
	projVar.Server.SetClientsAddrArray(clientIP.ClientIPs)

	projVar.Server.SetServerAddr(projType.Addr{
		Name: "ProxyServer",
		IP:   projConfig.Config.ProxyIP,
		Port: int(projConvert.ConvInt32(projConfig.Config.ProxyPort)),
	})

	projVar.Server.Start(c)
	log.Println("Socket start at " + projConfig.Config.ProxyIP + ":" + projConfig.Config.ProxyPort)
}
