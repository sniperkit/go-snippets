package main

import (
	"strings"
	"testing"
)

func TestTelephone(t *testing.T) {
	for k, v := range map[string]string{
		"4155230": "g1jjad0,g1jjae0,g1jjaf0,g1jjbd0,g1jjbe0,g1jjbf0,g1jjcd0,g1jjce0,g1jjcf0,g1jkad0,g1jkae0,g1jkaf0,g1jkbd0,g1jkbe0,g1jkbf0,g1jkcd0,g1jkce0,g1jkcf0,g1jlad0,g1jlae0,g1jlaf0,g1jlbd0,g1jlbe0,g1jlbf0,g1jlcd0,g1jlce0,g1jlcf0,g1kjad0,g1kjae0,g1kjaf0,g1kjbd0,g1kjbe0,g1kjbf0,g1kjcd0,g1kjce0,g1kjcf0,g1kkad0,g1kkae0,g1kkaf0,g1kkbd0,g1kkbe0,g1kkbf0,g1kkcd0,g1kkce0,g1kkcf0,g1klad0,g1klae0,g1klaf0,g1klbd0,g1klbe0,g1klbf0,g1klcd0,g1klce0,g1klcf0,g1ljad0,g1ljae0,g1ljaf0,g1ljbd0,g1ljbe0,g1ljbf0,g1ljcd0,g1ljce0,g1ljcf0,g1lkad0,g1lkae0,g1lkaf0,g1lkbd0,g1lkbe0,g1lkbf0,g1lkcd0,g1lkce0,g1lkcf0,g1llad0,g1llae0,g1llaf0,g1llbd0,g1llbe0,g1llbf0,g1llcd0,g1llce0,g1llcf0,h1jjad0,h1jjae0,h1jjaf0,h1jjbd0,h1jjbe0,h1jjbf0,h1jjcd0,h1jjce0,h1jjcf0,h1jkad0,h1jkae0,h1jkaf0,h1jkbd0,h1jkbe0,h1jkbf0,h1jkcd0,h1jkce0,h1jkcf0,h1jlad0,h1jlae0,h1jlaf0,h1jlbd0,h1jlbe0,h1jlbf0,h1jlcd0,h1jlce0,h1jlcf0,h1kjad0,h1kjae0,h1kjaf0,h1kjbd0,h1kjbe0,h1kjbf0,h1kjcd0,h1kjce0,h1kjcf0,h1kkad0,h1kkae0,h1kkaf0,h1kkbd0,h1kkbe0,h1kkbf0,h1kkcd0,h1kkce0,h1kkcf0,h1klad0,h1klae0,h1klaf0,h1klbd0,h1klbe0,h1klbf0,h1klcd0,h1klce0,h1klcf0,h1ljad0,h1ljae0,h1ljaf0,h1ljbd0,h1ljbe0,h1ljbf0,h1ljcd0,h1ljce0,h1ljcf0,h1lkad0,h1lkae0,h1lkaf0,h1lkbd0,h1lkbe0,h1lkbf0,h1lkcd0,h1lkce0,h1lkcf0,h1llad0,h1llae0,h1llaf0,h1llbd0,h1llbe0,h1llbf0,h1llcd0,h1llce0,h1llcf0,i1jjad0,i1jjae0,i1jjaf0,i1jjbd0,i1jjbe0,i1jjbf0,i1jjcd0,i1jjce0,i1jjcf0,i1jkad0,i1jkae0,i1jkaf0,i1jkbd0,i1jkbe0,i1jkbf0,i1jkcd0,i1jkce0,i1jkcf0,i1jlad0,i1jlae0,i1jlaf0,i1jlbd0,i1jlbe0,i1jlbf0,i1jlcd0,i1jlce0,i1jlcf0,i1kjad0,i1kjae0,i1kjaf0,i1kjbd0,i1kjbe0,i1kjbf0,i1kjcd0,i1kjce0,i1kjcf0,i1kkad0,i1kkae0,i1kkaf0,i1kkbd0,i1kkbe0,i1kkbf0,i1kkcd0,i1kkce0,i1kkcf0,i1klad0,i1klae0,i1klaf0,i1klbd0,i1klbe0,i1klbf0,i1klcd0,i1klce0,i1klcf0,i1ljad0,i1ljae0,i1ljaf0,i1ljbd0,i1ljbe0,i1ljbf0,i1ljcd0,i1ljce0,i1ljcf0,i1lkad0,i1lkae0,i1lkaf0,i1lkbd0,i1lkbe0,i1lkbf0,i1lkcd0,i1lkce0,i1lkcf0,i1llad0,i1llae0,i1llaf0,i1llbd0,i1llbe0,i1llbf0,i1llcd0,i1llce0,i1llcf0",
		"4041124": "g0g11ag,g0g11ah,g0g11ai,g0g11bg,g0g11bh,g0g11bi,g0g11cg,g0g11ch,g0g11ci,g0h11ag,g0h11ah,g0h11ai,g0h11bg,g0h11bh,g0h11bi,g0h11cg,g0h11ch,g0h11ci,g0i11ag,g0i11ah,g0i11ai,g0i11bg,g0i11bh,g0i11bi,g0i11cg,g0i11ch,g0i11ci,h0g11ag,h0g11ah,h0g11ai,h0g11bg,h0g11bh,h0g11bi,h0g11cg,h0g11ch,h0g11ci,h0h11ag,h0h11ah,h0h11ai,h0h11bg,h0h11bh,h0h11bi,h0h11cg,h0h11ch,h0h11ci,h0i11ag,h0i11ah,h0i11ai,h0i11bg,h0i11bh,h0i11bi,h0i11cg,h0i11ch,h0i11ci,i0g11ag,i0g11ah,i0g11ai,i0g11bg,i0g11bh,i0g11bi,i0g11cg,i0g11ch,i0g11ci,i0h11ag,i0h11ah,i0h11ai,i0h11bg,i0h11bh,i0h11bi,i0h11cg,i0h11ch,i0h11ci,i0i11ag,i0i11ah,i0i11ai,i0i11bg,i0i11bh,i0i11bi,i0i11cg,i0i11ch,i0i11ci",
		"9102341": "w10adg1,w10adh1,w10adi1,w10aeg1,w10aeh1,w10aei1,w10afg1,w10afh1,w10afi1,w10bdg1,w10bdh1,w10bdi1,w10beg1,w10beh1,w10bei1,w10bfg1,w10bfh1,w10bfi1,w10cdg1,w10cdh1,w10cdi1,w10ceg1,w10ceh1,w10cei1,w10cfg1,w10cfh1,w10cfi1,x10adg1,x10adh1,x10adi1,x10aeg1,x10aeh1,x10aei1,x10afg1,x10afh1,x10afi1,x10bdg1,x10bdh1,x10bdi1,x10beg1,x10beh1,x10bei1,x10bfg1,x10bfh1,x10bfi1,x10cdg1,x10cdh1,x10cdi1,x10ceg1,x10ceh1,x10cei1,x10cfg1,x10cfh1,x10cfi1,y10adg1,y10adh1,y10adi1,y10aeg1,y10aeh1,y10aei1,y10afg1,y10afh1,y10afi1,y10bdg1,y10bdh1,y10bdi1,y10beg1,y10beh1,y10bei1,y10bfg1,y10bfh1,y10bfi1,y10cdg1,y10cdh1,y10cdi1,y10ceg1,y10ceh1,y10cei1,y10cfg1,y10cfh1,y10cfi1,z10adg1,z10adh1,z10adi1,z10aeg1,z10aeh1,z10aei1,z10afg1,z10afh1,z10afi1,z10bdg1,z10bdh1,z10bdi1,z10beg1,z10beh1,z10bei1,z10bfg1,z10bfh1,z10bfi1,z10cdg1,z10cdh1,z10cdi1,z10ceg1,z10ceh1,z10cei1,z10cfg1,z10cfh1,z10cfi1",
		"0000000": "0000000",
		"1111111": "1111111",
		"2000000": "a000000,b000000,c000000",
		"3000000": "d000000,e000000,f000000",
		"4000000": "g000000,h000000,i000000",
		"5000000": "j000000,k000000,l000000",
		"6000000": "m000000,n000000,o000000",
		"7000000": "p000000,q000000,r000000,s000000",
		"8000000": "t000000,u000000,v000000",
		"9000000": "w000000,x000000,y000000,z000000"} {
		if r := strings.Join(telephone(k, ""), ","); r != v {
			t.Errorf("failed: telephone %s is %s, got %s",
				k, v, r)
		}
	}
}

