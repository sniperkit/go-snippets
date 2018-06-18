package statham_test

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/xoebus/statham"
)

func TestStatham(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Statham Suite")
}

type fakeRoundTripper struct {
	Calls int
}

func (rt *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.Calls++

	return nil, nil
}

var _ = Describe("A per-host transport mapping", func() {
	var (
		defaultTransport *fakeRoundTripper
		orgTransport     *fakeRoundTripper
		comTransport     *fakeRoundTripper

		transport http.RoundTripper
	)

	BeforeEach(func() {
		defaultTransport = &fakeRoundTripper{}
		orgTransport = &fakeRoundTripper{}
		comTransport = &fakeRoundTripper{}

		transport = statham.NewTransport(defaultTransport, statham.Mapping{
			"example.org": orgTransport,
			"example.com": comTransport,
		})
	})

	Context("when the host doesn't match any domain in the mapping", func() {
		It("uses the default transport", func() {
			req, _ := http.NewRequest("GET", "http://example.net", nil)

			transport.RoundTrip(req)

			Expect(defaultTransport.Calls).To(Equal(1))
			Expect(orgTransport.Calls).To(Equal(0))
			Expect(comTransport.Calls).To(Equal(0))
		})
	})

	Context("when the host matches a domain in the mapping", func() {
		It("uses the corresponding transport", func() {
			req, _ := http.NewRequest("GET", "http://example.com", nil)

			transport.RoundTrip(req)

			Expect(defaultTransport.Calls).To(Equal(0))
			Expect(orgTransport.Calls).To(Equal(0))
			Expect(comTransport.Calls).To(Equal(1))
		})

		It("uses the corresponding transport for another host", func() {
			req, _ := http.NewRequest("GET", "http://example.org", nil)

			transport.RoundTrip(req)

			Expect(defaultTransport.Calls).To(Equal(0))
			Expect(orgTransport.Calls).To(Equal(1))
			Expect(comTransport.Calls).To(Equal(0))
		})
	})
})
