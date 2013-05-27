package login

import (
	"flag"
	"github.com/Blackrush/gofus/shared"
	sio "github.com/Blackrush/gofus/shared/io"
	"log"
	"net"
)

var (
	laddr   = flag.String("laddr", ":5555", "the address to listen on")
	workers = flag.Int("workers", 1, "the number of workers to start")

	// number of bytes allocated for each client to receive data
	chunkLen = 64
	// message delimiter
	delimiter = []byte("\u0000")
)

type connectMsg struct {
	conn net.Conn
	new  bool
}

type rcvMsg struct {
	conn net.Conn
	data []byte
}

type network struct {
	run     bool
	connect chan connectMsg
	rcv     chan rcvMsg
}

func new_network() *network {
	return &network{
		false,
		make(chan connectMsg),
		make(chan rcvMsg),
	}
}

// This method returns a new login server
func NewNetworkServer() shared.Server {
	return new_network()
}

// This method implements Starter and automatically call network.Serve()
func (network *network) Start() error {
	network.run = true

	for i := 0; i < *workers; i++ {
		go network.create_worker()
	}
	go network.Serve()

	return nil
}

// This method implements Stopper
func (network *network) Stop() error {
	network.run = false

	return nil
}

// This method implements Server
func (network *network) Serve() {
	listener, err := net.Listen("net", *laddr)

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for network.run {
		conn, err := listener.Accept()

		if err != nil {
			panic(err)
		}

		network.connect <- connectMsg{conn, true}
		go network.serve_client(conn)
	}
}

func (network *network) serve_client(conn net.Conn) {
	defer func() {
		conn.Close()
		network.connect <- connectMsg{conn, false}
	}()

	buf := sio.BufferLimit(conn, chunkLen, delimiter)

	for network.run {
		if data, ok := <-buf; ok {
			network.rcv <- rcvMsg{conn, data}
		} else {
			break
		}
	}
}

func (network *network) create_worker() {
	for network.run {
		select {
		case <-network.connect: // todo

		case msg := <-network.rcv:
			if err := handle_rcv(msg); err != nil {
				if *debug {
					log.Println(err)
				} else {
					msg.conn.Close()
				}
			}
		}
	}
}
