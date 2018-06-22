package neogo

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
)

// OpCode script opcode
type OpCode byte

// OpCode
const (
	PUSH0           OpCode = 0x00 // An empty array of bytes is pushed onto the stack.
	PUSHF                  = PUSH0
	PUSHBYTES1             = 0x01 // 0x01-0x4B The next opcode bytes is data to be pushed onto the stack
	PUSHBYTES75            = 0x4B
	PUSHDATA1              = 0x4C // The next byte contains the number of bytes to be pushed onto the stack.
	PUSHDATA2              = 0x4D // The next two bytes contain the number of bytes to be pushed onto the stack.
	PUSHDATA4              = 0x4E // The next four bytes contain the number of bytes to be pushed onto the stack.
	PUSHM1                 = 0x4F // The number -1 is pushed onto the stack.
	PUSH1                  = 0x51 // The number 1 is pushed onto the stack.
	PUSHT                  = PUSH1
	PUSH2                  = 0x52 // The number 2 is pushed onto the stack.
	PUSH3                  = 0x53 // The number 3 is pushed onto the stack.
	PUSH4                  = 0x54 // The number 4 is pushed onto the stack.
	PUSH5                  = 0x55 // The number 5 is pushed onto the stack.
	PUSH6                  = 0x56 // The number 6 is pushed onto the stack.
	PUSH7                  = 0x57 // The number 7 is pushed onto the stack.
	PUSH8                  = 0x58 // The number 8 is pushed onto the stack.
	PUSH9                  = 0x59 // The number 9 is pushed onto the stack.
	PUSH10                 = 0x5A // The number 10 is pushed onto the stack.
	PUSH11                 = 0x5B // The number 11 is pushed onto the stack.
	PUSH12                 = 0x5C // The number 12 is pushed onto the stack.
	PUSH13                 = 0x5D // The number 13 is pushed onto the stack.
	PUSH14                 = 0x5E // The number 14 is pushed onto the stack.
	PUSH15                 = 0x5F // The number 15 is pushed onto the stack.
	PUSH16                 = 0x60 // The number 16 is pushed onto the stack.
	NOP                    = 0x61 // Does nothing.
	JMP                    = 0x62
	JMPIF                  = 0x63
	JMPIFNOT               = 0x64
	CALL                   = 0x65
	RET                    = 0x66
	APPCALL                = 0x67
	SYSCALL                = 0x68
	TAILCALL               = 0x69
	DUPFROMALTSTACK        = 0x6A
	TOALTSTACK             = 0x6B // Puts the input onto the top of the alt stack. Removes it from the main stack.
	FROMALTSTACK           = 0x6C // Puts the input onto the top of the main stack. Removes it from the alt stack.
	XDROP                  = 0x6D
	XSWAP                  = 0x72
	XTUCK                  = 0x73
	DEPTH                  = 0x74 // Puts the number of stack items onto the stack.
	DROP                   = 0x75 // Removes the top stack item.
	DUP                    = 0x76 // Duplicates the top stack item.
	NIP                    = 0x77 // Removes the second-to-top stack item.
	OVER                   = 0x78 // Copies the second-to-top stack item to the top.
	PICK                   = 0x79 // The item n back in the stack is copied to the top.
	ROLL                   = 0x7A // The item n back in the stack is moved to the top.
	ROT                    = 0x7B // The top three items on the stack are rotated to the left.
	SWAP                   = 0x7C // The top two items on the stack are swapped.
	TUCK                   = 0x7D // The item at the top of the stack is copied and inserted before the second-to-top item.
	CAT                    = 0x7E // Concatenates two strings.
	SUBSTR                 = 0x7F // Returns a section of a string.
	LEFT                   = 0x80 // Keeps only characters left of the specified point in a string.
	RIGHT                  = 0x81 // Keeps only characters right of the specified point in a string.
	SIZE                   = 0x82 // Returns the length of the input string.
	INVERT                 = 0x83 // Flips all of the bits in the input.
	AND                    = 0x84 // Boolean and between each bit in the inputs.
	OR                     = 0x85 // Boolean or between each bit in the inputs.
	XOR                    = 0x86 // Boolean exclusive or between each bit in the inputs.
	EQUAL                  = 0x87 // Returns 1 if the inputs are exactly equal 0 otherwise.
	INC                    = 0x8B // 1 is added to the input.
	DEC                    = 0x8C // 1 is subtracted from the input.
	SIGN                   = 0x8D
	NEGATE                 = 0x8F // The sign of the input is flipped.
	ABS                    = 0x90 // The input is made positive.
	NOT                    = 0x91 // If the input is 0 or 1 it is flipped. Otherwise the output will be 0.
	NZ                     = 0x92 // Returns 0 if the input is 0. 1 otherwise.
	ADD                    = 0x93 // a is added to b.
	SUB                    = 0x94 // b is subtracted from a.
	MUL                    = 0x95 // a is multiplied by b.
	DIV                    = 0x96 // a is divided by b.
	MOD                    = 0x97 // Returns the remainder after dividing a by b.
	SHL                    = 0x98 // Shifts a left b bits preserving sign.
	SHR                    = 0x99 // Shifts a right b bits preserving sign.
	BOOLAND                = 0x9A // If both a and b are not 0 the output is 1. Otherwise 0.
	BOOLOR                 = 0x9B // If a or b is not 0 the output is 1. Otherwise 0.
	NUMEQUAL               = 0x9C // Returns 1 if the numbers are equal 0 otherwise.
	NUMNOTEQUAL            = 0x9E // Returns 1 if the numbers are not equal 0 otherwise.
	LT                     = 0x9F // Returns 1 if a is less than b 0 otherwise.
	GT                     = 0xA0 // Returns 1 if a is greater than b 0 otherwise.
	LTE                    = 0xA1 // Returns 1 if a is less than or equal to b 0 otherwise.
	GTE                    = 0xA2 // Returns 1 if a is greater than or equal to b 0 otherwise.
	MIN                    = 0xA3 // Returns the smaller of a and b.
	MAX                    = 0xA4 // Returns the larger of a and b.
	WITHIN                 = 0xA5 // Returns 1 if x is within the specified range (left-inclusive) 0 otherwise.
	SHA1                   = 0xA7 // The input is hashed using SHA-1.
	SHA256                 = 0xA8 // The input is hashed using SHA-256.
	HASH160                = 0xA9
	HASH256                = 0xAA
	CHECKSIG               = 0xAC
	CHECKMULTISIG          = 0xAE
	ARRAYSIZE              = 0xC0
	PACK                   = 0xC1
	UNPACK                 = 0xC2
	PICKITEM               = 0xC3
	SETITEM                = 0xC4
	NEWARRAY               = 0xC5 //用作引用類型
	NEWSTRUCT              = 0xC6 //用作值類型
	THROW                  = 0xF0
	THROWIFNOT             = 0xF1
)

