package interfaces

type IHead interface {
	Marshal() []byte
	UnMarshal(buf []byte) error
	Check() error
	GetHeadLen() uint32
	SetHeadLen(uint32)
	GetBodyLen() uint32
	SetBodyLen(uint32)
	PreHandle(buf []byte)([]byte, error)
	PreSend(buf []byte) ([]byte, error)
	GetOpcode() uint32
	SetOpcode(uint32)
}
