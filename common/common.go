package common

import "sync/atomic"

type Saferef struct {
	ref uint32
}

func (s *Saferef) Add() {
	atomic.AddUint32(&s.ref, 1)
}

func (s *Saferef) Sub() {
	atomic.AddUint32(&s.ref, ^uint32(0))
}

func (s *Saferef) Get() uint32 {
	atomic.LoadUint32(&s.ref)
	return s.ref
}

func (s *Saferef) Set(uint32) {
	atomic.StoreUint32(&s.ref, 0)
}