// ScriptWriter neo script writer
type ScriptWriter struct {
	writer io.Writer
	Error  error
	depth  int
}

// NewScriptWriter create new script writer
func NewScriptWriter(writer io.Writer) *ScriptWriter {
	return &ScriptWriter{
		writer: writer,
	}
}

// Emit .
func (writer *ScriptWriter) Emit(opcode OpCode, arg []byte) *ScriptWriter {
	if writer.Error != nil {
		return writer
	}

	if arg != nil {
		writer.writer.Write(append([]byte{byte(opcode)}, arg...))
	} else {
		writer.writer.Write([]byte{byte(opcode)})
	}

	writer.depth++

	return writer
}

// EmitAPPCall .
func (writer *ScriptWriter) EmitAPPCall(scriptHash []byte, tailCall bool) *ScriptWriter {
	if len(scriptHash) != 20 {
		writer.Error = fmt.Errorf("[%d] EmitAPPCall scriptHash length must be 20 bytes", writer.depth)
		return writer
	}

	if tailCall {
		return writer.Emit(TAILCALL, scriptHash)
	}

	return writer.Emit(APPCALL, scriptHash)
}

// EmitJump .
func (writer *ScriptWriter) EmitJump(op OpCode, offset int16) *ScriptWriter {
	if op != JMP && op != JMPIF && op != JMPIFNOT && op != CALL {
		writer.Error = fmt.Errorf("[%d] EmitAPPCall scriptHash length must be 20 bytes", writer.depth)
		return writer
	}

	data := make([]byte, 2)

	binary.LittleEndian.PutUint16(data, uint16(offset))

	writer.Emit(op, data)

	return writer
}

