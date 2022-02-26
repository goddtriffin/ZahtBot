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
		log.Println("failed to initialize ZahtBot")
		panic(err)
	}
	defer func() {
		err = zb.StayConnectedUntilInterrupted(context.Background())
		if err != nil {
			log.Println(err)
		}
	}()
}
