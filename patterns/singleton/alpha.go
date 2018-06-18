package main

type alpha struct{}

func (alpha) String() string {
	return "alpha"
}

var (
	alphaPtr *alpha
)

func Alpha() *alpha {
	if alphaPtr == nil {
		alphaPtr = &alpha{}
	}
	return alphaPtr
}
