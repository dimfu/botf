package botf

type Balancer interface {
	// Returns all of the connections
	Get() Connection

	// Returns the connection from the balancer that are selected using round robin arrangement
	Connections() []Connection
}
