package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	px "github.com/thiagozs/go-proxy-audit"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	proxy := px.NewProxy(ctx, "127.0.0.1", ":3306", true)
	proxy.EnableDecoding()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Signal received %v, stopping and exiting...", sig)
			cancel()
		}
	}()

	err := proxy.Start("8888")
	if err != nil {
		log.Fatal(err)
	}
}
