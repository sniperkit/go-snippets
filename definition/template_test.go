package definition_test

import (
	"io/ioutil"
	"os"

	. "github.com/bronzdoc/orbi/definition"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template", func() {
	var template *Template

	BeforeEach(func() {
		content := []byte("Antonio Marga-reeeeeiiiiiiiiti.")
		var vars map[string]string
		template = NewTemplate("template_0", content, vars)
	})

	Describe("#Execute", func() {
		It("should generate a template with the correct content", func() {
			file, _ := os.Create("./tmp/some-file")
			template.Execute(file)

			content, _ := ioutil.ReadFile("./tmp/some-file")
			Expect(string(content)).To(Equal("Antonio Marga-reeeeeiiiiiiiiti."))
		})

		Context("Content has vars", func() {
			content := []byte(`
				{{.var_1}}: Gor-la… Gor-lo-mi? Per cortesia, me lo ripeti ancora?
				{{.var_2}}: Gorlami!
				{{  .var_3  }}: It's not about the money, is about sending a message
				{{  .var_3  }}: Why so serious?
				test {
					bronz: {doc},
					{
				{} test {fsdsfsdsd}}
				. test {}}}
				`,
			)

			It("should fill variables in template correctly", func() {
				template := NewTemplate("template_1", content, map[string]string{
					"var_1": "Col. Hans Landa",
					"var_2": "Lt. Aldo Raine",
					"var_3": "The Joker",
				})

				file, err := os.Create("./tmp/some-other-file")
				Expect(err).To(BeNil())

				_, err = template.Execute(file)
				Expect(err).To(BeNil())

				expected_content := `
				Col. Hans Landa: Gor-la… Gor-lo-mi? Per cortesia, me lo ripeti ancora?
				Lt. Aldo Raine: Gorlami!
				The Joker: It's not about the money, is about sending a message
				The Joker: Why so serious?
				test {
					bronz: {doc},
					{
				{} test {fsdsfsdsd}}
				. test {}}}
				`

				actual_content, err := ioutil.ReadFile("./tmp/some-other-file")
				Expect(err).To(BeNil())

				Expect(string(actual_content)).To(Equal(expected_content))
			})

			It("should fail if no vars are passed", func() {
				template := NewTemplate("template_2", content, map[string]string{})
				file, _ := os.Create("./tmp/some-other-other-file")

				_, err := template.Execute(file)
				Expect(err).ToNot(Equal(nil))
				Expect(err.Error()).To(Equal("validateVars: template_2: missing vars:\n  {{.var_1}}\n  {{.var_2}}\n  {{.var_3}}\n"))
			})
		})
	})

	Describe("#Name", func() {
		It("should return the correct name", func() {
			Expect(template.Name()).To(Equal("template_0"))
		})
	})

	Describe("#Content", func() {
		It("should return the correct content", func() {
			Expect(template.Content()).To(Equal([]byte("Antonio Marga-reeeeeiiiiiiiiti.")))
		})
	})

	Describe("#Vars", func() {
		It("should return the correct vars", func() {
			var nil_vars map[string]string
			Expect(template.Vars()).To(Equal(nil_vars))
		})
	})
})
