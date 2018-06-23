section .data
	SYS_WRITE equ 1
	STD_OUT   equ 1
	SYS_EXIT  equ 60
	EXIT_CODE equ 0

	NEW_LINE db 0xa
	INPUT db "Hello, welcome to china!"

section .bss
OUTPUT: 
 	align 0x4
	resb 12

section .text
	global _start

_start:
	mov rsi, INPUT
	xor rcx, rcx			;;counter for calculating length of the string
	cld						;;reset df flag to zero
	mov rdi, $ + 15			;;save the return address (15) is the offset
	call calculateStrLength
	xor rax, rax
	xor rdi, rdi
	jmp reverseStr

calculateStrLength:
	cmp byte [rsi], 0		;; check is it end of string
	je exitFromRoutine 		;; if yes exit from function
	lodsb					;; load 1 byte from rsi to al and inc rsi
	push rax				;; push symbol to stack
	inc rcx					;; increase counter
	jmp calculateStrLength	;; loop again

exitFromRoutine:
	push rdi				;; push return addres to stack again
	ret						;; return to _start

reverseStr:
	cmp rcx, 0
	je printResult
	pop rax
	mov [OUTPUT + rdi], rax	;;save to ouput buffer
	dec rcx
	inc rdi
	jmp reverseStr

printResult:
	mov rdx, rdi
	mov rax, 1
	mov rdi, 1
	mov rsi, OUTPUT
	syscall
	jmp printNewLine

printNewLine:
	mov rax, SYS_WRITE
	mov rdi, STD_OUT
	mov rsi, NEW_LINE
	mov rdx, 1
	syscall
	jmp exit

exit:
	mov rax, SYS_EXIT
	mov rdi, EXIT_CODE

;$$指向当前section相对于段基址的偏移地址，$指向当前行相对于段基址的偏移地址。
