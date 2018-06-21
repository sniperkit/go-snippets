import (
	"github.com/Workiva/go-datastructures/trie/ctrie"
	"log"
)

type NodeDB struct {
	DB *ctrie.Ctrie
}

type Node struct {
	Key  string
	ID   string
	Type string
}

func New() *NodeDB {
	return &NodeDB{
		DB: ctrie.New(nil),
	}
}

func (ndb *NodeDB) Add(n *Node) error {
	if ndb.IsExist(n) {
		return errors.New("Same case alread exist")
	}

	ndb.Insert(ndb.Key(n), n)
}

func (ndb *NodeDB) Del(n *Node) error {
	if ndb.IsExist(n) {
		return errors.New("Same case alread exist")
	}

	ndb.Remove(ndb.Key(c))
}

func (ndb *NodeDB) Get(n *Node) (*Node, error) {
	if ndb.IsCaseValid(n) {
		if i, ok := ndb.Lookup(ndb.Key(n)); ok {
			if r, ok := i.(*Node); ok {
				return r, nil
			}
			log.Println("Invalid result when Get case from DB")
		}
	}

	return nil, errors.New("Cannot find case")
}

func (ndb *NodeDB) IsExist(n *Nodee) error {
	if ndb.isNodeValid(n) {
		if _, ok := ndb.Lookup([]byte(ndb.Key(n))); ok {
			return true
		}
	}

	return false
}

func (ndb *NodeDB) isNodeValid(n *Node) bool {
	if n.Group == "" {
		return false
	}

	return true
}

func (ndb *NodeDB) Key(n *Node) []byte {
	var key = ""
	if ndb.isNodeValid(n) {
		key += n.Group
		if n.SubGroup != "" {
			key += n.SubGroup
		} else if n.Feature != "" {
			key += n.Feature
		} else {
			return []byte(key)
		}

		if n.Name != "" {
			key += n.Name
		}
	}

	return []byte(key)
}

func (cdb *CaseDB) GetAll() <-chan *Case {
	ch := make(chan *Case)
	go func(in chan<- *Case) {
		for i := range cdb.Iterator(nil) {
			if c, ok := i.(*Case); ok {
				ch <- c
			}
		}
		defer close(ch)
	}(ch)

	return ch
}
