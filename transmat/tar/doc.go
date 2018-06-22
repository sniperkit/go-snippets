/*
	The tar transmat packs filesystems into the widely-recognized "tar" format,
	and can use any k/v-styled warehouse for storage.
*/
package tartrans

import (
	"go.polydawn.net/go-timeless-api"
)

const PackType = api.PackType("tar")
