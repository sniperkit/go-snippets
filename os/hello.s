section .data			;data section
	msg db "Hello, world!"

section .text			;code section
	global _start		;Entry point

_start:
	mov rax, 1			;syscall number 1 (sys_write)
	mov rdi, 1;			;1st parameters(fd)
	mov rsi, msg		;2nd parameters(buff)
	mov rdx, 13			;3rd parameters(size)
	syscall
	mov rax, 60			;syscal number 60 (sys_exit)
	mov rdi, -1			;1st parameters (0)
	syscall

;nasm -f elf64 -o hello.o hello.asm
;ld -o hello hello.o
;echo $? 
