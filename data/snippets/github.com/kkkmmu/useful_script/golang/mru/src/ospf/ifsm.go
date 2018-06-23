package ospf

type InterfaceState interface {
	InterfaceUP(oi *Interface) (*InterfaceState, error)
	WaitTimer(oi *Interface) (*InterfaceState, error)
	NeighborChange(oi *Interface) (*InterfaceState, error)
	LoopInd(oi *Interface) (*InterfaceState, error)
	UnloopInd(oi *Interface) (*InterfaceState, error)
	InterfaceDown(oi *Interface) (*InterfaceState, error)
	NextState(oi *Interface) (*InterfaceState, error)
}

type IFSM_Down struct {
}

func (ifsm *IFSM_Down) InterfaceUP(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Down) WaitTimer(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Down) NeighborChange(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Down) LoopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Down) UnloopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Down) InterfaceDown(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Down) NextState(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}

type IFSM_Loopback struct {
}

func (ifsm *IFSM_Loopback) InterfaceUP(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Loopback) WaitTimer(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Loopback) NeighborChange(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Loopback) LoopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Loopback) UnloopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Loopback) InterfaceDown(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Loopback) NextState(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}

type IFSM_Waiting struct {
}

func (ifsm *IFSM_Waiting) InterfaceUP(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Waiting) WaitTimer(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Waiting) NeighborChange(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Waiting) LoopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Waiting) UnloopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Waiting) InterfaceDown(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Waiting) NextState(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}

type IFSM_PointToPoint struct {
}

func (ifsm *IFSM_PointToPoint) InterfaceUP(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_PointToPoint) WaitTimer(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_PointToPoint) NeighborChange(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_PointToPoint) LoopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_PointToPoint) UnloopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_PointToPoint) InterfaceDown(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_PointToPoint) NextState(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}

type IFSM_DROther struct {
}

func (ifsm *IFSM_DROther) InterfaceUP(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DROther) WaitTimer(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DROther) NeighborChange(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DROther) LoopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DROther) UnloopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DROther) InterfaceDown(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DROther) NextState(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}

type IFSM_Backup struct {
}

func (ifsm *IFSM_Backup) InterfaceUP(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Backup) WaitTimer(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Backup) NeighborChange(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Backup) LoopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Backup) UnloopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Backup) InterfaceDown(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_Backup) NextState(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}

type IFSM_DR struct {
}

func (ifsm *IFSM_DR) InterfaceUP(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DR) WaitTimer(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DR) NeighborChange(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DR) LoopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DR) UnloopInd(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DR) InterfaceDown(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
func (ifsm *IFSM_DR) NextState(oi *Interface) (*InterfaceState, error) {
	return nil, nil
}
