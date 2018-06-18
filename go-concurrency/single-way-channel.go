func (handler PersonHandlerImpl) Batch(origs <-chan Person) <- chan Person{
	dests := make(chan Person, 100)
	go func(){
		for{
			p,ok := <-origs
			if !ok {
				close(dests)
				break
			}
			handler.Handle(p)
			dests <- p
		}
	}()
	return dests
}