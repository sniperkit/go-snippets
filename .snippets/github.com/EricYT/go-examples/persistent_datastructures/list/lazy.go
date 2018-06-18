package list

import "sync"

type lazy struct {
	o    sync.Once
	f    func() (int, List)
	v    int
	next List
}

func (l *lazy) Value() int {
	l.o.Do(func() { l.v, l.next = l.f() })
	return l.v
}

func (l *lazy) Next() List {
	l.o.Do(func() { l.v, l.next = l.f() })
	return l.next
}
