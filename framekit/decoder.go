package framekit

import "sync"

type Decoder struct {
	Bytes []byte // 待处理数据
	sync.RWMutex
	TakeHandler func([]byte) ([]byte, []byte, bool)
}

func NewDecoder(initCapacity int, takeHandler func([]byte) ([]byte, []byte, bool)) *Decoder {
	return &Decoder{
		Bytes:       make([]byte, 0, initCapacity),
		TakeHandler: takeHandler,
	}
}

// Take 取一次
// return 未处理的bytes、成功出库的bytes、是否结束
func (th *Decoder) Take() ([]byte, bool) {
	th.Lock()
	defer th.Unlock()
	o, r, over := th.TakeHandler(th.Bytes)
	th.Bytes = o
	return r, over
}

// 等待接收的模式 todo

func (th *Decoder) Put(data []byte) {
	if len(data) > 0 {
		th.Lock()
		defer th.Unlock()
		th.Bytes = append(th.Bytes, data...)
	}
}
