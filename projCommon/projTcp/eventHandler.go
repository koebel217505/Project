package projTcp

import (
	"fmt"
	"github.com/cstockton/go-conv"
	"math"
	"reflect"
	"strings"
)

var eventHandler *EventHandler = nil

// EventFunc bla-bla
type EventFunc func(s *Session, e *EventBuffer)

//type BytesToStructFunc func(b *EventBuffer) (e *EventMsg)

type ProtocolFun struct {
	Event EventFunc
	//BytesToStruct BytesToStructFunc
}

// event bla-bla
type EventHandler [math.MaxUint16 + 1]ProtocolFun

// GetEventHandler bla-bla
func NewEventHandler(event any, initMethod EventFunc) (e *EventHandler) {
	if eventHandler != nil {
		e = eventHandler
		return
	}

	e = &EventHandler{}
	e.SetAll(ProtocolFun{Event: initMethod})

	rValue := reflect.ValueOf(event)
	rType := reflect.TypeOf(event)

	if rType.Kind() != reflect.Struct {
		return
	}

	args := make([]reflect.Value, 0)

	for i := 0; i < rValue.NumMethod(); i++ {
		//fmt.Println(runtime.FuncForPC(rValue.Method(i).Pointer()).Name())
		result := rValue.Method(i).Call(args)

		name := rType.Method(i).Name
		if strings.Index(name, `Event`) < 0 {
			continue
		}

		name = strings.Replace(rType.Method(i).Name, "Event", "", -1)
		var err error
		var msgNo uint16
		msgNo, err = conv.Uint16(name)
		if err != nil {
			fmt.Println("msgNo Error:", err)
			continue
		}

		// log.Println(result[0].Interface().(byte), result[1].Interface().(byte))
		//e[msgNo].BytesToStruct = result[0].Interface().(BytesToStructFunc)
		e[msgNo].Event = result[1].Interface().(EventFunc)
	}

	return
}

// Set bla-bla
func (e EventHandler) Set(msgNo byte, method ProtocolFun) {
	e[msgNo] = method
}

// SetAll bla-bla
func (e EventHandler) SetAll(method ProtocolFun) {
	for i := 0; i < len(e); i++ {
		e[i] = method
	}

}
