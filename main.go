package main

import (
	"bytes"
	"fmt"
	topfew "github.com/timbray/topfew/internal"
	"golang.org/x/exp/trace"
	"log"
	"os"
	"time"
)

func main() {
	var err error

	config, err := topfew.Configure(os.Args[1:]) // skip whatever go puts in os.Args[0]
	if err != nil {
		fmt.Println("Problem (tf -h for help): " + err.Error())
		os.Exit(1)
	}
	fr := trace.NewFlightRecorder()
	fr.SetPeriod(time.Minute * 10)
	fr.SetSize((1024 * 1024 * 1024) * 2)
	fr.Start()

	counts, err := topfew.Run(config, os.Stdin)
	var b bytes.Buffer
	_, err = fr.WriteTo(&b)
	if err != nil {
		log.Print(err)
		return
	}
	// Write it to a file.
	if err := os.WriteFile("trace.out", b.Bytes(), 0o755); err != nil {
		log.Print(err)
		return
	}

	if err != nil {
		os.Exit(1)
	}
	for _, kc := range counts {
		fmt.Printf("%d %s\n", *kc.Count, kc.Key)
	}
}
