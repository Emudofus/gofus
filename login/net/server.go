package net

import (
	"flag"
	sio "github.com/Nyasu/gofus/shared/io"
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

type baseMsg struct {
	net  *network
	conn net.Conn
}

func (msg *baseMsg) Write(b []byte) (int, error) {
	log.Printf("new outcoming data to %s (%d bytes)", msg.conn.RemoteAddr(), len(b))
	if msg.net.debug {
		println(string(b))
	}

	data := append(b, delimiter...)
	return msg.conn.Write(data)
}

type connectMsg struct {
	*baseMsg
	new bool
}

type rcvMsg struct {
	*baseMsg
	data []byte
}

type network struct {
	debug   bool
	run     bool
	connect chan connectMsg
	rcv     chan rcvMsg
}

// This method returns a new login server
func NewServer(debug bool) *network {
	return &network{
		debug,
		false,
		make(chan connectMsg),
		make(chan rcvMsg),
	}
}

// This method implements Starter and automatically call network.Serve()
func (network *network) Start() error {
	go network.Serve()

	for i := 0; i < *workers; i++ {
		go network.create_worker()
	}

	return nil
}

// This method implements Stopper
func (network *network) Stop() error {
	network.run = false

	return nil
}

// This method implements Server
func (network *network) Serve() {
	network.run = true

	listener, err := net.Listen("tcp", *laddr)

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for network.run {
		conn, err := listener.Accept()

		if err != nil {
			panic(err)
		}

		go network.serve_client(conn)
	}
}

func (network *network) serve_client(conn net.Conn) {
	log.Printf("new client %s", conn.RemoteAddr())

	msg := baseMsg{network, conn}
	buf := sio.BufferLimit(conn, chunkLen, delimiter)

	network.connect <- connectMsg{&msg, true}

	defer func() {
		conn.Close()
		network.connect <- connectMsg{&msg, false}
		log.Printf("client %s is gone", conn.RemoteAddr())
	}()

	for network.run {
		if data, ok := <-buf; ok {
			log.Printf("new incoming data from %s (%d bytes)", conn.RemoteAddr(), len(data))

			network.rcv <- rcvMsg{&msg, data}
		} else {
			break
		}
	}
}

func (network *network) create_worker() {
	for network.run {
		select {
		case msg := <-network.connect:
			if err := handle_connect(msg); err != nil {
				if network.debug {
					log.Println(err)
				} else if msg.new {
					msg.conn.Close()
				}
			}

		case msg := <-network.rcv:
			if err := handle_rcv(msg); err != nil {
				if network.debug {
					log.Println(err)
				} else {
					msg.conn.Close()
				}
			}
		}
	}
}
