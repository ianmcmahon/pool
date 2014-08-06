/*

   (c) 2014 Ian McMahon

   This code is based on the connection pooling example by Ryan Day, detailed here:
   http://www.ryanday.net/2012/09/12/golang-using-channels-for-a-connection-pool/

 */

package pool

type InitFunction func() (interface{}, error)

type ConnectionPoolWrapper struct {
	size int
	conn chan interface{}
}

/** 
	Call the init function 'size' times.  If the init function fails during any call, then 
	the creation of the pool is considered af ailure.  We don't return size because a nil
	return value indicates 'size' connections were successfully created.

	We call the init function 'size' times to make sure each connection shares the same
	state.  The init function should set defaults such as character encoding, timezone, 
	anything that needs to be the same in each connection.
*/
func (p *ConnectionPoolWrapper) InitPool(size int, initfn InitFunction) error {
	// create a buffered channel allowing 'size' senders
	p.conn = make(chan interface{}, size)
	for x := 0; x < size; x++ {
		conn, err := initfn()
		if err != nil {
			return err
		}

		p.conn <- conn
	}
	p.size = size

	return nil
}

/**
	Ask for a connection interface from our channel.  If there are no connections
	available, we block until a connection is ready.
*/
func (p *ConnectionPoolWrapper) GetConnection() interface{} {
	return <-p.conn
}

/**
	Return a connection we have used to the pool
*/
func (p *ConnectionPoolWrapper) ReleaseConnection(conn interface{}) {
	p.conn <- conn
}
