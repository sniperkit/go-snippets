package bitmap

import "errors"

var ErrOutOfRange = errors.New("bitmap: out of range.")

type Bitmap struct {
	size int
	data []byte
}

func New(size int) *Bitmap {
	if size < 0 {
		size = 0
	}
	return &Bitmap{
		size: size,
		data: make([]byte, size/8+1),
	}
}

func (bm *Bitmap) Size() int {
	return bm.size
}

func (bm *Bitmap) checkOutOfRange(i int) error {
	if (i < 0) || (i >= bm.size) {
		return ErrOutOfRange
	}
	return nil
}

func (bm *Bitmap) Set(i int, v bool) error {

	if err := bm.checkOutOfRange(i); err != nil {
		return err
	}

	index, shift := quoRem(i, 8)
	if v {
		bm.data[index] |= 1 << uint(shift)
	} else {
		bm.data[index] &= ^(1 << uint(shift))
	}

	return nil
}

func (bm *Bitmap) Get(i int) (v bool, err error) {

	if err := bm.checkOutOfRange(i); err != nil {
		return false, err
	}

	index, shift := quoRem(i, 8)
	v = ((bm.data[index] >> uint(shift)) & 1) == 1

	return v, nil
}

func quoRem(x, y int) (quo, rem int) {
	quo = x / y
	rem = x - quo*y
	return
}
