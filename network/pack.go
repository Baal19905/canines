package network

import (
	"github.com/Baal19905/canines/interfaces"
	"github.com/Baal19905/canines/utils"
)

type HeadBase struct {
}

func (h *HeadBase) Marshal() []byte                      { return nil }
func (h *HeadBase) UnMarshal(buf []byte) error           { return nil }
func (h *HeadBase) Check() error                         { return nil }
func (h *HeadBase) GetHeadLen() uint32                   { return 0 }
func (h *HeadBase) SetHeadLen(l uint32)                  {}
func (h *HeadBase) GetBodyLen() uint32                   { return 0 }
func (h *HeadBase) SetBodyLen(uint32)                    {}
func (h *HeadBase) PreHandle(buf []byte) ([]byte, error) { return buf, nil }
func (h *HeadBase) PreSend(buf []byte) ([]byte, error)   { return buf, nil }
func (h *HeadBase) GetOpcode() uint32                    { return 0 }
func (h *HeadBase) SetOpcode(uint32)                     {}

func RegistHead(head interfaces.IHead) {
	utils.HeadTemplate = head
}
