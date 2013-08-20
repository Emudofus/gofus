package network

import (
	"bytes"
	"net"
	"fmt"
	"github.com/Blackrush/gofus/protocol/msg"
	"github.com/Blackrush/gofus/shared"
	"io"
	"log"
	"math/rand"
	"time"
)

const (
	bufferLen = 64
	tasksQueueLen = 100
	eventsQueueLen = 100
	clientTicketLen = 32
	messageDelimiter = "\n\u0000"
	clientVersion = "1.29.1"
)

// TODO: Interface
type task struct {
	client Client
	data []byte
}

// TODO: Interface
type event struct {
	client Client
	login bool
}

type context struct {
	running bool
	tasks chan task
	events chan event
	nextClientId chan int
	nextClientTicket chan string
	clients map[int]Client

	config Configuration
}

type Configuration struct {
	Port uint16
	NbWorkers int
}

func New(config Configuration) shared.StartStopper {
	return &context{
		tasks: make(chan task, tasksQueueLen),
		events: make(chan event, eventsQueueLen),
		nextClientId: make(chan int),
		nextClientTicket: make(chan string),
		clients: make(map[int]Client),
		config: config,
	}
}

func (ctx *context) Start() {
	if ctx.running {
		// Should that panic? It could get over it or return an error
		log.Panic("network service already started")
	}
	ctx.running = true

	go client_id_generator(ctx)
	go client_ticket_generator(ctx)
	for i := 0; i < ctx.config.NbWorkers; i++ {
		log.Print("network service worker #", (i + 1), " started")
		go worker(ctx)
	}
	go start_server(ctx)

	log.Print("network service successfully started for Dofus ", clientVersion)
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
		// No need to stacktrace & shits, the error should be mentioned
		log.Fatal(fmt.Sprintf("can't listen on %d because: %s", ctx.config.Port, err.Error()))
	}

	defer listener.Close()
	defer stop_server(ctx)

	log.Print("network service listening on ", ctx.config.Port)

	for ctx.running {
		conn, err := listener.Accept()

		if err != nil {
			// No need to panic, the error concerns only new clients; just log and continue (might want to alert the admin though :<)
			log.Print(fmt.Sprintf("can't accept a connection on %d because: %s", ctx.config.Port, err.Error()))
			continue
		}

		go handle_conn(ctx, conn)
	}
}

func stop_server(ctx *context) {
	// nothing to do for now
}

func handle_conn(ctx *context, conn net.Conn) {
	var client Client = NewNetClient(conn, <-ctx.nextClientId, <-ctx.nextClientTicket)
	var chunk [bufferLen]byte
	var buffer []byte

	defer close_conn(ctx, client)

	ctx.events <- event { client: client, login: true }

	for ctx.running {
		n, err := conn.Read(chunk[0:])

		if n <= 0 || err == io.EOF { // no more data to read or end-of-file
			break
		}
		if err != nil {
			// Panic or Fatal? The error should be mentioned
			log.Panic(fmt.Sprintf("can't read data from %s because: %s", conn.RemoteAddr(), err.Error()))
		}

		received := chunk[:n]
		for len(received) > 0 {
			index := bytes.Index(received, []byte(messageDelimiter))
			if index < 0 {
				buffer = append(buffer, received...)
				log.Print("buffered ", len(received), " bytes (client #", client.Id(), ")")
				break
			}

			data := make([]byte, index + len(buffer))
			if len(buffer) > 0 {
				data = append(data, buffer...)
				buffer = nil
			}
			data = append(data, received[:index]...)

			ctx.tasks <- task { client, data }

			received = received[index+len(messageDelimiter):]
		}
	}
}

func close_conn(ctx *context, client Client) {
	client.Close()
	ctx.events <- event { client: client, login: false }
}

func worker(ctx *context) {
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

func handle_client_connection(ctx *context, client Client) {
	ctx.clients[client.Id()] = client
	log.Print("new client connected with id=", client.Id())

	client.Send(&msg.HelloConnect{ Ticket: client.Ticket() })
	client.SetState(VersionState)
}

func handle_client_disconnection(ctx *context, client Client) {
	delete(ctx.clients, client.Id())
	log.Print("client #", client.Id(), " disconnected")
}

func handle_client_data(ctx *context, client Client, data string) {
	log.Print("received ", len(data), " bytes `", data, "` from client #", client.Id())

	println("state=", client.State())

	switch client.State() {
	case VersionState:
		if clientVersion == data {
			client.SetState(LoginState)
		} else {
			client.Send(&msg.BadVersion{ Required: clientVersion })
			client.Close()
		}
	case LoginState:
	case RealmState:
	default:
		// No need to panic, it's only one client who got lost; just log and kick him out
		log.Print(fmt.Sprintf("unknown state %d of client #%d", client.State(), client.Id()))
		close_conn(ctx, client)
	}
}
