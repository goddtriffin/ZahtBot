package main

import (
	"context"
	"flag"
	"log"
	"os"
)

func main() {
	token := flag.String("token", "", "Discord Bot Token")
	flag.Parse()

	if *token == "" {
		flag.Usage()
		os.Exit(1)
	}

	zb, err := NewZahtBot(*token)
	if err != nil {
		log.Printf("New ZahtBot error")
		panic(err)
	}
	defer zb.StayConnectedUntilInterrupted(context.Background())
}
