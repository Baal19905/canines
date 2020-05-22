package handler

import (
	"container/heap"
	"github.com/baal19905/canines/interfaces"
	"github.com/baal19905/canines/utils"
	"sync"
)

type HandlerHeap []*JobHanelder
type RouterMap map[uint32]interfaces.RouterCallback
type HandlerMgr struct {
	heap    HandlerHeap
	routers RouterMap
	sm      interfaces.ISessionMgr
	lock    sync.Mutex
}

func NewHandlerMgr(sessionMgr interfaces.ISessionMgr) *HandlerMgr {
	return &HandlerMgr{
		heap:    make(HandlerHeap, utils.ConfInfo.MaxWorker),
		routers: make(RouterMap, utils.ConfInfo.MaxWorker),
		sm:      sessionMgr,
	}
}

func (hm *HandlerMgr) StartPool() {
	for i := 0; i < hm.heap.Len(); i++ {
		hm.heap[i] = NewJobHandler()
		hm.heap[i].hm = hm
		hm.heap[i].InitLoad()
		hm.heap[i].index = i
		hm.heap[i].sm = hm.sm
		go hm.heap[i].Start()
	}
	heap.Init(&hm.heap)
}

func (hm *HandlerMgr) StopPool() {
	for {
		jh := hm.Pop()
		if jh == nil {
			break
		}
		jh.Stop()
	}
}

func (hm *HandlerMgr) AddRouter(id uint32, router interfaces.RouterCallback) {
	hm.routers[id] = router
}

func (hm *HandlerMgr) GetRouter(id uint32) interfaces.RouterCallback {
	router, ok := hm.routers[id]
	if ok {
		return router
	}
	return nil
}

func (hm *HandlerMgr) SendToHandler(request interfaces.IRequest) {
	jh := hm.Pop()
	if jh == nil {
		// todo log
		return
	}
	jh.AddLoad()
	hm.Push(jh)
	jh.queue <- request
}

func (hm *HandlerMgr) Push(jh *JobHanelder) {
	hm.lock.Lock()
	defer hm.lock.Unlock()
	heap.Push(&hm.heap, jh)
}

func (hm *HandlerMgr) Pop() *JobHanelder {
	hm.lock.Lock()
	defer hm.lock.Unlock()
	return heap.Pop(&hm.heap).(*JobHanelder)
}

func (hm *HandlerMgr) Update(jh *JobHanelder) {
	hm.lock.Lock()
	defer hm.lock.Unlock()
	hm.heap.Update(jh)
}

func (hh HandlerHeap) Len() int {
	return len(hh)
}

func (hh HandlerHeap) Less(i, j int) bool {
	return hh[i].GetLoad() < hh[j].GetLoad()
}

func (hh HandlerHeap) Swap(i, j int) {
	hh[i], hh[j] = hh[j], hh[i]
	hh[i].index = i
	hh[j].index = j
}

func (hh *HandlerHeap) Push(x interface{}) {
	n := len(*hh)
	item := x.(*JobHanelder)
	item.index = n
	*hh = append(*hh, item)
}

func (hh *HandlerHeap) Pop() interface{} {
	old := *hh
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*hh = old[0 : n-1]
	return item
}

func (hh *HandlerHeap) Update(jh *JobHanelder) {
	heap.Fix(hh, jh.index)
}
