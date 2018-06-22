package definition

// Tree represents a definition tree structure
type Tree struct {
	root Resource
}

// NewTree Creates a new Tree
func NewTree(context string, definitionResources []map[interface{}]interface{}) *Tree {
	resource := &Directory{
		name:     context,
		id:       context,
		children: generate(context, definitionResources),
	}
	return &Tree{root: resource}
}

// Root gets a Tree root
func (t *Tree) Root() Resource {
	return t.root
}

// Traverse a tree and yield each node to a function
func (t *Tree) Traverse(action func(r Resource)) {
	traverse(t.root, action)
}

func traverse(r Resource, action func(r Resource)) {
	if r.Children() == nil {
		return
	}

	for _, node := range r.Children() {
		action(node)
		traverse(node, action)
	}
}

// Generates a definition Resource hierarchy
func generate(resourceID string, definitionResources []map[interface{}]interface{}) []Resource {
	var resources []Resource
	for _, resource := range definitionResources {
		for key, value := range resource {
			if key == "dir" {
				data := value.(map[interface{}]interface{})
				name := data["name"].(string)
				id := resourceID + "/" + name

				dirResources := []map[interface{}]interface{}{
					data,
				}

				resources = append(resources, &Directory{
					name:     name,
					id:       id,
					children: generate(id, dirResources),
				})

			} else if key == "files" {
				data := value.([]interface{})
				files := getFileResources(resourceID, filesStringify(data))
				resources = append(resources, files...)
			}
		}
	}
	return resources
}

// Convert a []string to []Resource
func getFileResources(resourceID string, fileNames []string) []Resource {
	var resources []Resource
	for _, file := range fileNames {
		id := resourceID + "/" + file
		resources = append(resources, &File{
			name: file,
			id:   id,
		})
	}
	return resources
}

// Convert a []interface{} to []string
func filesStringify(fileNames []interface{}) []string {
	var files []string
	for _, fileName := range fileNames {
		files = append(files, fileName.(string))
	}
	return files
}
