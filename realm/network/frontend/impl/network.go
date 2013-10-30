package impl

import (
	"bytes"
	"fmt"
	protocol "github.com/Blackrush/gofus/protocol/frontend"
	"github.com/Blackrush/gofus/protocol/frontend/msg"
	"github.com/Blackrush/gofus/realm/db"
	"github.com/Blackrush/gofus/realm/network/backend"
	"github.com/Blackrush/gofus/realm/network/frontend"
	"github.com/Blackrush/gofus/realm/network/frontend/handler"
	"github.com/Blackrush/gofus/shared"
	"log"
	"net"
)

const (
	input_msg_delimiter  = "\n\u0000"
	output_msg_delimiter = "\u0000"
	chunk_len            = 32
	queue_len            = 100
)

type context struct {
	config frontend.Configuration

	backend backend.Service
	players *db.Players

	running      bool
	nextClientId <-chan uint64
	tasks        chan task
	events       chan event
}

func New(backend backend.Service, players *db.Players, config frontend.Configuration) frontend.Service {
	return &context{
		config:  config,
		backend: backend,
		players: players,
		tasks:   make(chan task, queue_len),
		events:  make(chan event, queue_len),
	}
}

func (ctx *context) Config() frontend.Configuration {
	return ctx.config
}

func (ctx *context) Backend() backend.Service {
	return ctx.backend
}

func (ctx *context) Players() *db.Players {
	return ctx.players
}

func (ctx *context) Start() {
	if ctx.running {
		panic("frontend network service is already running")
	}
	ctx.running = true

	go client_id_generator(ctx)
	for i := 0; i < ctx.config.Workers; i++ {
		go worker_spawn(ctx)
	}
	go server_listen(ctx)

	log.Print("[frontend-net] successfully started")
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

	log.Printf("[frontend-net] listening on %d", ctx.config.Port)

	for ctx.running {
		conn, err := listener.Accept()
		if err != nil {
			panic(err.Error())
		}

		go server_conn_rcv(ctx, conn)
	}
}

func server_conn_rcv(ctx *context, conn net.Conn) {
	client := new_net_client(conn, <-ctx.nextClientId)
	defer server_conn_close(ctx, client)

	ctx.events <- event{client, true}

	buf := shared.Bufferize(conn, []byte(input_msg_delimiter), chunk_len)
	for ctx.running && client.alive {
		if data, ok := <-buf; ok {
			log.Printf("[frontend-net-client-%04d] RCV(%03d) %s", client.Id(), len(data), data)
			if msg, ok := msg.New(string(data[:2])); ok {
				in := bytes.NewReader(data)
				msg.Deserialize(in)

				ctx.tasks <- task{client, msg}
			} else {
				log.Printf("[frontend-net-client-%04d] unknown opcode %s", client.Id(), data[:2])
			}
		} else {
			break
		}
	}
}

func server_conn_close(ctx *context, client *net_client) {
	client.Close()

	ctx.events <- event{client, false}
}

type task struct {
	client *net_client
	data   protocol.MessageContainer
}

type event struct {
	client *net_client
	login  bool
}

func worker_spawn(ctx *context) {
	log.Print("[frontend-net-worker] spawned")

	for ctx.running {
		select {
		case task := <-ctx.tasks:
			handler.HandleClientData(ctx, task.client, task.data)
		case event := <-ctx.events:
			if event.login {
				handler.HandleClientConnection(ctx, event.client)
			} else {
				handler.HandleClientDisconnection(ctx, event.client)
			}
		}
	}
}
