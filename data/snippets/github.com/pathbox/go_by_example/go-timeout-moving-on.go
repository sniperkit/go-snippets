timeout := make(chan bool, 1)
go func() {
	time.Sleep(1 * time.Second)
	timeout <- true
}

select {
	case <-ch:
		// a read from ch has occurred
	case <-timeout:
		// read from ch has timed out
}

func Query(conns []Conn, query string) Result {
	ch := make(chan Result, 1)
	for _, conn := range conns {
		go func(c Conn) {
			select {
			case ch <- c.DoQuery(query):
			default:
			}
		}(conn)
	}
	return <-ch
}