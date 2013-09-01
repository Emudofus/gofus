package frontend

import (
	"database/sql"
	"fmt"
	"github.com/Blackrush/gofus/login/db"
	"github.com/Blackrush/gofus/shared"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	input_msg_delimiter  = "\n\u0000"
	output_msg_delimiter = "\u0000"
	client_version       = "1.29.1"
	chunk_len            = 64
	ticket_len           = 32
	queue_len            = 100
)

type Configuration struct {
	Port    uint16
	Workers int
}

type context struct {
	config Configuration
	db     *sql.DB
	users  *db.Users

	running          bool
	nextClientId     <-chan uint64
	nextClientTicket <-chan string

	tasks  chan task
	events chan event
}

func New(database *sql.DB, config Configuration) shared.StartStopper {
	return &context{
		config: config,
		db:     database,
		users:  &db.Users{database},
		tasks:  make(chan task, queue_len),
		events: make(chan event, queue_len),
	}
}

func (ctx *context) Start() {
	if ctx.running {
		panic("frontend network already running")
	}
	ctx.running = true

	go client_id_generator(ctx)
	go client_ticket_generator(ctx)
	for i := 0; i < ctx.config.Workers; i++ {
		go worker_spawn(ctx)
	}
	go server_listen(ctx)

	log.Printf("[frontend-net] successfully started")
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

func client_ticket_generator(ctx *context) {
	c := make(chan string)
	defer close(c)

	src := rand.NewSource(time.Now().UnixNano())

	ctx.nextClientTicket = c
	for ctx.running {
		ticket := shared.NextString(src, ticket_len)
		c <- ticket
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
	client := NewNetClient(conn, <-ctx.nextClientId, <-ctx.nextClientTicket)
	defer server_conn_close(ctx, client)

	ctx.events <- event{client: client, login: true}

	buf := shared.Bufferize(conn, []byte(input_msg_delimiter), chunk_len)
	for ctx.running && client.Alive() {
		if data, ok := <-buf; ok {
			ctx.tasks <- task{client, data}
		} else {
			break
		}
	}
}

func server_conn_close(ctx *context, client Client) {
	client.Close()
	ctx.events <- event{client: client, login: false}
}

type task struct {
	client Client
	data   []byte
}

type event struct {
	client Client
	login  bool
}

func worker_spawn(ctx *context) {
	log.Print("[frontend-net-worker] spawned")
	for ctx.running {
		select {
		case task := <-ctx.tasks:
			client_handle_data(ctx, task.client, task.data)
		case event := <-ctx.events:
			if event.login {
				client_connection(ctx, event.client)
			} else {
				client_disconnection(ctx, event.client)
			}
		}
	}
}
