package color

import (
	"testing"
)

func TestWrapAndDye(t *testing.T) {
	t.Log(Dye(FgRed, "redString"))
	t.Log(Dye(Underline, Dye(FgYellow, "yellowBoldString")))

	redStringWrap := New(FgRed).Wrap("redString")
	redStringDye := Dye(FgRed, "redString")
	if redStringDye != redStringWrap {
		t.Errorf("effort of wrap and dye should be same, but get wrap: %s; dye: %s",
			redStringWrap, redStringDye)
	}
}
