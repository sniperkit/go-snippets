package deque

import "container/list"

// inspired by juju/utils deque.
// read it and write it.

type Deque struct {
	maxLen            int
	blocks            list.List
	backIdx, frontIdx int
	len               int
}

const blockLen = 64
const blockCenter = (blockLen - 1) / 2

type blockT []interface{}

func newBlock() blockT {
	return make(blockT, blockLen)
}

func NewDeque() *Deque {
	return NewDequeWithMaxLen(0)
}

func NewDequeWithMaxLen(maxLen int) *Deque {
	deque := &Deque{maxLen: maxLen}
	deque.blocks.PushBack(newBlock())
	deque.recenter()
	return &deque
}

func (d *Deque) recenter() {
	d.frontIdx = blockCenter + 1
	d.backIdx = blockCenter
}

func (d *Deque) Len() int {
	return d.len
}

func (d *Deque) PushBack(item interface{}) {
	var block blockT
	if d.backIdx == blockLen {
		block = newBlock()
		d.blocks.PushBack(block)
		d.backIdx = -1
	} else {
		block = d.blocks.Back().Value.(blockT)
	}

	d.backIdx++
	block[d.backIdx] = item
	d.len++

	if d.maxLen > 0 && d.len > d.maxLen {
		d.PopBack()
	}
}

func (d *Deque) PushFront(item interface{}) {
	var block blockT
	if d.frontIdx == 0 {
		block = newBlock()
		d.blocks.PushFront(block)
		d.frontIdx = blockLen
	} else {
		block = d.blocks.Front().Value.(blockT)
	}

	d.frontIdx--
	block[d.frontIdx] = item
	d.len++

	if d.maxLen > 0 && d.len > d.maxLen {
		d.PopFront()
	}
}

func (d *Deque) PopBack() (interface{}, bool) {
	if d.len < 1 {
		return nil, false
	}

	elem := d.blocks.Back()
	block := elem.Value.(blockT)
	item := block[d.backIdx]
	block[d.backIdx] = nil
	d.backIdx--
	d.len--

	if d.backIdx == -1 {
		if d.len == 0 {
			d.recenter()
		} else {
			d.blocks.Remove(elem)
			d.backIdx = blockLen - 1
		}
	}

	return item, true
}

func (d *Deque) PopFront() (interface{}, bool) {
	if d.len < 1 {
		return nil, false
	}

	elem := d.blocks.Front()
	block := elem.Value.(blockT)
	item := block[d.frontIdx]
	block[d.frontIdx] = nil
	d.frontIdx++
	d.len--

	if d.frontIdx == blockLen {
		if d.len == 0 {
			d.recenter()
		} else {
			d.blocks.Remove(elem)
			d.frontIdx = 0
		}
	}

	return item, false
}
