package shared

// A Starter implements only one method 'Start'
type Starter interface {
	// Start should start a service without any error, only panic is acceptable
	Start()
}

// A Stopper implements only one method 'Stop'
type Stopper interface {
	// Stop should stop a service without any error, only panic is acceptable
	Stop()
}

// A StartStopper implements the two underlying interfaces Starter and Stopper
type StartStopper interface {
	Starter
	Stopper
}
