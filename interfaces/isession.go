package interfaces

type ISession interface {
	Start()
	Stop()
	GetSessionID() string
	PushMsg(opcode uint32, msg []byte)
	GetRemoteAddr() string
}
