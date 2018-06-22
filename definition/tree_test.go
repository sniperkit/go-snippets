package definition_test

import (
	. "github.com/bronzdoc/orbi/definition"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {
	var tree *Tree

	BeforeEach(func() {
		resources := []map[interface{}]interface{}{
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

		tree = NewTree("./test-resource", resources)
	})

	Describe("#Root", func() {
		It("should return the correct resource", func() {
			Expect(tree.Root().Name()).To(Equal("./test-resource"))
		})
	})

	Describe("#Traverse", func() {
		It("should yield each resource to a function", func() {
			resource_counter := 0
			resouce_names := []string{
				"tmp-dir",
				"tmp-dir-child",
				"tmp-file-child",
			}

			tree.Traverse(func(r Resource) {
				func(r Resource) {
					for _, name := range resouce_names {
						if r.Name() == name {
							resource_counter += 1
						}
					}
				}(r)
			})

			Expect(resource_counter).To(Equal(3))
		})
	})
})
