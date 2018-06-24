package internal

import (
	"bufio"
	"fmt"
	"io"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/iocount"
	"github.com/bsm/reason/internal/protoio"
	"github.com/gogo/protobuf/proto"
)

// NewOptimizer inits a new model with defaults
func NewOptimizer(model *core.Model, target string, size int) *Optimizer {
	return &Optimizer{
		Model:   model,
		Target:  target,
		Sums:    make([]float64, size),
		Weights: make([]float64, size),
	}
}

// WriteTo writes a tree to a Writer.
func (o *Optimizer) WriteTo(w io.Writer) (int64, error) {
	wc := &iocount.Writer{W: w}
	wp := &protoio.Writer{Writer: bufio.NewWriter(wc)}

	if o.Model != nil {
		if err := wp.WriteMessageField(1, o.Model); err != nil {
			return wc.N, err
		}
	}
	if o.Target != "" {
		if err := wp.WriteStringField(2, o.Target); err != nil {
			return wc.N, err
		}
	}
	if len(o.Sums) != 0 {
		if err := wp.WriteField(3, proto.WireBytes); err != nil {
			return wc.N, err
		}
		if err := wp.WriteVarint(uint64(len(o.Sums) * 8)); err != nil {
			return wc.N, err
		}
		for _, f := range o.Sums {
			if err := wp.WriteDouble(f); err != nil {
				return wc.N, err
			}
		}
	}
	if len(o.Weights) != 0 {
		if err := wp.WriteField(4, proto.WireBytes); err != nil {
			return wc.N, err
		}
		if err := wp.WriteVarint(uint64(len(o.Weights) * 8)); err != nil {
			return wc.N, err
		}
		for _, f := range o.Weights {
			if err := wp.WriteDouble(f); err != nil {
				return wc.N, err
			}
		}
	}
	return wc.N, wp.Flush()
}

// ReadFrom reads an optimizer from a Reader.
func (o *Optimizer) ReadFrom(r io.Reader) (int64, error) {
	rc := &iocount.Reader{R: r}
	rp := &protoio.Reader{Reader: bufio.NewReader(rc)}

	for {
		tag, wire, err := rp.ReadField()
		if err == io.EOF {
			return rc.N, nil
		} else if err != nil {
			return rc.N, err
		}

		switch tag {
		case 1: // model
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			model := new(core.Model)
			if err := rp.ReadMessage(model); err != nil {
				return rc.N, err
			}
			o.Model = model
		case 2: // target
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			str, err := rp.ReadString()
			if err != nil {
				return rc.N, err
			}
			o.Target = str
		case 3: // sums
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			slice, err := readFloatSlice(rp)
			if err != nil {
				return rc.N, err
			}
			o.Sums = slice
		case 4: // weights
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			slice, err := readFloatSlice(rp)
			if err != nil {
				return rc.N, err
			}
			o.Weights = slice
		default:
			return rc.N, fmt.Errorf("hoeffding: unexpected field tag %d", tag)
		}
	}
}

func readFloatSlice(rp *protoio.Reader) ([]float64, error) {
	u, err := rp.ReadVarint()
	if err != nil {
		return nil, err
	}
	n := int(u / 8)
	slice := make([]float64, 0, n)

	for i := 0; i < n; i++ {
		f, err := rp.ReadDouble()
		if err != nil {
			return nil, err
		}
		slice = append(slice, f)
	}
	return slice, nil
}
