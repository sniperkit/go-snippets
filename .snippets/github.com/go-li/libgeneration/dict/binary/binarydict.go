package binary

func Append(tree func(*) *[2]*, root **, key *, cmp func(*, *) int) * {
	var p *
	var q *

	var dir int

	if *root == nil {
		*root = key
		return key
	}

	for p = *root; p != nil; {
		var compare int = cmp(key, p)
		if compare == 0 {
			return p
		}
		if compare > 0 {
			dir = 1
		} else {
			dir = 0
		}

		q = p
		p = (*tree(p))[dir]
	}

	(*tree(q))[dir] = key

	return key

}

func Delete(tree func(*) *[2]*, root **, key *, cmp func(*, *) int) * {
	var p *
	var q *

	var dir int

	for p = *root; p != nil; {
		var compare int = cmp(key, p)
		if compare == 0 {
			if q == nil {
				q = *root
				if (*tree(q))[0] != nil {
					*root = (*tree(q))[0]
				} else {
					*root = (*tree(q))[1]
				}

				(*tree(q))[0] = nil
				(*tree(q))[1] = nil

				return p
			}

			if tree(p)[0] != nil && tree(p)[1] != nil {
				panic("the hard case of replacing this node with a successor")
			} else if tree(p)[0] != nil {
				(*tree(q))[dir] = tree(p)[0]
				tree(p)[0] = nil
			} else if tree(p)[1] != nil {
				(*tree(q))[dir] = tree(p)[1]
				tree(p)[1] = nil
			} else {
				(*tree(q))[dir] = nil
			}

			return p
		}
		if compare > 0 {
			dir = 1
		} else {
			dir = 0
		}

		q = p
		p = (*tree(p))[dir]
	}

	return nil

}

func Foreach(tree func(*) *[2]*, root *, apply func(*) bool) bool {
	if root == nil {
		return true
	}

	var paths = tree(root)

	if (paths[0] != nil) && !Foreach(tree, paths[0], apply) {
		return false
	}
	if !apply(root) {
		return false
	}
	if (paths[1] != nil) && !Foreach(tree, paths[1], apply) {
		return false
	}
	return true
}

func Find(tree func(*) *[2]*, root *, find func(*) byte) (f *) {
	if root == nil {
		return nil
	}
	f = root
	for {
		var children = tree(f)
		var dir byte
		if find != nil {
			dir = find(f)
		}
		if (*children)[dir] == nil {
			return f
		}
		f = (*children)[dir]
	}
}
