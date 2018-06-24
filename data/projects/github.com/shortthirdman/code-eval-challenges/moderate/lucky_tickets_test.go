package main

import (
	"fmt"
	"math/big"
	"testing"
)

func TestLucky(t *testing.T) {
	for k, v := range map[int64]string{
		2: "10", 4: "670", 6: "55252", 8: "4816030",
		100: "138681178063913146486663255108385891670476531416644888545033078503482282975641730091720919340564340"} {
		if r := fmt.Sprint(lucky(k)); r != v {
			t.Errorf("failed: lucky %d is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkLucky(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lucky(int64(i%99 + 2))
	}
}

var r, b, c *big.Int

func init() {
	r, b, c = new(big.Int), new(big.Int), new(big.Int)
}

func lucky(m int64) *big.Int {
	a := true
	r.SetInt64(0)
	for i := int64(0); i < m/2; i++ {
		b.Binomial(m, i)
		c.Binomial(11*m/2-1-10*i, m-1)
		if a {
			r.Add(r, b.Mul(b, c))
		} else {
			r.Sub(r, b.Mul(b, c))
		}
		a = !a
	}
	return r
}
