[BITS 16]
[ORG 0X7C00]
;program

boot:
	mov al, 'H'
	mov ah, 0x0e
	mov bh, 0x00
	mov bl, 0x07

	int 0x10
	jmp $

times 510-($-$$) db 0
db 0x55
db 0xaa

