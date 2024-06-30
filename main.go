package main

import (
	"fmt"
	"github.com/prusiksoftware/monorepo/obscurio/analytics"
	"github.com/prusiksoftware/monorepo/obscurio/http_server"
	"github.com/prusiksoftware/monorepo/obscurio/psql_proxy"
	"log"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	fmt.Println("starting proxy server")

	config, err := psql_proxy.GetConfig()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Println("config loaded")

	a := analytics.New(60)

	psqlserver, err := psql_proxy.NewServer(config, a)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Println("server created")

	httpserver := http_server.NewHTTPServer(a)
	go func() {
		httpserver.Start()
	}()

	go func(listeningChan <-chan psql_proxy.Status) {
		for {
			status := <-listeningChan
			fmt.Println("proxy status: ", status)
			if status == psql_proxy.Running {
				httpserver.SetLive(true)
				httpserver.SetReady(true)
			}
			if status == psql_proxy.Done {
				wg.Done()
				return
			}
		}
	}(psqlserver.StatusChan)

	psqlserver.Run()
	wg.Wait() // wait for the server to finish
}
