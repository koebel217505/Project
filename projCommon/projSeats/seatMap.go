package projSeats

import (
	"sync"
)

// SeatMap bla-bla
type SeatMap struct {
	rw   sync.RWMutex
	data map[int32]any
	max  int32
}

// isNil bla-bla
func (s *SeatMap) isNil(index int32) bool {
	return s.data[index] == nil
}

// IsNil bla-bla
func (s *SeatMap) IsNil(index int32) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.isNil(index)
}

// len bla-bla
func (s *SeatMap) len() int32 {
	return int32(len(s.data))
}

// findEmpty bla-bla
func (s *SeatMap) getEmpty() (index int32, err error) {
	return int32(len(s.data)), nil
}

// Add bla-bla
func (s *SeatMap) Add(data any) (index int32, err error) {
	s.rw.Lock()
	defer s.rw.Unlock()
	if data != nil {
		index, err = s.getEmpty()
		if err == nil {
			s.data[index] = data
		}
	}

	return index, err
}

// Del bla-bla
func (s *SeatMap) Del(index int32) (err error) {
	s.rw.Lock()
	defer s.rw.Unlock()
	delete(s.data, index)
	return nil
}

// Clear bla-bla
func (s *SeatMap) Clear() (err error) {
	s.rw.Lock()
	defer s.rw.Unlock()
	for key := range s.data {
		delete(s.data, key)
	}

	s.data = make(map[int32]any, 0)
	return nil
}

//Put 存儲操作
func (s *SeatMap) Put(index int32, data any) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if s.data[index] == nil {
		return nil
	}
	s.data[index] = data
	return nil
}

//Get 獲取操作
func (s *SeatMap) Get(index int32) any {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.data[index]
}

//Find bla-bla
func (s *SeatMap) Find(callback func(key any, value any) bool) (result any) {
	defer s.rw.RUnlock()
	s.rw.RLock()

	if callback == nil {
		return nil
	}

	for key, value := range s.data {
		if value == nil {
			continue
		}

		if ok := callback(key, value); ok == true {
			result = value
			return
		}
	}

	return
}

// ForEach bla-bla
func (s *SeatMap) ForEach(callback func(key any, value any)) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if callback == nil {
		return
	}

	for key, value := range s.data {
		if value != nil {
			callback(key, value)
		}
	}
}

// NewSeatMap bla-bla
func NewSeatMap(max int32) *SeatMap {
	if max == 0 {
		return nil
	}

	return &SeatMap{max: max, data: make(map[int32]any, max)}
}
