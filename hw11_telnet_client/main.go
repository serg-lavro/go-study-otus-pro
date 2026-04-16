package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("usage: go-telnet [--timeout=10s] host port]")
	}

	addr := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalf("...Connection failed: %v", err)
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", addr)
	defer client.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() { // stdin
		defer wg.Done()
		defer client.Close()

		if err := client.Send(); err == nil {
			fmt.Fprintln(os.Stderr, "...EOF")
		}
	}()

	go func() { // stdout
		defer wg.Done()
		defer client.Close()

		client.Receive()
		fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	}()

	wg.Wait()
}
