package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/hostpool"
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
	var hostPool = &hostpool.HostPool{}
	err = hostPool.Info(c)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(hostPool)
	for _, host := range hostPool.Hosts {
		fmt.Printf("%s (%s)\n", host, host.State)
	}
}
