package n2xsdk

type PacketSetList struct {
 Handler string
}

func(np *PacketSetLis) Add () error {
 //parameters: Type
 //AgtPacketSetList Add
 return nil
}

func(np *PacketSetLis) AddItems () error {
 //parameters: Type NumberOfItems
 //AgtPacketSetList AddItems
 return nil
}

func(np *PacketSetLis) Remove () error {
 //parameters: Handle
 //AgtPacketSetList Remove
 return nil
}

func(np *PacketSetLis) Copy () error {
 //parameters: Handle
 //AgtPacketSetList Copy
 return nil
}

func(np *PacketSetLis) SetName () error {
 //parameters: Handle Name
 //AgtPacketSetList SetName
 return nil
}

func(np *PacketSetLis) GetName ()(string, error) {
 //parameters: Handle
 //AgtPacketSetList GetName
 return "", nil
}

func(np *PacketSetLis) GetHandle ()(string, error) {
 //parameters: Name
 //AgtPacketSetList GetHandle
 return "", nil
}

func(np *PacketSetLis) GetType ()(string, error) {
 //parameters: Handle
 //AgtPacketSetList GetType
 return "", nil
}

func(np *PacketSetLis) LockItem () error {
 //parameters: Handle
 //AgtPacketSetList LockItem
 return nil
}

func(np *PacketSetLis) UnlockItem () error {
 //parameters: Handle
 //AgtPacketSetList UnlockItem
 return nil
}

func(np *PacketSetLis) ListHandles ()(string, error) {
 //parameters: 
 //AgtPacketSetList ListHandles
 return "", nil
}

func(np *PacketSetLis) ListNames ()(string, error) {
 //parameters: 
 //AgtPacketSetList ListNames
 return "", nil
}

func(np *PacketSetLis) ListHandlesByType ()(string, error) {
 //parameters: Type
 //AgtPacketSetList ListHandlesByType
 return "", nil
}

func(np *PacketSetLis) ListTypes ()(string, error) {
 //parameters: 
 //AgtPacketSetList ListTypes
 return "", nil
}

func(np *PacketSetLis) GetInterface ()(string, error) {
 //parameters: Type
 //AgtPacketSetList GetInterface
 return "", nil
}

func(np *PacketSetLis) GetLockCount ()(string, error) {
 //parameters: 
 //AgtPacketSetList GetLockCount
 return "", nil
}

func(np *PacketSetLis) GetHandleByRequestID ()(string, error) {
 //parameters: PacketRequestID
 //AgtPacketSetList GetHandleByRequestID
 return "", nil
}

func(np *PacketSetLis) WaitForPacketSets () error {
 //parameters: NumRequests psaPacketRequestIDs
 //AgtPacketSetList WaitForPacketSets
 return nil
}

