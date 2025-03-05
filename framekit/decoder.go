package framekit

import "sync"

type Decoder struct {
	Bytes []byte // 待处理数据
	sync.RWMutex
	TakeHandler func([]byte) []byte
}

func NewDecoder(initCapacity int, takeHandler func([]byte) []byte) *Decoder {
	return &Decoder{
		Bytes:       make([]byte, 0, initCapacity),
		TakeHandler: takeHandler,
	}
}

func (th *Decoder) Take() []byte {
	th.Lock()
	defer th.Unlock()
	data := th.TakeHandler(th.Bytes)
	return data
}

func (th *Decoder) Put(data []byte) {
	if len(data) > 0 {
		th.Lock()
		defer th.Unlock()
		th.Bytes = append(th.Bytes, data...)
	}
}
