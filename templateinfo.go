package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/templatepool"
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

	var tplPool = &templatepool.TemplatePool{}
	err = tplPool.Info(c)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(tplPool)

//	bla, err := tplPool.GetTemplateByName("BLA")
//	if err != nil {
//		log.Fatalln(err)
//	}

//	fmt.Println(bla)
//	vmName := "blackbox-foo-123-west"
//	err = bla.Instantiate(c, vmName)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Printf("instantiated %s from template %s (%d)\n", vmName, bla, bla.ID)

}