// EmitPushInteger .
func (writer *ScriptWriter) EmitPushInteger(number *big.Int) *ScriptWriter {
	if number.Int64() == -1 {
		return writer.Emit(PUSHM1, nil)
	}

	if number.Int64() == 0 {
		return writer.Emit(PUSH0, nil)
	}

	if number.Int64() > 0 && number.Int64() <= 16 {
		return writer.Emit(OpCode(byte(PUSH1)-1+byte(number.Int64())), nil)
	}

	data := reverseBytes(number.Bytes())

	return writer.EmitPushBytes(data)
}

func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

// EmitPushBool .
func (writer *ScriptWriter) EmitPushBool(data bool) *ScriptWriter {
	if data {
		return writer.Emit(PUSHT, nil)
	}

	return writer.Emit(PUSHF, nil)
}

// EmitPushString .
func (writer *ScriptWriter) EmitPushString(data string) *ScriptWriter {
	return writer.EmitPushBytes([]byte(data))
}

// EmitPushBytes .
func (writer *ScriptWriter) EmitPushBytes(data []byte) *ScriptWriter {

	if writer.Error != nil {
		return writer
	}

	if data == nil {
		writer.Error = fmt.Errorf("[%d] EmitPushBytes args can't be null", writer.depth)
		return writer
	}

	if len(data) < int(PUSHBYTES75) {
		return writer.Emit(OpCode(len(data)), data)
	}

	if len(data) < int(0x100) {
		writer.Emit(PUSHDATA1, nil)
		writer.writer.Write([]byte{byte(len(data))})
		writer.writer.Write(data)

		writer.depth++

		return writer
	}

	if len(data) < int(0x10000) {
		writer.Emit(PUSHDATA2, nil)
		bytesOfLength := make([]byte, 2)

		binary.LittleEndian.PutUint16(bytesOfLength, uint16(len(data)))

		writer.writer.Write(bytesOfLength)
		writer.writer.Write(data)

		writer.depth++

		return writer
	}

	writer.Emit(PUSHDATA4, nil)

	bytesOfLength := make([]byte, 4)

	binary.LittleEndian.PutUint32(bytesOfLength, uint32(len(data)))

	writer.writer.Write(bytesOfLength)
	writer.writer.Write(data)

	writer.depth++

	return writer
}

// EmitSysCall .
func (writer *ScriptWriter) EmitSysCall(api string) *ScriptWriter {
	if api == "" {
		writer.Error = fmt.Errorf("[%d] EmitSysCall api parameter can't be empty", writer.depth)
	}

	bytesOfAPI := []byte(api)

	if len(bytesOfAPI) > 252 {
		writer.Error = fmt.Errorf("[%d] EmitSysCall api name can't longer than 252", writer.depth)
	}

	return writer.Emit(SYSCALL, append([]byte{byte(len(bytesOfAPI))}, bytesOfAPI...))
}

// Offset .
func (writer *ScriptWriter) Offset() int {
	return writer.depth
}

// Reset reset script buffer
func (writer *ScriptWriter) Reset(newwriter io.Writer) *ScriptWriter {
	writer.depth = 0
	writer.Error = nil
	writer.writer = newwriter
	return writer
}
