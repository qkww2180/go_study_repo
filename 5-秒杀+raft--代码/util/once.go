package util

import (
	"sync"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	if o.done == 0 {
		// o.doSlow(f)
		o.m.Lock()
		defer o.m.Unlock()
		if o.done == 0 {
			f()
			o.done = 1
		}
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		f()
		o.done = 1
	}
}
