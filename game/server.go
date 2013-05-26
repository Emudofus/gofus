package game

type server struct {
}

// Create a new game server
func NewServer() *server {
	return &server{}
}

// Start game server without blocking
func (server *server) Start() error {
	return nil
}

// Stop game server without blocking
func (server *server) Stop() error {
	return nil
}
