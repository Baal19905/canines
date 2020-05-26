package handler

import (
	"github.com/Baal19905/canines/common"
	"github.com/Baal19905/canines/interfaces"
	"github.com/Baal19905/canines/utils"
)

type JobHanelder struct {
	load  common.Saferef
	index int
	hm    *HandlerMgr
	sm    interfaces.ISessionMgr
	queue chan interfaces.IRequest
	stop  chan bool
}

func NewJobHandler() *JobHanelder {
	return &JobHanelder{
		queue: make(chan interfaces.IRequest, utils.ConfInfo.MaxHandleQueue),
	}
}

func (jh *JobHanelder) AddLoad() {
	jh.load.Add()
}

func (jh *JobHanelder) SubLoad() {
	jh.load.Sub()
}

func (jh *JobHanelder) GetLoad() uint32 {
	return jh.load.Get()
}

func (jh *JobHanelder) InitLoad() {
	jh.load.Set(0)
}

func (jh *JobHanelder) Start() {
	for {
		select {
		case request := <-jh.queue:
			jh.SubLoad()
			jh.hm.Update(jh)
			router := jh.hm.GetRouter(request.GetOpcode())
			if router == nil {
				//todo log
				continue
			}
			opcode, msg := router(request)
			if opcode != 0 && msg != nil {
				session := jh.sm.GetSession(request.GetSessionID())
				if session == nil {
					//todo log
					continue
				}
				session.PushMsg(opcode, msg)
			}
		case <-jh.stop:
			close(jh.queue)
			return
		}
	}
}

func (jh *JobHanelder) Stop() {
	jh.stop <- true
}
