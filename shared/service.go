/*
  some common utilities
*/
package shared

// this interface represents an object that have to be started
type Starter interface {
	Start() error
}

// this interface represents an object that have to be stopped
type Stopper interface {
	Stop() error
}

// this interface represents an object that have to be started and stopped
type StartStopper interface {
	Starter
	Stopper
}

// this interface represents an action that will block a go-routine
type Server interface {
	// should be a blocking method
	Serve()
}
