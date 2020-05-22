package interfaces

type RouterCallback func (request IRequest) (opcode uint32, msg []byte)
type NotifyCallback func (session ISession)

type IServer interface {
	Stop()
	Serve()
	RegisterRouter(opcode uint32, router RouterCallback)
	Push(sid string, opcode uint32, msg []byte)
	RegisterOnConnect(connect NotifyCallback)
	RegisterOnDisConnect(disconnect NotifyCallback)
	CallOnConnect(sid string)
	CallOnDisConnect(sid string)
	GetHandleMgr() IHandleMgr
	GetSessionMgr() ISessionMgr
}
