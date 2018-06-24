package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func valid(a uint32) bool {
	return a >= 1<<24 && a < 1<<32-1
}

func o2d(a string) (ret uint32) {
	for i := 0; i < len(a); i++ {
		ret = ret<<3 + uint32(a[i]-'0')
	}
	return ret
}

func main() {
	dottedDecimal := regexp.MustCompile(`[12]?\d?\d\.[12]?\d?\d\.[12]?\d?\d\.[12]?\d?\d`)
	dottedHex := regexp.MustCompile(`0x[0-9a-fA-F]?[0-9a-fA-F]\.0x[0-9a-fA-F]?[0-9a-fA-F]\.0x[0-9a-fA-F]?[0-9a-fA-F]\.0x[0-9a-fA-F]?[0-9a-fA-F]`)
	dottedOctal := regexp.MustCompile(`0[0-3][0-7][0-7]\.0[0-3][0-7][0-7]\.0[0-3][0-7][0-7]\.0[0-3][0-7][0-7]`)
	dottedBinary := regexp.MustCompile(`[01]{8}\.[01]{8}\.[01]{8}\.[01]{8}`)
	binary := regexp.MustCompile(`[01]{32}`)
	octal := regexp.MustCompile(`0[0-3][0-7]{10}`)
	hexadecimal := regexp.MustCompile(`0x[0-9a-fA-F]{8}`)
	decimal := regexp.MustCompile(`[1-4]?\d{7,9}`)

	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	var (
		ips  []uint32
		ip   uint32
		nmax int
	)
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := scanner.Text()
		if t := dottedDecimal.FindAllString(s, -1); t != nil {
			for _, i := range t {
				var a, b, c, d, v uint32
				fmt.Sscanf(i, "%d.%d.%d.%d", &a, &b, &c, &d)
				v = a<<24 + b<<16 + c<<8 + d
				if valid(v) {
					ips = append(ips, v)
				}
			}
		}
		if t := dottedHex.FindAllString(s, -1); t != nil {
			for _, i := range t {
				var a, b, c, d, v uint32
				fmt.Sscanf(i, "0x%x.0x%x.0x%x.0x%x", &a, &b, &c, &d)
				v = a<<24 + b<<16 + c<<8 + d
				if valid(v) {
					ips = append(ips, v)
				}
			}
		}
		if t := dottedOctal.FindAllString(s, -1); t != nil {
			for _, i := range t {
				var v uint32
				vs := strings.Split(i, ".")
				v = o2d(vs[0])<<24 + o2d(vs[1])<<16 + o2d(vs[2])<<8 + o2d(vs[3])
				if valid(v) {
					ips = append(ips, v)
				}
			}
		}
		if t := dottedBinary.FindAllString(s, -1); t != nil {
			for _, i := range t {
				var a, b, c, d, v uint32
				fmt.Sscanf(i, "%b.%b.%b.%b", &a, &b, &c, &d)
				v = a<<24 + b<<16 + c<<8 + d
				if valid(v) {
					ips = append(ips, v)
				}
			}
		}
		if t := binary.FindAllString(s, -1); t != nil {
			for _, i := range t {
				var v uint32
				fmt.Sscanf(i, "%b", &v)
				if valid(v) {
					ips = append(ips, v)
				}
			}
		}
		if t := octal.FindAllString(s, -1); t != nil {
			for _, i := range t {
				v := o2d(i)
				if valid(v) {
					ips = append(ips, v)
				}
			}
		}
		if t := hexadecimal.FindAllString(s, -1); t != nil {
			for _, i := range t {
				var v uint32
				fmt.Sscanf(i, "0x%x", &v)
				if valid(v) {
					ips = append(ips, v)
				}
			}
		}
		if t := decimal.FindAllString(s, -1); t != nil {
			for _, i := range t {
				var v uint32
				fmt.Sscanf(i, "%d", &v)
				if valid(v) {
					ips = append(ips, v)
				}
			}
		}
	}
	nips := make(map[uint32]int)
	for _, i := range ips {
		if _, f := nips[i]; f {
			nips[i]++
		} else {
			nips[i] = 1
		}
		if nips[i] > nmax {
			nmax = nips[i]
			ip = i
		}
	}

	fmt.Printf("%d.%d.%d.%d\n", ip>>24&255, ip>>16&255, ip>>8&255, ip&255)
}
