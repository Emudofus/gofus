package backend

import (
	"database/sql"
	"fmt"
	"github.com/Blackrush/gofus/shared"
	"io"
	"log"
	"net"
)

const (
	chunk_len = 32
)

var (
	message_delimiter = []byte{0}
)

type Configuration struct {
	Port uint16
}

type context struct {
	config Configuration

	db *sql.DB

	running      bool
	nextClientId <-chan uint64
}

func New(database *sql.DB, config Configuration) shared.StartStopper {
	return &context{
		config: config,
		db:     database,
	}
}

func (ctx *context) Start() {
	if ctx.running {
		panic("backend network service already running")
	}
	ctx.running = true

	go client_id_generator(ctx)
	go server_listen(ctx)

	log.Print("[backend-net] successfully started")
}

func (ctx *context) Stop() {
	ctx.running = false
}

func client_id_generator(ctx *context) {
	c := make(chan uint64)
	defer close(c)

	var nextId uint64

	ctx.nextClientId = c
	for ctx.running {
		nextId++
		c <- nextId
	}
}

func server_listen(ctx *context) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ctx.config.Port))
	if err != nil {
		panic(err.Error())
	}

	defer listener.Close()

	log.Printf("[backend-net] listening on %d", ctx.config.Port)

	for ctx.running {
		conn, err := listener.Accept()
		if err != nil {
			panic(err.Error())
		}

		go server_conn_rcv(ctx, conn)
	}
}

func server_conn_rcv(ctx *context, conn net.Conn) {
	client := &Client{
		WriteCloser: conn,
		id:          <-ctx.nextClientId,
		alive:       true,
	}
	defer server_conn_close(ctx, client)

	log.Printf("[backend-net-client-%04d] CONN", client.id)

	buf := shared.Bufferize(conn, message_delimiter, chunk_len)
	for ctx.running && client.alive {
		if data, ok := <-buf; ok {
			log.Printf("[backend-net-client-%04d] RCV(%03d)", client.id, len(data))
			client_handle_data(ctx, client, data)
		} else {
			break
		}
	}
}

func server_conn_close(ctx *context, client *Client) {
	client.Close()

	log.Printf("[backend-net-client-%04d] DCONN", client.id)
}

type Client struct {
	io.WriteCloser
	id    uint64
	alive bool
}

func (client *Client) Write(b []byte) (int, error) {
	log.Printf("[backend-net-client-%04d] SND(%03d)", client.id, len(b))
	return client.WriteCloser.Write(b)
}

func (client *Client) Close() error {
	client.alive = false
	return client.WriteCloser.Close()
}

func (client *Client) Id() uint64 {
	return client.id
}

func (client *Client) Alive() bool {
	return client.alive
}
