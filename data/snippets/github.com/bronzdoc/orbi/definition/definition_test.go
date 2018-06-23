package definition_test

import (
	"fmt"
	"io/ioutil"
	"sort"

	. "github.com/bronzdoc/orbi/definition"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Definition", func() {
	var definition *Definition
	var resources []map[interface{}]interface{}

	BeforeEach(func() {
		resources = []map[interface{}]interface{}{
			{
				"dir": map[interface{}]interface{}{
					"name": "tmp-dir",
					"dir": map[interface{}]interface{}{
						"name": "tmp-dir-child",
					},
					"files": []interface{}{
						"tmp-file-child",
					},
				},
			},
		}

		map_definition := map[interface{}]interface{}{
			"context":   "./test-resource",
			"resources": resources,
		}

		options := map[string]interface{}{}

		definition = New(map_definition, options)
	})

	Context("When creating a new definition", func() {
		It("should create the correct definition from a yml file or a map", func() {
			definition_content := []byte(`---
context: ./test-resource
resources:
  - dir:
     name: tmp-dir
     dir:
      name: tmp-dir-child
     files:
      - tmp-file-child
 `)

			definition_file := "./tmp/definition.yml"
			ioutil.WriteFile(definition_file, definition_content, 0777)

			map_definition := definition

			var options map[string]interface{}
			file_definition := New(definition_file, options)

			result := equal(map_definition, file_definition)
			Expect(result).To(Equal(true))
		})
	})

	Describe("#Create", func() {
		It("should create the defined resources", func() {
			err := definition.Create()
			Expect(err).To(BeNil())

			definition.ResourceTree.Traverse(func(r Resource) {
				file_exists, _ := exists(r.ID())
				Expect(file_exists).To(Equal(true))
			})
		})

		Context("Definition context doesn't exists", func() {
			It("should return error", func() {
				map_definition := map[interface{}]interface{}{
					"context":   "./bad-context",
					"resources": resources,
				}

				options := map[string]interface{}{}

				definition = New(map_definition, options)

				err := definition.Create()
				Expect(err).ToNot(BeNil())
			})

		})
	})

	Describe("#Search", func() {
		It("should return the correct resource", func() {
			resource := definition.Search("tmp-dir/tmp-file-child")
			Expect(resource.Name()).To(Equal("tmp-file-child"))
		})
	})
})

func equal(a, b *Definition) bool {
	var a_names, b_names []string

	a.ResourceTree.Traverse(func(r Resource) {
		a_names = append(a_names, r.Name())
	})

	b.ResourceTree.Traverse(func(r Resource) {
		b_names = append(b_names, r.Name())
	})

	if a_names == nil && b_names == nil {
		return true
	}

	if a_names == nil || b_names == nil {
		return false
	}

	if len(a_names) != len(b_names) {
		return false
	}

	sort.Strings(a_names)
	sort.Strings(b_names)

	for i := range a_names {
		if a_names[i] != b_names[i] {
			fmt.Printf("%s not equal %s\n", a_names[i], b_names[i])
			return false
		}
	}

	return true
}
