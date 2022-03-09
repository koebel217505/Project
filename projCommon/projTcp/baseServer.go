package projTcp

import (
	"context"
	"fmt"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/koebel217505/Project/projCommon/projType"
	"log"
	"net"
	"sync"
)

type BaseServer Server

// Server 描述一个服务器的结构
type baseServer struct {
	listener        net.Listener   // 监听句柄
	wg              sync.WaitGroup // 等待所有goroutine结束
	Sessions        Sessions
	serverHandler   ServerHandler
	userHandler     UserHandler
	eventHandler    *EventHandler
	data            any
	serverAddr      projType.Addr
	clientAddrArray []projType.Addr
	//UIStatus      uint8
}

//var serverlogger *log.Logger

//func init() {
//	os.Mkdir("log", os.ModePerm)
//	file, err := os.OpenFile("./log/Server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 666)
//	if err != nil {
//		serverlogger.Fatal(err)
//	}
//	serverlogger = log.New(file, "", log.LstdFlags)
//	serverlogger.SetFlags(log.LstdFlags | log.Lshortfile)
//	file.Close()
//}

type Server interface {
	Start(c context.Context)
	Stop()
	Close()
	Kick(index int32)
	GetSessions() Sessions
	SetSessions(n int32)
	SetUserHandler(value UserHandler)
	SetServerHandler(value ServerHandler)
	SetEventHandler(*EventHandler)
	SetServerAddr(value projType.Addr)

	GetClientsAddrArray() []projType.Addr
	SetClientsAddrArray(value []projType.Addr)

	GetSessionID(ip string) int32

	OnUserConnect(s Session)
	OnUserDisconnect(s Session)

	OnServerInit(s BaseServer)
	OnServerDestroy(s BaseServer)

	SendMsgToAll(msgNo uint16, b []byte)
	SendMsgToSomeOne(ID int32, msgNo uint16, b []byte)
	SendMsgExclude(excludes []int32, msgNo uint16, b []byte)
	FindSessionByID(ID int32) (p Session)

	//GetUIStatus() uint8
	//SetUIStatus(value uint8)
}

func (bs *baseServer) SetUserHandler(value UserHandler) {
	bs.userHandler = value
}

func (bs *baseServer) SetServerHandler(value ServerHandler) {
	bs.serverHandler = value
}

func (bs *baseServer) SetEventHandler(value *EventHandler) {
	bs.eventHandler = value
}

func (bs *baseServer) GetSessions() Sessions {
	return bs.Sessions
}

func (bs *baseServer) SetSessions(n int32) {
	bs.Sessions = NewSessions(n)
}

func (bs *baseServer) SetServerAddr(value projType.Addr) {
	bs.serverAddr = value
}

func (bs *baseServer) GetClientsAddrArray() []projType.Addr {
	return bs.clientAddrArray
}

func (bs *baseServer) SetClientsAddrArray(value []projType.Addr) {
	bs.clientAddrArray = value
}

func (bs *baseServer) GetSessionID(ip string) int32 {
	host, _, _ := net.SplitHostPort(ip)
	for key, value := range bs.clientAddrArray {
		if host == value.IP {
			return int32(key)
		}
	}

	return -1
}

func (bs *baseServer) OnUserConnect(s Session) {
	//if bs.GetUIStatus() == projType.UIStatus_None {
	//	s.Close()
	//	fmt.Println("UIStart False")
	//	return
	//}

	if bs.clientAddrArray != nil && len(bs.clientAddrArray) != 0 {
		if id := bs.GetSessionID(s.GetConn().RemoteAddr().String()); id >= 0 {
			if bs.GetSessions().IsNil(int32(id)) {
				if e := bs.GetSessions().Put(int32(id), s); e != nil {
					s.Close()
					return
				}
			}
		} else {
			s.Close()
			return
		}
	} else {
		if _, err := bs.GetSessions().Add(s); err != nil {
			s.Close()
			return
		} /*else {
			bs.GetSessionMgr().Put(uint16(id), s)
		}*/
	}

	if bs.userHandler != nil {
		bs.userHandler.OnUserConnect(s)
	}

	log.Printf("Client [%d] %s connect LocalServer %s", s.GetID(), s.RemoteAddr(), s.LocalAddr())
}

func (bs *baseServer) OnUserDisconnect(s Session) {
	if !bs.GetSessions().IsNil(s.GetID()) {
		bs.GetSessions().Del(s.GetID())
	}

	if bs.userHandler != nil {
		bs.userHandler.OnUserDisconnect(s)
	}

	log.Printf("Client [%d] %s disconnect LocalServer %s", s.GetID(), s.RemoteAddr(), s.LocalAddr())
	s.Close()
}

func (bs *baseServer) OnServerInit(s BaseServer) {
	if bs.serverHandler != nil {
		bs.serverHandler.OnServerInit(s)
	}
}

func (bs *baseServer) OnServerDestroy(s BaseServer) {
	if bs.serverHandler != nil {
		bs.serverHandler.OnServerDestroy(s)
	}
}

// Start 开始服务
func (bs *baseServer) Start(c context.Context) {
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bs.serverAddr.IP, bs.serverAddr.Port))
		if err != nil {
			log.Println("Server Listen err:", err)
		} else {
			bs.OnServerInit(bs)

			log.Println("Server OnServerInit")
			bs.wg.Add(1)
			go func() {
				defer bs.wg.Done()

				for {
					conn, err := listener.Accept()
					if err != nil {
						log.Println("Server Accept err:", err)
						break
					}

					NewSession(conn, bs, &bs.wg, *gtype.NewBool(false), bs.eventHandler, make(chan func(), 10000)).Start()
				}
			}()

			bs.wg.Wait()
			bs.OnServerDestroy(bs)
			log.Println("Server OnServerDestroy")
			listener.Close()
			return
		}

		for {
			select {
			case <-c.Done():
				return
			}
		}
	}()
}

// Stop 停止服务
func (bs *baseServer) Stop() {
	bs.Close()
	bs.listener.Close()
}

// Close bla-bla
func (bs *baseServer) Close() {
	bs.GetSessions().Close()
}

func (bs *baseServer) Kick(ID int32) {

}

//func (s *baseServer) GetUIStatus() uint8 {
//	return s.UIStatus
//}
//
//func (s *baseServer) SetUIStatus(value uint8) {
//	s.UIStatus = value
//}

func (bs *baseServer) SendMsgToAll(msgNo uint16, b []byte) {
	bs.Sessions.SendMsgToAll(msgNo, b)
}

func (bs *baseServer) SendMsgToSomeOne(ID int32, msgNo uint16, b []byte) {
	bs.Sessions.SendMsgToOne(ID, msgNo, b)
}

func (bs *baseServer) SendMsgExclude(excludes []int32, msgNo uint16, b []byte) {
	bs.Sessions.SendMsgExclude(excludes, msgNo, b)
}

func (bs *baseServer) FindSessionByID(ID int32) (p Session) {
	return bs.Sessions.FindSessionByID(ID)
}

// NewBaseServer 创建一个Server, 返回*Server
func NewBaseServer() BaseServer { return &baseServer{} }
