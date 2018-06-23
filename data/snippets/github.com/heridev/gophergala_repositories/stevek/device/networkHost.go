package device

type networkHost interface {
	Filter(in <-chan message) <-chan message
}
