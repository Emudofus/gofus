package network

import (
	"fmt"
	"github.com/Blackrush/gofus/shared"
	"log"
	"net"
)

const (
	chunk_len            = 64
	input_msg_delimiter  = "\n\u0000"
	output_msg_delimiter = "\u0000"
	events_queue_len     = 100
	works_queue_len      = 100
)

type Configuration struct {
	Port    uint16
	Workers int
}

type work struct {
	client Client
	data   []byte
}

type event struct {
	client Client
	login  bool
}

type context struct {
	config Configuration

	runnning     bool
	nextClientId chan uint64
	events       chan event
	works        chan work
}

func New(config Configuration) shared.StartStopper {
	return &context{
		config:       config,
		nextClientId: make(chan uint64),
		events:       make(chan event, events_queue_len),
		works:        make(chan work, works_queue_len),
	}
}

func client_id_generator(ctx *context) {
	defer close(ctx.nextClientId)
	var nextId uint64

	for ctx.runnning {
		nextId++
		ctx.nextClientId <- nextId
	}
}

func (ctx *context) Start() {
	if ctx.runnning {
		panic("network service is already running")
	}
	ctx.runnning = true

	go client_id_generator(ctx)
	for i := 0; i < ctx.config.Workers; i++ {
		go spawn_worker(ctx)
	}
	go start_server(ctx)

	log.Print("[network] successfully started")
}

func (ctx *context) Stop() {
	if !ctx.runnning {
		return
	}
	ctx.runnning = false
}

func start_server(ctx *context) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ctx.config.Port))

	if err != nil {
		panic(fmt.Sprintf("can't listen on %d because: %s", ctx.config.Port, err))
	}

	defer listener.Close()
	defer stop_server(ctx)

	log.Print("[network] listening on ", ctx.config.Port)

	for ctx.runnning {
		conn, err := listener.Accept()

		if err != nil {
			panic(fmt.Sprint("can't accept a connection because: ", err))
		}

		go handle_conn(ctx, conn)
	}
}

func stop_server(ctx *context) {
	close(ctx.events)
}

func handle_conn(ctx *context, conn net.Conn) {
	client := NewNetClient(conn, <-ctx.nextClientId)
	defer close_client(ctx, client)

	ctx.events <- event{client: client, login: true}

	buf := shared.Bufferize(conn, []byte(input_msg_delimiter), chunk_len)
	for ctx.runnning && client.Alive() {
		data := <-buf
		handle_client_data(ctx, client, data) // FIXME dispatch work (aegis won)
	}
}

func close_client(ctx *context, client Client) {
	client.Close()
	ctx.events <- event{client: client, login: false}
}

func spawn_worker(ctx *context) {
	log.Print("[network-worker] spawned")

	for ctx.runnning {
		select {
		case work := <-ctx.works:
			handle_client_data(ctx, work.client, work.data)
		case event := <-ctx.events:
			if event.login {
				handle_client_connection(ctx, event.client)
			} else {
				handle_client_disconnection(ctx, event.client)
			}
		}
	}
}
