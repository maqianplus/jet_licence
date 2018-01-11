package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"flag"
	"fmt"
)

var gPrivateKey []byte

func getPort() int {
	port := flag.Int("port", 8011, "Listen port")
	flag.Parse()
	return *port
}

func main() {
	port := getPort()

	var err error
	gPrivateKey, err = ioutil.ReadFile("key.pem")
	if err != nil {
		log.Println(err)
		return
	}

	tm := new(ticketMux)
	tm.RegisterRouter()

	log.Printf("Jetbrain licence server start and listen port %d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
