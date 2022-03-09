package projTcp

import (
	"context"
	"fmt"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/jianfengye/collection"
	"github.com/koebel217505/Project/projCommon/projType"
	"github.com/spf13/cast"
	"log"
	"net"
	"sync"
	"time"
)

//const (
//	stateNone byte = iota
//	stateConnect
//	stateDisConnect
//	stateReConnect
//)

type tcpClientState struct {
	remoteAddr string
	remotePort int
	connected  bool
}

// Client TCP客户端描述
type baseClient struct {
	tcpClientState
	//session         Session
	//wg              sync.WaitGroup
	//state           byte
	reConnectSecond time.Duration
	userHandler     UserHandler
	eventHandler    *EventHandler
	//closeCh         projChannel.Channel
	serverAddrArray []projType.Addr
	sAddrArray      *collection.ObjCollection
	Sessions        Sessions
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

type BaseClient interface {
	Connect(c context.Context)
	Close()
	GetSessions() Sessions
	SetSessions(n int32)
	Kick(ID int32)
	IsConnect(ID int32) bool
	OnUserConnect(s Session)
	OnUserDisconnect(s Session)
	SetUserHandler(UserHandler)
	SetReConnectSecond(time.Duration)
	SetEventHandler(*EventHandler)
	GetServerAddrArray() []projType.Addr
	SetServerAddrArray(value []projType.Addr)

	SendMsgToAll(msgNo uint16, b []byte)
	SendMsgToSomeOne(ID int32, msgNo uint16, b []byte)
	SendMsgExclude(excludes []int32, msgNo uint16, b []byte)
	FindSessionByID(ID int32) (p Session)

	GetServerAddrIndex(ip string) int32
}

func (bc *baseClient) GetServerAddrIndex(ip string) int32 {
	host, port, _ := net.SplitHostPort(ip)
	for key, value := range bc.serverAddrArray {
		if host == value.IP && port == cast.ToString(value.Port) {
			return int32(key)
		}
	}

	return -1
}

func (bc *baseClient) OnUserConnect(s Session) {
	if bc.serverAddrArray != nil && len(bc.serverAddrArray) != 0 {
		if index := bc.GetServerAddrIndex(s.GetConn().RemoteAddr().String()); index >= 0 {
			sessionByID := bc.FindSessionByID(index)
			if sessionByID != nil {
				//sessionByID.Close()
			}

			if e := bc.GetSessions().Put(int32(index), s); e != nil {
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

	if _, err := bc.GetSessions().Add(s); err != nil {
		s.Close()
		return
	}

	if bc.userHandler != nil {
		bc.userHandler.OnUserConnect(s)
	}

	log.Printf("LocalClient [%d] %s connect Server %s", s.GetID(), s.LocalAddr(), s.RemoteAddr())
}

func (bc *baseClient) OnUserDisconnect(s Session) {
	if !bc.GetSessions().IsNil(s.GetID()) {
		bc.GetSessions().Del(s.GetID())
	}

	if bc.userHandler != nil {
		bc.userHandler.OnUserDisconnect(s)
	}

	log.Printf("LocalClient [%d] %s disconnect Server %s", s.GetID(), s.LocalAddr(), s.RemoteAddr())
	s.Close()
}

func (bc *baseClient) GetSessions() Sessions {
	return bc.Sessions
}

func (bc *baseClient) SetSessions(n int32) {
	bc.Sessions = NewSessions(n)
}

func (bc *baseClient) SetUserHandler(value UserHandler) {
	bc.userHandler = value
}

func (bc *baseClient) SetReConnectSecond(value time.Duration) {
	bc.reConnectSecond = value
}

func (bc *baseClient) SetEventHandler(value *EventHandler) {
	bc.eventHandler = value
}

func (bc *baseClient) GetServerAddrArray() []projType.Addr {
	return bc.serverAddrArray
}

func (bc *baseClient) SetServerAddrArray(value []projType.Addr) {
	bc.serverAddrArray = value
}

// Connect 连接到服务器
func (bc *baseClient) Connect(c context.Context) {
	for _, value := range bc.serverAddrArray {
		go func(addr projType.Addr) {
			time.Sleep(1 * time.Second)
			var s Session
			var wg sync.WaitGroup
			isReConn := gtype.NewBool(false)
			var reConnectSecond = time.Second * 3
			for {
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr.IP, addr.Port))
				if err != nil {
					log.Printf("net.Dial Error: %s\n\n", err)
					continue
				}

				if s != nil {
					s.Close()
					s = nil
				}
				s = NewSession(conn, bc, &wg, *isReConn, bc.eventHandler, make(chan func(), 10000))
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

//// Send 发送数据
//func (c *baseClient) Send(msgNo uint16, b []byte) {
//	pa := projPacket.PacketPool.Get()
//	defer projPacket.PacketPool.Put(pa)
//	pa.WriteUint16(uint16(2 + 2 + len(b)))
//	pa.WriteUint16(msgNo)
//	pa.WriteBytes(b)
//	c.session.Send(pa.CopyBytes())
//
//	s := c.FindSessionByIndex(index)
//}

// Close 关闭连接
func (bc *baseClient) Close() {
	bc.GetSessions().Close()
}

// Kick 关闭连接
func (bc *baseClient) Kick(ID int32) {
	s := bc.FindSessionByID(ID)
	if s != nil {
		s.Close()
	}
}

// IsConnect 关闭连接
func (bc *baseClient) IsConnect(ID int32) bool {
	return bc.FindSessionByID(ID) != nil
}

func (bc *baseClient) SendMsgToAll(msgNo uint16, b []byte) {
	bc.Sessions.SendMsgToAll(msgNo, b)
}

func (bc *baseClient) SendMsgToSomeOne(ID int32, msgNo uint16, b []byte) {
	bc.Sessions.SendMsgToOne(ID, msgNo, b)
}

func (bc *baseClient) SendMsgExclude(excludes []int32, msgNo uint16, b []byte) {
	bc.Sessions.SendMsgExclude(excludes, msgNo, b)
}

func (bc *baseClient) FindSessionByID(ID int32) (p Session) {
	return bc.Sessions.FindSessionByID(ID)
}

// NewBaseClient 创建一个TCPClient实例
func NewBaseClient() BaseClient { return &baseClient{} }
