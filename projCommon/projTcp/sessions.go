package projTcp

import (
	"github.com/koebel217505/Project/projCommon/projChannel"
	"github.com/koebel217505/Project/projCommon/projPacket"
	"github.com/koebel217505/Project/projCommon/projSeats"
)

type Sessions interface {
	Add(c Session) (ID int32, err error)
	Del(ID int32) (err error)
	Put(ID int32, c Session) (err error)
	Start()
	Clear()
	Close()
	SendMsgToAll(msgNo uint16, b []byte)
	SendMsgToOne(ID int32, msgNo uint16, b []byte)
	SendMsgExclude(excludes []int32, msgNo uint16, b []byte)
	FindSessionByID(ID int32) (p Session)
	IsNil(id int32) bool
}

// sessions bla-bla
type sessions struct {
	seats  *projSeats.SeatArray
	sendCh *projChannel.Channel
}

// NewSessions bla-bla
func NewSessions(max int32) (result Sessions) {
	result = &sessions{seats: projSeats.NewSeatArray(max), sendCh: projChannel.NewChannel(1000)}
	result.Start()
	return
}

// Add bla-bla
func (s *sessions) Add(c Session) (ID int32, err error) {
	i, e := s.seats.Add(c)
	c.SetID(ID)
	return int32(i), e
}

// Del bla-bla
func (s *sessions) Del(ID int32) (err error) {
	return s.seats.Del(int32(ID))
}

// Put bla-bla
func (s *sessions) Put(ID int32, c Session) (err error) {
	c.SetID(ID)
	return s.seats.Put(ID, c)
}

// Clear bla-bla
func (s *sessions) Clear() {
	s.seats.Clear()
}

func (s sessions) Encode(msgNo uint16, b []byte) []byte {
	pa := projPacket.PacketPool.Get()
	defer projPacket.PacketPool.Put(pa)
	pa.WriteUint16(uint16(2 + 2 + len(b)))
	pa.WriteUint16(msgNo)
	pa.WriteBytes(b)
	return pa.CopyBytes()
}

// SendMsgToAll bla-bla
func (s *sessions) SendMsgToAll(msgNo uint16, b []byte) {
	s.sendCh.Push(func() {
		v := s.Encode(msgNo, b)
		s.seats.ForEach(
			func(key any, value any) {
				if value == nil {
					return
				}

				if session, ok := value.(Session); ok {
					session.Send(v)
				}
			})
	})
}

// SendMsgToOne bla-bla
func (s *sessions) SendMsgToOne(ID int32, msgNo uint16, b []byte) {
	s.sendCh.Push(func() {
		v := s.Encode(msgNo, b)
		s.seats.Find(
			func(key any, value any) bool {
				if value == nil {
					return false
				}

				if session, ok := value.(Session); ok {
					if session.GetID() == ID {
						session.Send(v)
						return true
					}
				}

				return false
			})
	})
}

// SendMsgExclude bla-bla
func (s *sessions) SendMsgExclude(excludes []int32, msgNo uint16, b []byte) {
	s.sendCh.Push(func() {
		v := s.Encode(msgNo, b)
		s.seats.ForEach(
			func(key any, value any) {
				if value == nil {
					return
				}
				if session, ok := value.(Session); ok {
					for ID := range excludes {
						if session.GetID() == int32(ID) {
							return
						}
					}

					session.Send(v)
				}
			})
	})
}

// FindSessionByID bla-bla
func (s *sessions) FindSessionByID(id int32) (p Session) {
	find := s.seats.Find(
		func(key any, value any) bool {
			if value == nil {
				return false
			}

			if session, ok := value.(Session); ok {
				if int32(session.GetID()) == id {
					return true
				}
			}

			return false
		})

	if find == nil {
		return nil
	}

	return find.(Session)
}

// IsNil bla-bla
func (s *sessions) IsNil(id int32) bool {
	return s.seats.IsNil(int32(id))
}

func (s *sessions) Close() {
	s.seats.ForEach(
		func(key any, value any) {
			if value == nil {
				return
			}
			if session, ok := value.(Session); ok {
				session.Close()
			}
		})

	s.seats.Clear()
}

func (s *sessions) Start() {
	go s.sendThread()
}

func (s *sessions) sendThread() {
	for {
		if v := s.sendCh.Pop(); v != nil {
			if v, ok := v.(func()); ok {
				v()
			}
		}
	}
}
