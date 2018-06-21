	org 07c00h ; load this program to 7c00

	mov ax, cs
	mov ds, ax
	mov es, ax

call DispStr
	jmp $ ; infinite loop

DispStr:
	mov ax, BootMessage
	mov bp, ax
	mov cx, 16; str length
	mov ax, 01301h
	mov bx, 000bh; page number 0, background: black, color red
	mov dl, 0
	int 10h; interrupt 10
	ret
BootMessage:
	db "Hello, OS world!"
	times 510-($-$$) db 0 ;; fill the remain space 
	dw 0xaa55; this is the identifier for end.
