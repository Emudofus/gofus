package login

type server struct {
}

// Create a new login server
func NewServer() *server {
	return &server{}
}

// Start login server without blocking
func (server *server) Start() error {
	return nil
}

// Stop login server without blocking
func (server *server) Stop() error {
	return nil
}
