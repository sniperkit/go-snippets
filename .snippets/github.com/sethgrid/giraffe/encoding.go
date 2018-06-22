package giraffe

import (
	"bytes"
	"encoding/gob"
	"errors"
)

// GobEncode satisfies the gob encoder interface
func (g *Graph) GobEncode() ([]byte, error) {
	g.Lock()
	defer g.Unlock()

	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(g.Name)
	if err != nil {
		return nil, errors.New("1")
	}
	err = encoder.Encode(g.duplicateKeys)
	if err != nil {
		return nil, errors.New("2")
	}
	err = encoder.Encode(g.circularRelationship)
	if err != nil {
		return nil, errors.New("3")
	}
	err = encoder.Encode(g.Nodes)
	if err != nil {
		return nil, errors.New("4")
	}
	err = encoder.Encode(g.keys)
	if err != nil {
		return nil, errors.New("5")
	}
	err = encoder.Encode(g.topNodeID)
	if err != nil {
		return nil, errors.New("6")
	}
	return w.Bytes(), nil
}

// GobDecode satisfies the gob encoder interface
func (g *Graph) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&g.Name)
	if err != nil {
		return err
	}
	err = decoder.Decode(&g.duplicateKeys)
	if err != nil {
		return err
	}
	err = decoder.Decode(&g.circularRelationship)
	if err != nil {
		return err
	}
	err = decoder.Decode(&g.Nodes)
	if err != nil {
		return err
	}
	err = decoder.Decode(&g.keys)
	if err != nil {
		return err
	}
	err = decoder.Decode(&g.topNodeID)
	if err != nil {
		return err
	}

	// node.sources causes a locking issue
	// on encode/decode. to work around this,
	// we use node.sourceIDs and rebuild node.sources
	for _, node := range g.Nodes {
		for _, id := range node.sourceIDs {
			node.sources = append(node.sources, g.Nodes[id])
		}
		node.sourceIDs = make([]uint64, 0)
	}

	return nil
}

// GobEncode satisfies the gob encoder interface
func (n *Node) GobEncode() ([]byte, error) {
	n.Lock()
	defer n.Unlock()

	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(n.ID)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(n.Key)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(n.Value)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(n.circularRelationship)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(n.destinations)
	if err != nil {
		return nil, err
	}

	// poulate sourceIDs
	for _, n := range n.sources {
		n.sourceIDs = append(n.sourceIDs, n.ID)
	}

	err = encoder.Encode(n.sourceIDs)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// GobDecode satisfies the gob encoder interface
func (n *Node) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&n.ID)
	if err != nil {
		return err
	}
	err = decoder.Decode(&n.Key)
	if err != nil {
		return err
	}
	err = decoder.Decode(&n.Value)
	if err != nil {
		return err
	}
	err = decoder.Decode(&n.circularRelationship)
	if err != nil {
		return err
	}
	err = decoder.Decode(&n.destinations)
	if err != nil {
		return err
	}
	err = decoder.Decode(&n.sourceIDs)
	if err != nil {
		return err
	}
	return nil
}
