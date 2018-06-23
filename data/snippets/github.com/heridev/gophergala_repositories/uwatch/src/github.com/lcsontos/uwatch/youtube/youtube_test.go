package youtube

import (
	"fmt"
	"testing"

	"github.com/lcsontos/uwatch/catalog"
)

func TestSearchByID(t *testing.T) {
	service := newService(t)

	// Call with existing video ID

	service.doTestSearchByID("3M3iK_a-azM", false, t)

	// Call with existing video ID

	service.doTestSearchByID("a", true, t)
}

func TestSearchByTitle(t *testing.T) {

}

func (service *Service) doTestSearchByID(videoId string, wantNoSuchVideoError bool, t *testing.T) {
	videoRecord, err := service.SearchByID(videoId)

	nsve, ok := err.(*catalog.NoSuchVideoError)

	if !ok && err != nil {
		t.Fatal(err)
	}

	if (nsve != nil && !wantNoSuchVideoError) || (nsve == nil && wantNoSuchVideoError) {
		fmt.Printf(
			"err=%v, ok=%v, videoRecord=%v, wantNoSuchVideoError=%v\n",
			err, ok, videoRecord, wantNoSuchVideoError)

		t.Fatal(err)
	}
}

func newService(t *testing.T) *Service {
	service, err := New()

	if err != nil {
		t.Fatal(err)
	}

	return service
}
