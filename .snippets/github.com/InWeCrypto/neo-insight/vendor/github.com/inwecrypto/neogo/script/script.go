package script

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	"github.com/dynamicgo/slf4go"
	"golang.org/x/crypto/ripemd160"
)

// Script .
type Script struct {
	slf4go.Logger
	Ops   []*Op
	Name  string
	Error error
}

// New create new script with display name
func New(name string) *Script {
	return &Script{
		Logger: slf4go.Get("neogo-script"),
		Name:   name,
	}
}

// Reset .
func (script *Script) Reset() {
	script.Ops = nil
	script.Error = nil
}

// Emit emit one op
func (script *Script) Emit(opcode OpCode, arg []byte) *Script {

	if !script.checkEmit() {
		return script
	}

	op := &Op{
		Code: opcode,
		Arg:  arg,
	}

	script.DebugF("emit %s", op)

	script.Ops = append(script.Ops, op)

	return script
}

// EmitJump .
func (script *Script) EmitJump(op OpCode, offset int16) *Script {
	if op != JMP && op != JMPIF && op != JMPIFNOT && op != CALL {
		script.Error = fmt.Errorf("[%d] invalid EmitJump opcode %s", len(script.Ops), op2Strings[op])
		return script
	}

	data := make([]byte, 2)

	binary.LittleEndian.PutUint16(data, uint16(offset))

	script.Emit(op, data)

	return script
}

// EmitAPPCall .
func (script *Script) EmitAPPCall(scriptHash []byte, tailCall bool) *Script {
	if len(scriptHash) != 20 {
		script.Error = fmt.Errorf("[%d] EmitAPPCall scriptHash length must be 20 bytes", len(script.Ops))
		return script
	}

	if tailCall {
		return script.Emit(TAILCALL, scriptHash)
	}

	return script.Emit(APPCALL, scriptHash)
}

// EmitPushInteger .
func (script *Script) EmitPushInteger(number *big.Int) *Script {
	if number.Int64() == -1 {
		return script.Emit(PUSHM1, nil)
	}

	if number.Int64() == 0 {
		return script.Emit(PUSH0, nil)
	}

	if number.Int64() > 0 && number.Int64() <= 16 {
		return script.Emit(OpCode(byte(PUSH1)-1+byte(number.Int64())), nil)
	}

	data := reverseBytes(number.Bytes())

	if number.Int64() > 0 {
		data = append(data, 0x00)
	} else {
		data = append(data, 0x80)
	}

	return script.EmitPushBytes(data)
}

func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

// EmitPushBool .
func (script *Script) EmitPushBool(data bool) *Script {
	if data {
		return script.Emit(PUSHT, nil)
	}

	return script.Emit(PUSHF, nil)
}

// EmitPushString .
func (script *Script) EmitPushString(data string) *Script {
	return script.EmitPushBytes([]byte(data))
}

// EmitPushBytes .
func (script *Script) EmitPushBytes(data []byte) *Script {

	if script.Error != nil {
		return script
	}

	if data == nil {
		script.Error = fmt.Errorf("[%d] EmitPushBytes args can't be null", len(script.Ops))
		return script
	}

	if len(data) < int(PUSHBYTES75) {
		return script.Emit(OpCode(len(data)), data)
	}

	if len(data) < int(0x100) {

		var buff bytes.Buffer
		buff.Write([]byte{byte(len(data))})
		buff.Write(data)

		script.Emit(PUSHDATA1, buff.Bytes())

		return script
	}

	if len(data) < int(0x10000) {
		var buff bytes.Buffer

		bytesOfLength := make([]byte, 2)
		binary.LittleEndian.PutUint16(bytesOfLength, uint16(len(data)))

		buff.Write(bytesOfLength)
		buff.Write(data)

		script.Emit(PUSHDATA2, buff.Bytes())

		return script
	}

	var buff bytes.Buffer

	bytesOfLength := make([]byte, 4)

	binary.LittleEndian.PutUint32(bytesOfLength, uint32(len(data)))

	buff.Write(bytesOfLength)
	buff.Write(data)

	script.Emit(PUSHDATA4, buff.Bytes())

	return script
}

// EmitSysCall .
func (script *Script) EmitSysCall(api string) *Script {
	if api == "" {
		script.Error = fmt.Errorf("[%d] EmitSysCall api parameter can't be empty", len(script.Ops))
	}

	bytesOfAPI := []byte(api)

	if len(bytesOfAPI) > 252 {
		script.Error = fmt.Errorf("[%d] EmitSysCall api name can't longer than 252", len(script.Ops))
	}

	return script.Emit(SYSCALL, append([]byte{byte(len(bytesOfAPI))}, bytesOfAPI...))
}

func (script *Script) checkEmit() bool {
	return script.Error == nil
}

func (script *Script) Write(writer io.Writer) error {

	if script.Error != nil {
		return script.Error
	}

	for _, op := range script.Ops {
		_, err := writer.Write(append([]byte{byte(op.Code)}, op.Arg...))

		if err != nil {
			return err
		}
	}

	return nil
}

// Bytes get script bytes
func (script *Script) Bytes() ([]byte, error) {
	var buff bytes.Buffer

	if err := script.Write(&buff); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// Hash get script hash
func (script *Script) Hash() ([]byte, error) {

	buff, err := script.Bytes()

	if err != nil {
		return nil, err
	}

	/* SHA256 Hash */
	sha256h := sha256.New()
	sha256h.Reset()
	sha256h.Write(buff)
	pubhash1 := sha256h.Sum(nil)

	/* RIPEMD-160 Hash */
	ripemd160h := ripemd160.New()
	ripemd160h.Reset()
	ripemd160h.Write(pubhash1)
	pubhash2 := ripemd160h.Sum(nil)

	programhash := pubhash2

	return programhash, nil
}

// Hash .
func Hash(script []byte) []byte {
	/* SHA256 Hash */
	sha256h := sha256.New()
	sha256h.Reset()
	sha256h.Write(script)
	pubhash1 := sha256h.Sum(nil)

	/* RIPEMD-160 Hash */
	ripemd160h := ripemd160.New()
	ripemd160h.Reset()
	ripemd160h.Write(pubhash1)
	pubhash2 := ripemd160h.Sum(nil)

	programhash := pubhash2

	return programhash
}

// JSON .
func (script *Script) JSON() string {
	var ops []string

	for _, op := range script.Ops {
		ops = append(ops, op.String())
	}

	jsondata, _ := json.Marshal(ops)

	return string(jsondata)
}
