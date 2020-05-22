package interfaces

type IHandleMgr interface {
	StartPool()
	StopPool()
	AddRouter(id uint32, router RouterCallback)
	SendToHandler(request IRequest)
}
