package fs

import (
	"sort"

	"go.polydawn.net/rio/lib/treewalk"
)

type WalkFunc func(filenode *FilewalkNode) error

/*
	Walks a filesystem.

	This is much like the standard library's `path/filepath.Walk`,
	except it's based on `treewalk`, which means it supports both pre- and post-order traversals;
	and, if uses fs.RelPath (of course) to normalize path names.

	If walking directories, implicitly the first path will always be `./`;
	if the basePath is a file however, the first (and only) path with be `.`.
	This retains the same invarients from the perspective of the visit funcs
	(namely, that `filepath.Join(basePath, node.Path)` must be a correct path),
	but may also require additional understanding from the calling code to handle
	single files correctly.

	In order to get a name for the file in special case that basePath is a single
	file, use `node.Info.Name()`.

	Symlinks are not followed.

	The traversal order of siblings is *not* guaranteed, and is *not* necessarily
	stable.

	Caveat: calling `node.NextChild()` during your walk results in undefined behavior.
*/
func Walk(afs FS, preVisit WalkFunc, postVisit WalkFunc) error {
	return treewalk.Walk(
		newFileWalkNode(afs, RelPath{}),
		func(node treewalk.Node) error {
			filenode := node.(*FilewalkNode)
			if preVisit != nil {
				if err := preVisit(filenode); err != nil {
					return err
				}
			}
			return filenode.prepareChildren(afs)
		},
		func(node treewalk.Node) error {
			filenode := node.(*FilewalkNode)
			var err error
			if postVisit != nil {
				err = postVisit(filenode)
			}
			filenode.forgetChildren()
			return err
		},
	)
}

var _ treewalk.Node = &FilewalkNode{}

type FilewalkNode struct {
	Info *Metadata
	Err  error

	children []*FilewalkNode // note we didn't sort this
	itrIndex int             // next child offset
}

func (t *FilewalkNode) NextChild() treewalk.Node {
	if t.itrIndex >= len(t.children) {
		return nil
	}
	t.itrIndex++
	return t.children[t.itrIndex-1]
}

func newFileWalkNode(afs FS, path RelPath) (filenode *FilewalkNode) {
	// Fill in attributes.
	//  We could leave it to the user's code to do this, but when expanding
	//  the children list in the previsit func, we need to know if we're dealing
	//  with a dir, so, might as well keep and expose that info.
	filenode = &FilewalkNode{}
	filenode.Info, filenode.Err = afs.LStat(path)
	// We don't expand the children until the previsit function,
	//  because we don't want them all crashing into memory at once.
	return
}

/*
	Expand next subtree.  Used in the pre-order visit step so we don't walk
	every dir up front.  `Walk()` wraps the user-defined pre-visit function
	to do this at the end.
*/
func (t *FilewalkNode) prepareChildren(afs FS) error {
	if t.Info.Type != Type_Dir {
		return nil
	}
	names, err := afs.ReadDirNames(t.Info.Name)
	if err != nil {
		return err
	}
	sort.Strings(names)
	t.children = make([]*FilewalkNode, len(names))
	for i, name := range names {
		t.children[i] = newFileWalkNode(afs, t.Info.Name.Join(RelPath{name, -1}))
	}
	return nil
}

/*
	Used in the post-order visit step so we don't continuously consume more
	memory as we walk.  `Walk()` wraps the user-defined post-visit function
	to do this at the end.
*/
func (t *FilewalkNode) forgetChildren() {
	t.children = nil
}
