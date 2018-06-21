package ospf

type NeighborState interface {
	HelloReceived(n *Neighbor) (*NeighborState, error)
	Start(n *Neighbor) (*NeighborState, error)
	TwoWayReceived(n *Neighbor) (*NeighborState, error)
	NegotiationDone(n *Neighbor) (*NeighborState, error)
	ExchangeDone(n *Neighbor) (*NeighborState, error)
	BadLSReq(n *Neighbor) (*NeighborState, error)
	LoadingDone(n *Neighbor) (*NeighborState, error)
	AdjOK(n *Neighbor) (*NeighborState, error)
	SeqNumberMismatch(n *Neighbor) (*NeighborState, error)
	OneWayReceived(n *Neighbor) (*NeighborState, error)
	KillNbr(n *Neighbor) (*NeighborState, error)
	InactivityTimer(n *Neighbor) (*NeighborState, error)
	LLDown(n *Neighbor) (*NeighborState, error)
}

type NFSM_Down struct {
}

func (nfsm *NFSM_Down) HelloReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) Start(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) TwoWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) NegotiationDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) ExchangeDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) BadLSReq(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) LoadingDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) AdjOK(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) SeqNumberMismatch(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) OneWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) KillNbr(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) InactivityTimer(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Down) LLDown(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}

type NFSM_Attempt struct {
}

func (nfsm *NFSM_Attempt) HelloReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) Start(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) TwoWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) NegotiationDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) ExchangeDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) BadLSReq(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) LoadingDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) AdjOK(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) SeqNumberMismatch(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) OneWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) KillNbr(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) InactivityTimer(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Attempt) LLDown(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}

type NFSM_Init struct {
}

func (nfsm *NFSM_Init) HelloReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) Start(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) TwoWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) NegotiationDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) ExchangeDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) BadLSReq(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) LoadingDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) AdjOK(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) SeqNumberMismatch(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) OneWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) KillNbr(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) InactivityTimer(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Init) LLDown(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}

type NFSM_2Way struct {
}

func (nfsm *NFSM_2Way) HelloReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) Start(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) TwoWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) NegotiationDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) ExchangeDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) BadLSReq(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) LoadingDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) AdjOK(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) SeqNumberMismatch(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) OneWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) KillNbr(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) InactivityTimer(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_2Way) LLDown(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}

type NFSM_Exstart struct {
}

func (nfsm *NFSM_Exstart) HelloReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) Start(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) TwoWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) NegotiationDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) ExchangeDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) BadLSReq(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) LoadingDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) AdjOK(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) SeqNumberMismatch(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) OneWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) KillNbr(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) InactivityTimer(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exstart) LLDown(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}

type NFSM_Exchange struct {
}

func (nfsm *NFSM_Exchange) HelloReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) Start(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) TwoWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) NegotiationDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) ExchangeDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) BadLSReq(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) LoadingDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) AdjOK(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) SeqNumberMismatch(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) OneWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) KillNbr(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) InactivityTimer(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Exchange) LLDown(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}

type NFSM_Loading struct {
}

func (nfsm *NFSM_Loading) HelloReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) Start(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) TwoWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) NegotiationDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) ExchangeDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) BadLSReq(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) LoadingDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) AdjOK(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) SeqNumberMismatch(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) OneWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) KillNbr(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) InactivityTimer(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Loading) LLDown(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}

type NFSM_Full struct {
}

func (nfsm *NFSM_Full) HelloReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) Start(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) TwoWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) NegotiationDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) ExchangeDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) BadLSReq(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) LoadingDone(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) AdjOK(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) SeqNumberMismatch(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) OneWayReceived(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) KillNbr(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) InactivityTimer(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
func (nfsm *NFSM_Full) LLDown(n *Neighbor) (*NeighborState, error) {
	return nil, nil
}
