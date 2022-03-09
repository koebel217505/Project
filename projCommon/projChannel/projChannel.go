package projChannel

import (
	"errors"

	"github.com/gogf/gf/g/container/gtype"
)

// Channel bla-bla
type Channel struct {
	ch     chan interface{}
	closed *gtype.Bool
}

// Push bla-bla
func (c *Channel) Push(value interface{}) error {
	if c.closed.Val() {
		return errors.New("ch is closed")
	}
	c.ch <- value
	return nil
}

// Pop bla-bla
func (c *Channel) Pop() interface{} {
	return <-c.ch
}

func (c *Channel) Self() chan interface{} {
	return c.ch
}

// Close bla-bla
func (c *Channel) Close() {
	if !c.closed.Set(true) {
		close(c.ch)
	}
}

// Size bla-bla
func (c *Channel) Size() int {
	return c.Len()
}

// Len bla-bla
func (c *Channel) Len() int {
	return len(c.ch)
}

// Cap bla-bla
func (c *Channel) Cap() int {
	return cap(c.ch)
}

// New bla-bla
func New(limit int) *Channel {
	return &Channel{
		ch:     make(chan interface{}, limit),
		closed: gtype.NewBool(),
	}
}
