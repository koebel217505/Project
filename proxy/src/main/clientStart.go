package main

import (
	"context"
	"github.com/koebel217505/Project/projCommon/projChannel"
	"github.com/koebel217505/Project/projCommon/projTcp"
	"github.com/koebel217505/Project/projCommon/projVar"
	clientHandler "github.com/koebel217505/Project/proxy/src/client"
	"net"
)

func clientStart(c context.Context, closeCh *projChannel.Channel, sendCh *projChannel.Channel, eventCh *projChannel.Channel) {
	//tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:7099")
	//fmt.Println(tcpAddr.IP)
	//fmt.Println(tcpAddr.Port)
	//fmt.Println(tcpAddr.Zone)
	//fmt.Println(tcpAddr.String())
	//fmt.Println(tcpAddr.Network())
	//fmt.Println(err)
	//time.Sleep(1 * time.Second)
	projVar.Client = projTcp.NewBaseClient()

	EventHandler := projTcp.NewEventHandler(clientHandler.Event{}, nil)
	ConnectHandler := &clientHandler.Connect{}
	projVar.Client.SetEventHandler(EventHandler)
	projVar.Client.SetUserHandler(ConnectHandler)
	projVar.Server.SetSendCh(sendCh)
	projVar.Server.SetEventCh(eventCh)
	projVar.Client.SetSessions(5)

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:7099")

	projVar.Client.SetPermissions(append([]projTcp.AddrInfo{},
		projTcp.AddrInfo{
			ID:      1,
			Name:    "ProxyClient",
			TCPAddr: tcpAddr,
		}))

	projVar.Client.Connect(c)
}
