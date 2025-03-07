package framekit

import "sync"

type Decoder struct {
	Bytes []byte // 待处理数据
	sync.RWMutex
	TakeHandler func([]byte) ([]byte, []byte, bool) // return 未处理的bytes、成功出库的bytes、是否结束
	enableRecv  bool
	recvHandler func([]byte)
}

func NewDecoder(initCapacity int, takeHandler func([]byte) ([]byte, []byte, bool)) *Decoder {
	return &Decoder{
		Bytes:       make([]byte, 0, initCapacity),
		TakeHandler: takeHandler,
	}
}

// Take 取一次
func (th *Decoder) Take() ([]byte, bool) {
	th.Lock()
	defer th.Unlock()
	o, r, over := th.TakeHandler(th.Bytes)
	th.Bytes = o
	return r, over
}

// Recv 等待接收的模式 todo
func (th *Decoder) Recv(f func([]byte)) {
	th.enableRecv = true
	th.recvHandler = f
}

func (th *Decoder) Put(data []byte) {
	if len(data) > 0 {
		th.Lock()
		defer th.Unlock()
		th.Bytes = append(th.Bytes, data...)
		// 触发接收模式
		if th.enableRecv {
			for {
				d, _ := th.Take()
				if len(d) == 0 {
					break
				} else {
					th.recvHandler(d)
				}
			}
		}
	}
}
