package projSeats

import (
	"errors"
	"sync"
)

// SeatArray bla-bla
type SeatArray struct {
	rw   sync.RWMutex
	data []any
	max  int32
}

// isNil bla-bla
func (s *SeatArray) isNil(index int32) bool {
	return s.isRange(index) || s.data[index] == nil
}

// IsNil bla-bla
func (s *SeatArray) IsNil(index int32) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.isNil(index)
}

// isRange bla-bla
func (s *SeatArray) isRange(index int32) bool {
	return index <= 0 || index > s.max
}

// len bla-bla
func (s *SeatArray) len() int32 {
	return int32(len(s.data))
}

// findEmpty bla-bla
func (s *SeatArray) getEmpty() (index int32, err error) {
	for i := int32(0); i < s.max; i++ {
		if s.data[i] == nil {
			return i, nil
		}
	}

	return -1, errors.New("no Empty")
}

// Add bla-bla
func (s *SeatArray) Add(data any) (index int32, err error) {
	s.rw.Lock()
	defer s.rw.Unlock()
	if data == nil {
		return -1, errors.New("data Err")
	}

	index, err = s.getEmpty()
	if err == nil {
		s.data[index] = data
	}

	return index, err
}

// Del bla-bla
func (s *SeatArray) Del(index int32) (err error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	if s.isRange(index) == false {
		return errors.New("index Err")
	}

	s.data = append(s.data[:index], s.data[index+1:]...)
	return nil
}

// Clear bla-bla
func (s *SeatArray) Clear() (err error) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.data = make([]any, 0)
	return nil
}

//Put 存儲操作
func (s *SeatArray) Put(index int32, data any) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if s.isRange(index) == false {
		return errors.New("index Err")
	}

	s.data[index] = data
	return nil
}

//Get 獲取操作
func (s *SeatArray) Get(index int32) (r any, err error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if s.isRange(index) == false {
		return nil, errors.New("index Err")
	}

	return s.data[index], nil
}

//Find bla-bla
func (s *SeatArray) Find(callback func(key any, value any) bool) (r any) {
	s.rw.Lock()
	defer s.rw.Unlock()

	if callback == nil {
		return nil
	}

	for key, value := range s.data {
		if value == nil {
			continue
		}

		if ok := callback(key, value); ok == true {
			r = value
			return
		}
	}

	return
}

// ForEach bla-bla
func (s *SeatArray) ForEach(callback func(key any, value any)) {
	s.rw.Lock()
	defer s.rw.Unlock()

	if callback == nil {
		return
	}

	for key, value := range s.data {
		if value != nil {
			callback(key, value)
		}
	}
}

// NewSeatArray bla-bla
func NewSeatArray(max int32) *SeatArray {
	if max == 0 {
		return nil
	}

	return &SeatArray{max: max, data: make([]any, max)}
}
