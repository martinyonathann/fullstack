package main

import (
	"fmt"
	"log"
	"time"

	"github.com/martinyonathann/restapi_golang_postgres/api"
	"github.com/martinyonathann/restapi_golang_postgres/api/redisbroker"
)

type myMessage string

func main() {
	api.Run()
	var b redisbroker.Broker

	//Trying the in memory broker.
	b = redisbroker.NewMemoryBroker()

	//subCh is a readony channel that we will
	//receive message publisher on "ch1"

	subCh, err := b.Subscribe("ch1")
	if err != nil {
		log.Fatalln(err)
	}

	//start a publish loop
	//publish a message every second.

	go func() {
		defer b.Close()

		i := 0
		for {
			fmt.Println("masuk loop")
			i++
			if err := b.Publish("ch1", fmt.Sprintf("message %d", i)); err != nil {
				log.Fatalln(err)
			}

			time.Sleep(time.Second)

			// stop after 5 iterations
			if i == 5 {
				if err := b.Unsubscribe("ch1"); err != nil {
					log.Fatalln(err)
				}
				return
			}
		}
	}()

	//read message from subCh published on "ch1".
	for m := range subCh {
		fmt.Printf("got message : %s\n", m)
	}
}
