package projTcp

import (
	"encoding/binary"
	"fmt"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/koebel217505/Project/projCommon/projPacket"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Test struct {
	A int8
	B int16
	C int32
}

type Session interface {
	RemoteAddr() string
	LocalAddr() string
	Start()
	Close()
	Send(b []byte)
	PushEvent(event func())

	GetID() int32
	SetID(value int32)

	GetBuffer() []byte

	GetConn() net.Conn

	//SetGMap(value *gmap.Map)
	//GetGMap() *gmap.Map

	SetData(value any)
	GetData() any

	//Get(k string) any
	//GetOrSetFunc(key any, f func() interface{}) any
	//SetSessionID(value int32)
	//GetSessionID() int32
}

// Session 代表一个连接会话
type session struct {
	id           int32
	conn         net.Conn
	sendChan     chan []byte
	buffer       []byte
	wg           *sync.WaitGroup
	eventChan    chan func()
	eventHandler *EventHandler
	userHandler  UserHandler

	data any
	//gMap          *gmap.Map
	isReConn      gtype.Bool
	closeOnce     sync.Once
	readDeadline  time.Time
	writeDeadline time.Time
	//sessionID     gtype.Int32
	isClosed gtype.Bool
}

func (s *session) GetConn() net.Conn {
	return s.conn
}

func (s *session) GetBuffer() []byte {
	return s.buffer
}

func (s *session) GetID() int32 {
	return s.id
}

func (s *session) SetID(value int32) {
	s.id = value
}

func (s *session) SetData(value any) {
	s.data = value
}

func (s *session) GetData() any {
	return s.data
}

//func (s *session) SetGMap(value *gmap.Map) {
//	if s.gMap != nil {
//		s.gMap.Clear()
//		s.gMap = nil
//	}
//
//	s.gMap = value
//}
//
//func (s *session) GetGMap() *gmap.Map {
//	return s.gMap
//}
//
//func (s *session) Get(key string) any {
//	return s.gMap.Get(key)
//}
//
//func (s *session) GetOrSetFunc(key any, f func() interface{}) any {
//	return s.gMap.GetOrSetFunc(key, f)
//}

//func (s *session) SetSessionID(value int32) {
//	s.sessionID.Set(value)
//}
//
//func (s *session) GetSessionID() int32 {
//	return s.sessionID.Val()
//}

func (s *session) SetReadDeadline(value time.Time) {
	s.readDeadline = value
}

func (s *session) GetReadDeadline() time.Time {
	return s.readDeadline
}

func (s *session) SetWriteDeadline(value time.Time) {
	s.writeDeadline = value
}

func (s *session) GetWriteDeadline() time.Time {
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
func (s *session) RemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

// LocalAddr 返回本機地址和端口
func (s *session) LocalAddr() string {
	return s.conn.LocalAddr().String()
}

// Start 開始
func (s *session) Start() {
	func() {
		s.wg.Add(1)
		go s.receiveThread()
		s.wg.Add(1)
		go s.sendThread()
		s.wg.Add(1)
		go s.eventThread()
	}()
}

// Close 關閉連接
func (s *session) Close() {
	s.closeOnce.Do(s.close)
}

func (s *session) close() {
	//s.conn.Close()
	s.isClosed.Set(true)
	close(s.eventChan)
	//s.SetReadDeadline(time.Now().Add(0))
}

func (s *session) receiveThread() {
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

func (s *session) sendThread() {
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

// eventThread bla-bla
func (s *session) eventThread() {
	defer s.wg.Done()

	for event := range s.eventChan {
		event()
	}

	close(s.sendChan)
}

// Send 發送數據
func (s *session) Send(b []byte) {
	if s.isClosed.Val() == false {
		s.sendChan <- b
	}
}

//PushEvent 使用者事件
func (s *session) PushEvent(event func()) {
	if s.eventChan != nil {
		s.eventChan <- event
	}
}

// Decode bla-bla
func (s *session) Decode() (e *EventBuffer, err error) {
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
func (s *session) Encode(b []byte) (r *projPacket.Packet, err error) {

	return
}

// NewSession 生成一个新的Session
func NewSession(conn net.Conn, userHandler UserHandler, wg *sync.WaitGroup, isReConn gtype.Bool, eventHandler *EventHandler, eventChan chan func()) (result Session) {
	result = &session{
		conn:     conn,
		sendChan: make(chan []byte, 100),
		buffer:   make([]byte, 1024*8),
		wg:       wg,
		//eventChan:    make(projChannel func(), 10000),
		eventChan:    eventChan,
		userHandler:  userHandler,
		eventHandler: eventHandler,
		isReConn:     isReConn,
	}

	return result
}

func NewSessionNoCon(value any) (result Session) {
	result = &session{
		data: value,
	}

	return result
}
