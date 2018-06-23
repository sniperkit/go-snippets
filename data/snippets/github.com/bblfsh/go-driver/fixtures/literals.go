package fixtures

var (
	_ = nil

	_ = true
	_ = false

	_ = 42
	_ = 42.0
	_ = 42e0
	_ = 0x42
	_ = 042

	_ = "next\\\nline"
	_ = `next\
line`

	_ = 'a'
	_ = 'λ'

	_ = []byte{0}
	_ = [1]byte{0}
	_ = map[string]string{"foo": "bar"}
	_ = struct{ name string }{"foo"}
)
