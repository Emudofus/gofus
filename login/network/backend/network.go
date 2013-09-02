package backend

import (
	digest "crypto/sha512"
	"database/sql"
	"fmt"
	"github.com/Blackrush/gofus/login/db"
	"github.com/Blackrush/gofus/protocol/backend"
	"github.com/Blackrush/gofus/shared"
	"hash"
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	chunk_len = 32
	salt_len  = 100
)

type Configuration struct {
	Port     uint16
	Password string
}

type Service interface {
	shared.StartStopper

	GetRealm(id int) (*Realm, bool)
	GetRealms() []*Realm
}

type context struct {
	config Configuration

	db    *sql.DB
	users *db.Users

	running        bool
	nextClientId   <-chan uint64
	nextClientSalt <-chan string

	realms map[int]*Realm
}

func New(database *sql.DB, config Configuration) Service {
	return &context{
		config: config,
		db:     database,
		users:  &db.Users{database},
		realms: make(map[int]*Realm),
	}
}

func (ctx *context) Start() {
	if ctx.running {
		panic("backend network service already running")
	}
	ctx.running = true

	go client_id_generator(ctx)
	go client_salt_generator(ctx)
	go server_listen(ctx)

	log.Print("[backend-net] successfully started")
}

func (ctx *context) Stop() {
	ctx.running = false
}

func (ctx *context) GetRealm(id int) (*Realm, bool) {
	realm, ok := ctx.realms[id]
	return realm, ok
}

// TODO perform some caching
func (ctx *context) GetRealms() []*Realm {
	realms := make([]*Realm, len(ctx.realms))
	i := 0
	for _, realm := range ctx.realms {
		realms[i] = realm
		i++
	}
	return realms
}

func shexdigest(digest hash.Hash, input string) []byte {
	return digest.Sum([]byte(input))
}

func (ctx *context) get_password_hash(salt string) []byte {
	d := digest.New()
	return shexdigest(d, fmt.Sprintf("%x%s", shexdigest(d, ctx.config.Password), salt))
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

func client_salt_generator(ctx *context) {
	c := make(chan string)
	defer close(c)

	src := rand.NewSource(time.Now().UnixNano())

	ctx.nextClientSalt = c
	for ctx.running {
		salt := shared.NextString(src, salt_len)
		c <- salt
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

func conn_rcv(conn net.Conn) (backend.Message, bool) {
	var opcode uint16

	if n, err := backend.Read(conn, &opcode); n <= 0 || err == io.EOF {
		return nil, false
	} else if err != nil {
		panic(err.Error())
	}

	if msg, ok := backend.NewMsg(opcode); ok {
		if err := msg.Deserialize(conn); err == io.EOF {
			return nil, false
		} else if err != nil {
			panic(err.Error())
		}

		return msg, true
	}

	log.Printf("[backend-net-client-%04d] unknown opcode %d", opcode)

	return nil, false
}

func server_conn_rcv(ctx *context, conn net.Conn) {
	client := &Client{
		WriteCloser: conn,
		id:          <-ctx.nextClientId,
		salt:        <-ctx.nextClientSalt,
		alive:       true,
	}
	defer server_conn_close(ctx, client)

	log.Printf("[backend-net-client-%04d] CONN", client.id)

	client_connection(ctx, client)

	for ctx.running && client.alive {
		if msg, ok := conn_rcv(conn); ok {
			log.Printf("[backend-net-client-%04d] RCV(%d)", client.id, msg.Opcode())

			client_handle_data(ctx, client, msg)
		} else {
			break
		}
	}
}

func server_conn_close(ctx *context, client *Client) {
	client.Close()

	client_disconnection(ctx, client)

	log.Printf("[backend-net-client-%04d] DCONN", client.id)
}

type Client struct {
	io.WriteCloser
	id    uint64
	salt  string
	alive bool
	realm *Realm
}

func (client *Client) Close() error {
	client.alive = false
	return client.WriteCloser.Close()
}

func (client *Client) Send(msg backend.Message) error {
	log.Printf("[backend-net-client-%04d] SND(%d)", client.id, msg.Opcode())

	backend.Put(client, msg.Opcode())
	return msg.Serialize(client)
}

func (client *Client) Id() uint64 {
	return client.id
}

func (client *Client) Salt() string {
	return client.salt
}

func (client *Client) Alive() bool {
	return client.alive
}

func (client *Client) Realm() *Realm {
	return client.realm
}
