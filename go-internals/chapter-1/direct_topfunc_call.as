"".add t=1 size=20 args=0x10 locals=0x0
	0x0000 00000 (direct_topfunc_call.go:4)	TEXT	"".add(SB), $0-16
	0x0000 00000 (direct_topfunc_call.go:4)	FUNCDATA	$0, gclocals·f207267fbf96a0178e8758c6e3e0ce28(SB)
	0x0000 00000 (direct_topfunc_call.go:4)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (direct_topfunc_call.go:4)	MOVL	"".b+12(FP), AX
	0x0004 00004 (direct_topfunc_call.go:4)	MOVL	"".a+8(FP), CX
	0x0008 00008 (direct_topfunc_call.go:4)	ADDL	CX, AX
	0x000a 00010 (direct_topfunc_call.go:4)	MOVL	AX, "".~r2+16(FP)
	0x000e 00014 (direct_topfunc_call.go:4)	MOVB	$1, "".~r3+20(FP)
	0x0013 00019 (direct_topfunc_call.go:4)	RET
	0x0000 8b 44 24 0c 8b 4c 24 08 01 c8 89 44 24 10 c6 44  .D$..L$....D$..D
	0x0010 24 14 01 c3                                      $...
"".main t=1 size=1 args=0x0 locals=0x0
	0x0000 00000 (direct_topfunc_call.go:6)	TEXT	"".main(SB), $0-0
	0x0000 00000 (direct_topfunc_call.go:6)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (direct_topfunc_call.go:6)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (direct_topfunc_call.go:8)	RET
	0x0000 c3                                               .
"".init t=1 size=79 args=0x0 locals=0x8
	0x0000 00000 (direct_topfunc_call.go:9)	TEXT	"".init(SB), $8-0
	0x0000 00000 (direct_topfunc_call.go:9)	MOVQ	(TLS), CX
	0x0009 00009 (direct_topfunc_call.go:9)	CMPQ	SP, 16(CX)
	0x000d 00013 (direct_topfunc_call.go:9)	JLS	72
	0x000f 00015 (direct_topfunc_call.go:9)	SUBQ	$8, SP
	0x0013 00019 (direct_topfunc_call.go:9)	MOVQ	BP, (SP)
	0x0017 00023 (direct_topfunc_call.go:9)	LEAQ	(SP), BP
	0x001b 00027 (direct_topfunc_call.go:9)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x001b 00027 (direct_topfunc_call.go:9)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x001b 00027 (direct_topfunc_call.go:9)	MOVBLZX	"".initdone·(SB), AX
	0x0022 00034 (direct_topfunc_call.go:9)	CMPB	AL, $1
	0x0024 00036 (direct_topfunc_call.go:9)	JLS	$0, 47
	0x0026 00038 (direct_topfunc_call.go:9)	MOVQ	(SP), BP
	0x002a 00042 (direct_topfunc_call.go:9)	ADDQ	$8, SP
	0x002e 00046 (direct_topfunc_call.go:9)	RET
	0x002f 00047 (direct_topfunc_call.go:9)	JNE	$0, 56
	0x0031 00049 (direct_topfunc_call.go:9)	PCDATA	$0, $0
	0x0031 00049 (direct_topfunc_call.go:9)	CALL	runtime.throwinit(SB)
	0x0036 00054 (direct_topfunc_call.go:9)	UNDEF
	0x0038 00056 (direct_topfunc_call.go:9)	MOVB	$2, "".initdone·(SB)
	0x003f 00063 (direct_topfunc_call.go:9)	MOVQ	(SP), BP
	0x0043 00067 (direct_topfunc_call.go:9)	ADDQ	$8, SP
	0x0047 00071 (direct_topfunc_call.go:9)	RET
	0x0048 00072 (direct_topfunc_call.go:9)	NOP
	0x0048 00072 (direct_topfunc_call.go:9)	PCDATA	$0, $-1
	0x0048 00072 (direct_topfunc_call.go:9)	CALL	runtime.morestack_noctxt(SB)
	0x004d 00077 (direct_topfunc_call.go:9)	JMP	0
	0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 39 48  eH..%....H;a.v9H
	0x0010 83 ec 08 48 89 2c 24 48 8d 2c 24 0f b6 05 00 00  ...H.,$H.,$.....
	0x0020 00 00 3c 01 76 09 48 8b 2c 24 48 83 c4 08 c3 75  ..<.v.H.,$H....u
	0x0030 07 e8 00 00 00 00 0f 0b c6 05 00 00 00 00 02 48  ...............H
	0x0040 8b 2c 24 48 83 c4 08 c3 e8 00 00 00 00 eb b1     .,$H...........
	rel 5+4 t=16 TLS+0
	rel 30+4 t=15 "".initdone·+0
	rel 50+4 t=8 runtime.throwinit+0
	rel 58+4 t=15 "".initdone·+-1
	rel 73+4 t=8 runtime.morestack_noctxt+0
gclocals·33cdeccccebe80329f1fdbee7f5874cb t=8 dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
gclocals·f207267fbf96a0178e8758c6e3e0ce28 t=8 dupok size=9
	0x0000 01 00 00 00 02 00 00 00 00                       .........
go.info."".add t=45 size=91
	0x0000 02 22 22 2e 61 64 64 00 00 00 00 00 00 00 00 00  ."".add.........
	0x0010 00 00 00 00 00 00 00 00 01 05 61 00 01 9c 00 00  ..........a.....
	0x0020 00 00 00 00 00 00 05 62 00 04 9c 11 04 22 00 00  .......b....."..
	0x0030 00 00 00 00 00 00 05 7e 72 32 00 04 9c 11 08 22  .......~r2....."
	0x0040 00 00 00 00 00 00 00 00 05 7e 72 33 00 04 9c 11  .........~r3....
	0x0050 0c 22 00 00 00 00 00 00 00 00 00                 .".........
	rel 8+8 t=1 "".add+0
	rel 16+8 t=1 "".add+20
	rel 30+8 t=28 go.info.int32+0
	rel 46+8 t=28 go.info.int32+0
	rel 64+8 t=28 go.info.int32+0
	rel 82+8 t=28 go.info.bool+0
go.info."".main t=45 size=27
	0x0000 02 22 22 2e 6d 61 69 6e 00 00 00 00 00 00 00 00  ."".main........
	0x0010 00 00 00 00 00 00 00 00 00 01 00                 ...........
	rel 9+8 t=1 "".main+0
	rel 17+8 t=1 "".main+1
go.info."".init t=45 size=27
	0x0000 02 22 22 2e 69 6e 69 74 00 00 00 00 00 00 00 00  ."".init........
	0x0010 00 00 00 00 00 00 00 00 00 01 00                 ...........
	rel 9+8 t=1 "".init+0
	rel 17+8 t=1 "".init+79
"".initdone· t=32 size=1
