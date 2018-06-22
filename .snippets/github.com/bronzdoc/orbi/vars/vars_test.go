package vars_test

import (
	"reflect"

	. "github.com/bronzdoc/orbi/vars"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vars", func() {
	Describe("Parse", func() {
		It("Should parse a string with vars correctly", func() {
			strvars := `
			ninja_turtle=donatello logan = wolverine
			gandalf= magneto
			cat-dog = nickelodeon
			homer.simpson=Doh!
			`

			expectedMap := map[string]string{
				"ninja_turtle":  "donatello",
				"logan":         "wolverine",
				"gandalf":       "magneto",
				"cat-dog":       "nickelodeon",
				"homer.simpson": "Doh!",
			}

			vars, err := Parse(strvars)
			Expect(err).To(BeNil())

			eql := reflect.DeepEqual(vars, expectedMap)
			Expect(eql).To(BeTrue())
		})
	})
})
