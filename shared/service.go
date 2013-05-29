/*
  some common utilities
*/
package shared

import (
	"errors"
)

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
	Stopper
	// should be a blocking method
	Serve()
}

// this method will call Starter.Start() or Server.Serve() without blocking
func Start(any interface{}) error {
	switch service := any.(type) {
	case Starter:
		return service.Start()
	case Server:
		go service.Serve()
		return nil
	default:
		return errors.New("`any` is not a Starter nor Server")
	}
}

func Stop(any interface{}) error {
	if service, ok := any.(Stopper); ok {
		return service.Stop()
	} else {
		return errors.New("`any` is not a Stopper")
	}
}
