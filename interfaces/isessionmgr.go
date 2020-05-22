package interfaces

type ISessionMgr interface {
	Size() uint32
	Len() uint32
	GetSession(sid string) ISession
	Add(s ISession)
	Remove(sid string)
	RemoveAll()
	Start()
	Stop()
}
