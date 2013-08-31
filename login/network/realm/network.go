package realm

import (
	"fmt"
	"github.com/Blackrush/gofus/shared"
	"log"
	"net"
)

const (
	output_msg_delimiter = "\u0000"
	input_msg_delimiter  = "\u0000"
	chunk_len            = 64
)

type Configuration struct {
	Port uint16
}

type context struct {
	config Configuration

	running      bool
	nextClientId chan uint64
}

func New(config Configuration) shared.StartStopper {
	return &context{
		config: config,
	}
}

func (ctx *context) Start() {
	if ctx.running {
		panic("realm network service already running")
	}
	ctx.running = true

	go client_id_generator(ctx)
	go server_listen(ctx)

	log.Print("[realm-net] successfully started")
}

func (ctx *context) Stop() {
	ctx.running = false
}

func client_id_generator(ctx *context) {
	defer close(ctx.nextClientId)
	ctx.nextClientId = make(chan uint64)

	var nextId uint64

	for ctx.running {
		nextId++
		ctx.nextClientId <- nextId
	}
}

func server_listen(ctx *context) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ctx.config.Port))
	if err != nil {
		panic(err.Error())
	}

	defer listener.Close()

	log.Printf("[realm-net] successfully listening on %d", ctx.config.Port)

	for ctx.running {
		conn, err := listener.Accept()
		if err != nil {
			panic(err.Error())
		}

		go server_conn_rcv(ctx, conn)
	}
}

type client struct {
	net.Conn

	id    uint64
	alive bool
}

func (client *client) Close() error {
	client.alive = false
	return client.Conn.Close()
}

func (client *client) Write(data []byte) (int, error) {
	log.Printf("[realm-net-client-%04d] SND(%03d) %s", client.id, len(data), data)
	data = append(data, []byte(output_msg_delimiter)...)
	return client.Conn.Write(data)
}

func server_conn_rcv(ctx *context, conn net.Conn) {
	client := &client{
		Conn:  conn,
		id:    <-ctx.nextClientId,
		alive: true,
	}
	defer server_conn_close(ctx, client)

	log.Printf("[realm-net-client-%04d] connected")

	buf := shared.Bufferize(conn, []byte(input_msg_delimiter), chunk_len)
	for ctx.running && client.alive {
		data := <-buf
		log.Printf("[realm-net-client-%04d] RCV(%03d) %s", client.id, len(data), data)
	}
}

func server_conn_close(ctx *context, client *client) {
	client.Close()

	log.Printf("[realm-net-client-%04d] disconnected")
}
