---
title: "Python Handy Conversions"
date: 2018-03-09T02:35:48+08:00
tags: ["programming", "python"]
draft: true
---

`TODO: write content` 


- string to hex (import binascii or codecs): binascii.hexlify == binascii.b2a_hex == codecs.decode("...", "hex") == 'hello world'.decode('hex')
- hex to string (import binascii or codecs): binascii.unhexlify == binascii.a2b_hex == codecs.encode("...", "hex") == 'hex_string'.decode('hex')

- base64 encoding: 'hello world'.encode('base64')
- base64 decoding: 'hello world'.decode('base64')

- ascii to binary: bin(int('hello'.encode('hex'),16))
- binary to ascii

- binary to int: `int('0b110100001100101011011000110110001101111', 2)`
- int to binary: `bin(123123435234)`

- int to hex: `hex(23423434)`
- int to oct: `oct(23423434)`

note: the `'bla bla'.encode|decode` pattern is removed in python 3

if on unix, can use built-in `base64` like so: `echo "hello world" | base64` for encoding, and `echo "aGVsbG8gd29ybGQK" | base64 -D` for decoding

6236343a20615735305a584a755a58526659323975646d567963326c76626c3930623239736331397962324e72

- python playground online: https://repl.it


<p class="text-center">***</p>

*Outline:* (all in python 2)

1. decode hex: codecs.decode("707974686f6e2d666f72756d2e696f", "hex")
2. ascii to hex: binascii.a2b_hex (ascii to binary)
3. how to identify base64 string?
4. how to identify hex string?

codecs.decode(0x11, "binary")

codecs.decode("110110", "bit")


crypto challenge: https://github.com/VulnHub/ctf-writeups/blob/master/2016/angstrom-ctf/what-the-hex.md

---

6557393149475a766457356b4947526c4948646c61513d

- 'flag{YoU_h4v3_f0und_de_w31}'[::-1] ← reverse string

ctf: NjY1NDQ1N2E2NDMxMzk2YzVhNDYzOTZiNjI2ZTU1Nzc1YTZjMzg3YTY0NmE1MjZmNTgzMTU2NzY1\nNzU4NzQ2ZTU5NTc3ODZkMGE

https://github.com/VulnHub/ctf-writeups/blob/master/2016/angstrom-ctf/supersecure.md


---
hex to binary: `bin(int('hello'.encode('hex'), 16))`
binary to hex: `hex(0b110100001100101011011000110110001101111)[2:].decode('hex')`


octal to int (or int from octal): `int('211212323', 8)`
int to octal: `oct(121241)`

octal to string: `hex(int('06414533066157', 8))[2:].decode('hex')`

int from hex: `int(0xfffabcdde)` ← think of "casting" to int (hex to int)
int from bin: `int(0b111010101100101010)`

hex from octal: `hex(023234)` ← think of "casting" to hex (octal to hex)
bin from int: `bin(234234234234234)` ← think of "casting" to binary (int to bin)