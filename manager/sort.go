// Implementing sort.Interface to enable easy sorting of manager structs

package manager

////////////////// For sorting workers from fastest to slowest

// bySpeed implements sort.Interface
type bySpeed []*workerStats

func (bs bySpeed) Len() int {
	return len(bs)
}
func (bs bySpeed) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}
func (bs bySpeed) Less(i, j int) bool {
	// Faster ones should be first
	return bs[i].speed > bs[j].speed

}
