package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/thiagozs/go-mysql-audit/proxy"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	px := proxy.NewProxy(ctx, "127.0.0.1", ":3306", true)
	px.EnableDecoding()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Signal received %v, stopping and exiting...", sig)
			cancel()
		}
	}()

	err := px.Start("8888")
	if err != nil {
		log.Fatal(err)
	}
}
