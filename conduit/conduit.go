package conduit

// Conduit is a channel of structs
type Conduit chan struct{}

func New() Conduit {
	return Conduit(make(chan struct{}))
}

// Either will return a channel which blocks on both this and that.
// If something is received from either it is sent down the conduit.
// Otherwise, Either will need a subsequent call to try again.
// Useful in a select statement.
func (c Conduit) Either(this, that <-chan struct{}) Conduit {
	go func() {
		select {
		case t, ok := <-this:
			if !ok {
				close(c)
				return
			}
			c <- t
		case t, ok := <-that:
			if !ok {
				close(c)
				return
			}
			c <- t
		}
	}()
	return c
}

// Trigger will return a channel which blocks until it receives
// a value from `pipe` channel UNLESS it receives on `trigger`
// channel, in which case it will close the conduit.
func (c Conduit) Trigger(pipe, trigger <-chan struct{}) Conduit {
	go func() {
		select {
		case <-trigger:
			close(c)
		case p := <-pipe:
			c <- p
		}
	}()
	return c
}
