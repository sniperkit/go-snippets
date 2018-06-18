package proto

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"bytes"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestMarshalUnmarshalBar1(t *testing.T) {
	v := Bar1{}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

func BenchmarkMarshalMsgBar1(b *testing.B) {
	v := Bar1{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.MarshalMsg(nil)
	}
}

func BenchmarkAppendMsgBar1(b *testing.B) {
	v := Bar1{}
	bts := make([]byte, 0, v.Msgsize())
	bts, _ = v.MarshalMsg(bts[0:0])
	b.SetBytes(int64(len(bts)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bts, _ = v.MarshalMsg(bts[0:0])
	}
}

func BenchmarkUnmarshalBar1(b *testing.B) {
	v := Bar1{}
	bts, _ := v.MarshalMsg(nil)
	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncodeDecodeBar1(t *testing.T) {
	v := Bar1{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)

	m := v.Msgsize()
	if buf.Len() > m {
		t.Logf("WARNING: Msgsize() for %v is inaccurate", v)
	}

	vn := Bar1{}
	err := msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	msgp.Encode(&buf, &v)
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkEncodeBar1(b *testing.B) {
	v := Bar1{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)
	b.SetBytes(int64(buf.Len()))
	en := msgp.NewWriter(msgp.Nowhere)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.EncodeMsg(en)
	}
	en.Flush()
}

func BenchmarkDecodeBar1(b *testing.B) {
	v := Bar1{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)
	b.SetBytes(int64(buf.Len()))
	rd := msgp.NewEndlessReader(buf.Bytes(), b)
	dc := msgp.NewReader(rd)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := v.DecodeMsg(dc)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMarshalUnmarshalBar2(t *testing.T) {
	v := Bar2{}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

func BenchmarkMarshalMsgBar2(b *testing.B) {
	v := Bar2{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.MarshalMsg(nil)
	}
}

func BenchmarkAppendMsgBar2(b *testing.B) {
	v := Bar2{}
	bts := make([]byte, 0, v.Msgsize())
	bts, _ = v.MarshalMsg(bts[0:0])
	b.SetBytes(int64(len(bts)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bts, _ = v.MarshalMsg(bts[0:0])
	}
}

func BenchmarkUnmarshalBar2(b *testing.B) {
	v := Bar2{}
	bts, _ := v.MarshalMsg(nil)
	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncodeDecodeBar2(t *testing.T) {
	v := Bar2{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)

	m := v.Msgsize()
	if buf.Len() > m {
		t.Logf("WARNING: Msgsize() for %v is inaccurate", v)
	}

	vn := Bar2{}
	err := msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	msgp.Encode(&buf, &v)
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkEncodeBar2(b *testing.B) {
	v := Bar2{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)
	b.SetBytes(int64(buf.Len()))
	en := msgp.NewWriter(msgp.Nowhere)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.EncodeMsg(en)
	}
	en.Flush()
}

func BenchmarkDecodeBar2(b *testing.B) {
	v := Bar2{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)
	b.SetBytes(int64(buf.Len()))
	rd := msgp.NewEndlessReader(buf.Bytes(), b)
	dc := msgp.NewReader(rd)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := v.DecodeMsg(dc)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMarshalUnmarshalFooWithPointer(t *testing.T) {
	v := FooWithPointer{}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

func BenchmarkMarshalMsgFooWithPointer(b *testing.B) {
	v := FooWithPointer{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.MarshalMsg(nil)
	}
}

func BenchmarkAppendMsgFooWithPointer(b *testing.B) {
	v := FooWithPointer{}
	bts := make([]byte, 0, v.Msgsize())
	bts, _ = v.MarshalMsg(bts[0:0])
	b.SetBytes(int64(len(bts)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bts, _ = v.MarshalMsg(bts[0:0])
	}
}

func BenchmarkUnmarshalFooWithPointer(b *testing.B) {
	v := FooWithPointer{}
	bts, _ := v.MarshalMsg(nil)
	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncodeDecodeFooWithPointer(t *testing.T) {
	v := FooWithPointer{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)

	m := v.Msgsize()
	if buf.Len() > m {
		t.Logf("WARNING: Msgsize() for %v is inaccurate", v)
	}

	vn := FooWithPointer{}
	err := msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	msgp.Encode(&buf, &v)
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkEncodeFooWithPointer(b *testing.B) {
	v := FooWithPointer{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)
	b.SetBytes(int64(buf.Len()))
	en := msgp.NewWriter(msgp.Nowhere)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.EncodeMsg(en)
	}
	en.Flush()
}

func BenchmarkDecodeFooWithPointer(b *testing.B) {
	v := FooWithPointer{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)
	b.SetBytes(int64(buf.Len()))
	rd := msgp.NewEndlessReader(buf.Bytes(), b)
	dc := msgp.NewReader(rd)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := v.DecodeMsg(dc)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMarshalUnmarshalOper1(t *testing.T) {
	v := Oper1{}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

func BenchmarkMarshalMsgOper1(b *testing.B) {
	v := Oper1{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.MarshalMsg(nil)
	}
}

func BenchmarkAppendMsgOper1(b *testing.B) {
	v := Oper1{}
	bts := make([]byte, 0, v.Msgsize())
	bts, _ = v.MarshalMsg(bts[0:0])
	b.SetBytes(int64(len(bts)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bts, _ = v.MarshalMsg(bts[0:0])
	}
}

func BenchmarkUnmarshalOper1(b *testing.B) {
	v := Oper1{}
	bts, _ := v.MarshalMsg(nil)
	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncodeDecodeOper1(t *testing.T) {
	v := Oper1{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)

	m := v.Msgsize()
	if buf.Len() > m {
		t.Logf("WARNING: Msgsize() for %v is inaccurate", v)
	}

	vn := Oper1{}
	err := msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	msgp.Encode(&buf, &v)
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkEncodeOper1(b *testing.B) {
	v := Oper1{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)
	b.SetBytes(int64(buf.Len()))
	en := msgp.NewWriter(msgp.Nowhere)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.EncodeMsg(en)
	}
	en.Flush()
}

func BenchmarkDecodeOper1(b *testing.B) {
	v := Oper1{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)
	b.SetBytes(int64(buf.Len()))
	rd := msgp.NewEndlessReader(buf.Bytes(), b)
	dc := msgp.NewReader(rd)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := v.DecodeMsg(dc)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMarshalUnmarshalOper2(t *testing.T) {
	v := Oper2{}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

func BenchmarkMarshalMsgOper2(b *testing.B) {
	v := Oper2{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.MarshalMsg(nil)
	}
}

func BenchmarkAppendMsgOper2(b *testing.B) {
	v := Oper2{}
	bts := make([]byte, 0, v.Msgsize())
	bts, _ = v.MarshalMsg(bts[0:0])
	b.SetBytes(int64(len(bts)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bts, _ = v.MarshalMsg(bts[0:0])
	}
}

func BenchmarkUnmarshalOper2(b *testing.B) {
	v := Oper2{}
	bts, _ := v.MarshalMsg(nil)
	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncodeDecodeOper2(t *testing.T) {
	v := Oper2{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)

	m := v.Msgsize()
	if buf.Len() > m {
		t.Logf("WARNING: Msgsize() for %v is inaccurate", v)
	}

	vn := Oper2{}
	err := msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	msgp.Encode(&buf, &v)
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkEncodeOper2(b *testing.B) {
	v := Oper2{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)
	b.SetBytes(int64(buf.Len()))
	en := msgp.NewWriter(msgp.Nowhere)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.EncodeMsg(en)
	}
	en.Flush()
}

func BenchmarkDecodeOper2(b *testing.B) {
	v := Oper2{}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)
	b.SetBytes(int64(buf.Len()))
	rd := msgp.NewEndlessReader(buf.Bytes(), b)
	dc := msgp.NewReader(rd)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := v.DecodeMsg(dc)
		if err != nil {
			b.Fatal(err)
		}
	}
}
