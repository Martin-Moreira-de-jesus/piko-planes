package main

import (
	"net"
	"os"
	"strings"
	"sync"
)

type SafeCounter struct {
	mu          sync.Mutex
	instruction []string
}

type ClientInfos struct {
	conn     *net.TCPConn
	servAddr string
}

var State = ClientInfos{
	conn: &net.TCPConn{},
}

var wg sync.WaitGroup

func (ci ClientInfos) Client() {
	InitConfig()
	ci.servAddr = Cfg.Client.Address
	tcpAddr, err := net.ResolveTCPAddr("tcp", ci.servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	State.conn, _ = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	wg.Add(1)
	go tcpRead()
	//go tcpWrite(conn) Don't need it here, old commit
	wg.Wait()
}

func SendButtonPressed(upKeyState string, downKeyState string) {
	message := "up=" + upKeyState + "," + "down=" + downKeyState + ","
	go tcpWrite(message)
}

func tcpRead() {

	c := SafeCounter{}
	defer wg.Done()
	for {
		received := make([]byte, 1024)
		n, err := State.conn.Read(received)
		if err != nil {
			println("Read from server failed:", err.Error())
			os.Exit(1)
		}
		c.Lock(strings.TrimRight(string(received[:n]), "\n"))
	}
}

func tcpWrite(message string) {
	println(message)
	_, err := State.conn.Write([]byte(message))
	println("done")
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
}
