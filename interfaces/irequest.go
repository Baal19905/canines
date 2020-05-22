package interfaces

type IRequest interface {
	GetSessionID() string
	GetOpcode() uint32
	GetMsg() []byte
}
