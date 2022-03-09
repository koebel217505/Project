package clientHandler

import "github.com/koebel217505/Project/projCommon/projTcp"

// Event10000 bla-bla
func (e *Event) Event10000() projTcp.EventFunc {
	return func(s *projTcp.Session, e *projTcp.EventBuffer) {
		//fmt.Println(000)
	}
}
