package projTcp

import (
	"context"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/koebel217505/Project/projCommon/projChannel"
	"log"
	"net"
	"sync"
	"time"
)

// BaseClient 客户端描述
type BaseClient struct {
	//tcpClientState
	reConnectSecond time.Duration
	userHandler     UserHandler
	eventHandler    *EventHandler
	//closeCh         projChannel.Channel
	permissions []AddrInfo
	Sessions    *Sessions
	eventCh     *projChannel.Channel
	sendCh      *projChannel.Channel
}

//var clientlogger *log.Logger

//func init() {
//	os.Mkdir("log", os.ModePerm)
//	file, err := os.OpenFile("./log/Client.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 666)
//	if err != nil {
//		clientlogger.Fatal(err)
//	}
//	clientlogger = log.New(file, "", log.LstdFlags)
//	clientlogger.SetFlags(log.LstdFlags | log.Lshortfile)
//	file.Close()
//}

func (bc *BaseClient) GetPermission(ip string) AddrInfo {
	for _, value := range bc.permissions {
		if ip == value.TCPAddr.String() {
			return value
		}
	}

	return AddrInfo{}
}

func (bc *BaseClient) OnUserConnect(s *Session) {
	if len(bc.permissions) != 0 {
		permission := bc.GetPermission(s.GetConn().RemoteAddr().String())
		if permission.ID != 0 {
			if bc.GetSessions().IsNil(int32(permission.ID)) {
				if e := bc.GetSessions().Put(int32(permission.ID), s); e != nil {
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
		if _, err := bc.GetSessions().Add(s); err != nil {
			s.Close()
			return
		}
	}

	if bc.userHandler != nil {
		bc.userHandler.OnUserConnect(s)
	}

	log.Printf("LocalClient [%d] %s connect Server %s", s.GetID(), s.LocalAddr(), s.RemoteAddr())
}

func (bc *BaseClient) OnUserDisconnect(s *Session) {
	if !bc.GetSessions().IsNil(s.GetID()) {
		bc.GetSessions().Del(s.GetID())
	}

	if bc.userHandler != nil {
		bc.userHandler.OnUserDisconnect(s)
	}

	log.Printf("LocalClient [%d] %s disconnect Server %s", s.GetID(), s.LocalAddr(), s.RemoteAddr())
	s.Close()
}

func (bc *BaseClient) GetSessions() *Sessions {
	return bc.Sessions
}

func (bc *BaseClient) SetSessions(n int32) {
	bc.Sessions = NewSessions(n, bc.sendCh)
}

func (bc *BaseClient) SetUserHandler(value UserHandler) {
	bc.userHandler = value
}

func (bc *BaseClient) SetReConnectSecond(value time.Duration) {
	bc.reConnectSecond = value
}

func (bc *BaseClient) SetEventHandler(value *EventHandler) {
	bc.eventHandler = value
}

func (bc *BaseClient) GetPermissions() []AddrInfo {
	return bc.permissions
}

func (bc *BaseClient) SetPermissions(value []AddrInfo) {
	bc.permissions = value
}

func (bc *BaseClient) SetSendCh(sendCh *projChannel.Channel) {
	bc.sendCh = sendCh
}

func (bc *BaseClient) SetEventCh(eventCh *projChannel.Channel) {
	bc.eventCh = eventCh
}

// Connect 连接到服务器
func (bc *BaseClient) Connect(c context.Context) {
	for _, value := range bc.permissions {
		go func(addr AddrInfo) {
			time.Sleep(1 * time.Second)
			var s *Session
			var wg sync.WaitGroup
			isReConn := gtype.NewBool(false)
			var reConnectSecond = time.Second * 3
			for {
				conn, err := net.Dial("tcp", addr.TCPAddr.String())
				if err != nil {
					log.Printf("net.Dial Error: %s\n\n", err)
					continue
				}

				if s != nil {
					s.Close()
					s = nil
				}
				s = NewSession(conn, bc, &wg, *isReConn, bc.eventHandler, bc.eventCh)
				s.Start()
				wg.Wait()
				if bc.reConnectSecond == 0 {
					break
				}
				isReConn.Set(true)
				reConnectSecond = bc.reConnectSecond
				time.Sleep(reConnectSecond)
			}

			for {
				select {
				case <-c.Done():
					return
				}
			}
		}(value)
	}
	return
}

// Close 关闭连接
func (bc *BaseClient) Close() {
	bc.GetSessions().Close()
}

// Kick 关闭连接
func (bc *BaseClient) Kick(ID int32) {
	s := bc.FindSessionByID(ID)
	if s != nil {
		s.Close()
	}
}

// IsConnect 关闭连接
func (bc *BaseClient) IsConnect(ID int32) bool {
	return bc.FindSessionByID(ID) != nil
}

func (bc *BaseClient) SendMsgToAll(msgNo uint16, b []byte) {
	bc.Sessions.SendMsgToAll(msgNo, b)
}

func (bc *BaseClient) SendMsgToOne(ID int32, msgNo uint16, b []byte) {
	bc.Sessions.SendMsgToOne(ID, msgNo, b)
}

func (bc *BaseClient) SendMsgExclude(excludes []int32, msgNo uint16, b []byte) {
	bc.Sessions.SendMsgExclude(excludes, msgNo, b)
}

func (bc *BaseClient) FindSessionByID(ID int32) (p *Session) {
	return bc.Sessions.FindSessionByID(ID)
}

// NewBaseClient 创建一个TCPClient实例
func NewBaseClient() *BaseClient { return &BaseClient{} }
