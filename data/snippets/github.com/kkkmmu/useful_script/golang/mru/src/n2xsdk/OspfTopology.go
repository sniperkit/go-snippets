package n2xsdk

type OspfTopology struct {
 Handler string
}

func(np *OspfTopology) Add () error {
 //parameters: DatabaseHandle Type
 //AgtOspfTopology Add
 return nil
}

func(np *OspfTopology) GetType ()(string, error) {
 //parameters: ObjectHandle
 //AgtOspfTopology GetType
 return "", nil
}

func(np *OspfTopology) Remove () error {
 //parameters: ObjectHandle
 //AgtOspfTopology Remove
 return nil
}

func(np *OspfTopology) RemoveAllObjects () error {
 //parameters: DatabaseHandle
 //AgtOspfTopology RemoveAllObjects
 return nil
}

func(np *OspfTopology) ListHandles ()(string, error) {
 //parameters: DatabaseHandle
 //AgtOspfTopology ListHandles
 return "", nil
}

func(np *OspfTopology) ListHandlesByType ()(string, error) {
 //parameters: DatabaseHandle Type
 //AgtOspfTopology ListHandlesByType
 return "", nil
}

func(np *OspfTopology) GetMaximumObjectConnections ()(string, error) {
 //parameters: 
 //AgtOspfTopology GetMaximumObjectConnections
 return "", nil
}

