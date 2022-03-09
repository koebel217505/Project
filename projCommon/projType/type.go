package projType

//const UIStatus_None = 0
//const UIStatus_Start = 1
//const UIStatus_NotLoad = 2

//// Checkfunc bla-bla
//type Checkfunc func(key any, value any) bool
//
//// CallBackfunc bla-bla
//type CallBackfunc func(key any, value any)

type Addr struct {
	//ID   int
	Name string
	IP   string
	Port int
}

type IP struct {
	Local  Addr
	Server Addr
	Client []Addr
}

type AccountData struct {
	//PlayerId int64
	Account  string
	Password string
}
