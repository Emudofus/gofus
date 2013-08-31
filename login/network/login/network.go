package login

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
	bufferLen           = 64
	tasksQueueLen       = 100
	eventsQueueLen      = 100
	clientTicketLen     = 32
	inMessageDelimiter  = "\n\u0000"
	outMessageDelimiter = "\u0000"
	clientVersion       = "1.29.1"
)

type context struct {
	running          bool
	tasks            chan task
	events           chan event
	nextClientId     chan int
	nextClientTicket chan string
	clients          map[int]Client

	db     *sql.DB
	users  db.Users
	config Configuration
}

type task struct {
	client Client
	data   []byte
}

type event struct {
	client Client
	login  bool
}

type Configuration struct {
	Port      uint16
	NbWorkers int
}

func New(database *sql.DB, config Configuration) shared.StartStopper {
	return &context{
		tasks:            make(chan task, tasksQueueLen),
		events:           make(chan event, eventsQueueLen),
		nextClientId:     make(chan int),
		nextClientTicket: make(chan string),
		clients:          make(map[int]Client),
		db:               database,
		users:            db.Users{database},
		config:           config,
	}
}

func (ctx *context) Start() {
	if ctx.running {
		panic("network service already started")
	}
	ctx.running = true

	go client_id_generator(ctx)
	go client_ticket_generator(ctx)
	for i := 0; i < ctx.config.NbWorkers; i++ {
		go spawn_worker(ctx)
	}
	go start_server(ctx)

	log.Print("[login-net] successfully started for Dofus ", clientVersion)
}

func (ctx *context) Stop() {
	if !ctx.running {
		return // just get over it, application is stopping, don't mess up other's stop
	}
	ctx.running = false
}

func client_id_generator(ctx *context) {
	defer close(ctx.nextClientId)

	var id int

	for ctx.running {
		id++
		ctx.nextClientId <- id
	}
}

func client_ticket_generator(ctx *context) {
	defer close(ctx.nextClientTicket)

	src := rand.NewSource(time.Now().UnixNano())

	for ctx.running {
		ticket := shared.NextString(src, clientTicketLen)
		ctx.nextClientTicket <- ticket
	}
}

func start_server(ctx *context) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ctx.config.Port))

	if err != nil {
		panic(fmt.Sprintf("can't listen on %d because: %s", ctx.config.Port, err.Error()))
	}

	defer listener.Close()
	defer stop_server(ctx)

	log.Print("[login-net] listening on ", ctx.config.Port)

	for ctx.running {
		conn, err := listener.Accept()

		if err != nil {
			panic(fmt.Sprintf("can't accept a connection on %d because: %s", ctx.config.Port, err.Error()))
		}

		go handle_conn(ctx, conn)
	}
}

func stop_server(ctx *context) {
	// nothing to do for now
}

func handle_conn(ctx *context, conn net.Conn) {
	client := NewNetClient(conn, <-ctx.nextClientId, <-ctx.nextClientTicket)
	defer close_conn(ctx, client)

	ctx.events <- event{client: client, login: true}

	buffer := shared.Bufferize(conn, []byte(inMessageDelimiter), bufferLen)
	for ctx.running && client.Alive() {
		data := <-buffer
		handle_client_data(ctx, client, string(data))
	}
}

func close_conn(ctx *context, client Client) {
	client.Close()
	ctx.events <- event{client: client, login: false}
}

func spawn_worker(ctx *context) {
	log.Print("[login-net-worker] spawned")

	for ctx.running {
		select {
		case task := <-ctx.tasks:
			handle_client_data(ctx, task.client, string(task.data))
		case event := <-ctx.events:
			if event.login {
				handle_client_connection(ctx, event.client)
			} else {
				handle_client_disconnection(ctx, event.client)
			}
		}
	}
}
