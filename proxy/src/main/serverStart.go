package main

import (
	"context"
	"github.com/koebel217505/Project/Proxy/src/server"
	"github.com/koebel217505/Project/projCommon/projChannel"
	"github.com/koebel217505/Project/projCommon/projConfig"
	"github.com/koebel217505/Project/projCommon/projTcp"
	"github.com/koebel217505/Project/projCommon/projVar"
	"github.com/koebel217505/Project/proxy/src/clientIP"
	"log"
	"net"
)

func serverStart(c context.Context, closeCh *projChannel.Channel, sendCh *projChannel.Channel, eventCh *projChannel.Channel) {
	projVar.Server = projTcp.NewBaseServer()

	EventHandler := projTcp.NewEventHandler(serverHandler.Event{}, nil)
	ConnectHandler := &serverHandler.Connect{}
	projVar.Server.SetEventHandler(EventHandler)
	projVar.Server.SetUserHandler(ConnectHandler)
	projVar.Server.SetServerHandler(ConnectHandler)
	projVar.Server.SetSendCh(sendCh)
	projVar.Server.SetEventCh(eventCh)
	projVar.Server.SetSessions(5)
	projVar.Server.SetPermissions(clientIP.ClientIPs)

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:7099")

	projVar.Server.SetServerAddr(
		projTcp.AddrInfo{
			ID:      1,
			Name:    "ProxyServer",
			TCPAddr: tcpAddr,
		})

	tcpAddr, _ = net.ResolveTCPAddr("tcp4", "127.0.0.1:7099")
	projVar.Server.SetPermissions(append([]projTcp.AddrInfo{},
		projTcp.AddrInfo{
			ID:      1,
			Name:    "ProxyClient",
			TCPAddr: tcpAddr,
		}))

	projVar.Server.Start(c)
	log.Println("Socket start at " + projConfig.Config.ProxyIP + ":" + projConfig.Config.ProxyPort)
}
