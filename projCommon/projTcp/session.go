package projTcp

import (
	"encoding/binary"
	"fmt"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/koebel217505/Project/projCommon/projChannel"
	"github.com/koebel217505/Project/projCommon/projPacket"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type AddrInfo struct {
	ID      uint64
	Name    string
	TCPAddr *net.TCPAddr
}

// Session 代表一个连接会话
type Session struct {
	id            int32
	conn          net.Conn
	sendChan      chan []byte
	buffer        []byte
	wg            *sync.WaitGroup
	eventCh       *projChannel.Channel
	eventHandler  *EventHandler
	userHandler   UserHandler
	data          any
	isReConn      gtype.Bool
	closeOnce     sync.Once
	readDeadline  time.Time
	writeDeadline time.Time
	isClosed      gtype.Bool
}

func (s *Session) GetConn() net.Conn {
	return s.conn
}

func (s *Session) GetBuffer() []byte {
	return s.buffer
}

func (s *Session) GetID() int32 {
	return s.id
}

func (s *Session) SetID(value int32) {
	s.id = value
}

func (s *Session) SetData(value any) {
	s.data = value
}

func (s *Session) GetData() any {
	return s.data
}

func (s *Session) SetReadDeadline(value time.Time) {
	s.readDeadline = value
}

func (s *Session) GetReadDeadline() time.Time {
	return s.readDeadline
}

func (s *Session) SetWriteDeadline(value time.Time) {
	s.writeDeadline = value
}

func (s *Session) GetWriteDeadline() time.Time {
	return s.writeDeadline
}

// EventBuffer bla-bla
type EventBuffer struct {
	MsgNo  uint16
	Buffer []byte
}

// EventMsg bla-bla
type EventMsg struct {
	MsgNo uint16
	Data  any
}

// RemoteAddr 返回客户端的地址和端口
func (s *Session) RemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

// LocalAddr 返回本機地址和端口
func (s *Session) LocalAddr() string {
	return s.conn.LocalAddr().String()
}

// Start 開始
func (s *Session) Start() {
	func() {
		s.wg.Add(1)
		go s.receiveThread()
		s.wg.Add(1)
		go s.sendThread()
		//s.wg.Add(1)
		//go s.eventThread()
	}()
}

// Close 關閉連接
func (s *Session) Close() {
	s.closeOnce.Do(s.close)
}

func (s *Session) close() {
	//s.conn.Close()
	s.isClosed.Set(true)
	//close(s.eventChan)
	//s.SetReadDeadline(time.Now().Add(0))
}

func (s *Session) receiveThread() {
	defer s.wg.Done()

	if s.isReConn.Val() == false {
		s.userHandler.OnUserConnect(s)
	}

	for {
		if s.isClosed.Val() == true {
			break
		}

		if err := s.conn.SetReadDeadline(time.Now().Add(time.Second * 1 * 60 * 60)); err != nil {
			log.Printf("SetReadDeadline TimeOut:%v\n", err)
			break
		}

		if eventBuffer, err := s.Decode(); err != nil {
			if eventBuffer == nil {
				break
			}
			if err != io.EOF {
				break
			}
			break
		} else {
			if s.eventHandler[eventBuffer.MsgNo].Event == nil {
				log.Printf("event[%d].Event nil \n", eventBuffer.MsgNo)
				continue
			}

			s.PushEvent(func() {
				s.eventHandler[eventBuffer.MsgNo].Event(s, eventBuffer)
			})
		}
	}

	s.userHandler.OnUserDisconnect(s)
	s.Close()
}

func (s *Session) sendThread() {
	defer s.wg.Done()

	for msg := range s.sendChan {
		if err := s.conn.SetWriteDeadline(time.Now().Add(time.Second * 60)); err != nil {
			log.Println("SetWriteDeadline TimeOut")
			break
		}

		if _, err := s.conn.Write(msg); err != nil {
			break
		}
	}

	//s.Close()
	err := s.conn.Close()
	if err != nil {
		return
	}
}

//// eventThread bla-bla
//func (s *Session) eventThread() {
//	defer s.wg.Done()
//
//	for event := range s.eventChan {
//		event()
//	}
//
//	close(s.sendChan)
//}

// Send 發送數據
func (s *Session) Send(b []byte) {
	if s.isClosed.Val() == false {
		s.sendChan <- b
	}
}

//PushEvent 使用者事件
func (s *Session) PushEvent(event func()) {
	s.eventCh.Push(event)
}

// Decode bla-bla
func (s *Session) Decode() (e *EventBuffer, err error) {
	reader := s.GetConn().(io.Reader)
	buffer := s.GetBuffer()

	_, err = io.ReadFull(reader, buffer[0:2])
	if err != nil {
		fmt.Println(err)
		return &EventBuffer{}, err
	}

	pa := projPacket.PacketPool.Get()
	defer projPacket.PacketPool.Put(pa)
	size := binary.LittleEndian.Uint16(buffer[0:2])
	pa.WriteUint16(size)

	_, err = io.ReadFull(reader, buffer[0:2])
	if err != nil {
		fmt.Println(err)
		return &EventBuffer{}, err
	}
	msgNo := binary.LittleEndian.Uint16(buffer[0:2])
	pa.WriteUint16(msgNo)

	for size > pa.Len() {
		n := size - pa.Len()
		if n > uint16(len(buffer)) {
			if _, err := io.ReadFull(reader, buffer[:]); err != nil {
				return &EventBuffer{}, err
			}

			pa.WriteBytes(buffer[:])
		} else {
			if _, err := io.ReadFull(reader, buffer[0:n]); err != nil {
				return &EventBuffer{}, err
			}

			pa.WriteBytes(buffer[0:n])
		}
	}

	// log.Println("Buffer:", string(projPacket.bytes()))
	// log.Println("Buffer:", string(c.recvMsg[0:size]))

	// time.Sleep(time.Second * 10)
	fmt.Println(msgNo, pa.CopyBytes())
	return &EventBuffer{msgNo, pa.CopyBytes()}, nil
}

// Encode bla-bla
func (s *Session) Encode(b []byte) (r *projPacket.Packet, err error) {

	return
}

// NewSession 生成一个新的Session
func NewSession(conn net.Conn, userHandler UserHandler, wg *sync.WaitGroup, isReConn gtype.Bool, eventHandler *EventHandler, eventCh *projChannel.Channel) (result *Session) {
	result = &Session{
		conn:     conn,
		sendChan: make(chan []byte, 100),
		buffer:   make([]byte, 1024*8),
		wg:       wg,
		//eventChan:    make(projChannel func(), 10000),
		eventCh:      eventCh,
		userHandler:  userHandler,
		eventHandler: eventHandler,
		isReConn:     isReConn,
	}

	return result
}

func NewSessionNoCon(value any) (result *Session) {
	result = &Session{
		data: value,
	}

	return result
}
