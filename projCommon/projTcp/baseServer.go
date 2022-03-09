package projTcp

import (
	"context"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/koebel217505/Project/projCommon/projChannel"
	"github.com/spf13/cast"
	"log"
	"net"
	"sync"
)

// BaseServer 描述一个服务器的结构
type BaseServer struct {
	listener      net.Listener   // 监听句柄
	wg            sync.WaitGroup // 等待所有goroutine结束
	Sessions      *Sessions
	serverHandler ServerHandler
	userHandler   UserHandler
	eventHandler  *EventHandler
	data          any
	addr          AddrInfo
	permissions   []AddrInfo
	eventCh       *projChannel.Channel
	sendCh        *projChannel.Channel
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

func (bs *BaseServer) SetUserHandler(value UserHandler) {
	bs.userHandler = value
}

func (bs *BaseServer) SetServerHandler(value ServerHandler) {
	bs.serverHandler = value
}

func (bs *BaseServer) SetEventHandler(value *EventHandler) {
	bs.eventHandler = value
}

func (bs *BaseServer) GetSessions() *Sessions {
	return bs.Sessions
}

func (bs *BaseServer) SetSessions(n int32) {
	bs.Sessions = NewSessions(n, bs.sendCh)
}

func (bs *BaseServer) SetServerAddr(value AddrInfo) {
	bs.addr = value
}

func (bs *BaseServer) GetServerAddr(value AddrInfo) {
	bs.addr = value
}

func (bs *BaseServer) GetPermissions() []AddrInfo {
	return bs.permissions
}

func (bs *BaseServer) SetPermissions(value []AddrInfo) {
	if value != nil {
		bs.permissions = value
	}
}

func (bs *BaseServer) GetPermission(ip string) AddrInfo {
	host, _, _ := net.SplitHostPort(ip)
	for _, value := range bs.permissions {
		if host == cast.ToString(value.TCPAddr.IP) {
			return value
		}
	}

	return AddrInfo{}
}

func (bs *BaseServer) SetSendCh(sendCh *projChannel.Channel) {
	bs.sendCh = sendCh
}

func (bs *BaseServer) SetEventCh(eventCh *projChannel.Channel) {
	bs.eventCh = eventCh
}

func (bs *BaseServer) OnUserConnect(s *Session) {
	if len(bs.permissions) == 0 {
		permission := bs.GetPermission(s.GetConn().RemoteAddr().String())
		if permission.ID != 0 {
			if bs.GetSessions().IsNil(int32(permission.ID)) {
				if e := bs.GetSessions().Put(int32(permission.ID), s); e != nil {
					s.Close()
					return
				}
			} else {
				s.Close()
				return
			}
		} else {
			s.Close()
			return
		}
	} else {
		if _, err := bs.GetSessions().Add(s); err != nil {
			s.Close()
			return
		}
	}

	if bs.userHandler != nil {
		bs.userHandler.OnUserConnect(s)
	}

	log.Printf("Client [%d] %s connect LocalServer %s", s.GetID(), s.RemoteAddr(), s.LocalAddr())
}

func (bs *BaseServer) OnUserDisconnect(s *Session) {
	if !bs.GetSessions().IsNil(s.GetID()) {
		bs.GetSessions().Del(s.GetID())
	}

	if bs.userHandler != nil {
		bs.userHandler.OnUserDisconnect(s)
	}

	log.Printf("Client [%d] %s disconnect LocalServer %s", s.GetID(), s.RemoteAddr(), s.LocalAddr())
	s.Close()
}

func (bs *BaseServer) OnServerInit(s *BaseServer) {
	if bs.serverHandler != nil {
		bs.serverHandler.OnServerInit(s)
	}
}

func (bs *BaseServer) OnServerDestroy(s *BaseServer) {
	if bs.serverHandler != nil {
		bs.serverHandler.OnServerDestroy(s)
	}
}

// Start 开始服务
func (bs *BaseServer) Start(c context.Context) {
	go func() {
		listener, err := net.Listen("tcp", bs.addr.TCPAddr.String())
		if err != nil {
			log.Println("BaseServer Listen err:", err)
		} else {
			bs.OnServerInit(bs)

			log.Println("BaseServer OnServerInit")
			bs.wg.Add(1)
			go func() {
				defer bs.wg.Done()

				for {
					conn, err := listener.Accept()
					if err != nil {
						log.Println("Server Accept err:", err)
						break
					}

					NewSession(conn, bs, &bs.wg, *gtype.NewBool(false), bs.eventHandler, bs.eventCh).Start()
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
func (bs *BaseServer) Stop() {
	bs.Close()
	bs.listener.Close()
}

// Close bla-bla
func (bs *BaseServer) Close() {
	bs.GetSessions().Close()
}

func (bs *BaseServer) Kick(ID int32) {

}

func (bs *BaseServer) SendMsgToAll(msgNo uint16, b []byte) {
	bs.Sessions.SendMsgToAll(msgNo, b)
}

func (bs *BaseServer) SendMsgToOne(ID int32, msgNo uint16, b []byte) {
	bs.Sessions.SendMsgToOne(ID, msgNo, b)
}

func (bs *BaseServer) SendMsgExclude(excludes []int32, msgNo uint16, b []byte) {
	bs.Sessions.SendMsgExclude(excludes, msgNo, b)
}

func (bs *BaseServer) FindSessionByID(ID int32) (p *Session) {
	return bs.Sessions.FindSessionByID(ID)
}

// NewBaseServer 创建一个Server, 返回*IBaseServer
func NewBaseServer() *BaseServer { return &BaseServer{} }
