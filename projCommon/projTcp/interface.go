package projTcp

// UserHandler bla-bla
type UserHandler interface {
	OnUserConnect(s *Session)
	OnUserDisconnect(s *Session)
	//OnUserReConnect(s Session)
}

// ServerHandler bla-bla
type ServerHandler interface {
	OnServerInit(s *BaseServer)
	OnServerDestroy(s *BaseServer)
}
