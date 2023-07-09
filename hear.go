package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
)

//This is an async server with state store (cache - map[string]string)
//will not scale well. Could use a buffer/ or some sort of queue

/*cache with a mutex to prevent a data race (concurrent read and writes)
synthetic type*/
type cache struct {
	data map[string](string)
	*sync.RWMutex
}

//do we want to use a map as a data structure for caching
var c = cache{
	data:    make(map[string]string),
	RWMutex: &sync.RWMutex{},
}

//server only processes raw bytes - stay DRY
var InvalidCommand = []byte("Invalid Command")

//handler command
func handleCommand(inp string, conn net.Conn) {
	str := strings.Split(inp, " ")

	if len(str) <= 0 {
		conn.Write(InvalidCommand)
		return
	}

	command := str[0]

	//mux
	switch command {
	case "GET":
		get(str[1:], conn)
	case "SET":
		set(str[1:], conn)
	default:
		conn.Write(InvalidCommand)
	}
	//is this a stream => slice of bytes
	conn.Write([]byte("\n"))
}

/*handler per Connections
will take a connection, returns nothing*/
func handleConnection(conn net.Conn) {
	s := bufio.NewScanner(conn)
	for s.Scan() {
		data := s.Text()

		if data == " " {
			//is this a stream => slice of bytes
			conn.Write([]byte(">"))
			continue
		}
		if data == "exit" {
			return
		}
		handleCommand(data, conn)
	}
}

/*client command set, takes a slice of string and a connection to the server
will set a key value pair; returns nothihg */
func set(cmd []string, conn net.Conn) {
	if len(cmd) > 2 {
		conn.Write(InvalidCommand)
		return
	}
	key := cmd[0]
	val := cmd[1]
	//write locks
	c.Lock()
	c.data[key] = val
	c.Unlock()

	conn.Write([]byte("Ok"))
}

/*client command get, takes a slice of string and a connection to the server
will get a key/value pair */
func get(cmd []string, conn net.Conn) {
	if len(cmd) < 1 {
		conn.Write(InvalidCommand)
		return
	}

	val := cmd[0]
	//read locks
	c.RLock()
	ret, ok := c.data[val]
	c.RUnlock()

	if !ok {
		conn.Write([]byte("Nil"))
		return
	}
	conn.Write([]byte(ret))
}

func main() {
	//start listener
	listner, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	//handle connections
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Print("Accept Error", err)
			continue
		}

		log.Println("Accepted", conn.RemoteAddr())
		conn.Write([]byte(">"))

		//non-blocking call (mux call)
		go handleConnection(conn)
	}
}
