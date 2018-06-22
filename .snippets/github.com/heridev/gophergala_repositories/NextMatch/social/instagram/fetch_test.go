package instagram

import (
	"testing"
)

func TestByTag(t *testing.T) {
	o, err := ByTag(`bullsvsmavericks`)
	if err != nil {
		t.Error(`Failed to fetch recent photos by tag`, err)
	}

	if o.Meta.Code != 200 {
		t.Error(`Wrong status code`)
	}

	t.Log(`OK`)
}
