package color

import (
	"errors"
	"fmt"
)

type ColorRGBA interface {
	GetColor() (r, g, b, a uint16)
	SetColor(r, g, b, a uint16)
}

type Chanel int32

const (
	chanelSize Chanel = (1 << 16)
)

func ChanelIsValid(ch Chanel) bool {

	return (ch >= 0) && (ch < chanelSize)
}

type Chanels struct {
	R, G, B, A Chanel
}

func (cs Chanels) IsValid() (res bool) {

	if ChanelIsValid(cs.R) {
		if ChanelIsValid(cs.G) {
			if ChanelIsValid(cs.B) {
				if ChanelIsValid(cs.A) {
					res = true
				}
			}
		}
	}

	return
}

type Color interface {
	GetChanels() (ch Chanels, ok bool)
	SetChanels(Chanels) error
}

type RgbColor uint32

const (
	redShift   = 0
	greenShift = 8
	blueShift  = 16
	validShift = 24
)

const (
	redMask   = 0xFF << redShift
	greenMask = 0xFF << greenShift
	blueMask  = 0xFF << blueShift
	validMask = 0xFF << validShift
)

const (
	validSign RgbColor = 0x8A << validShift
)

const (
	White  RgbColor = validSign | 0xFFFFFF
	Black  RgbColor = validSign | 0x000000
	Red    RgbColor = validSign | 0x0000FF
	Green  RgbColor = validSign | 0x00FF00
	Blue   RgbColor = validSign | 0xFF0000
	Yellow RgbColor = validSign | 0x00FFFF
)

func (c RgbColor) IsValid() bool {
	return (c & validMask) == validSign
}

func getError(c RgbColor) (err error) {

	if !c.IsValid() {
		err = errors.New("RgbColor is not valid")
	}
	return
}

func (c *RgbColor) GetChanels() (cs Chanels, ok bool) {

	if c != nil {
		cs.R = Chanel((((*c) & redMask) >> redShift) << 8)
		cs.G = Chanel((((*c) & greenMask) >> greenShift) << 8)
		cs.B = Chanel((((*c) & blueMask) >> blueShift) << 8)
		cs.A = chanelSize - 1
		ok = true
	}
	return
}

func (c *RgbColor) SetChanels(cs Chanels) {
	if c != nil {
		if cs.IsValid() {

			var r int

			r |= ((int(cs.R) >> 8) << redShift) & redMask
			r |= ((int(cs.G) >> 8) << greenShift) & greenMask
			r |= ((int(cs.B) >> 8) << blueShift) & blueMask

			*c = RgbColor(r)
		}
	}
}

func (c RgbColor) Negative() (resColor RgbColor) {

	if cs, ok := c.GetChanels(); ok {

		cs.r = (chanelSize - 1) - cs.r
		cs.g = (chanelSize - 1) - cs.g
		cs.b = (chanelSize - 1) - cs.b

		resColor.SetChanels(cs)
	}

	return
}

func (c RgbColor) Grayscale() (resColor RgbColor) {

	if err, v := c.GetRgbValues(); err == nil {

		gray := byte((30*int(v.r) + 59*int(v.g) + 11*int(v.b)) / 100)

		resColor.SetRgbValues(gray, gray, gray)
	}

	return
}

/*
func (c RgbColor) GetRValue() (err error, r byte) {

	if err = getError(c); err == nil {
		r = byte((c & redMask) >> redShift)
	}

	return
}

func (c RgbColor) GetGValue() (err error, r byte) {

	if err = getError(c); err == nil {
		r = byte((c & greenMask) >> greenShift)
	}

	return
}

func (c RgbColor) GetBValue() (err error, r byte) {

	if err = getError(c); err == nil {
		r = byte((c & blueMask) >> blueShift)
	}

	return
}

func (c *RgbColor) SetRValue(r byte) (err error) {

	if c != nil {

	}
}

func (c *RgbColor) SetGValue(g Chanel) {

	c.g = g
}

func (c *RgbColor) SetBValue(b Chanel) {

	c.b = b
}
*/

func (c RgbColor) ToString() (str string) {

	if err, v := c.GetRgbValues(); err == nil {
		str = fmt.Sprintf("RGB(%v, %v, %v)", v.r, v.g, v.b)
	} else {
		str = fmt.Sprintln("RgbColor is invalid")
	}

	return
}

func Mix2(c1, c2 RgbColor) (c3 RgbColor) {

	err1, v1 := c1.GetRgbValues()
	if err1 == nil {
		err2, v2 := c2.GetRgbValues()
		if err2 == nil {
			r3 := byte((int(v1.r) + int(v2.r)) / 2)
			g3 := byte((int(v1.g) + int(v2.g)) / 2)
			b3 := byte((int(v1.b) + int(v2.b)) / 2)
			c3.SetRgbValues(r3, g3, b3)
		}
	}

	return
}

func Mix3(c1, c2, c3 RgbColor) (c4 RgbColor) {

	err1, v1 := c1.GetRgbValues()
	if err1 == nil {
		err2, v2 := c2.GetRgbValues()
		if err2 == nil {
			err3, v3 := c3.GetRgbValues()
			if err3 == nil {
				r4 := byte((int(v1.r) + int(v2.r) + int(v3.r)) / 3)
				g4 := byte((int(v1.g) + int(v2.g) + int(v3.g)) / 3)
				b4 := byte((int(v1.b) + int(v2.b) + int(v3.b)) / 3)
				c4.SetRgbValues(r4, g4, b4)
			}
		}
	}

	return
}

func Test1() {

	var c1 RgbColor
	var c2 RgbColor

	c1.SetRgbValues(34, 150, 255)
	c2.SetRgbValues(0, 200, 100)

	c3 := Mix2(c1, c2)

	fmt.Println(c3.ToString())
}
