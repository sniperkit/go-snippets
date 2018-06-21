package busybody

import (
	"fmt"
	"hash/crc32"
)

func crc32hash(input string) string {
	crchash := crc32.NewIEEE()
	crchash.Write([]byte(input))

	return fmt.Sprintf("%x", crchash.Sum(nil))
}
