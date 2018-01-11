package main

import (
	"net/http"
	"io/ioutil"
	"log"
)

var gPrivateKey []byte

func main() {
	var err error

	gPrivateKey, err = ioutil.ReadFile("key.pem")
	if err != nil{
		log.Println(err)
		return
	}

	tm := new(ticketMux)
	tm.RegisterRouter()

	log.Println("Jetbrain licence server start and listen port 8011")
	err = http.ListenAndServe(":8011", nil)
	if err != nil {
		panic(err)
	}
}
