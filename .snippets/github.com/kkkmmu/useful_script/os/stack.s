section .data
	SYS_WRITE 	equ 1
	STD_IN	  	equ 1
	SYS_EXIT  	equ 60
	EXIT_CODE 	equ 0
	NEW_LINE  	db 0xa
	WRONG_ARGC 	db "Must assign two input parameters!"

section .text
	global _start

;If  we runn application with comand line arguments:
	;[rsp] -	top of stack will contain arguments count.
	;[rsp+8] - 	will contain argv[0]
	;[rsp+16] - will contain argv[1]
	;...and so on...
_start:
	pop rcx			;get argument count
	cmp rcx, 3
	jne argcError

	add rsp, 8		;first argument address
	pop rsi			;get first argument
	call str_to_int

	mov r10, rax
	pop rsi			;get second argument
	call str_to_int
	mov r11, rax

	add r10, r11
	mov rax, r10
	
	xor r12, r12
	jmp int_to_str


argcError:
	;;sys_write syscall
	mov rax, 1
	;; file descriptor, standard output
	mov rdi, 1
	;; message buffer address
	mov rsi, WRONG_ARGC
	;; length of message
	mov rdx, 34
	;; call write syscall
	syscall
	;;exit from program
	jmp exit

str_to_int:
	xor rax, rax
	mov rcx, 10

next:
	cmp [rsi], byte 0
	je return_str
	mov bl, [rsi]
	sub bl, 48
	mul rcx
	add rax, rbx
	inc rsi
	jmp next

return_str:
	ret

int_to_str:
	mov rdx, 0
	mov rbx, 10
	div rbx			;divide rax with rbx(10) and get the reminder in rdx and whole part in rax
	add rdx, 48		;get ascii symbol
	add rdx, 0x0
	push rdx		;Save to stack
	inc r12
	cmp rax, 0x0
	jne int_to_str
	jmp print

print:
	;;; calculate number length
	mov rax, 1
	mul r12
	mov r12, 8
	mul r12
	mov rdx, rax

	;;; print sum
	mov rax, SYS_WRITE
	mov rdi, STD_IN
	mov rsi, rsp
	;; call sys_write
	syscall
	
	jmp exit

exit:
	mov rax, SYS_EXIT
	;;exit code
	mov rdi, EXIT_CODE
	syscall
