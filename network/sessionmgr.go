package network

import (
	"github.com/Baal19905/canines/interfaces"
	"sync"
	"github.com/Baal19905/canines/utils"
)

type SessionMgr struct {
	size      uint32
	len       uint32
	sessions  map[string]interfaces.ISession
	stopqueue chan string
	exit      chan bool
	lock      sync.RWMutex
}

func NewSessionManage() *SessionMgr {
	return &SessionMgr{
		size:      utils.ConfInfo.MaxCli,
		len:       0,
		sessions:  make(map[string]interfaces.ISession, utils.ConfInfo.MaxCli),
		stopqueue: make(chan string, utils.ConfInfo.MaxCli),
	}
}

func (sm *SessionMgr) Size() uint32 {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	ret := sm.size
	return ret
}

func (sm *SessionMgr) Len() uint32 {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	ret := sm.len
	return ret
}

func (sm *SessionMgr) GetSession(sid string) interfaces.ISession {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	s, ok := sm.sessions[sid]
	if !ok {
		return nil
	}
	return s
}

func (sm *SessionMgr) Add(s interfaces.ISession) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if sm.len < sm.size {
		sm.sessions[s.GetSessionID()] = s
		sm.len++
	}
}

func (sm *SessionMgr) Remove(sid string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	s, ok := sm.sessions[sid]
	if !ok {
		return
	}
	s.Stop()
	delete(sm.sessions, sid)
	sm.len--
}

func (sm *SessionMgr) RemoveAll() {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	for id, s := range sm.sessions {
		s.Stop()
		delete(sm.sessions, id)
	}
	sm.len = 0
}

func (sm *SessionMgr) Start() {
	go func() {
		select {
		case sessionid := <-sm.stopqueue:
			sm.Remove(sessionid)
		case <-sm.exit:
			return
		}
	}()
}

func (sm *SessionMgr) Stop() {
	sm.exit <- true
}