func telephone(s, t string) []string {
	if len(s) == 0 {
		return []string{t}
	}
	var r []string
	switch s[0] {
	case '0':
		r = telephone(s[1:len(s)], t+"0")
	case '1':
		r = telephone(s[1:len(s)], t+"1")
	case '2':
		r = append(telephone(s[1:len(s)], t+"a"),
			append(telephone(s[1:len(s)], t+"b"),
				telephone(s[1:len(s)], t+"c")...)...)
	case '3':
		r = append(telephone(s[1:len(s)], t+"d"),
			append(telephone(s[1:len(s)], t+"e"),
				telephone(s[1:len(s)], t+"f")...)...)
	case '4':
		r = append(telephone(s[1:len(s)], t+"g"),
			append(telephone(s[1:len(s)], t+"h"),
				telephone(s[1:len(s)], t+"i")...)...)
	case '5':
		r = append(telephone(s[1:len(s)], t+"j"),
			append(telephone(s[1:len(s)], t+"k"),
				telephone(s[1:len(s)], t+"l")...)...)
	case '6':
		r = append(telephone(s[1:len(s)], t+"m"),
			append(telephone(s[1:len(s)], t+"n"),
				telephone(s[1:len(s)], t+"o")...)...)
	case '7':
		r = append(telephone(s[1:len(s)], t+"p"),
			append(telephone(s[1:len(s)], t+"q"),
				append(telephone(s[1:len(s)], t+"r"),
					telephone(s[1:len(s)], t+"s")...)...)...)
	case '8':
		r = append(telephone(s[1:len(s)], t+"t"),
			append(telephone(s[1:len(s)], t+"u"),
				telephone(s[1:len(s)], t+"v")...)...)
	case '9':
		r = append(telephone(s[1:len(s)], t+"w"),
			append(telephone(s[1:len(s)], t+"x"),
				append(telephone(s[1:len(s)], t+"y"),
					telephone(s[1:len(s)], t+"z")...)...)...)
	}
	return r
}
