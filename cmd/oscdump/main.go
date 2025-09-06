// oscdump provides a server that listens for OSC messages and dumps text
// representations of them to STDOUT, for testing

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hypebeast/go-osc/osc"
)

func fatalIfError(err error, msg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error ", msg, ": ", err)
		os.Exit(-1)
	}
}

func main() {
	listen := "localhost:54321"
	if len(os.Args) > 1 {
		listen = os.Args[1]
	}

	dispatcher := osc.NewStandardDispatcher()
	dispatcher.AddMsgHandler("*", func(msg *osc.Message) {
		fmt.Println(msg)
	})

	srv := osc.Server{
		Addr:       listen,
		Dispatcher: dispatcher,
	}

	log.Print("listening on ", listen)
	fatalIfError(srv.ListenAndServe(), "listening")
}
