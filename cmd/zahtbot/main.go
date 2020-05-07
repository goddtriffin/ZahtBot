package main

import (
	"context"
	"flag"
	"os"
)

func main() {
	token := flag.String("token", "", "Discord Bot Token")
	flag.Parse()

	if *token == "" {
		flag.Usage()
		os.Exit(1)
	}

	defer NewZahtBot(*token).StayConnectedUntilInterrupted(context.Background())
}
