package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nats-io/nats.go"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

var (
	NATSClient *nats.Conn
	Policy = bluemonday.UGCPolicy()
)

func main() {
	nc, err := nats.Connect(os.Getenv("NATS_HOST"))
	if err != nil {
		panic(err)
	}
	NATSClient = nc

	router := fasthttprouter.New()
	HTTPInit(router)
	FSInit()

	log.Print("Always listening.")
	log.Fatal(fasthttp.ListenAndServe("0.0.0.0:8000", router.Handler))
}
