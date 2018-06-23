package n2xsdk

type IsisTopology struct {
 Handler string
}

func(np *IsisTopology) Add () error {
 //parameters: DatabaseHandle Type
 //AgtIsisTopology Add
 return nil
}

func(np *IsisTopology) GetType ()(string, error) {
 //parameters: ObjectHandle
 //AgtIsisTopology GetType
 return "", nil
}

func(np *IsisTopology) Remove () error {
 //parameters: ObjectHandle
 //AgtIsisTopology Remove
 return nil
}

func(np *IsisTopology) ListHandles ()(string, error) {
 //parameters: DatabaseHandle
 //AgtIsisTopology ListHandles
 return "", nil
}

func(np *IsisTopology) ListHandlesByType ()(string, error) {
 //parameters: DatabaseHandle Type
 //AgtIsisTopology ListHandlesByType
 return "", nil
}

func(np *IsisTopology) IsTopologyChangedOnRestore () error {
 //parameters: DatabaseHandle
 //AgtIsisTopology IsTopologyChangedOnRestore
 return nil
}

