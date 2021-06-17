package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/vmpool"
)

const (
	endpoint = "http://localhost:2633/RPC2"
	user     = "oneadmin"
	pass     = "opennebula"
)

func main() {
	c, err := api.NewClient(endpoint, user, pass, &http.Transport{}, 30*time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	var vmPool = &vmpool.VMPool{}
	err = vmPool.Info(c)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(vmPool)
	bla, err := vmPool.GetVMsByName("bla")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(bla)
}
