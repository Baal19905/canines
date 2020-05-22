package network

type Request struct {
	sessionid string
	opcode  uint32
	msg     []byte
}

type Response struct {
	Opcode uint32
	Msg    []byte
}

func (r *Request) GetSessionID() string {
	return r.sessionid
}

func (r *Request) GetOpcode() uint32 {
	return r.opcode
}

func (r *Request) GetMsg() []byte {
	return r.msg
}
