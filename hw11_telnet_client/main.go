package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()

	args := os.Args[2:]
	if len(args) < 2 {
		fmt.Printf("Invalid arguments number: want - 2, got - %d", len(args))
		os.Exit(1)
	}

	address := net.JoinHostPort(args[0], args[1])
	tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := tc.Connect(); err != nil {
		fmt.Printf("%s: failed to connect: %v", address, err)
		os.Exit(1)
	}
	defer tc.Close()

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	exitCtx, exit := context.WithCancel(context.Background())
	defer exit()

	go func() {
		<-signalCtx.Done()
		exit()
	}()

	go func() {
		tc.Send()
		exit()
	}()

	go func() {
		tc.Receive()
		exit()
	}()

	<-exitCtx.Done()
}
