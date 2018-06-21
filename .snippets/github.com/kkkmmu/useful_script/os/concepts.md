1. Register - register is a small amount of storage inside processor. Main point of processor is data processing. Processor can get data from memory, but it is slow operation. That’s why processor has own internal restricted set of data storage which name is - register.
2. Little-endian - we can imagine memory as one large array. It contains bytes. Each address stores one element of the memory “array”. Each element is one byte. For example we have 4 bytes: AA 56 AB FF. In little-endian the least significant byte has the smallest address:
	 0 FF
	 1 AB
	 2 56
	 3 AA
	 where 0,1,2 and 3 are memory addresses.

	 Big-endian - big-endian stores bytes in opposite order than little-endian. So if we have AA 56 AB FF bytes sequence it will be:
	 0 AA
	 1 56
	 2 AB
	 3 FF
3. Syscall - is the way a user level program asks the operating system to do something for it.
4. Section - every assembly program consists from sections. There are following sections:
	data - section is used for declaring initialized data or constants
	bss - section is used for declaring non initialized variables
	text - section is used for code

5.  NASM supports a number of pseudo-instructions:
	1. DB, DW, DD, DQ, DT, DO, DY and DZ - are used for declaring initialized data. For example:
		;; Initialize 4 bytes 1h, 2h, 3h, 4h
		db 0x01,0x02,0x03,0x04

		;; Initialize word to 0x12 0x34
		dw    0x1234
	2. RESB, RESW, RESD, RESQ, REST, RESO, RESY and RESZ - are used for declaring non initialized variables
	3. INCBIN - includes External Binary Files
	4. EQU - defines constant. For example:
		;; now one is 1
		one equ 1
	5. TIMES - Repeating Instructions or Data. 
	6. ADD - integer add
	7. SUB - substract
	8. MUL - unsigned multiply
	9. IMUL - signed multiply
	10. DIV - unsigned divide
	11. IDIV - signed divide
	12. INC - increment
	13. DEC - decrement
	14. NEG - negate
	15. JE - if equal
	16. JZ - if zero
	17. JNE - if not equal
	18. JNZ - if not zero
	19. JG - if first operand is greater than second
	20. JGE - if first operand is greater or equal to second
	21. JA - the same that JG, but performs unsigned comparison
	22. JAE - the same that JGE, but performs unsigned comparison
	23. JMP - Unconditional Jump
	24. CMP - compare two value and set the coresponding flags.
6.  System V AMD64 ABI, the first six function arguments passed in registers. They are:

	rdi - first argument
	rsi - second argument
	rdx - third argument
	rcx - fourth argument
	r8 - fifth argument
	r9 - sixth
	Next arguments will be passed in stack. So if we have function like this:
		int foo(int a1, int a2, int a3, int a4, int a5, int a6, int a7)
		{
	    	return (a1 + a2 - a3 - a4 + a5 - a6) * a7;
		}
	Then first six arguments will be passed in registers, but 7 argument will be passed in stack.
7. RBP is the base pointer register. It points to the base of the current stack frame. RSP is the stack pointer, which points to the top of current stack frame.
8. We have two commands for work with stack:

	push argument - increments stack pointer (RSP) and stores argument in location pointed by stack pointer
	pop argument - copied data to argument from location pointed by stack pointer
9. NASM supports two form of macro:

	1. single-line
	2. multiline
	All single-line macro must start from %define directive. It form is following:
		%define macro_name(parameter) value
	Nasm macro behaves and looks very similar as in C. For example, we can create following single-line macro:
		%define argc rsp + 8
		%define cliArg1 rsp + 24
	and than use it in code:

		mov rax, [argc]
		cmp rax, 3
		jne .mustBe3args
	Multiline macro starts with %macro nasm directive and end with %endmacro. It general form is following:

		%macro number_of_parameters
    		instruction
	   		instruction
	    	instruction
		%endmacro
	For example:
			%macro bootstrap 1
		          push ebp
		          mov ebp,esp
			%endmacro
	And we can use it:
			_start:
			      bootstrap
